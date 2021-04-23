/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trade

import (
	"fmt"
	"github.com/Micahel754267513/pkg/hb/client"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
)

// 最新的交易
//	id	integer	唯一交易id（将被废弃）
//	trade-id	integer	唯一成交ID（NEW）
//	amount	float	以基础币种为单位的交易量
//	price	float	以报价币种为单位的成交价格
//	ts	integer	调整为新加坡时间的时间戳，单位毫秒
//	direction	string	交易方向：“buy” 或 “sell”, “buy” 即买，“sell” 即卖
func LastetTrade(coin string) (mt *market.TradeTick) {
	mt, err := client.MarketClient().GetLatestTrade(coin)
	if err != nil {
		fmt.Println(err)
	}

	return
}

//id	long	调整为新加坡时间的时间戳，单位秒，并以此作为此K线柱的id
//amount	float	以基础币种计量的交易量
//count	integer	交易次数
//open	float	本阶段开盘价
//close	float	本阶段收盘价
//low	float	本阶段最低价
//high	float	本阶段最高价
//vol	float	以报价币种计量的交易量
func HistoryTrade(coin string, size int) {
	client.MarketClient().GetHistoricalTrade(coin, market.GetHistoricalTradeOptionalRequest{Size: size})
}
