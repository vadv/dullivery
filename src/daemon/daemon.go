package daemon

// https://github.com/golang/go/issues/227
// https://habrahabr.ru/post/187668/
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	envDaemonName  = "_runned_is_daemon"
	envDaemonValue = "1"
)

func Daemonize(fd *os.File) error {
	return reborn(0002, fd)
}
func reborn(umask uint32, fd *os.File) (err error) {
	if !IsDaemon() {
		var path string
		if path, err = filepath.Abs(os.Args[0]); err != nil {
			return
		}
		cmd := exec.Command(path, os.Args[1:]...)
		envVar := fmt.Sprintf("%s=%s", envDaemonName, envDaemonValue)
		cmd.Env = append(os.Environ(), envVar)
		cmd.Stdout = fd
		cmd.Stderr = fd
		if err = cmd.Start(); err != nil {
			return
		}
		os.Exit(0)
	}
	return
}
func IsDaemon() bool {
	return os.Getenv(envDaemonName) == envDaemonValue
}
