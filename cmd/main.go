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

package main

import (
	"github.com/Micahel754267513/pkg/hb/ai"
	"time"
)

func main() {
	//ms,err  := trade.GetAllTicks()
	//if err !=nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("BCH USDT: ",fmt.Sprint(trade.GetCoinPercent("bchusdt",ms)))
	//fmt.Println("DOGE USDT: ",fmt.Sprint(trade.GetCoinPercent("dogeusdt",ms)))
	//fmt.Println("HIT USDT: ",fmt.Sprint(trade.GetCoinPercent("hitusdt",ms)))
	//fmt.Println("BCH3l USDT: ",fmt.Sprint(trade.GetCoinPercent("bch3lusdt",ms)))
	//// 做空
	//fmt.Println("BCH3s USDT: ",fmt.Sprint(trade.GetCoinPercent("bch3susdt",ms)))
	//
	//account.GetBalance(config.AccountId,"doge")
	//
	//trade.BuyCoin()

	//trade.SellCoin()

	//cc := client.PrivateAccount()
	//aa,_ := cc.GetAccountInfo()
	//for _,a :=range aa {
	//	fmt.Println(a.Id)
	//	 coins,_ := cc.GetAccountBalance(fmt.Sprint(a.Id))
	//	 for _,c1 := range coins.List  {
	//	 	if  c1.Currency == "usdt" && c1.Type == "trade"{
	//	 		fmt.Println(c1)
	//		}
	//	 }
	//}
	var cc ai.TradeCoins
	cc.Coin = "doge"
	for i := 0; i < 20; i++ {
		cc.Init()
		time.Sleep(1 * time.Second)
	}

}
