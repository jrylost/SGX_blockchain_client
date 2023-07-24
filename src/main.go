package main

import (
	"SGX_blockchain_client/src/crypto"
	"SGX_blockchain_client/src/requests"
	"SGX_blockchain_client/src/utils"
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

var functionInput = `[
    {
        "input_name" : "context",
        "input_type" : "context", 
        "input_value" : ""
    },
    {
        "input_name" : "data",
        "input_type" : "string",
        "input_value" : "{\"evidenceId\":\"aa\",\"uploaderSign\":\"cc\",\"content\":\"dd\"}"
    }
]`

var functionInput2 = `[
    {
        "input_name" : "context",
        "input_type" : "context", 
        "input_value" : ""
    },
    {
        "input_name" : "evidenceId",
        "input_type" : "string",
        "input_value" : "Evi_aa"
    }
]`

func main() {
	pk := "0x7b079a78348a4ed27b69f5c02ab4a944f61f117405d008a566a85ffbb4fce1d5"
	//pk := "0xd9b29c4dc2b3202d8b43ae6677ac02e754a1217f94a193c53f21e89d26d06685"
	account := requests.CreateNewSingleAccount("http://47.102.115.171:8888", pk)
	//account := requests.CreateNewSingleAccount("http://localhost:8888", pk)
	fmt.Println(utils.EncodeBytesToHexStringWith0x(account.Keypair.PublicKey.SerializeUncompressed()))
	fmt.Println(utils.EncodeBytesToHexStringWith0x(account.Keypair.PublicKey.SerializeCompressed()))
	//IP地址和端口
	fmt.Println(hex.EncodeToString(account.Keypair.PrivateKey.Serialize()))
	accountInfo := account.GetAccountsInfo()
	fmt.Println(accountInfo)
	//获取账户信息

	fileHash, resp, filetxhash, fileBlockNumber := account.StoreFile([]byte("File content is this!"))
	fmt.Println(utils.EncodeBytesToHexStringWith0x(fileHash))
	fmt.Println(resp)
	//存储文件

	//fmt.Println(hex.EncodeToString(account.Keypair.PrivateKey.Serialize()))
	accountInfo = account.GetAccountsInfo()
	fmt.Println(accountInfo)
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

	fileb, err := os.ReadFile(`./src/example_contracts/evidence_v6/evidence_v6.go2`)
	if err != nil {
		fmt.Println(err)
	}
	abi, err := os.ReadFile(`./src/example_contracts/evidence_v6/test_abi.json`)
	if err != nil {
		fmt.Println(err)
	}
	contracttxhash, contractrespstring, contractBlockNumber, contractaddress := account.DeployContract("evidence_v6", string(fileb), string(abi))
	fmt.Println(string(contracttxhash), contractrespstring)

	contracttxinfo := account.GetTransactionInfo(string(contracttxhash))
	fmt.Println(contracttxinfo)

	contractcodehash := utils.EncodeBytesToHexStringWith0x(crypto.Keccak256(fileb))

	contracttxhash, contractrespstring, contractBlockNumber, resultstr := account.CallContract(contractcodehash, contractaddress, "AddEvidence", functionInput)
	fmt.Println(resultstr)
	fmt.Println(contractrespstring)

	contracttxhash, contractrespstring, contractBlockNumber, resultstr = account.CallContract(contractcodehash, contractaddress, "QueryEvidenceById", functionInput2)
	fmt.Println(resultstr)
	fmt.Println(contractrespstring)

	fileblock := account.GetBlockInfo(fileBlockNumber / 1000)
	kvblock := account.GetBlockInfo(KVBlockNumber / 1000)
	contractblock := account.GetBlockInfo(contractBlockNumber / 1000)
	fmt.Println(fileblock, kvblock, contractblock)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadRune()
}
