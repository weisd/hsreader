package dbf

// http://www.clicketyclick.dk/databases/xbase/format/dbf.html#DBF_NOTE_9_TARGET
// show2003.dbf
/*
字段	类型	长度	小数位数	为空
字段名	字段说明	类型	长度	小数位数
S1	证券代码	Char	6
S2	证券名称	Char	8
S3	前收盘价格	Num	8	3
S4	今开盘价格	Num	8	3
S5	今成交金额	Num	12
S6	最高价	Num	8	3
S7	最低价	Num	8	3
S8	最新价	Num	8	3
S9	当前买入价	Num	8	3
S10	当前卖出价	Num	8	3
S11	成交数量	Num	10
S13	市盈率	Num	8	3
S15	申买量一	Num	10
S16	申买价二	Num	8	3
S17	申买量二	Num	10
S18	申买价三	Num	8	3
S19	申买量三	Num	10
S21	申卖量一	Num	10
S22	申卖价二	Num	8	3
S23	申卖量二	Num	10
S24	申卖价三	Num	8	3
S25	申卖量三	Num	10
S26	申买价四	Num	8	3
S27	申买量四	Num	10
S28	申买价五	Num	8	3
S29	申买量五	Num	10
S30	申卖价四	Num	8	3
S31	申卖量四	Num	10
S32	申卖价五	Num	8	3
S33	申卖量五	Num	10
*/
/*
数据库内容

1. 整体行情：第1条记录

备注：
1除上述字段外，其他字段为空。
2.记录示例：2002年9月11日下午15时05分13秒时，表中的第一条记录如下：



2. 分类指数（证券代码为000001－000013，测试数据为000300）：第2－15条记录

备注：
分类指数包括：上证指数、A股指数、B股指数、工业指数、商业指数、地产指数、公用指数、综合指数、上证180、基金指数、国债指数、企债指数和测试数据。
参与计算相应指数的交易数量（S11）的单位和参与计算的证券类型相关。证券类型是股票的指数交易数量单位是100股，基金指数的交易数量单位是100份，债券指数的交易数量单位是手。
除上述字段外，其他字段内容为空。
*/

// sjshq.dbf
/*
序号	字段名	字段描述	类型	长度	备注
1	HQZQDM	证券代码	C	6
2	HQZQJC	证券简称	C	8
3	HQZRSP	昨日收盘价	N	9,3
4	HQJRKP	今日开盘价	N	9,3
5	HQZJCJ	最近成交价	N	9,3
6	HQCJSL	成交数量	N	12,0
7	HQCJJE	成交金额	N	17,3
8	 HQCJBS	成交笔数	N	9,0
9	HQZGCJ	最高成交价	N	9,3
10	HQZDCJ	最低成交价	N	9,3
11	HQSYL1	市盈率1	N	7,2
12	HQSYL2	市盈率2	N	7,2
13	 HQJSD1	价格升跌1	N	9,3
14	 HQJSD2	价格升跌2	N	9,3
15	HQHYCC	合约持仓量	N	12,0
16	HQSJW4	卖价位四	N	9,3
17	HQSSL4	卖数量四	N	12,0
18	HQSJW3	卖价位三	N	9,3
19	HQSSL3	卖数量三	N	12,0
20	HQSJW2	卖价位二	N	9,3
21	HQSSL2	卖数量二	N	12,0
22	HQSJW1	卖价位一/叫卖揭示价	N	9,3
23	HQSSL1	卖数量一	N	12,0
24	HQBJW1	买价位一/叫买揭示价	N	9,3
25	HQBSL1	买数量一	N	12,0
26	HQBJW2	买价位二	N	9,3
27	HQBSL2	买数量二	N	12,0
28	HQBJW3	买价位三	N	9,3
29	HQBSL3	买数量三	N	12,0
30	HQBJW4	买价位四	N	9,3
31	HQBSL4	买数量四	N	12,0
*/

/*
每个交易日本库的第一条记录为特殊记录，HQZQDM为“000000”，HQZQJC存放当前日期CCYYMMDD，HQCJBS存放当前时间HHMMSS，HQZRSP存放“指数因子”，HQCJSL存放行情状态信息。当本库的记录为指数记录时（HQZQDM的最左两位为39），相应的字段HQZRSP、HQJRKP、HQZJCJ、HQZGCJ、HQZDCJ等都必须乘上该“指数因子”，计算出的结果才为实际的指数值。HQCJSL个位数存放收市行情标志（0：非收市行情；1：表示收市行情），HQCJSL十位数存放正式行情与测试行情标志(0：正式行情；1：表示测试行情)，即HQCJSL值为0时表示正式非收市行情，为1时表示正式收市行情，为10时表示测试的非收市行情；为11时表示测试的收市行情。
HQZQDM(字段1：证券代码)为关键字。
HQJSD1(字段13：升跌一) = HQZJCJ(字段5：最近成交价) - HQZRSP(字段3：昨收盘价)。
HQJSD2(字段14：升跌二) = HQZJCJ(最近成交价) - 上笔成交价。
对于当天的第一笔成交：
HQJSD2(升跌二)= HQZJCJ(最近成交价) - HQZRSP(昨收盘价)。
HQHYCC(字段15：合约持仓量)目前未用。
卖盘三至卖盘一（HQSJW3，HQSSL3，HQSJW2，HQSSL2，HQSJW1，HQSSL1）为实时最低三个价位卖出申报价和数量，买盘一至买盘三（HQBJW1，HQBSL1，HQBJW2，HQBSL2，HQBJW3，HQBSL3）为实时最高三个价位买入申报价和数量。    HQSJW4(字段16：卖价位四) 、HQSSL4(字段17：卖数量四)为有效竞价范围内除卖价位一、二、三之外的所有卖委托的加权平均价和总量；HQBJW4(字段30：买价位四) 、HQBSL4(字段31：买数量四)为有效竞价范围内除买价位一、二、三之外的所有买委托的加权平均价和总量。
如买卖盘价格字段都不为空（都大于零），字段值有如下关系：HQSJW4>HQSJW3>HQSJW2>HQSJW1，HQBJW4<HQBJW3<HQBJW2< HQBJW1。
*/
