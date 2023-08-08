package requests

import (
	"SGX_blockchain_client/src/utils"
	"bytes"
	"fmt"
	"github.com/tidwall/pretty"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type BlockInfoRequest struct {
	Data struct {
		Number int64 `json:"number"`
		Ts     int64 `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type TransactionInfoRequest struct {
	Data struct {
		Hash string `json:"hash"`
		Ts   int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

func (account *SingleAccount) GetBlockInfo(number int64) string {

	blockinfoRequest := &BlockInfoRequest{
		Data: struct {
			Number int64 `json:"number"`
			Ts     int64 `json:"ts"`
		}{
			Number: number,
			Ts:     time.Now().UnixMilli(),
		},
		Signature: "",
	}

	bodyBytes, err := utils.SignJsonWithData(blockinfoRequest, account.Keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(account.Url+"/block/info", jsonContentType, buffer)

	if err != nil {
		fmt.Println(err)
		return "Wrong!"
	} else {
		//fmt.Println("区块信息读取正常")
		body, _ := io.ReadAll(resp.Body)
		result := pretty.Pretty(body)

		return string(result)
	}
	defer resp.Body.Close()
	return ""
}

func (account *SingleAccount) GetTransactionInfo(hash string) string {

	transactioninfoRequest := &TransactionInfoRequest{
		Data: struct {
			Hash string `json:"hash"`
			Ts   int64  `json:"ts"`
		}{
			Hash: hash,
			Ts:   time.Now().UnixMilli(),
		},
		Signature: "",
	}

	bodyBytes, err := utils.SignJsonWithData(transactioninfoRequest, account.Keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer := bytes.NewBuffer(bodyBytes)
	//fmt.Println(string(bodyBytes))
	resp, err := http.Post(account.Url+"/transaction/info", jsonContentType, buffer)

	if err != nil {
		fmt.Println(err)
		return "Wrong!"
	} else {
		//fmt.Println("交易信息读取正常")
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(string(body))
		result := pretty.Pretty(body)

		return string(result)
	}
	defer resp.Body.Close()
	return ""
}

func (account *SingleAccount) Storetxinfo() {
	fmt.Println("生成交易请求:")
	//var Hashrandom []byte
	randomhash := make([]byte, 32)
	t1 := time.Now()
	ttt := time.Now().UnixMilli()
	for i := 0; i < 100000; i++ {

		rand.Read(randomhash)
		transactioninfoRequest := &TransactionInfoRequest{
			Data: struct {
				Hash string `json:"hash"`
				Ts   int64  `json:"ts"`
			}{
				Hash: utils.EncodeBytesToHexStringWith0x(randomhash),
				Ts:   ttt,
			},
			Signature: "",
		}
		fmt.Println("交易查询请求序号：", i+1, transactioninfoRequest)
	}
	t2 := time.Now()
	fmt.Println("耗时：", t2.Sub(t1))
	fmt.Println("性能：", 100000/t2.Sub(t1).Seconds())
}

func (account *SingleAccount) Storecontextinfo() {
	fmt.Println("生成交易请求:")
	//var Hashrandom []byte
	randomhash := make([]byte, 32)
	t1 := time.Now()
	ttt := time.Now().UnixMilli()

	for i := 0; i < 100000; i++ {

		rand.Read(randomhash)
		transactioninfoRequest := &TransactionInfoRequest{
			Data: struct {
				Hash string `json:"hash"`
				Ts   int64  `json:"ts"`
			}{
				Hash: utils.EncodeBytesToHexStringWith0x(randomhash),
				Ts:   ttt,
			},
			Signature: "",
		}
		fmt.Println("context存储请求序号：", i+1, transactioninfoRequest)
	}
	t2 := time.Now()
	fmt.Println("耗时：", t2.Sub(t1))
	fmt.Println("性能：", 100000/t2.Sub(t1).Seconds())
}
