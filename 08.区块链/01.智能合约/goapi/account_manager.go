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
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

//import "github.com/ethereum/go-ethereum/common"

func main() {
	ks := keystore.NewKeyStore("./key.store", keystore.StandardScryptN, keystore.StandardScryptP)
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)
	fmt.Println(am.Accounts())

	//// Create a new account with the specified encryption passphrase.
	//newAcc, _ := ks.NewAccount("Creation password")
	//fmt.Println(newAcc)
	//
	//// Export the newly created account with a different passphrase. The returned
	//// data from this method invocation is a JSON encoded, encrypted key-file.
	//jsonAcc, _ := ks.Export(newAcc, "Creation password", "Export password")
	//
	//// Update the passphrase on the account created above inside the local keystore.
	//_ = ks.Update(newAcc, "Creation password", "Update password")
	//
	//// Delete the account updated above from the local keystore.
	//_ = ks.Delete(newAcc, "Update password")
	//
	//// Import back the account we've exported (and then deleted) above with yet
	//// again a fresh passphrase.
	//impAcc, _ := ks.Import(jsonAcc, "Export password", "Import password")
	//impAcc.Address.Value()

}
