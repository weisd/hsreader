package txt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	// "strconv"
	// "strings"
	// "time"

	// log "code.google.com/p/log4go"
	"../mahonia"
)

var (
	HeaderTab = "HEADER"
	FooterTab = "TRAILER"
	IndexTab  = "MD001" // 指数
	StockTab  = "MD002" // 股票 A、B股
	BondTab   = "MD003" // 债券
	FundTab   = "MD004" // 基金
)

type TxtContent struct {
	Header    map[string]string
	Footer    map[string]string
	StockList []map[string]string
	BondList  []map[string]string
	FundList  []map[string]string
	IndexList []map[string]string
}

func ReadFile(fpath string) (contents TxtContent, err error) {
	fp, err := os.OpenFile(fpath, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer fp.Close()

	header := map[string]string{}
	footer := map[string]string{}
	stockList := make([]map[string]string, 0)
	bondList := make([]map[string]string, 0)
	fundList := make([]map[string]string, 0)
	indexList := make([]map[string]string, 0)

	contents = TxtContent{}

	breader := bufio.NewReader(fp)
	lines := make([][]byte, 0)
	for {
		line, err := breader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return contents, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// |线分割
		// 去space

		lines = append(lines, line)
	}

	fmt.Println(len(lines))

	for i, _ := range lines {
		info := splitLine(lines[i])

		flag := info[0]

		switch string(flag) {
		case HeaderTab:
			header = parseFields(info, HeaderFields)
		case FooterTab:
			footer = parseFields(info, FooterFields)
		case IndexTab:
			item := parseFields(info, IndexFields)
			indexList = append(indexList, item)
		case StockTab:
			item := parseFields(info, StockFields)
			stockList = append(stockList, item)
		case BondTab:
			item := parseFields(info, StockFields)
			bondList = append(bondList, item)
		case FundTab:
			item := parseFields(info, FoundFields)
			fundList = append(fundList, item)
		}

	}

	contents.Header = header
	contents.Footer = footer
	contents.StockList = stockList
	contents.BondList = bondList
	contents.FundList = fundList
	contents.IndexList = indexList

	return
}

var HeaderFields = []string{
	"BeginString",
	"Version",
	"BodyLength",
	"TotNumTradeReports",
	"MDReportID",
	"SenderCompID",
	"XSHG01",
	"MDTime",
	"MDUpdateType",
	"MDSesStatus",
}

var FooterFields = []string{
	"EndString",
	"CheckSum",
}

var IndexFields = []string{
	"MDStreamID",
	"SecurityID",
	"Symbol",
	"TradeVolume",
	"TotalValueTraded",
	"PreClosePx",
	"OpenPrice",
	"HighPrice",
	"LowPrice",
	"TradePrice",
	"ClosePx",
	"TradingPhaseCode",
	"Timestamp",
}

var StockFields = []string{
	"MDStreamID",
	"SecurityID",
	"Symbol",
	"TradeVolume",
	"TotalValueTraded",
	"PreClosePx",
	"OpenPrice",
	"HighPrice",
	"LowPrice",
	"TradePrice",
	"ClosePx",
	"BuyPrice1",
	"BuyVolume1",
	"SellPrice1",
	"SellVolume1",
	"BuyPrice2",
	"BuyVolume2",
	"SellPrice2",
	"SellVolume2",
	"BuyPrice3",
	"BuyVolume3",
	"SellPrice3",
	"SellVolume3",
	"BuyPrice4",
	"BuyVolume4",
	"SellPrice4",
	"SellVolume4",
	"BuyPrice5",
	"BuyVolume5",
	"SellPrice5",
	"SellVolume5",
	"TradingPhaseCode",
	"Timestamp",
}

var FoundFields = []string{
	"MDStreamID",
	"SecurityID",
	"Symbol",
	"TradeVolume",
	"TotalValueTraded",
	"PreClosePx",
	"OpenPrice",
	"HighPrice",
	"LowPrice",
	"TradePrice",
	"ClosePx",
	"BuyPrice1",
	"BuyVolume1",
	"SellPrice1",
	"SellVolume1",
	"BuyPrice2",
	"BuyVolume2",
	"SellPrice2",
	"SellVolume2",
	"BuyPrice3",
	"BuyVolume3",
	"SellPrice3",
	"SellVolume3",
	"BuyPrice4",
	"BuyVolume4",
	"SellPrice4",
	"SellVolume4",
	"BuyPrice5",
	"BuyVolume5",
	"SellPrice5",
	"SellVolume5",
	"PreCloseIOPV",
	"IOPV",
	"TradingPhaseCode",
	"Timestamp",
}

func parseFields(vals [][]byte, fileds []string) map[string]string {
	info := map[string]string{}
	for k, _ := range vals {
		key := fileds[k]
		if key == "Symbol" {
			// gbk to utf-8
			info[key] = Gkb2Utf8(string(vals[k]))
			continue
		}
		info[key] = string(vals[k])
	}
	return info
}

// |线分割 去space
func splitLine(line []byte) [][]byte {
	vals := bytes.Split(line, []byte("|"))
	for k, _ := range vals {
		vals[k] = bytes.TrimSpace(vals[k])
	}
	return vals
}

func Gkb2Utf8(s string) string {
	d := mahonia.NewDecoder("gbk")
	return d.ConvertString(s)
}

// func ReadMktdtFromBytes(src []byte) (list map[string]map[string]string, err error) {
// 	// list = make(map[string]map[string]string)

// 	lines := bytes.Split(src, []byte("\n"))
// 	// log.Info("lines %d", len(lines))

// 	header := parseMktdtHeader(lines[0])

// 	// log.Info("header %v", header)

// 	dTime, err := time.Parse("20060102-15:04:05", string(header.MDTime[0:17]))
// 	if err != nil {
// 		return nil, err
// 	}

// 	bodys := parseMKtdtBody(lines[1:header.TotNumTradeReports+1], dTime.Format("20060102"))

// 	// footer := parseMktdtFooter(lines[header.TotNumTradeReports+1])

// 	// log.Info("footer %v", footer)

// 	return bodys, nil
// }

// var MktdtKeysMap = map[string]string{
// 	"SecurityID": "HQZQDM", // 证券代码
// 	"Symbol":     "HQZQJC", // 证券简称

// 	"OpenPrice":  "HQJRKP", // 今日开盘价
// 	"PreClosePx": "HQZRSP", // 昨日收盘价

// 	"TradePrice": "HQZJCJ", // 最近成交价
// 	"HighPrice":  "HQZGCJ", // 最高成交价
// 	"LowPrice":   "HQZDCJ", // 最低成交价

// 	"TradeVolume":      "HQCJSL", // 成交数量
// 	"TotalValueTraded": "HQCJJE", // 成交金额

// 	"BuyPrice1":   "HQBJW1", // 买价位一
// 	"BuyVolume1":  "HQBSL1", // 买数量一
// 	"SellPrice1":  "HQSJW1", // 卖价位一
// 	"SellVolume1": "HQSSL1", // 卖数量一

// 	"BuyPrice2":   "HQBJW2", // 买价位一
// 	"BuyVolume2":  "HQBSL2", // 买数量一
// 	"SellPrice2":  "HQSJW2", // 卖价位一
// 	"SellVolume2": "HQSSL2", // 卖数量一

// 	"BuyPrice3":   "HQBJW3", // 买价位一
// 	"BuyVolume3":  "HQBSL3", // 买数量一
// 	"SellPrice3":  "HQSJW3", // 卖价位一
// 	"SellVolume3": "HQSSL3", // 卖数量一

// 	"BuyPrice4":   "HQBJW4", // 买价位一
// 	"BuyVolume4":  "HQBSL4", // 买数量一
// 	"SellPrice4":  "HQSJW4", // 卖价位一
// 	"SellVolume4": "HQSSL4", // 卖数量一

// 	"BuyPrice5":   "HQBJW5", // 买价位一
// 	"BuyVolume5":  "HQBSL5", // 买数量一
// 	"SellPrice5":  "HQSJW5", // 卖价位一
// 	"SellVolume5": "HQSSL5", // 卖数量一

// 	"S13": "HQSYL1", // 市盈率1
// }

// type MktdtMd001 struct {
// 	MDStreamID       string // MD001 表示指数行情数据格式类型，其中指数目前实际精度为4位小数；
// 	SecurityID       string // 指数代码
// 	Symbol           string // 指数简称
// 	TradeVolume      string // 成交数量 参与计算相应指数的交易数量，股票指数交易数量单位是100股，基金指数的交易数量单位是100份，债券指数的交易数量单位是手。
// 	TotalValueTraded string // 成交金额
// 	PreClosePx       string // 昨收
// 	OpenPrice        string // 今开盘价
// 	HighPrice        string // 最高价
// 	LowPrice         string // 最低价
// 	TradePrice       string // 最新价
// 	ClosePx          string // 今日收盘价 无取值取空格
// 	TradingPhaseCode string // 指数实时阶段及标志 该字段为8位字符串，左起每位表示特定的含义，无定义则填空格。（预留）
// 	Timestamp        string // 时间 HH:MM:SS.000
// }

// func (m *MktdtMd001) ToMap() map[string]string {
// 	info := map[string]string{}
// 	info["MDStreamID"] = m.MDStreamID
// 	info["SecurityID"] = m.SecurityID
// 	info["Symbol"] = m.Symbol
// 	info["TradeVolume"] = m.TradeVolume
// 	info["TotalValueTraded"] = m.TotalValueTraded
// 	info["PreClosePx"] = m.PreClosePx
// 	info["OpenPrice"] = m.OpenPrice
// 	info["HighPrice"] = m.HighPrice
// 	info["LowPrice"] = m.LowPrice
// 	info["TradePrice"] = m.TradePrice
// 	info["ClosePx"] = m.ClosePx
// 	info["TradingPhaseCode"] = m.TradingPhaseCode
// 	info["Timestamp"] = m.Timestamp

// 	// info["code"] = "sh" + m.SecurityID
// 	// info["date"] = ""
// 	// info["time"] = strings.Replace(string(m.Timestamp[0:8]), ":", "", -1)
// 	// info["delete"] = "1"

// 	oldInfo := map[string]string{}

// 	for k, v := range MktdtKeysMap {
// 		val, ok := info[k]
// 		if ok {
// 			oldInfo[v] = val
// 		} else {
// 			oldInfo[v] = ""
// 		}

// 	}

// 	oldInfo["code"] = "sh" + m.SecurityID
// 	oldInfo["date"] = ""
// 	oldInfo["time"] = strings.Replace(string(m.Timestamp[0:8]), ":", "", -1)
// 	oldInfo["delete"] = "1"

// 	return oldInfo

// 	return info
// }

// type MktdtMd002 struct {
// 	MDStreamID       string // MD001 表示指数行情数据格式类型，其中指数目前实际精度为4位小数；
// 	SecurityID       string // 指数代码
// 	Symbol           string // 指数简称
// 	TradeVolume      string // 成交数量 参与计算相应指数的交易数量，股票指数交易数量单位是100股，基金指数的交易数量单位是100份，债券指数的交易数量单位是手。
// 	TotalValueTraded string // 成交金额
// 	PreClosePx       string // 昨收
// 	OpenPrice        string // 今开盘价
// 	HighPrice        string // 最高价
// 	LowPrice         string // 最低价
// 	TradePrice       string // 最新价
// 	ClosePx          string // 今日收盘价 无取值取空格

// 	BuyPrice1   string // 买一价
// 	BuyVolume1  string // 买一量
// 	SellPrice1  string // 卖一价
// 	SellVolume1 string // 卖一量

// 	BuyPrice2   string // 买一价
// 	BuyVolume2  string // 买一量
// 	SellPrice2  string // 卖一价
// 	SellVolume2 string // 卖一量

// 	BuyPrice3   string // 买一价
// 	BuyVolume3  string // 买一量
// 	SellPrice3  string // 卖一价
// 	SellVolume3 string // 卖一量

// 	BuyPrice4   string // 买一价
// 	BuyVolume4  string // 买一量
// 	SellPrice4  string // 卖一价
// 	SellVolume4 string // 卖一量

// 	BuyPrice5   string // 买一价
// 	BuyVolume5  string // 买一量
// 	SellPrice5  string // 卖一价
// 	SellVolume5 string // 卖一量

// 	PreCloseIOPV string // 基金T-1日收盘时刻IOPV 可选字段，仅当MDStreamID=MD004时存在该字段。
// 	IOPV         string // 基金IOPV 可选字段，仅当MDStreamID=MD004时存在该字段。

// 	TradingPhaseCode string // 指数实时阶段及标志 该字段为8位字符串，左起每位表示特定的含义，无定义则填空格。（预留）
// 	Timestamp        string // 时间 HH:MM:SS.000

// }

// func (m *MktdtMd002) ToMap() map[string]string {
// 	info := map[string]string{}
// 	info["MDStreamID"] = m.MDStreamID
// 	info["SecurityID"] = m.SecurityID
// 	info["Symbol"] = m.Symbol
// 	info["TradeVolume"] = m.TradeVolume
// 	info["TotalValueTraded"] = m.TotalValueTraded
// 	info["PreClosePx"] = m.PreClosePx
// 	info["OpenPrice"] = m.OpenPrice
// 	info["HighPrice"] = m.HighPrice
// 	info["LowPrice"] = m.LowPrice
// 	info["TradePrice"] = m.TradePrice
// 	info["ClosePx"] = m.ClosePx

// 	info["BuyPrice1"] = m.BuyPrice1
// 	info["BuyVolume1"] = m.BuyVolume1
// 	info["SellPrice1"] = m.SellPrice1
// 	info["SellVolume1"] = m.SellVolume1

// 	info["BuyPrice2"] = m.BuyPrice2
// 	info["BuyVolume2"] = m.BuyVolume2
// 	info["SellPrice2"] = m.SellPrice2
// 	info["SellVolume2"] = m.SellVolume2

// 	info["BuyPrice3"] = m.BuyPrice3
// 	info["BuyVolume3"] = m.BuyVolume3
// 	info["SellPrice3"] = m.SellPrice3
// 	info["SellVolume3"] = m.SellVolume3

// 	info["BuyPrice4"] = m.BuyPrice4
// 	info["BuyVolume4"] = m.BuyVolume4
// 	info["SellPrice4"] = m.SellPrice4
// 	info["SellVolume4"] = m.SellVolume4

// 	info["BuyPrice5"] = m.BuyPrice5
// 	info["BuyVolume5"] = m.BuyVolume5
// 	info["SellPrice5"] = m.SellPrice5
// 	info["SellVolume5"] = m.SellVolume5

// 	info["PreCloseIOPV"] = m.PreCloseIOPV
// 	info["IOPV"] = m.IOPV

// 	info["TradingPhaseCode"] = m.TradingPhaseCode
// 	info["Timestamp"] = m.Timestamp

// 	oldInfo := map[string]string{}

// 	for k, v := range MktdtKeysMap {
// 		val, ok := info[k]
// 		if ok {
// 			oldInfo[v] = val
// 		} else {
// 			oldInfo[v] = ""
// 		}

// 	}

// 	oldInfo["code"] = "sh" + m.SecurityID
// 	oldInfo["date"] = ""
// 	oldInfo["time"] = strings.Replace(string(m.Timestamp[0:8]), ":", "", -1)
// 	oldInfo["delete"] = "1"

// 	return oldInfo
// }

// func parseMKtdtBody(lines [][]byte, date string) map[string]map[string]string {
// 	l := len(lines)

// 	list := map[string]map[string]string{}

// 	for i := 0; i < l; i++ {
// 		line := lines[i]

// 		items := bytes.Split(line, []byte("|"))

// 		if len(items) < 13 {
// 			continue
// 		}

// 		info := map[string]string{}

// 		mdStreamID := strings.TrimSpace(string(items[0]))
// 		switch mdStreamID {
// 		case "MD001":
// 			md001 := parseMktdtMd001(line)
// 			info = md001.ToMap()
// 		case "MD002", "MD003", "MD004":
// 			md002 := parseMktdtMd002(line)
// 			info = md002.ToMap()
// 		default:
// 			log.Warn("unkown body type %s", mdStreamID)
// 			continue
// 		}

// 		info["date"] = date
// 		code := info["code"]

// 		list[code] = info

// 	}

// 	return list
// }

// func parseMktdtMd001(data []byte) *MktdtMd001 {
// 	items := bytes.Split(data, []byte("|"))

// 	if len(items) < 13 {
// 		log.Warn("parseMktdtMd002 len err")
// 		return nil
// 	}

// 	item := &MktdtMd001{}

// 	item.MDStreamID = strings.TrimSpace(string(items[0]))
// 	item.SecurityID = strings.TrimSpace(string(items[1]))
// 	// item.Symbol = Gkb2Utf8(strings.TrimSpace(string(items[2])))
// 	item.TradeVolume = strings.TrimSpace(string(items[3]))
// 	item.TotalValueTraded = strings.TrimSpace(string(items[4]))
// 	item.PreClosePx = strings.TrimSpace(string(items[5]))
// 	item.OpenPrice = strings.TrimSpace(string(items[6]))
// 	item.HighPrice = strings.TrimSpace(string(items[7]))
// 	item.LowPrice = strings.TrimSpace(string(items[8]))
// 	item.TradePrice = strings.TrimSpace(string(items[9]))
// 	item.ClosePx = strings.TrimSpace(string(items[10]))
// 	item.TradingPhaseCode = strings.TrimSpace(string(items[11]))
// 	item.Timestamp = strings.TrimSpace(string(items[12]))

// 	return item
// }

// func parseMktdtMd002(data []byte) *MktdtMd002 {
// 	items := bytes.Split(data, []byte("|"))

// 	if len(items) < 33 {
// 		log.Warn("parseMktdtMd002 len err")
// 		return nil
// 	}

// 	item := &MktdtMd002{}

// 	item.MDStreamID = strings.TrimSpace(string(items[0]))
// 	item.SecurityID = strings.TrimSpace(string(items[1]))
// 	// item.Symbol = Gkb2Utf8(strings.TrimSpace(string(items[2])))
// 	item.TradeVolume = strings.TrimSpace(string(items[3]))
// 	item.TotalValueTraded = strings.TrimSpace(string(items[4]))
// 	item.PreClosePx = strings.TrimSpace(string(items[5]))
// 	item.OpenPrice = strings.TrimSpace(string(items[6]))
// 	item.HighPrice = strings.TrimSpace(string(items[7]))
// 	item.LowPrice = strings.TrimSpace(string(items[8]))
// 	item.TradePrice = strings.TrimSpace(string(items[9]))
// 	item.ClosePx = strings.TrimSpace(string(items[10]))

// 	item.BuyPrice1 = strings.TrimSpace(string(items[11]))
// 	item.BuyVolume1 = strings.TrimSpace(string(items[12]))
// 	item.SellPrice1 = strings.TrimSpace(string(items[13]))
// 	item.SellVolume1 = strings.TrimSpace(string(items[14]))

// 	item.BuyPrice2 = strings.TrimSpace(string(items[15]))
// 	item.BuyVolume2 = strings.TrimSpace(string(items[16]))
// 	item.SellPrice2 = strings.TrimSpace(string(items[17]))
// 	item.SellVolume2 = strings.TrimSpace(string(items[18]))

// 	item.BuyPrice3 = strings.TrimSpace(string(items[19]))
// 	item.BuyVolume3 = strings.TrimSpace(string(items[20]))
// 	item.SellPrice3 = strings.TrimSpace(string(items[21]))
// 	item.SellVolume3 = strings.TrimSpace(string(items[22]))

// 	item.BuyPrice4 = strings.TrimSpace(string(items[23]))
// 	item.BuyVolume4 = strings.TrimSpace(string(items[24]))
// 	item.SellPrice4 = strings.TrimSpace(string(items[25]))
// 	item.SellVolume4 = strings.TrimSpace(string(items[26]))

// 	item.BuyPrice5 = strings.TrimSpace(string(items[27]))
// 	item.BuyVolume5 = strings.TrimSpace(string(items[28]))
// 	item.SellPrice5 = strings.TrimSpace(string(items[29]))
// 	item.SellVolume5 = strings.TrimSpace(string(items[30]))

// 	if item.MDStreamID != "MD004" {
// 		item.TradingPhaseCode = strings.TrimSpace(string(items[31]))
// 		item.Timestamp = strings.TrimSpace(string(items[32]))
// 		return item
// 	}

// 	item.PreCloseIOPV = strings.TrimSpace(string(items[31]))
// 	item.IOPV = strings.TrimSpace(string(items[32]))

// 	item.TradingPhaseCode = strings.TrimSpace(string(items[33]))
// 	item.Timestamp = strings.TrimSpace(string(items[34]))

// 	return item
// }

// // header
// type MktdtHeader struct {
// 	BeginString        string // 起始标识符
// 	Version            string // 版本
// 	BodyLength         int    // 数据长度
// 	TotNumTradeReports int    // 文件体记录数
// 	MDReportID         int    // 行情序号
// 	SenderCompID       string // 发送方
// 	MDTime             string // 行情时间 格式为YYYYMMDD-HH:MM:SS.000
// 	MDUpdateType       int    // 发送方式 0 = 快照 Full refresh , 1 = 增量Incremental（暂不支持）
// 	MDSesStatus        string // 市场行情状态 第1位：‘S’表示全市场启动期间（开市前），‘T’表示全市场处于交易期间（含中间休市）， ‘E’表示全市场处于闭市期间。第2位：‘1’表示开盘集合竞价结束标志，未结束取‘0’。第3位：‘1’表示市场行情结束标志，未结束取‘0’。
// }

// func parseMktdtHeader(data []byte) *MktdtHeader {
// 	items := bytes.Split(data, []byte("|"))

// 	if len(items) < 9 {
// 		return nil
// 	}

// 	header := &MktdtHeader{}
// 	header.BeginString = strings.TrimSpace(string(items[0]))
// 	header.Version = strings.TrimSpace(string(items[1]))
// 	bodyLength, err := strconv.Atoi(strings.TrimSpace(string(items[2])))
// 	if err != nil {
// 		log.Error("parse bodyLength err %v ", err)
// 	}
// 	header.BodyLength = bodyLength
// 	totNumTradeReports, err := strconv.Atoi(strings.TrimSpace(string(items[3])))
// 	if err != nil {
// 		log.Error("parse totNumTradeReports err %v ", err)
// 	}
// 	header.TotNumTradeReports = totNumTradeReports

// 	mDReportIDStr := strings.TrimSpace(string(items[4]))
// 	if len(mDReportIDStr) == 0 {
// 		header.MDReportID = 0
// 	} else {
// 		mDReportID, err := strconv.Atoi(mDReportIDStr)
// 		if err != nil {
// 			log.Error("parse mDReportID err %v ", err)
// 		}
// 		header.MDReportID = mDReportID
// 	}

// 	header.SenderCompID = strings.TrimSpace(string(items[5]))
// 	header.MDTime = strings.TrimSpace(string(items[6]))
// 	mDUpdateType, err := strconv.Atoi(strings.TrimSpace(string(items[7])))
// 	if err != nil {
// 		log.Error("parse mDReportID err %v ", err)
// 	}
// 	header.MDUpdateType = mDUpdateType
// 	header.MDSesStatus = strings.TrimSpace(string(items[8]))

// 	return header
// }

// // footer
// type MktdtFooter struct {
// 	EndString string
// 	CheckSum  string
// }

// func parseMktdtFooter(data []byte) *MktdtFooter {
// 	items := bytes.Split(data, []byte("|"))

// 	if len(items) < 2 {
// 		return nil
// 	}

// 	footer := &MktdtFooter{}
// 	footer.EndString = strings.TrimSpace(string(items[0]))
// 	footer.CheckSum = strings.TrimSpace(string(items[1]))

// 	return footer

// }
