package dbf

import (
	"fmt"
	"strings"
)

type Reader interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	ReadAt(b []byte, off int64) (n int, err error)
}

type DbfHead struct {
	Version    []byte
	Updatedate string
	Records    int64
	Headerlen  int64
	Recordlen  int64
}
type Field struct {
	Name             string
	Fieldtype        string
	FieldDataaddress []byte
	FieldLen         int64
	DecimalCount     []byte
	Workareaid       []byte
}
type Record struct {
	Delete bool
	//Data   string
	Data map[string]string
}

func GetDbfHead(reader Reader) (dbfhead DbfHead) {
	buf := make([]byte, 16)
	reader.Seek(0, 0)
	_, err := reader.Read(buf)
	if err != nil {
		panic(err)
	}
	dbfhead.Version = buf[0:1]
	dbfhead.Updatedate = fmt.Sprintf("%d", buf[1:4])
	dbfhead.Headerlen = Changebytetoint(buf[8:10])
	dbfhead.Recordlen = Changebytetoint(buf[10:12])
	dbfhead.Records = Changebytetoint(buf[4:8])
	return dbfhead
}
func RemoveNullfrombyte(b []byte) (s string) {
	for _, val := range b {
		if val == 0 {
			continue
		}
		s = s + string(val)
	}
	return
}

func GetFields(reader Reader) []Field {
	dbfhead := GetDbfHead(reader)

	// off := dbfhead.Headerlen - 32 - 264
	off := dbfhead.Headerlen - 8
	fieldlist := make([]Field, off/32)
	buf := make([]byte, off)
	_, err := reader.ReadAt(buf, 32)
	if err != nil {
		panic(err)
	}
	curbuf := make([]byte, 32)
	for i, val := range fieldlist {
		a := i * 32

		curbuf = buf[a:]
		val.Name = RemoveNullfrombyte(curbuf[0:11])
		val.Fieldtype = fmt.Sprintf("%s", curbuf[11:12])
		val.FieldDataaddress = curbuf[12:16]
		val.FieldLen = Changebytetoint(curbuf[16:17])
		val.DecimalCount = curbuf[17:18]
		val.Workareaid = curbuf[20:21]
		fieldlist[i] = val

	}
	return fieldlist
}

func Changebytetoint(b []byte) (x int64) {
	for i, val := range b {
		if i == 0 {
			x = x + int64(val)
		} else {
			x = x + int64(2<<7*int64(i)*int64(val))
		}
		//fmt.Println(x)
	}
	//fmt.Println(fieldlist)

	return
}

func GetRecords(fp Reader) []map[string]string {
	dbfhead := GetDbfHead(fp)
	fp.Seek(0, 0)
	fields := GetFields(fp)
	recordlen := dbfhead.Recordlen
	start := dbfhead.Headerlen
	buf := make([]byte, recordlen)
	i := 0
	temp := make([]map[string]string, dbfhead.Records)
	for {
		_, err := fp.ReadAt(buf, start)
		if err != nil {
			return temp
			panic(err)
		}
		record := map[string]string{}
		if string(buf[0:1]) == " " {
			record["delete"] = "0"
		} else if string(buf[0:1]) == "*" {
			record["delete"] = "1"
		}

		a := int64(1)
		for _, val := range fields {
			fieldlen := val.FieldLen
			record[val.Name] = strings.TrimSpace(fmt.Sprintf("%s", buf[a:a+fieldlen]))
			a = a + fieldlen
		}
		temp[i] = record
		start = start + recordlen
		i = i + 1
	}
}
