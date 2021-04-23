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

package client

import (
	"github.com/Micahel754267513/pkg/hb/config"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
)

func PrivateAccount() (c *client.AccountClient) {
	c = new(client.AccountClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}

func WalletClient() (c *client.WalletClient) {
	c = new(client.WalletClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}

func OrderClient() (c *client.OrderClient) {
	c = new(client.OrderClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}

func IsolatedMarginClient() (c *client.IsolatedMarginClient) {
	c = new(client.IsolatedMarginClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}

func CrossMarginClient() (c *client.CrossMarginClient) {
	c = new(client.CrossMarginClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}

func ETFClient() (c *client.ETFClient) {
	c = new(client.ETFClient).Init(config.AccessKey, config.SecretKey, config.Host)
	return
}
