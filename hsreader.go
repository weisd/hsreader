package hsreader

import (
	"./dbf"
	"./txt"
	"fmt"
	"path"
)

var quoteFields = []string{
	"market",
	"code",
	"name",
	"turnRate", //转手率
	"turnVol",
	"vol",
	"sum",
	"preClose",
	"open",
	"close",
	"high",
	"low",
	"price",
	"status",
	"timeStamp",

	"buyp1",
	"buyv1",
	"buyp2",
	"buyv2",
	"buyp3",
	"buyv3",
	"buyp4",
	"buyv4",
	"buyp5",
	"buyv5",

	"sellp1",
	"sellv1",
	"sellp2",
	"sellv2",
	"sellp3",
	"sellv3",
	"sellp4",
	"sellv4",
	"sellp5",
	"sellv5",
}

// var txtFieldsMap = map[string]string{
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
	data, err := dbf.FromFile()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, err
	}

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

	return nil, nil
}

func parseSzData(data []map[string]string) ([]map[string]string, error) {

	return nil, nil
}

func ReadTxtFromFile(filename string) ([]map[string]string, error) {
	contents, err := txt.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	list := []map[string]string{}

	// 处理指数
	for _, item := range contents.IndexList {
		info := map[string]string{}
		info["market"] = "sh"
		info["code"] = info["market"] + item["SecurityID"]
		info["name"] = item["Symbol"]
		info["turnRate"] = ""
		info["turnVol"] = ""
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
		info["turnRate"] = ""
		info["turnVol"] = ""
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
