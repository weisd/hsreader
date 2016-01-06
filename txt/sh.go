package txt

import (
	"bufio"
	"bytes"
	// "fmt"
	"io"
	"os"

	"github.com/weisd/hsreader/mahonia"
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

func ReadBytes(data []byte) (contents TxtContent, err error) {
	return Read(bytes.NewReader(data))
}

func ReadFile(fpath string) (contents TxtContent, err error) {
	fp, err := os.OpenFile(fpath, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer fp.Close()

	return Read(fp)

}

func Read(fp io.Reader) (contents TxtContent, err error) {

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
	d := mahonia.NewDecoder("gbk")
	for k, _ := range vals {
		key := fileds[k]
		if key == "Symbol" {
			// gbk to utf-8
			info[key] = d.ConvertString(string(vals[k]))
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
