package dbf

import (
	"bytes"
	"fmt"
	"os"
)

func FromBytes(b []byte) (data []map[string]string, err error) {
	defer func(err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}(err)

	rd := bytes.NewReader(b)
	data = GetRecords(rd)

	return
}

func FromFile(fpath string) (data []map[string]string, err error) {
	defer func(err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}(err)

	var fp *os.File
	fp, err = os.OpenFile(fpath, os.O_RDONLY, 0)
	if err != nil {
		return
	}

	defer fp.Close()

	data = GetRecords(fp)

	return
}
