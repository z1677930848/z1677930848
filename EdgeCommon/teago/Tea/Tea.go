package Tea

import (
	"os"
	"path/filepath"
	"strings"
)

// Root represents the project root. Default to current working directory if not set.
var Root string = func() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}()

// DS is the system path separator.
var DS = string(os.PathSeparator)

// Env indicates the current environment (prod/test/dev).
var Env = "dev"

var (
	publicDir string
	viewsDir  string
	tmpDir    string
	configDir string
	logDir    string
)

// IsTesting reports whether the process is running in test mode.
// Detect common go test launcher patterns and a GO_TEST hint.
func IsTesting() bool {
	if os.Getenv("GO_TEST") != "" {
		return true
	}
	exe := filepath.Base(os.Args[0])
	return strings.HasSuffix(exe, ".test") || strings.HasSuffix(exe, ".test.exe")
}

// ConfigFile joins Root with a config file name under "configs".
func ConfigFile(name string) string {
	base := configDir
	if base == "" {
		base = filepath.Join(Root, "configs")
	}
	return filepath.Join(base, name)
}

// LogDir returns the directory for logs.
func LogDir() string {
	if logDir != "" {
		return logDir
	}
	return filepath.Join(Root, "logs")
}

// LogFile returns the path to a log file in the log dir.
func LogFile(name string) string {
	return filepath.Join(LogDir(), name)
}

// SetConfigDir sets config directory.
func SetConfigDir(dir string) {
	configDir = dir
}

// SetPublicDir sets public assets directory.
func SetPublicDir(dir string) {
	publicDir = dir
}

// SetViewsDir sets template/views directory.
func SetViewsDir(dir string) {
	viewsDir = dir
}

// SetTmpDir sets temporary directory.
func SetTmpDir(dir string) {
	tmpDir = dir
}

// SetLogDir sets log directory.
func SetLogDir(dir string) {
	logDir = dir
}

// UpdateRoot updates project root.
func UpdateRoot(root string) {
	Root = root
}
