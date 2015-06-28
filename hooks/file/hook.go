package file

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gogap/logrus"
	"github.com/gogap/logrus/hooks/caller"
)

func NewHook(file string) (f *FileHook) {
	path := strings.Split(file, "/")
	if len(path) > 1 {
		exec.Command("mkdir", path[0]).Run()
	}
	w := NewFileWriter()
	config := fmt.Sprintf(`{"filename":"%s","maxdays":7}`, file)
	w.Init(config)
	return &FileHook{w}
}

type FileHook struct {
	W LoggerInterface
}

func (hook *FileHook) Fire(entry *logrus.Entry) (err error) {
	file, lineNumber := caller.GetCallerIgnoringLogMulti(1)
	if file != "" {
		sep := fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
		fileName := strings.Split(file, sep)
		if len(fileName) >= 2 {
			file = fileName[1]
		}
	}
	entry.Data["file"] = file
	entry.Data["line"] = lineNumber
	message, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[ERROR] %s", message), LevelError)
	case logrus.FatalLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[ERROR] %s", message), LevelError)
	case logrus.ErrorLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[ERROR] %s", message), LevelError)
	case logrus.WarnLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[WARN] %s", message), LevelWarn)
	case logrus.InfoLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[INFO] %s", message), LevelInfo)
	case logrus.DebugLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[DEBUG] %s", message), LevelDebug)
	default:
		return nil
	}
	return
}

func (hook *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
