package main

import (
	"SGX_blockchain_client/src/crypto"
	"SGX_blockchain_client/src/requests"
	"SGX_blockchain_client/src/utils"
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
	fmt.Println("测试可信执行环境：")
	fmt.Println("测试生成远程认证报告、获取远程认证报告、验证远程认证报告、建立TLS传输通道、加入分布式网络、同步节点状态：")
	account := requests.CreateNewSingleAccount("http://47.102.115.171:8888", pk)
	//account := requests.CreateNewSingleAccount("http://localhost:8888", pk)
	fmt.Println("账户公钥：", utils.EncodeBytesToHexStringWith0x(account.Keypair.PublicKey.SerializeUncompressed()))
	fmt.Println("账户压缩公钥：", utils.EncodeBytesToHexStringWith0x(account.Keypair.PublicKey.SerializeCompressed()))
	//IP地址和端口
	fmt.Println("账户私钥：", hex.EncodeToString(account.Keypair.PrivateKey.Serialize()))
	accountInfo := account.GetAccountsInfo()
	fmt.Println("下面测试账户信息：")
	fmt.Println("测试获取账户交易数量、文件数量、智能合约数量：")
	fmt.Println(accountInfo)
	//获取账户信息

	fmt.Println("下面测试存储功能：")
	fmt.Println("测试存储文件内容、生成分布式锁、锁定存储信息、读取存储信息、写入存储信息、解锁存储信息、删除分布式锁：")
	fileHash, resp, filetxhash, fileBlockNumber := account.StoreFile([]byte("File content is this!"))
	fmt.Println("测试文件哈希为：", utils.EncodeBytesToHexStringWith0x(fileHash))
	fmt.Println(resp)
	//存储文件
	fmt.Println("测试查询区块大小、查询区块交易数量、查询区块生成时间：")
	accountblock := account.GetBlockInfo(fileBlockNumber / 1000)
	//fmt.Println(hex.EncodeToString(account.Keypair.PrivateKey.Serialize()))
	//fmt.Println("下面账户交易信息查询：")
	//accountInfo = account.GetAccountsInfo()
	//fmt.Println(accountInfo)
	fmt.Println("下面测试交易信息查询：")
	fmt.Println("测试交易事务凭据、交易事务类型、交易事务对手方：")
	filetxinfo := account.GetTransactionInfo(filetxhash)
	fmt.Println(filetxinfo)

	fmt.Println("下面测试存储功能读取查询：")
	fmt.Println("测试查询文件内容、验证文件所有权：")
	content, respbody := account.RetrieveFile(fileHash)
	fmt.Println(string(content), respbody)
	//获取文件

	key := "this is key"
	value := "this is value"
	fmt.Println("下面测试键值对值存取：")
	fmt.Println("测试存储键值对信息：")
	//kvtxhash, resp2, KVBlockNumber := account.StoreKV(key, value)
	_, resp2, _ := account.StoreKV(key, value)
	//_, resp2, KVBlockNumber := account.StoreKV(key, value)
	fmt.Println(resp2)
	fmt.Println("测试查询键值对信息、验证键值对所有权：")
	val, respbody2 := account.RetrieveKV(key)
	fmt.Println(string(val), respbody2)

	//kvtxinfo := account.GetTransactionInfo(string(kvtxhash))
	//fmt.Println(kvtxinfo)
	//
	//accountInfo = account.GetAccountsInfo()
	//fmt.Println(accountInfo)

	fmt.Println("下面测试虚拟机和智能合约功能：")
	fileb, err := os.ReadFile(`./src/example_contracts/evidence_v6/evidence_v6.go2`)
	if err != nil {
		fmt.Println(err)
	}
	abi, err := os.ReadFile(`./src/example_contracts/evidence_v6/test_abi.json`)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("智能合约编译部署：")
	fmt.Println("测试部署智能合约、分析智能合约合法性、生成智能合约AST、存取智能合约context：")
	//contracttxhash, contractrespstring, contractBlockNumber, contractaddress := account.DeployContract("evidence_v6", string(fileb), string(abi))
	contracttxhash, contractrespstring, _, contractaddress := account.DeployContract("evidence_v6", string(fileb), string(abi))
	fmt.Println(string(contracttxhash), contractrespstring)

	contracttxinfo := account.GetTransactionInfo(string(contracttxhash))
	fmt.Println(contracttxinfo)

	contractcodehash := utils.EncodeBytesToHexStringWith0x(crypto.Keccak256(fileb))
	fmt.Println("智能合约执行：")

	fmt.Println("测试执行智能合约、返回智能合约执行结果、返回智能合约执行错误信息：")
	//contracttxhash, contractrespstring, contractBlockNumber, resultstr := account.CallContract(contractcodehash, contractaddress, "AddEvidence", functionInput)
	contracttxhash, contractrespstring, _, resultstr := account.CallContract(contractcodehash, contractaddress, "AddEvidence", functionInput)
	fmt.Println(resultstr)
	fmt.Println(contractrespstring)

	//contracttxhash, contractrespstring, contractBlockNumber, resultstr = account.CallContract(contractcodehash, contractaddress, "QueryEvidenceById", functionInput2)
	contracttxhash, contractrespstring, _, resultstr = account.CallContract(contractcodehash, contractaddress, "QueryEvidenceById", functionInput2)
	fmt.Println(resultstr)
	fmt.Println(contractrespstring)

	fmt.Println("获取账户内交易、文件、智能合约列表：")
	accountblock = account.GetBlockInfo(fileBlockNumber / 1000)
	//kvblock := account.GetBlockInfo(KVBlockNumber / 1000)
	//contractblock := account.GetBlockInfo(contractBlockNumber / 1000)
	fmt.Println(accountblock)
	//fmt.Println("下面测试存储交易事务信息性能：")
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadRune()
	//account.Storetxinfo()
	//fmt.Println("下面测试存取智能合约context性能：")
	//reader = bufio.NewReader(os.Stdin)
	//reader.ReadRune()
	//account.Storecontextinfo()
	//fmt.Println("下面测试执行智能合约执行性能：")
	//reader = bufio.NewReader(os.Stdin)
	//reader.ReadRune()
	//account.Storeexecutecontract()
	//reader = bufio.NewReader(os.Stdin)
	//reader.ReadRune()

	fmt.Println("测试退出分布式网络：")
	//reader = bufio.NewReader(os.Stdin)
	//reader.ReadRune()
}
