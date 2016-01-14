package hsreader

import (
	"fmt"
	"github.com/weisd/hsreader/dbf"
	"github.com/weisd/hsreader/mahonia"
	"github.com/weisd/hsreader/txt"
	"path"
	"strings"
	"time"
)

// var quoteFields = []string{
// 	"market",
// 	"code",
// 	"name",
// 	"turnRate", //转手率
// 	"turnVol",
// 	"vol",
// 	"sum",
// 	"preClose",
// 	"open",
// 	"close",
// 	"high",
// 	"low",
// 	"price",
// 	"status",
// 	"timeStamp",

// 	"buyp1",
// 	"buyv1",
// 	"buyp2",
// 	"buyv2",
// 	"buyp3",
// 	"buyv3",
// 	"buyp4",
// 	"buyv4",
// 	"buyp5",
// 	"buyv5",

// 	"sellp1",
// 	"sellv1",
// 	"sellp2",
// 	"sellv2",
// 	"sellp3",
// 	"sellv3",
// 	"sellp4",
// 	"sellv4",
// 	"sellp5",
// 	"sellv5",
// }

func ReadFromBytes(filename string, data []byte) ([]map[string]string, error) {
	// dbf or txt
	ext := path.Ext(filename)
	switch ext {
	case ".txt":
		return ReadTxtFromBytes(data)
	case ".dbf":
		return ReadDbfFromBytes(data)
	default:
		return nil, fmt.Errorf("file ext not supper")
	}

	return nil, fmt.Errorf("file ext not supper")

}

func ReadFromFile(filename string) ([]map[string]string, error) {
	// dbf or txt
	ext := path.Ext(filename)
	switch ext {
	case ".txt":
		return ReadTxtFromFile(filename)
	case ".dbf":
		return ReadDbfFromFile(filename)
	default:
		return nil, fmt.Errorf("file ext not supper")
	}

	return nil, fmt.Errorf("file ext not supper")

}

func ReadDbfFromFile(filename string) ([]map[string]string, error) {
	data, err := dbf.FromFile(filename)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, err
	}

	return parseDbfData(data)
}

func ReadDbfFromBytes(d []byte) ([]map[string]string, error) {
	data, err := dbf.FromBytes(d)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, err
	}

	return parseDbfData(data)
}

func parseDbfData(data []map[string]string) ([]map[string]string, error) {
	// 是否sh
	tmp := data[0]
	if _, ok := tmp["S1"]; ok {
		return parseShData(data)
	}

	if _, ok := tmp["HQZQDM"]; ok {
		return parseSzData(data)
	}

	return nil, fmt.Errorf("unknow dbf file")
}

func parseShData(data []map[string]string) ([]map[string]string, error) {

	statusInfo := map[string]string{}

	for _, item := range data {
		if item["S1"] == "000000" {
			statusInfo = item
			break
		}
	}

	if len(statusInfo) == 0 {
		return nil, fmt.Errorf("parseShData status info no found")
	}

	if len(statusInfo["S2"]) < 6 {
		statusInfo["S2"] = "0" + statusInfo["S2"]
	}

	t, err := time.Parse("150405", statusInfo["S2"])
	if err != nil {
		fmt.Println("time.Parse(150405, statusInfo[S2]) %v", err)
		err = nil
	}

	ctime := t.Format("15:04:05")

	list := []map[string]string{}

	d := mahonia.NewDecoder("gbk")

	for _, item := range data {
		if item["S1"] == "000000" {
			continue
		}

		code := item["S1"]
		// return nil, nil
		switch string(code[0:3]) {
		case "000", "600", "601", "603", "900":
		default:
			continue
		}

		info := parsesFields(item, shfiledsMap)
		info["market"] = "sh"
		info["code"] = info["market"] + info["code"]
		info["timeStamp"] = ctime

		info["name"] = d.ConvertString(info["name"])

		info["status"] = ""
		if item["delete"] == "1" {
			info["status"] = "P"
		}

		list = append(list, info)
	}

	return list, nil
}

func parseSzData(data []map[string]string) ([]map[string]string, error) {

	statusInfo := map[string]string{}

	for _, item := range data {
		if item["HQZQDM"] == "000000" {
			statusInfo = item
			break
		}
	}

	if len(statusInfo) == 0 {
		return nil, fmt.Errorf("parseShData status info no found")
	}

	if len(statusInfo["HQCJBS"]) < 6 {
		statusInfo["HQCJBS"] = "0" + statusInfo["HQCJBS"]
	}

	t, err := time.Parse("150405", statusInfo["HQCJBS"])
	if err != nil {
		fmt.Println("time.Parse(150405, statusInfo[HQCJBS])", len(statusInfo["HQCJBS"]), err)
		err = nil
	}

	ctime := t.Format("15:04:05")

	list := []map[string]string{}

	d := mahonia.NewDecoder("gbk")

	for _, item := range data {
		if item["HQZQDM"] == "000000" {
			continue
		}

		code := item["HQZQDM"]
		// return nil, nil
		switch string(code[0:3]) {
		case "000", "300", "002", "399", "200", "001", "400":
		default:
			continue
		}

		info := parsesFields(item, szfiledsMap)
		info["market"] = "sz"
		info["code"] = info["market"] + info["code"]
		info["timeStamp"] = ctime

		info["name"] = d.ConvertString(info["name"])

		info["status"] = ""
		if item["delete"] == "1" {
			info["status"] = "P"
		}

		list = append(list, info)
	}

	return list, nil
}

var shfiledsMap = map[string]string{
	"code":     "S1",
	"name":     "S2",
	"vol":      "S11",
	"sum":      "S5",
	"preClose": "S3",
	"open":     "S4",
	"close":    "S8",
	"high":     "S6",
	"low":      "S7",
	"price":    "S8",

	"buyp1": "S14",
	"buyv1": "S15",
	"buyp2": "S16",
	"buyv2": "S17",
	"buyp3": "S18",
	"buyv3": "S19",
	"buyp4": "S26",
	"buyv4": "S27",
	"buyp5": "S28",
	"buyv5": "S29",

	"sellp1": "S20",
	"sellv1": "S21",
	"sellp2": "S22",
	"sellv2": "S23",
	"sellp3": "S24",
	"sellv3": "S25",
	"sellp4": "S30",
	"sellv4": "S31",
	"sellp5": "S32",
	"sellv5": "S33",
}

var szfiledsMap = map[string]string{
	"code":     "HQZQDM",
	"name":     "HQZQJC",
	"vol":      "HQCJSL",
	"sum":      "HQCJJE",
	"preClose": "HQZRSP",
	"open":     "HQJRKP",
	"close":    "HQZJCJ",
	"high":     "HQZGCJ",
	"low":      "HQZDCJ",
	"price":    "HQZJCJ",

	"buyp1": "HQBJW1",
	"buyv1": "HQBSL1",
	"buyp2": "HQBJW2",
	"buyv2": "HQBSL2",
	"buyp3": "HQBJW3",
	"buyv3": "HQBSL3",
	"buyp4": "HQBJW4",
	"buyv4": "HQBSL4",
	"buyp5": "S28",
	"buyv5": "S29",

	"sellp1": "HQSJW1",
	"sellv1": "HQSSL1",
	"sellp2": "HQSJW2",
	"sellv2": "HQSSL2",
	"sellp3": "HQSJW3",
	"sellv3": "HQSSL3",
	"sellp4": "HQSJW4",
	"sellv4": "HQSSL4",
	"sellp5": "S32",
	"sellv5": "S33",
}

func parsesFields(item map[string]string, fieldsMap map[string]string) map[string]string {
	info := map[string]string{}
	for k, v := range fieldsMap {
		if _, ok := item[v]; !ok {
			info[k] = ""
			continue
		}

		info[k] = strings.Replace(item[v], "-", "0", -1)
	}
	return info
}

func ReadTxtFromBytes(data []byte) ([]map[string]string, error) {
	contents, err := txt.ReadBytes(data)
	if err != nil {
		return nil, err
	}

	return parseTxtData(contents)
}

func ReadTxtFromFile(filename string) ([]map[string]string, error) {
	contents, err := txt.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parseTxtData(contents)
}

func parseTxtData(contents txt.TxtContent) ([]map[string]string, error) {

	list := []map[string]string{}

	// 处理指数
	for _, item := range contents.IndexList {
		info := map[string]string{}
		info["market"] = "sh"
		info["code"] = info["market"] + item["SecurityID"]
		info["name"] = item["Symbol"]
		info["vol"] = item["TradeVolume"]
		info["sum"] = item["TotalValueTraded"]
		info["preClose"] = item["PreClosePx"]
		info["open"] = item["OpenPrice"]
		info["close"] = item["ClosePx"]
		info["high"] = item["HighPrice"]
		info["low"] = item["LowPrice"]
		info["price"] = item["TradePrice"]

		status := item["TradingPhaseCode"]
		if len(status) > 0 {
			info["status"] = string(status[0])
		} else {
			info["status"] = ""
		}

		ts := item["Timestamp"]
		info["timeStamp"] = string(ts[0:8])

		info["buyp1"] = ""
		info["buyv1"] = ""
		info["buyp2"] = ""
		info["buyv2"] = ""
		info["buyp3"] = ""
		info["buyv3"] = ""
		info["buyp4"] = ""
		info["buyv4"] = ""
		info["buyp5"] = ""
		info["buyv5"] = ""

		info["sellp1"] = ""
		info["sellv1"] = ""
		info["sellp2"] = ""
		info["sellv2"] = ""
		info["sellp3"] = ""
		info["sellv3"] = ""
		info["sellp4"] = ""
		info["sellv4"] = ""
		info["sellp5"] = ""
		info["sellv5"] = ""

		list = append(list, info)
	}

	// 处理指数
	for _, item := range contents.StockList {
		info := map[string]string{}
		info["market"] = "sh"
		info["code"] = info["market"] + item["SecurityID"]
		info["name"] = item["Symbol"]
		info["vol"] = item["TradeVolume"]
		info["sum"] = item["TotalValueTraded"]
		info["preClose"] = item["PreClosePx"]
		info["open"] = item["OpenPrice"]
		info["close"] = item["ClosePx"]
		info["high"] = item["HighPrice"]
		info["low"] = item["LowPrice"]
		info["price"] = item["TradePrice"]

		status := item["TradingPhaseCode"]
		if len(status) > 0 {
			info["status"] = string(status[0])
		} else {
			info["status"] = ""
		}
		ts := item["Timestamp"]

		info["timeStamp"] = string(ts[0:8])

		info["buyp1"] = item["BuyPrice1"]
		info["buyv1"] = item["BuyVolume1"]
		info["buyp2"] = item["BuyPrice2"]
		info["buyv2"] = item["BuyVolume2"]
		info["buyp3"] = item["BuyPrice3"]
		info["buyv3"] = item["BuyVolume3"]
		info["buyp4"] = item["BuyPrice4"]
		info["buyv4"] = item["BuyVolume4"]
		info["buyp5"] = item["BuyPrice5"]
		info["buyv5"] = item["BuyVolume5"]

		info["sellp1"] = item["SellPrice1"]
		info["sellv1"] = item["SellVolume1"]
		info["sellp2"] = item["SellPrice2"]
		info["sellv2"] = item["SellVolume2"]
		info["sellp3"] = item["SellPrice3"]
		info["sellv3"] = item["SellVolume3"]
		info["sellp4"] = item["SellPrice4"]
		info["sellv4"] = item["SellVolume4"]
		info["sellp5"] = item["SellPrice5"]
		info["sellv5"] = item["SellVolume5"]

		list = append(list, info)
	}

	return list, nil

}
