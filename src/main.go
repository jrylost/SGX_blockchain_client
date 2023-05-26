package main

import (
	"SGX_blockchain_client/src/requests"
	"SGX_blockchain_client/src/utils"
	"fmt"
)

func main() {
	//account := requests.CreateNewSingleAccount("http://47.102.115.171:8888")
	account := requests.CreateNewSingleAccount("http://localhost:8888")
	//IP地址和端口

	accountInfo := account.GetAccountsInfo()
	fmt.Println(accountInfo)
	//获取账户信息

	fileHash, resp, filetxhash, fileBlockNumber := account.StoreFile([]byte("File content is this!"))
	fmt.Println(utils.EncodeBytesToHexStringWith0x(fileHash))
	fmt.Println(resp)
	//存储文件

	filetxinfo := account.GetTransactionInfo(filetxhash)
	fmt.Println(filetxinfo)

	content, respbody := account.RetrieveFile(fileHash)
	fmt.Println(string(content), respbody)
	//获取文件

	key := "this is key"
	value := "this is value"
	kvtxhash, resp2, KVBlockNumber := account.StoreKV(key, value)
	fmt.Println(resp2)
	val, respbody2 := account.RetrieveKV(key)
	fmt.Println(string(val), respbody2)

	kvtxinfo := account.GetTransactionInfo(string(kvtxhash))
	fmt.Println(kvtxinfo)

	accountInfo = account.GetAccountsInfo()
	fmt.Println(accountInfo)

	fileblock := account.GetBlockInfo(fileBlockNumber / 1000)
	kvblock := account.GetBlockInfo(KVBlockNumber / 1000)
	fmt.Println(fileblock, kvblock)

}
