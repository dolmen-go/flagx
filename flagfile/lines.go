package flagfile

import (
	"bufio"
	"io"
)

// Lines is a Loader that reads a text file and takes each line as an argument.
func Lines(r io.Reader, fragment string) (interface{}, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
