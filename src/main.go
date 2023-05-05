package main

import (
	"SGX_blockchain_client/src/requests"
	"SGX_blockchain_client/src/utils"
	"fmt"
)

func main() {
	account := requests.CreateNewSingleAccount("http://47.102.115.171:8888")
	//IP地址和端口

	accountInfo := account.GetAccountsInfo()
	fmt.Println(accountInfo)
	//获取账户信息

	fileHash, resp := account.StoreFile([]byte("File content is this!"))
	fmt.Println(utils.EncodeBytesToHexStringWith0x(fileHash))
	fmt.Println(resp)
	//存储文件

	content, respbody := account.RetrieveFile(fileHash)
	fmt.Println(string(content), respbody)
	//获取文件

	accountInfo = account.GetAccountsInfo()
	fmt.Println(accountInfo)

}
