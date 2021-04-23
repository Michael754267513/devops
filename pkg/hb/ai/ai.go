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

package ai

import (
	"fmt"
	"github.com/Micahel754267513/pkg/hb/client"
	"github.com/Micahel754267513/pkg/hb/config"
	"github.com/Micahel754267513/pkg/hb/trade"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/huobirdcenter/huobi_golang/pkg/model/order"
	"github.com/shopspring/decimal"
	"strconv"
)

type TradeCoins struct {
	AccountID      string                     `json:"account_id"`      // 账号id
	Coin           string                     `json:"coin"`            // 币种
	CoinType       string                     `json:"coin_type"`       // 币类型
	CoinBalance    float64                    `json:"coin_balance"`    // 余额
	UsdtBalance    float64                    `json:"usdt_balance"`    // USDT 余额
	AllTicks       []market.SymbolCandlestick `json:"all_ticks"`       // 所有币种当前数据
	Min1Percent    []decimal.Decimal          `json:"min_1_percent"`   // 每分钟获取比重当前百分比
	LastestPercent decimal.Decimal            `json:"lastest_percent"` // 最新数据
	SellPercent    decimal.Decimal            `json:"sell_percent"`
	BuyPercent     decimal.Decimal            `json:"buy_percent"`
	Status         bool                       `json:"status"`
}

// 获取当前正负百分比
func (t *TradeCoins) Get1Percent() {
	if len(t.Min1Percent) > 10 {
		t.Min1Percent = t.Min1Percent[len(t.Min1Percent)-10:]
	}
	pp := trade.GetCoinPercent(fmt.Sprintf("%s%s", t.Coin, t.CoinType), t.AllTicks)
	t.LastestPercent = pp
	t.Min1Percent = append(t.Min1Percent, pp)
	fmt.Println(t.Min1Percent)
}

// 获取当前币种可交易余额
func (t *TradeCoins) GetCoinBalance() {
	cp := client.PrivateAccount()
	coins, _ := cp.GetAccountBalance(t.AccountID)
	for _, c := range coins.List {
		if c.Currency == t.Coin && c.Type == "trade" {
			t.CoinBalance, _ = strconv.ParseFloat(c.Balance, 64)
		}
	}

}

func (t *TradeCoins) GetUSDTBalance() {
	cp := client.PrivateAccount()
	coins, _ := cp.GetAccountBalance(fmt.Sprint(t.AccountID))
	for _, c := range coins.List {
		if c.Currency == "usdt" && c.Type == "trade" {
			t.UsdtBalance, _ = strconv.ParseFloat(c.Balance, 64)
		}
	}
}

// 获取所有ticks
func (t *TradeCoins) GetAllTicks() (err error) {
	t.AllTicks, err = client.MarketClient().GetAllSymbolsLast24hCandlesticksAskBid()
	return
}

func (t *TradeCoins) Init() {
	t.CoinType = "usdt"
	t.AccountID = config.AccountId
	if t.GetAllTicks() != nil {
		return
	}
	t.Get1Percent()
	t.GetCoinBalance()
	t.GetUSDTBalance()
	t.AIBuySell()
}

// 按照USDT 数量购买
func (t *TradeCoins) BuyCoin(usdt int) {
	c := client.OrderClient()
	r := &order.PlaceOrderRequest{
		AccountId: config.AccountId,
		Symbol:    fmt.Sprintf("%s%s", t.Coin, t.CoinType),
		Type:      "buy-market", // 买
		Amount:    fmt.Sprint(usdt),
	}
	opr, _ := c.PlaceOrder(r)
	fmt.Println(opr)
}

// 根据所持有货币的数量售卖
func (t *TradeCoins) SellCoin(coin int) {
	c := client.OrderClient()
	r := &order.PlaceOrderRequest{
		AccountId: config.AccountId,
		Symbol:    fmt.Sprintf("%s%s", t.Coin, t.CoinType),
		Type:      "buy-market", // 买
		Amount:    fmt.Sprint(coin),
	}
	opr, _ := c.PlaceOrder(r)
	fmt.Println(opr)
}

// 计算买卖
func (t *TradeCoins) AIBuySell() {
	if len(t.Min1Percent) < 5 {
		return
	}
	var b0 decimal.Decimal

	if t.LastestPercent.GreaterThanOrEqual(b0) {
		b10, _ := t.LastestPercent.Float64()
		// 10% 不关心
		if b10 < 0.10 {
			return
		}
		avg := decimal.Avg(b0, t.Min1Percent...)
		// 判断平均数是否小于最后的数
		if !avg.GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		// 判断最大数是否小于等于最新值
		if !decimal.Max(b0, t.Min1Percent...).GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		sub, _ := t.Min1Percent[len(t.Min1Percent)-2].Sub(t.LastestPercent).Float64()
		if sub > 0.0005 {
			fmt.Println(fmt.Sprintf("涨 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}
		if sub < 0.0005 {
			fmt.Println(fmt.Sprintf("跌 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}
	} else {
		b_10, _ := t.LastestPercent.Float64()
		sub, _ := t.Min1Percent[len(t.Min1Percent)-2].Sub(t.LastestPercent).Float64()
		if sub > 0.0005 {
			fmt.Println(fmt.Sprintf(" 涨 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}
		if sub < 0.0005 {
			fmt.Println(fmt.Sprintf(" 跌 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}

		// -10% 不关心
		if b_10 > -0.10 {
			return
		}
		avg := decimal.Avg(b0, t.Min1Percent...)
		if avg.GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		if !decimal.Max(b0, t.Min1Percent...).GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		sub, _ = t.Min1Percent[len(t.Min1Percent)-2].Sub(t.LastestPercent).Float64()
		// 跌的多
		if sub > 0.0005 {
			fmt.Println(fmt.Sprintf(" 涨 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}
		if sub < 0.0005 {
			fmt.Println(fmt.Sprintf(" 跌 USDT: %s  %s:%s %s", t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			fmt.Println("sell")
			t.Status = true
		}

	}

}
