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
	"github.com/Micahel754267513/pkg/cronjob"
	"github.com/Micahel754267513/pkg/hb"
	"github.com/robfig/cron/v3"
)

func main() {
	hb.RunHB("doge")
	hb.RunHB("bch3l")
	hb.RunHB("pvt")
	select {}
}

func init() {
	cronjob.J.C = cron.New(cron.WithSeconds())
	cronjob.J.C.Start()
}
