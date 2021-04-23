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
	"github.com/Micahel754267513/pkg/hb/config"
	"github.com/huobirdcenter/huobi_golang/pkg/model/order"
)

func BuyCoin() {
	c := client.OrderClient()
	r := &order.PlaceOrderRequest{
		AccountId: config.AccountId,
		Symbol:    "dogeusdt",   // 币种
		Type:      "buy-market", // 买
		Amount:    "5",
	}
	opr, _ := c.PlaceOrder(r)

	fmt.Println(opr)
	//834
}

func SellCoin() {
	c := client.OrderClient()
	r := &order.PlaceOrderRequest{
		AccountId: config.AccountId,
		Symbol:    "dogeusdt",    // 币种
		Type:      "sell-market", // 卖
		Amount:    "20",
	}
	opr, err := c.PlaceOrder(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(opr)
	//834
}
