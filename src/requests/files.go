package requests

import (
	"SGX_blockchain_client/src/crypto"
	"SGX_blockchain_client/src/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"io"
	"net/http"
	"time"
)

type FileStoreRequest struct {
	Data struct {
		Content  string `json:"content"`
		FileHash string `json:"fileHash"`
		From     string `json:"from"`
		Nonce    int64  `json:"nonce"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type FileRetrieveRequest struct {
	Data struct {
		From     string `json:"from"`
		FileHash string `json:"fileHash"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

func (account *SingleAccount) StoreFile(contentBytes []byte) ([]byte, string, string, int64) {
	filehash := crypto.Keccak256(contentBytes)
	filecontentbase64 := base64.StdEncoding.EncodeToString(contentBytes)

	filestorerequest := &FileStoreRequest{
		Data: struct {
			Content  string `json:"content"`
			FileHash string `json:"fileHash"`
			From     string `json:"from"`
			Nonce    int64  `json:"nonce"`
			Ts       int64  `json:"ts"`
		}{
			Content:  filecontentbase64,
			FileHash: utils.EncodeBytesToHexStringWith0x(filehash),
			From:     utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			Nonce:    1,
			Ts:       time.Now().UnixMilli(),
		},
		Signature: "",
	}
	bodyBytes, err := utils.SignJsonWithData(filestorerequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	resp, error := http.Post(account.Url+"/files/store", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!", "", 0
	} else {
		//fmt.Println("文件存储正常")
		body, _ := io.ReadAll(resp.Body)
		hashstr := gjson.GetBytes(body, "data.hash").String()
		blockNumber := gjson.GetBytes(body, "ts").Int()
		result := pretty.Pretty(body)
		return filehash, string(result), hashstr, blockNumber
	}
}

func (account *SingleAccount) RetrieveFile(fileHash []byte) ([]byte, string) {
	fileretrieverequest := &FileRetrieveRequest{
		Data: struct {
			From     string `json:"from"`
			FileHash string `json:"fileHash"`
			Ts       int64  `json:"ts"`
		}{
			From:     utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			FileHash: utils.EncodeBytesToHexStringWith0x(fileHash),
			Ts:       time.Now().UnixMilli(),
		},
		Signature: "",
	}

	bodyBytes, err := utils.SignJsonWithData(fileretrieverequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	resp, error := http.Post(account.Url+"/files/retrieve", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!"
	} else {
		//fmt.Println("文件读取正常")
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println("body is here", string(body))
		result := pretty.Pretty(body)

		base64str := gjson.GetBytes(body, "data.content").String()
		fcontent, _ := base64.StdEncoding.DecodeString(base64str)
		return fcontent, string(result)
	}
}
