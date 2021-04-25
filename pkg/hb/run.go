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

package hb

import (
	"fmt"
	"github.com/Micahel754267513/pkg/cronjob"
	ai2 "github.com/Micahel754267513/pkg/hb/ai"
)

func RunHB(ctype string) {
	_, err := cronjob.J.AddCronJob(GetFunc(ctype), "* * * * * *")
	if err != nil {
		fmt.Println(err)
	}
}

func GetFunc(ctype string) (f func()) {
	var ai ai2.TradeCoins
	ai.Coin = ctype
	f = func() {
		ai.Coin = ctype
		ai.Init()
	}

	return
}
