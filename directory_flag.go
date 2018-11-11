package main

import (
	"fmt"
	"os"
)

type DirectoryFlag string

func (d *DirectoryFlag) UnmarshalFlag(value string) (err error) {
	if value == "" {
		return nil
	}

	finfo, err := os.Stat(value)
	if err != nil {
		err = fmt.Errorf("couldn't retrieve info regarding directory '%s'", value)
		return
	}

	if !finfo.IsDir() {
		err = fmt.Errorf("value '%s' is not a directory", value)
		return
	}

	*d = DirectoryFlag(value)
	return
}
