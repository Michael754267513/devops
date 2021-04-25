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
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/huobirdcenter/huobi_golang/pkg/model/order"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
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
	Price          decimal.Decimal            `json:"price"`      // 当前价格
	Price24H       decimal.Decimal            `json:"price_24_h"` // 24h最高价
	Price24L       decimal.Decimal            `json:"price_24_l"` // 24h最低价
	Price24O       decimal.Decimal            `json:"price_24_o"` // 开盘价
}

func (t *TradeCoins) GetCoinPercent() (p decimal.Decimal) {
	symbol := fmt.Sprintf("%s%s", t.Coin, t.CoinType)
	for _, c := range t.AllTicks {
		if c.Symbol == symbol {
			p = c.Close.Sub(c.Open).Div(c.Open)
			t.Price = c.Close
			t.Price24H = c.High
			t.Price24L = c.Low
			t.Price24O = c.Open
		}
	}
	return
}

// 获取当前正负百分比
func (t *TradeCoins) Get1Percent() {
	if len(t.Min1Percent) > 10 {
		t.Min1Percent = t.Min1Percent[len(t.Min1Percent)-10:]
	}
	pp := t.GetCoinPercent()
	t.LastestPercent = pp
	t.Min1Percent = append(t.Min1Percent, pp)

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
	fmt.Println(fmt.Sprintf("USDT余额：%s  %s余额：%s 当前价:%s 开盘价:%s 最高价:%s 最低价:%s ",
		fmt.Sprint(t.UsdtBalance),
		t.Coin, fmt.Sprint(t.CoinBalance),
		fmt.Sprint(t.Price),
		fmt.Sprint(t.Price24O),
		fmt.Sprint(t.Price24H),
		fmt.Sprint(t.Price24L)))
}

// 按照USDT 数量购买
func (t *TradeCoins) BuyCoin(usdt float64) {
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
func (t *TradeCoins) SellCoin(coin float64) {
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
		//10% 不关心
		if b10 < 0.010 {
			return
		}
		avg := decimal.Avg(b0, t.Min1Percent...) // 平均数
		sub, _ := (t.Min1Percent[len(t.Min1Percent)-2].Sub(t.LastestPercent)).Float64()
		sub, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", sub), 64) // 最后两次差值
		if !avg.GreaterThanOrEqual(b0) && t.LastestPercent.GreaterThanOrEqual(b0) {
			if t.Status == true {
				if t.UsdtBalance >= 6 {
					fmt.Println(fmt.Sprintf("%s 购买:%s 数量:%s USDT", fmt.Sprint(time.Now()), fmt.Sprint(t.Coin), fmt.Sprint(t.UsdtBalance)))
					t.BuyCoin(t.UsdtBalance)
					t.Status = false
				}
			}
		}
		// 判断平均数是否小于最后的数
		if !avg.GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		// 判断最大数是否小于等于最新值
		if !decimal.Max(b0, t.Min1Percent...).GreaterThanOrEqual(t.LastestPercent) {
			return
		}

		// 涨多少
		if sub > 0.010 {
			fmt.Println(fmt.Sprintf("%s 涨 USDT: %s  %s:%s %s", fmt.Sprint(time.Now()), fmt.Sprint(t.UsdtBalance), fmt.Sprint(t.Coin),
				fmt.Sprint(t.CoinBalance), fmt.Sprint(sub)))
			//t.Status = true
		}
		// 跌多少
		if b10 <= 0.50 {
			if sub < 0.010 {
				fmt.Println(fmt.Sprintf("%s 出售:%s 数量:%s", fmt.Sprint(time.Now()), fmt.Sprint(t.Coin), fmt.Sprint(t.CoinBalance)))
				t.SellCoin(t.CoinBalance)
				t.Status = true
			}
		}

		if b10 >= 0.50 {
			if sub < 0.010*2.0 {
				fmt.Println(fmt.Sprintf("%s 出售:%s 数量:%s", fmt.Sprint(time.Now()), fmt.Sprint(t.Coin), fmt.Sprint(t.CoinBalance)))
				t.SellCoin(t.CoinBalance)
				t.Status = true
			}
		}

	} else {
		b_10, _ := t.LastestPercent.Float64()
		avg := decimal.Avg(b0, t.Min1Percent...) // 平均数
		sub, _ := (t.Min1Percent[len(t.Min1Percent)-2].Sub(t.LastestPercent)).Float64()
		sub, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", sub), 64) // 最后两次差值
		// 平均值小于0 啥事不干
		if avg.GreaterThanOrEqual(b0) {
			return
		}

		// -10% 不关心
		if b_10 > -0.10 {
			return
		}
		// allin
		if b_10 < -0.20 {
			if sub > 0.03 {
				return
			}
			if !avg.GreaterThanOrEqual(t.LastestPercent) {
				return
			}
			if t.Status == true {
				if t.UsdtBalance >= 6 {
					fmt.Println(fmt.Sprintf("%s 购买:%s 数量:%s USDT", fmt.Sprint(time.Now()), fmt.Sprint(t.Coin), fmt.Sprint(t.UsdtBalance)))
					t.BuyCoin(t.UsdtBalance)
					t.Status = false
				}
			}
		}
		if avg.GreaterThanOrEqual(t.LastestPercent) {
			return
		}
		if !decimal.Max(b0, t.Min1Percent...).GreaterThanOrEqual(t.LastestPercent) {
			return
		}

		if sub > 0.05 {
			fmt.Println(fmt.Sprintf("%s 涨 USDT: %s  %s:%s %s", fmt.Sprint(time.Now()), t.UsdtBalance, t.Coin, t.CoinBalance, sub))
			//
			//t.Status = true
		}

		if sub < 0.05 {
			fmt.Println(fmt.Sprintf("%s 出售:%s 数量:%s", fmt.Sprint(time.Now()), fmt.Sprint(t.Coin), fmt.Sprint(t.CoinBalance)))
			t.SellCoin(t.CoinBalance)
			t.Status = true
		}

	}

}
