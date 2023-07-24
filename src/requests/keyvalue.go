package requests

import (
	"SGX_blockchain_client/src/utils"
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"io"
	"net/http"
	"time"
)

type KVStoreRequest struct {
	Data struct {
		From  string `json:"from"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Nonce int64  `json:"nonce"`
		Ts    int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type KVRetrieveRequest struct {
	Data struct {
		From string `json:"from"`
		Key  string `json:"key"`
		Ts   int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

func (account *SingleAccount) StoreKV(key, value string) ([]byte, string, int64) {
	kvstorerequest := &KVStoreRequest{
		Data: struct {
			From  string `json:"from"`
			Key   string `json:"key"`
			Value string `json:"value"`
			Nonce int64  `json:"nonce"`
			Ts    int64  `json:"ts"`
		}{
			From:  utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			Key:   key,
			Value: value,
			Nonce: 1,
			Ts:    time.Now().UnixMilli(),
		},
		Signature: "",
	}
	bodyBytes, err := utils.SignJsonWithData(kvstorerequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	resp, error := http.Post(account.Url+"/kv/store", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!", 0
	} else {
		fmt.Println("kv存储正常")
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(body)
		result := pretty.Pretty(body)
		res := gjson.GetBytes(body, "data.hash")
		blockNumber := gjson.GetBytes(body, "ts").Int()
		return []byte(res.String()), string(result), blockNumber
	}
}

func (account *SingleAccount) RetrieveKV(key string) ([]byte, string) {
	kvretrieverequest := &KVRetrieveRequest{
		Data: struct {
			From string `json:"from"`
			Key  string `json:"key"`
			Ts   int64  `json:"ts"`
		}{
			From: utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			Key:  key,
			Ts:   time.Now().UnixMilli(),
		},
		Signature: "",
	}

	bodyBytes, err := utils.SignJsonWithData(kvretrieverequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	resp, error := http.Post(account.Url+"/kv/retrieve", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!"
	} else {
		fmt.Println("kv读取正常")
		body, _ := io.ReadAll(resp.Body)
		result := pretty.Pretty(body)

		rawstr := gjson.GetBytes(body, "data.value").String()

		return []byte(rawstr), string(result)
	}
}
