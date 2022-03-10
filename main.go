package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/cuigh/auxo/config"
)

func main() {
	config.SetDefaultValue("banner", false)
	config.SetDefaultValue("app_path", "/usr/share/nginx/html")
	config.SetDefaultValue("var_name", "config")

	profile := config.GetString("profile")
	varName := config.GetString("var_name")
	dirs := strings.Split(config.GetString("app_path"), ",")
	for _, dir := range dirs {
		ensure(injectConfig(dir, profile, varName))
	}

	ensure(startNginx())
}

func ensure(err error) {
	if err != nil {
		writeLog("error", err.Error())
		os.Exit(1)
	}
}

func injectConfig(dir, profile, varName string) error {
	c, err := loader.Load(dir, profile)
	if err != nil {
		return err
	} else if c == "" {
		return nil
	}
	writeLog("notice", "app: %s, profile: %s, config: %s", dir, profile, c)

	filename := filepath.Join(dir, "index.html")
	d, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	f, err := Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// place config variables before first script tag
	index := bytes.Index(d, []byte("<script"))
	if index == -1 {
		return errors.New("no script tag in index.html")
	}

	f.Truncate()
	f.Write(d[:index])
	f.WriteString(fmt.Sprintf(`<script type="text/javascript">window.%s=%s</script>`, varName, c))
	f.Write([]byte{'\n'})
	f.WriteString("    ")
	f.Write(d[index:])
	return f.Error()
}

func startNginx() error {
	cmd := exec.Command("nginx", "-g", "daemon off;")
	return syscall.Exec(cmd.Path, cmd.Args, os.Environ())
}

func writeLog(level, format string, args ...interface{}) {
	log.Println(fmt.Sprintf("[%s]", level), fmt.Sprintf(format, args...))
}

type FileWrapper struct {
	file *os.File
	err  error
}

func Open(filename string) (*FileWrapper, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileWrapper{file: f}, nil
}

func (w *FileWrapper) Close() {
	_ = w.file.Close()
}

func (w *FileWrapper) Truncate() {
	if w.err == nil {
		w.err = w.file.Truncate(0)
	}
}

func (w *FileWrapper) Write(b []byte) {
	if w.err == nil {
		_, w.err = w.file.Write(b)
	}
}

func (w *FileWrapper) WriteString(s string) {
	if w.err == nil {
		_, w.err = w.file.WriteString(s)
	}
}

func (w *FileWrapper) Error() error {
	return w.err
}
