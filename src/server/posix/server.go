package posix

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api "api/posix"
	auth "auth"
	ut "utils"
)

const (
	POSIX_FILE_SERVER_TEMPORARY_SUFFIX      = ".tmp"
	POSIX_FILE_SERVER_RESCAN_LIST_EVERY_MIN = 5
)

type PosixServer struct {
	sync.RWMutex
	ContentDir         string                    `json:"-"`
	List               []*posixFileInfo          `json:"list"`
	NeedRescanByLastOp bool                      `json:"need_rescan_by_last_op"`
	ListUpdatedAt      int64                     `json:"list_updated_at"`
	ChecksumCache      map[string]*posixCheckSum `json:"checksum_cache"`
	stateFile          string
	updateInProgress   bool
	forceRescanChan    chan bool
}

type posixFileInfo struct {
	Path string `json:"path"` // real path
	Size int64  `json:"size"`
}

type posixCheckSum struct {
	Size int64  `json:"size"`
	Md5  string `json:"md5"`
}

// создание нового сервера
func NewPosixServer(contentDir, stateFile string) (*PosixServer, error) {
	p := &PosixServer{
		stateFile:       stateFile,
		ContentDir:      contentDir,
		List:            make([]*posixFileInfo, 0),
		ChecksumCache:   make(map[string]*posixCheckSum, 0),
		forceRescanChan: make(chan bool),
	}
	// загружаем cache
	if err := os.MkdirAll(filepath.Dir(stateFile), 0750); err != nil {
		return nil, err
	}
	go p.rescanBackground() // в принципе её уже можно запустить
	log.Printf("[INFO] read state file: %s\n", p.stateFile)
	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return p, p.SaveState()
	}
	if err := json.Unmarshal(data, p); err != nil {
		log.Printf("[ERROR] load old state: %s\n", err.Error())
	}
	return p, p.SaveState()
}

// сохраняем статус нового сервера
func (p *PosixServer) SaveState() error {
	p.Lock()
	defer p.Unlock()
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("[FATAL] json marshal object posix server: %s\n", err.Error())
	}
	return ioutil.WriteFile(p.stateFile, data, 0640)
}

// рутина по фоновому обновлению списка файлов
func (p *PosixServer) rescanBackground() {

	rescan := func() {
		stat, err := os.Stat(p.ContentDir)
		if err != nil {
			log.Printf("[ERROR] state of content dir: %s\n", err.Error())
			return
		}
		if stat.ModTime().Unix() > p.ListUpdatedAt {
			if err := p.fetchList(); err != nil {
				log.Printf("[ERROR] scan: %s\n", err.Error())
			}
			p.ListUpdatedAt = time.Now().Unix()
			p.NeedRescanByLastOp = false
			p.SaveState()
		} else {
			log.Printf("[INFO] skip file list update, because modify time is older, when last update\n")
		}
	}

	log.Printf("[INFO] start background routine for upgrade list of file in dir: %s\n", p.ContentDir)
	rescan()
	ticker := time.NewTicker(POSIX_FILE_SERVER_RESCAN_LIST_EVERY_MIN * time.Minute)
	select {

	case <-ticker.C:
		p.updateInProgress = true
		rescan()
		p.updateInProgress = false

	case <-p.forceRescanChan:
		p.updateInProgress = true
		rescan()
		p.updateInProgress = false

	}
}

// собственно само обновление листа
func (p *PosixServer) fetchList() error {
	p.RLock()
	log.Printf("[INFO] start scan content dir\n")
	defer func() {
		p.RUnlock()
		log.Printf("[INFO] scan content dir complete\n")
	}()
	p.List = make([]*posixFileInfo, 0)
	return p.fetchFromDir(p.ContentDir)
}

// рекурсивная функция-хелпер по обновлению листа
func (p *PosixServer) fetchFromDir(dir string) error {
	dh, err := os.Open(dir)
	if err != nil {
		log.Printf("[ERROR] open directory `%s`: %s\n", dir, err.Error())
		return err
	}
	defer dh.Close()
	for {
		fileInfoList, err := dh.Readdir(10)
		if err == io.EOF {
			break
		}
		for _, fileInfo := range fileInfoList {
			path := filepath.Join(dir, fileInfo.Name())
			if fileInfo.IsDir() {
				if p.fetchFromDir(path); err != nil {
					return err
				}
			} else {
				p.List = append(p.List, &posixFileInfo{Path: path, Size: fileInfo.Size()})
			}
		}
	}
	return nil
}

func (p *PosixServer) getCachedMd5(filename string) (string, error) {
	state, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[INFO] md5 calculate, bad file %s: %s\n", filename, err)
			return "", err
		}
		log.Printf("[INFO] file state %s: not exists\n", filename)
		delete(p.ChecksumCache, filename)
		p.SaveState()
		return "", nil
	}
	if info, ok := p.ChecksumCache[filename]; ok {
		// информацию нашли в кэше, проверим её валидность
		if state.Size() == info.Size {
			return info.Md5, nil
		}
	}
	// значит пришло время вычислить cache
	log.Printf("[INFO] start calculate checksum of %s\n", filename)
	md5, err := ut.Md5(filename)
	if err != nil {
		log.Printf("[ERROR] calculate checksum of %s: %s\n", filename, err.Error())
		return "", err
	}
	log.Printf("[INFO] calculate checksum of %s completed\n", filename)
	p.ChecksumCache[filename] = &posixCheckSum{Md5: md5, Size: state.Size()}
	p.SaveState()
	return md5, nil
}

// реализация поиска
func (p *PosixServer) Find(ctx context.Context, filter *api.Filter) (*api.List, error) {
	log.Printf("[INFO] find request: %v\n", filter)
	result := &api.List{Files: make([]*api.Info, 0)}
	regMatch := filter.PathMatch
	if regMatch == "*" {
		regMatch = ".*"
	}
	reg, err := regexp.Compile(regMatch)
	if err != nil {
		result.State = api.List_ERROR
		err = fmt.Errorf("build regexp %s: %s", filter.PathMatch, err.Error())
		result.Error = err.Error()
		log.Printf("[ERROR] reply on find request: %s\n", err.Error())
		return result, nil
	}
	// если это файл, попробуем найти его это в лоб
	if md5, err := p.getCachedMd5(filepath.Join(p.ContentDir, filter.PathMatch)); err == nil {
		result.State = api.List_OK
		if stat, err := os.Stat(filepath.Join(p.ContentDir, filter.PathMatch)); err == nil {
			result.Files = append(result.Files, &api.Info{Path: filter.PathMatch, Md5: md5, Size: stat.Size()})
			log.Printf("[INFO] find completed (found PathMatch on disk), count %d\n", len(result.Files))
			return result, nil
		}
	}
	log.Printf("[INFO] find in cache list\n")
	// иначе, считаем что это regexp и пробуем искать его на диске
	if p.NeedRescanByLastOp {
		// если необходимо отсылаем инфу о том, что необходимо посканить
		go func() { p.forceRescanChan <- true }()
		time.Sleep(100 * time.Millisecond)
	}

	// ждем, если необходимо обновления
	for i := 0; i < 600; i++ {
		if !p.updateInProgress {
			break
		}
		log.Printf("[INFO] wait update...\n")
		time.Sleep(time.Second)
	}
	// не дождались обновления
	if p.updateInProgress {
		result.State = api.List_RETRY
		log.Printf("[INFO] reply on find request: RETRY (not ready)\n")
		return result, nil
	}

	for _, info := range p.List {
		if info == nil {
			continue
		}
		if strings.HasSuffix(info.Path, POSIX_FILE_SERVER_TEMPORARY_SUFFIX) {
			continue
		}
		virtualFileName := strings.TrimPrefix(info.Path, p.ContentDir)
		virtualFileName = strings.TrimPrefix(virtualFileName, string(filepath.Separator))
		if !reg.MatchString(virtualFileName) {
			continue
		}
		md5, err := p.getCachedMd5(info.Path)
		if err != nil {
			result.State = api.List_ERROR
			result.Error = fmt.Sprintf("File found, but calculate md5: %s", err.Error())
			log.Printf("[ERROR] reply on find request: %v\n", result)
			return result, nil
		}
		resultInfo := &api.Info{Path: virtualFileName, Md5: md5, Size: info.Size}
		result.Files = append(result.Files, resultInfo)
	}
	result.State = api.List_OK
	result.SnaphotTime = p.ListUpdatedAt
	log.Printf("[INFO] find completed, count %d\n", len(result.Files))
	return result, nil
}

// реализация получение файла на сервер
func (p *PosixServer) Receive(stream api.File_ReceiveServer) error {

	uploadErrorLogAndReport := func(stream api.File_ReceiveServer, format string, a ...interface{}) error {
		message := fmt.Sprintf(format, a...)
		log.Printf("[ERROR] %s\n", message)
		return stream.SendAndClose(&api.Info{State: api.Info_ERROR, Error: message})
	}

	writeChunk := func(fd *os.File, chunk *api.Chunk) error {
		fd.Seek(chunk.Offset, 0)
		count, err := fd.Write(chunk.Data)
		if err != nil {
			return err
		} else {
			if count != len(chunk.Data) {
				return fmt.Errorf("only %d bytes of the %d was written", count, len(chunk.Data))
			}
		}
		return nil
	}

	log.Printf("[INFO] start receive file\n")

	chunk, err := stream.Recv()
	if err != nil && err != io.EOF {
		return uploadErrorLogAndReport(stream, "receive first chunk: %s", err.Error())
	}
	fileInfo := chunk.File
	if fileInfo == nil {
		return uploadErrorLogAndReport(stream, "empty file information if first chunk")
	}
	if fileInfo.Md5 == "" {
		return uploadErrorLogAndReport(stream, "empty md5 information in first chunk")
	}
	if fileInfo.Path == "" {
		return uploadErrorLogAndReport(stream, "empty file path information in first chunk")
	}
	if fileInfo.Size == 0 {
		return uploadErrorLogAndReport(stream, "empty file size information in first chunk")
	}

	// открываем дескриптор
	path := filepath.Join(p.ContentDir, fileInfo.Path+POSIX_FILE_SERVER_TEMPORARY_SUFFIX)
	fd, err := os.Create(path)
	if err != nil {
		return uploadErrorLogAndReport(stream, "create file %s: %s", fileInfo.Path, err.Error())
	}
	defer fd.Close()
	defer os.Remove(path)

	if err := writeChunk(fd, chunk); err != nil {
		return uploadErrorLogAndReport(stream, "write chunk: %s", err.Error())
	}

	p.NeedRescanByLastOp = true
	result := &api.Info{Path: fileInfo.Path}

	chunkCount := 1
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			stat, _ := fd.Stat()
			if stat.Size() != fileInfo.Size {
				return uploadErrorLogAndReport(stream, "file received, but size expected %d, got %d", fileInfo.Size, stat.Size())
			}
			log.Printf("[INFO] file %s receive complete (%s)\n", path, humanize.Bytes(uint64(stat.Size())))
			fd.Seek(0, 0)
			md5, err := ut.Md5(path)
			if err != nil {
				return uploadErrorLogAndReport(stream, "calculate md5: %s", err.Error())
			}
			if md5 != fileInfo.Md5 {
				return uploadErrorLogAndReport(stream, "md5 expected %s, got %s", fileInfo.Md5, md5)
			}
			if err := fd.Close(); err != nil {
				return uploadErrorLogAndReport(stream, "close temp file: %s", err.Error())
			}
			if err := os.Rename(path, filepath.Join(p.ContentDir, fileInfo.Path)); err != nil {
				return uploadErrorLogAndReport(stream, "rename %s: %s", path, err.Error())
			}
			result.Md5 = md5
			result.State = api.Info_OK
			log.Printf("[INFO] receive file %s done\n", fileInfo.Path)
			return stream.SendAndClose(result)
		} // end EOF

		if err := writeChunk(fd, chunk); err != nil {
			log.Printf("[ERROR] write chunk: %s\n", err.Error())
			return err
		}

		// print info
		chunkCount++
		if chunkCount%50 == 0 {
			if stat, err := fd.Stat(); err == nil {
				log.Printf("[INFO] file %s: download %d chunks, size: %s\n", path, chunkCount, humanize.Bytes(uint64(stat.Size())))
			}
		}
	} // end for receive
}

// реализация стриминг файла с сервера
func (p *PosixServer) Stream(info *api.Info, stream api.File_StreamServer) error {

	log.Printf("[INFO] start stream file {%v}\n", info)

	fd, err := os.Open(filepath.Join(p.ContentDir, info.Path))
	if err != nil {
		log.Printf("[ERROR] open file: %s\n", err.Error())
		resultFile := &api.Info{Path: info.Path, State: api.Info_ERROR, Error: err.Error()}
		if send := stream.Send(&api.Chunk{File: resultFile}); send != nil {
			return send
		}
		return err
	}
	defer fd.Close()

	offset := int64(0)
	reader := bufio.NewReader(fd)
	buf := make([]byte, 0, POSIX_FILE_BUFFER_SIZE)
	for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Printf("[ERROR] read file: %s\n", err.Error())
			resultFile := &api.Info{Path: info.Path, State: api.Info_ERROR, Error: err.Error()}
			return stream.Send(&api.Chunk{File: resultFile})
		}
		chunk := &api.Chunk{Data: buf, Offset: offset, File: &api.Info{Path: info.Path, State: api.Info_OK}}
		offset += int64(len(buf))
		if err := stream.Send(chunk); err != nil {
			log.Printf("[ERROR] send chunk: %s\n", err.Error())
		}
	}

	log.Printf("[INFO] file %s (%s) stream complete\n", info.Path, humanize.Bytes(uint64(offset)))
	return nil
}

// реализация локальных операций
func (p *PosixServer) LocalOps(ctx context.Context, op *api.LocalOperation) (*api.LocalOperation, error) {

	copyFromTo := func(src, dst string) error {
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		os.MkdirAll(filepath.Dir(dst), 0750)
		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		return out.Sync()
	}

	log.Printf("[INFO] local operation: {%v}\n", op)

	// валидация
	if op.File == nil || op.File.Path == "" {
		log.Printf("[ERROR] empty file information\n")
		return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: "empty file information"}, nil
	}
	switch op.Type {
	case api.LocalOperation_COPY, api.LocalOperation_MOVE:
		if op.DstFile == nil || op.DstFile.Path == "" {
			log.Printf("[ERROR] empty destination file info\n")
			return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: "empty dst_file for operation LOCAL_{COPY,MOVE}"}, nil
		}
	case api.LocalOperation_DELETE:
		if op.DstFile != nil {
			log.Printf("[ERROR] not empty destination file info\n")
			return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: "not empty destination file info for operation LOCAL_DELETE"}, nil
		}
	default:
		log.Printf("[ERROR] unknown local operation: %v\n", op.Type)
		return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: fmt.Sprintf("unknown local operation: %v", op.Type)}, nil
	}

	// запоминаем что нам нужно обновить информацию
	p.NeedRescanByLastOp = true

	source := filepath.Join(p.ContentDir, op.File.Path)
	md5, err := p.getCachedMd5(source)
	if err != nil {
		log.Printf("[ERROR] calculate md5 %s: %s\n", source, err.Error())
		return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: fmt.Sprintf("calculate md5 %s: %s\n", source, err.Error())}, nil
	}
	if op.Type == api.LocalOperation_DELETE {
		os.MkdirAll(filepath.Dir(source), 0750)
		if err := os.Remove(source); err == nil {
			log.Printf("[INFO] file %s removed\n", source)
			return &api.LocalOperation{State: api.LocalOperation_OK}, nil
		} else {
			log.Printf("[ERROR] remove file %s: %s\n", source, err.Error())
			return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: err.Error()}, err
		}
	}
	dest := filepath.Join(p.ContentDir, op.DstFile.Path)
	if op.Type == api.LocalOperation_MOVE {
		os.MkdirAll(filepath.Dir(dest), 0750)
		if err := os.Rename(source, dest); err == nil {
			log.Printf("[INFO] file %s moved to %s\n", source, dest)
			return &api.LocalOperation{State: api.LocalOperation_OK, DstFile: &api.Info{Md5: md5}}, nil
		} else {
			log.Printf("[ERROR] file %s move to %s: %s\n", source, dest, err.Error())
			return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: err.Error()}, err
		}
	}

	// op.Type == api.LocalOperation_COPY
	dstMd5, _ := p.getCachedMd5(dest)
	if dstMd5 == md5 {
		// операция не нужна
		log.Printf("[INFO] source file %s and destination file %s are have equal md5, skip copy\n", source, dest)
		return &api.LocalOperation{State: api.LocalOperation_OK, File: &api.Info{Md5: md5}, DstFile: &api.Info{Md5: md5}}, nil
	}
	if err := copyFromTo(source, dest); err == nil {
		dstMd5, _ := p.getCachedMd5(dest)
		state := api.LocalOperation_OK
		if dstMd5 != md5 {
			state = api.LocalOperation_ERROR
		}
		log.Printf("[INFO] copy from %s to %s done\n", source, dest)
		return &api.LocalOperation{State: state, File: &api.Info{Md5: md5}, DstFile: &api.Info{Md5: dstMd5}}, nil
	} else {
		log.Printf("[ERROR] copy from %s to %s: %s\n", source, dest, err.Error())
		return &api.LocalOperation{State: api.LocalOperation_ERROR, Error: err.Error()}, err
	}
}

// реализация приказов для выполнений операци типа: "скопируй файл с сервера"
func (p *PosixServer) RemoteOps(ctx context.Context, op *api.RemoteOperation) (*api.RemoteOperation, error) {

	log.Printf("[INFO] remote operation: {%v}\n", op)

	if op.Type != api.RemoteOperation_COPY_FROM {
		log.Printf("[ERROR] unknown remote operation: %v\n", op.Type)
		return &api.RemoteOperation{State: api.RemoteOperation_ERROR, Error: "unknown remote opertaion"}, nil
	}

	// валидация
	if (op.RemoteFile == nil || op.RemoteFile.Path == "") || (op.ToFile == nil || op.ToFile.Path == "") || op.RemoteServer == "" {
		log.Printf("[ERROR] empty information\n")
		return &api.RemoteOperation{State: api.RemoteOperation_ERROR, Error: "empty information"}, nil
	}

	// создаем коннект, что бы получить информацию с удаленного сервера
	p.NeedRescanByLastOp = true
	server := op.RemoteServer
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", server, op.RemotePort),
		grpc.WithPerRPCCredentials(auth.GrpcAuthFor(server)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[ERROR] dial %s: %s\n", server, err.Error())
		return &api.RemoteOperation{State: api.RemoteOperation_ERROR, Error: fmt.Sprintf("dial %s: %s\n", server, err.Error())}, err
	}
	defer conn.Close()

	client := api.NewFileClient(conn)
	localfilepath := filepath.Join(p.ContentDir, op.ToFile.Path)
	// проверим нужно ли выполнять операцию
	sourceInfos, err := client.Find(context.Background(), &api.Filter{PathMatch: op.RemoteFile.Path})
	if err == nil && sourceInfos.State == api.List_OK && len(sourceInfos.Files) == 1 {
		remoteFile := sourceInfos.Files[0]
		if localMd5, _ := p.getCachedMd5(localfilepath); localMd5 == remoteFile.Md5 {
			// скачивать ничего не надо
			log.Printf("[INFO] remote {%v} and local file {%v} has equal md5 sum\n", op.RemoteFile, op.ToFile)
			return &api.RemoteOperation{State: api.RemoteOperation_OK}, nil
		}
	}

	// откроем файл на запись
	os.MkdirAll(filepath.Dir(localfilepath), 0750)
	fd, err := os.Create(localfilepath)
	if err != nil {
		log.Printf("[ERROR] create %s: %s\n", localfilepath, err.Error())
		return &api.RemoteOperation{State: api.RemoteOperation_ERROR, Error: fmt.Sprintf("create file %s: %s\n", localfilepath, err.Error())}, err
	}

	// попросим удаленный сервер застримить нам файл
	log.Printf("[INFO] send command receive remote file: %s\n", op.RemoteFile.Path)
	if err := ReceiveFileClient(client, op.RemoteFile.Path, fd); err != nil {
		log.Printf("[ERROR] recive %s: %s\n", localfilepath, err.Error())
		return &api.RemoteOperation{State: api.RemoteOperation_ERROR, Error: fmt.Sprintf("receive file %s: %s\n", localfilepath, err.Error())}, err
	}

	log.Printf("[INFO] file %s received\n", localfilepath)
	return &api.RemoteOperation{State: api.RemoteOperation_OK}, nil
}
