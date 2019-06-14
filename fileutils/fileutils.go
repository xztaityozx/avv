package fileutils

import (
	"io/ioutil"
	"os"
)


// TryMakeDirAll is try make directory
// params:
//  - path: path to directory
func TryMakeDirAll(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

// WriteFile is write string to file
// params:
//  - path: path to target file
//  - data: string to write
func WriteFile(path,data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

