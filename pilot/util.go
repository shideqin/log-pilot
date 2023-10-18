package pilot

import (
	"os"
	"strings"
)

// ReadFile return string list separated by separator
func ReadFile(path string, separator string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), separator), nil
}
