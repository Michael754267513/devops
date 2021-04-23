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
	"github.com/Micahel754267513/pkg/hb/client"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/shopspring/decimal"
)

func GetCoinPercent(symbol string, ms []market.SymbolCandlestick) (p decimal.Decimal) {
	for _, c := range ms {
		if c.Symbol == symbol {
			p = c.Close.Sub(c.Open).Div(c.Open)
		}
	}
	return
}

func GetAllTicks() (ms []market.SymbolCandlestick, err error) {
	ms, err = client.MarketClient().GetAllSymbolsLast24hCandlesticksAskBid()
	return
}
