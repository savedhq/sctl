package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	saved "github.com/savedhq/sdk-go"
)

// PrintAPIError prints a friendly error message for SDK errors
func PrintAPIError(err error) error {
	if apiErr, ok := err.(*saved.GenericOpenAPIError); ok {
		// If we can extract a body, try to show it
		if len(apiErr.Body()) > 0 {
			return fmt.Errorf("API Error: %s", string(apiErr.Body()))
		}
		// If there's an internal model, maybe show that
		if apiErr.Model() != nil {
			return fmt.Errorf("API Error: %v", apiErr.Model())
		}
	}

	return fmt.Errorf("API Error: %w", err)
}

func CheckErr(err error) {
	if err != nil {
		color.Red("Error: %s", err)
		os.Exit(1)
	}
}

func PrintJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		CheckErr(fmt.Errorf("failed to marshal JSON: %w", err))
	}
	fmt.Println(string(data))
}
