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
			record[val.Name] = strings.Trim(fmt.Sprintf("%s", buf[a:a+fieldlen]), " ")
			a = a + fieldlen
		}
		temp[i] = record
		start = start + recordlen
		i = i + 1
	}
}

// var CodeKeys = map[string]string{
// 	"S1":  "HQZQDM",
// 	"S2":  "HQZQJC",
// 	"S4":  "HQJRKP",
// 	"S3":  "HQZRSP",
// 	"S8":  "HQZJCJ",
// 	"S6":  "HQZGCJ",
// 	"S7":  "HQZDCJ",
// 	"S9":  "HQBJW1",
// 	"S10": "HQSJW1",
// 	"S11": "HQCJSL",
// 	"S5":  "HQCJJE",
// 	"S15": "HQBSL1",
// 	"S17": "HQBSL2",
// 	"S16": "HQBJW2",
// 	"S19": "HQBSL3",
// 	"S18": "HQBJW3",
// 	"S27": "HQBSL4",
// 	"S26": "HQBJW4",
// 	"S29": "HQBSL5",
// 	"S28": "HQBJW5",
// 	"S21": "HQSSL1",
// 	"S23": "HQSSL2",
// 	"S22": "HQSJW2",
// 	"S25": "HQSSL3",
// 	"S24": "HQSJW3",
// 	"S31": "HQSSL4",
// 	"S30": "HQSJW4",
// 	"S33": "HQSSL5",
// 	"S32": "HQSJW5",
// 	"S13": "HQSYL1",
// }

// // 兼容两个文件
// func ReadStocksFromDbf(fp *bytes.Reader, is_sh bool) (map[string]map[string]string, error) {

// 	records := GetRecords(fp)
// 	// 表中更新时间
// 	updateTime := ""
// 	updateDate := ""

// 	code_key := "HQZQDM"
// 	time_key := "HQCJBS"
// 	date_key := "HQZQJC"
// 	pre := "sz"
// 	if is_sh {
// 		code_key = "S1"
// 		time_key = "S2"
// 		date_key = "S6"
// 		pre = "sh"
// 	}

// 	for _, val := range records {
// 		if val.Data[code_key] == "000000" {
// 			updateTime = val.Data[time_key]
// 			updateDate = val.Data[date_key]
// 			break
// 		}
// 	}

// 	// 取时间

// 	list := map[string]map[string]string{}

// 	for _, val := range records {
// 		info := map[string]string{}

// 		for sh, sz := range CodeKeys {
// 			kk := sz

// 			if is_sh {
// 				kk = sh
// 			}
// 			info[sz] = strings.Replace(val.Data[kk], "-", "0", -1)
// 			// cod
// 			if sz == "HQZQDM" {
// 				info["code"] = pre + info[sz]
// 			} else if sz == "HQZQJC" {

// 				// 转不了就不转
// 				info[sz] = Mahonia(info[sz])
// 			}
// 		}

// 		info["date"] = updateDate
// 		info["time"] = updateTime

// 		// 是否删除
// 		if val.Delete {
// 			info["delete"] = "1"
// 		} else {
// 			info["delete"] = "0"
// 		}

// 		list[info["code"]] = info
// 	}

// 	return list, nil
// }

// func Float64frombytes(bytes []byte) float64 {
// 	bits := binary.LittleEndian.Uint64(bytes)
// 	f := math.Float64frombits(bits)
// 	return f
// }

// func Mahonia(s string) string {
// 	d := mahonia.NewDecoder("gbk")
// 	return d.ConvertString(s)
// }

// func FindStock(fp *bytes.Reader, is_sh bool, searchField, searchValue string) (map[string]string, error) {

// 	records := GetRecords(fp)
// 	// 表中更新时间
// 	updateTime := ""
// 	updateDate := ""

// 	code_key := "HQZQDM"
// 	time_key := "HQCJBS"
// 	date_key := "HQZQJC"
// 	pre := "sz"
// 	if is_sh {
// 		code_key = "S1"
// 		time_key = "S2"
// 		date_key = "S6"
// 		pre = "sh"
// 	}

// 	// 取时间
// 	for _, val := range records {
// 		if val.Data[code_key] == "000000" {
// 			updateTime = val.Data[time_key]
// 			updateDate = val.Data[date_key]
// 			break
// 		}
// 	}

// 	info := map[string]string{}
// 	for _, val := range records {

// 		for sh, sz := range CodeKeys {
// 			kk := sz

// 			if is_sh {
// 				kk = sh
// 			}
// 			info[sz] = strings.Replace(val.Data[kk], "-", "0", -1)
// 			// cod
// 			if sz == "HQZQDM" {
// 				info["code"] = pre + info[sz]
// 			} else if sz == "HQZQJC" {

// 				// 转不了就不转
// 				info[sz] = Mahonia(info[sz])
// 			}
// 		}

// 		info["date"] = updateDate
// 		info["time"] = updateTime

// 		// 是否删除
// 		if val.Delete {
// 			info["delete"] = "1"
// 		} else {
// 			info["delete"] = "0"
// 		}

// 		if info[searchField] == searchValue {
// 			return info, nil
// 		}

// 	}

// 	return nil, errors.New("value in field not found")
// }
