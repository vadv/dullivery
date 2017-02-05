package posix

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	api "api/posix"
)

var POSIX_FILE_BUFFER_SIZE = 2 * 1024 * 1024

func ReceiveFileClient(client api.FileClient, remotePath string, fd *os.File) error {

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

	// init connection
	stream, err := client.Stream(context.Background(), &api.Info{Path: remotePath})
	if err != nil {
		log.Printf("[ERROR] open stream: %s\n", err.Error())
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[ERROR] recieve: %s\n", err.Error())
			return err
		}
		if err := writeChunk(fd, chunk); err != nil {
			log.Printf("[ERROR] write chunk: %s\n", err.Error())
			return err
		}
	}

	return nil
}

func StreamFileClient(client api.FileClient, fd *os.File, filename, md5 string) (*api.Info, error) {

	stream, err := client.Receive(context.Background())
	if err != nil {
		log.Printf("[ERROR] open stream: %s\n", err.Error())
		return nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		log.Printf("[ERROR] stat: %s\n", err.Error())
		return nil, err
	}

	fileInfo := &api.Info{Path: filename, Md5: md5, Size: stat.Size(), State: api.Info_OK}

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
			return nil, err
		}
		chunk := &api.Chunk{Data: buf, Offset: offset, File: fileInfo}
		offset += int64(len(buf))

		// send
		if err := stream.Send(chunk); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("[ERROR] send chunk: %s\n", err.Error())
			return nil, err
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("[ERROR] send close stream error: %s\n", err.Error())
		return reply, err
	}
	return reply, nil
}
