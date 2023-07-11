package requests

import (
	"SGX_blockchain_client/src/crypto"
	"SGX_blockchain_client/src/utils"
	"bytes"
	"fmt"
	"github.com/tidwall/pretty"
	"io"
	"net/http"
	"time"
)

const jsonContentType string = "application/json"

type AccountInfoRequest struct {
	Data struct {
		Address string `json:"address"`
		Ts      int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type SingleAccount struct {
	Keypair *crypto.KeyPair
	Url     string
}

func CreateNewSingleAccount(url string, pk string) *SingleAccount {
	keypair := crypto.Initialize(utils.DecodeHexStringToBytesWith0x(pk))

	return &SingleAccount{Keypair: keypair, Url: url}
}

func (account *SingleAccount) GetAccountsInfo() string {

	accountsRequest := &AccountInfoRequest{
		Data: struct {
			Address string `json:"address"`
			Ts      int64  `json:"ts"`
		}{
			Address: utils.EncodeBytesToHexStringWith0x(account.Keypair.PublicKey.SerializeUncompressed()),
			Ts:      time.Now().UnixMilli(),
		},
		Signature: "",
	}

	bodyBytes, err := utils.SignJsonWithData(accountsRequest, account.Keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(account.Url+"/account/info", jsonContentType, buffer)

	if err != nil {
		fmt.Println(err)
		return "Wrong!"
	} else {
		fmt.Println("账户信息读取正常")
		body, _ := io.ReadAll(resp.Body)
		result := pretty.Pretty(body)

		return string(result)
	}
	defer resp.Body.Close()
	return ""

}
