package exit

import (
	"fmt"
	"os"

	"github.com/savedhq/sctl/internal/render"
)

// Success exits with a 0 exit code.
func Success() {
	os.Exit(0)
}

// Error prints an error message and exits with a 1 exit code.
func Error(err error) {
	render.Error(err)
	os.Exit(1)
}

// ErrorMessage prints a formatted error message and exits with a 1 exit code.
func ErrorMessage(format string, a ...interface{}) {
	err := fmt.Errorf(format, a...)
	render.Error(err)
	os.Exit(1)
}
