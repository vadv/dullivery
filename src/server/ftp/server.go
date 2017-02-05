package ftp

// данный файл отвечает за общение с ftp-сервером

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	// github.com/jlaffaye/ftp
	"github.com/secsy/goftp"
	"golang.org/x/net/context"

	api "api/ftp"
)

type FtpServer struct {
	ContentDir string
}

func NewFtpServer(contentDir string) *FtpServer {
	return &FtpServer{
		ContentDir: contentDir,
	}
}

// рекурсивная функция поиска по директории
func ftpWalkDir(ftpClient *goftp.Client, dir string) ([]os.FileInfo, error) {
	result := make([]os.FileInfo, 0)
	files, err := ftpClient.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file)
		} else {
			if newFiles, err := ftpWalkDir(ftpClient, file.Name()); err != nil {
				return nil, err
			} else {
				result = append(result, newFiles...)
			}
		}
	}
	return files, nil
}

// манипуляции с файлами
type ftpOperation string

const (
	ftpOpDownload ftpOperation = "download"
	ftpOpUpload   ftpOperation = "upload"
	ftpOpDelete   ftpOperation = "delete"
)

func (f *FtpServer) manipulation(ctx context.Context, file *api.Info, op ftpOperation) (*api.Info, error) {
	log.Printf("[INFO] request to %s: %v\n", op, file)
	result := &api.Info{}

	// open file
	filename := filepath.Join(f.ContentDir, file.Path)
	switch op {
	case ftpOpDownload, ftpOpUpload:
		os.MkdirAll(filepath.Dir(filename), 0640)
	}

	// open or create file
	var err error
	var fd *os.File
	switch op {
	case ftpOpDownload:
		fd, err = os.Create(filename)
	case ftpOpUpload:
		fd, err = os.Open(filename)
	}

	if err != nil {
		log.Printf("[ERROR] open file: %s\n", err.Error())
		result.State = api.Info_ERROR
		result.Error = err.Error()
		return result, err
	}

	// build client
	if file.Connection == nil {
		info := "empty connection info"
		log.Printf("[ERROR] file: %s\n", info)
		result.State = api.Info_ERROR
		result.Error = info
		return result, fmt.Errorf(info)
	}

	ftpConfig := goftp.Config{
		User:     file.Connection.User,
		Password: file.Connection.Password,
		Timeout:  5 * time.Second,
	}

	ftpClient, err := goftp.DialConfig(ftpConfig, file.Connection.Host)
	if err != nil {
		log.Printf("[ERROR] dial to ftp %s: %s\n", file.Connection.Host, err.Error())
		result.State = api.Info_ERROR
		result.Error = err.Error()
		return result, err
	}

	// operation
	switch op {
	case ftpOpDownload:
		err = ftpClient.Retrieve(file.Path, fd)
	case ftpOpUpload:
		err = ftpClient.Store(file.Path, fd)
	case ftpOpDelete:
		err = ftpClient.Delete(file.Path)
	}

	if err != nil {
		log.Printf("[ERROR] %s from ftp %s: %s\n", op, file.Connection.Host, err.Error())
		result.State = api.Info_ERROR
		result.Error = err.Error()
		return result, err
	}

	log.Printf("[INFO] %s complete\n", op)
	result.State = api.Info_OK
	return result, nil

}

// реализация поиска
func (f *FtpServer) Find(ctx context.Context, filter *api.Filter) (*api.List, error) {

	log.Printf("[INFO] find request: %v\n", filter)
	result := &api.List{Files: make([]*api.Info, 0)}
	// build regexp
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
	// build client
	if filter.Connection == nil {
		info := "empty connection info"
		log.Printf("[ERROR] filter: %s\n", info)
		result.State = api.List_ERROR
		result.Error = info
		return result, fmt.Errorf(info)
	}
	ftpConfig := goftp.Config{
		User:     filter.Connection.User,
		Password: filter.Connection.Password,
		Timeout:  5 * time.Second,
	}
	ftpClient, err := goftp.DialConfig(ftpConfig, filter.Connection.Host)
	if err != nil {
		log.Printf("[ERROR] dial to ftp %s: %s\n", filter.Connection.Host, err.Error())
		result.State = api.List_ERROR
		result.Error = err.Error()
		return result, err
	}
	// find
	files, err := ftpWalkDir(ftpClient, "/")
	if err != nil {
		log.Printf("[ERROR] walk directories %s: %s\n", filter.Connection.Host, err.Error())
		result.State = api.List_ERROR
		result.Error = err.Error()
		return result, err
	}
	// filter
	for _, file := range files {
		if reg.MatchString(file.Name()) {
			result.Files = append(result.Files, &api.Info{
				Path:    file.Name(),
				Size:    file.Size(),
				ModTime: file.ModTime().Unix(),
			})
		}
	}
	result.State = api.List_OK
	return result, nil
}

// реализация загрузки файла
func (f *FtpServer) Download(ctx context.Context, file *api.Info) (*api.Info, error) {
	return f.manipulation(ctx, file, ftpOpDownload)
}

// реализация выгрузки файла
func (f *FtpServer) Upload(ctx context.Context, file *api.Info) (*api.Info, error) {
	return f.manipulation(ctx, file, ftpOpUpload)
}

// реализация удаления файла с ftp
func (f *FtpServer) Delete(ctx context.Context, file *api.Info) (*api.Info, error) {
	return f.manipulation(ctx, file, ftpOpDelete)
}
