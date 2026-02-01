package render

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// TTY is true if the current process is running in a terminal.
var TTY = isatty.IsTerminal(os.Stdout.Fd())

// JSON is true if the --json flag is set.
var JSON = false

// init checks if the current process is running in a terminal and disables color if not.
func init() {
	if !TTY {
		color.NoColor = true
	}
}

// Message prints a message to the console.
func Message(msg string) {
	if JSON {
		// In JSON mode, we don't print messages unless they are part of an error.
		return
	}
	fmt.Println(msg)
}

// Object prints a struct as JSON if --json is set, otherwise it prints a
// formatted string.
func Object(v interface{}) {
	if JSON {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(b))
		return
	}

	// If not in JSON mode, we expect the object to have a String() method.
	fmt.Println(v)
}

// Error prints an error message to the console.
func Error(err error) {
	if JSON {
		v := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
		return
	}
	fmt.Fprintf(os.Stderr, "%s %s\n", color.RedString("Error:"), err.Error())
}
