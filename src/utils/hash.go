package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func Md5FD(fd *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, fd); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}

func Md5(file string) (string, error) {
	fd, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return Md5FD(fd)
}
