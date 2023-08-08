package requests

import (
	"SGX_blockchain_client/src/crypto"
	"SGX_blockchain_client/src/utils"
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type FunctionInput struct {
	InputName       string          `json:"input_name"`
	InputType       string          `json:"input_type"`
	InputComponents []FunctionInput `json:"input_components,omitempty"`
}

type FunctionOutput struct {
	OutputName       string           `json:"output_name"`
	OutputType       string           `json:"output_type"`
	OutputComponents []FunctionOutput `json:"output_components,omitempty"`
}

type ContractFunction struct {
	FunctionName    string           `json:"function_name"`
	FunctionInputs  []FunctionInput  `json:"function_inputs,omitempty"`
	FunctionOutputs []FunctionOutput `json:"function_outputs,omitempty"`
}

type ContractABI struct {
	ContractName      string             `json:"contract_name"`
	ContractFunctions []ContractFunction `json:"contract_functions"`
}

type ContractDeployRequest struct {
	Data struct {
		Name     string `json:"name"`
		ABI      string `json:"abi"`
		From     string `json:"from"`
		Code     string `json:"code"`
		CodeHash string `json:"codeHash"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type ContractDeployResponse struct {
	Status string `json:"status"`
	Data   struct {
		From            string `json:"from"`
		Hash            string `json:"hash"`
		ContractAddress string `json:"contractAddress"`
		Nonce           int64  `json:"nonce"`
		CodeHash        string `json:"codeHash"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type ContractCallRequest struct {
	Data struct {
		CodeHash        string `json:"codeHash"`
		ContractAddress string `json:"contractAddress"`
		From            string `json:"from"`
		FunctionName    string `json:"functionName"`
		FuntionInputs   string `json:"functionInputs"`
		Ts              int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type ContractCallResponse struct {
	Status string `json:"status"`
	Data   struct {
		CodeHash string      `json:"codeHash"`
		From     string      `json:"from"`
		Hash     string      `json:"hash"`
		Nonce    int         `json:"nonce"`
		Result   interface{} `json:"result"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

func (account *SingleAccount) DeployContract(name, code, abi string) ([]byte, string, int64, string) {

	contractdeployrequest := &ContractDeployRequest{
		Data: struct {
			Name     string `json:"name"`
			ABI      string `json:"abi"`
			From     string `json:"from"`
			Code     string `json:"code"`
			CodeHash string `json:"codeHash"`
			Ts       int64  `json:"ts"`
		}{
			Name:     name,
			ABI:      abi,
			From:     utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			Code:     code,
			CodeHash: utils.EncodeBytesToHexStringWith0x(crypto.Keccak256([]byte(code))),
			Ts:       time.Now().UnixMilli(),
		},
		Signature: "",
	}
	fmt.Println("contract from", utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK))
	bodyBytes, err := utils.SignJsonWithData(contractdeployrequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	//fmt.Println("合约depoly")
	fmt.Println(string(bodyBytes))
	resp, error := http.Post(account.Url+"/contract/deploy", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!", 0, ""
	} else {
		//fmt.Println("合约部署正常")
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(body)
		result := pretty.Pretty(body)
		res := gjson.GetBytes(body, "data.hash")
		contractaddress := gjson.GetBytes(body, "data.contractAddress").String()
		blockNumber := gjson.GetBytes(body, "ts").Int()
		return []byte(res.String()), string(result), blockNumber, contractaddress
	}
}

func (account *SingleAccount) CallContract(codeHash, contractAddress, functionName, functionInputs string) ([]byte, string, int64, string) {

	contractcallrequest := &ContractCallRequest{
		Data: struct {
			CodeHash        string `json:"codeHash"`
			ContractAddress string `json:"contractAddress"`
			From            string `json:"from"`
			FunctionName    string `json:"functionName"`
			FuntionInputs   string `json:"functionInputs"`
			Ts              int64  `json:"ts"`
		}{
			CodeHash:        codeHash,
			ContractAddress: contractAddress,
			From:            utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
			FunctionName:    functionName,
			FuntionInputs:   functionInputs,
			Ts:              time.Now().UnixMilli(),
		},
		Signature: "",
	}
	bodyBytes, err := utils.SignJsonWithData(contractcallrequest, account.Keypair)
	if err != nil {
		fmt.Println("error in signing json!")
	}
	buffer := bytes.NewBuffer(bodyBytes)
	fmt.Println("合约call")
	fmt.Println(string(bodyBytes))
	resp, error := http.Post(account.Url+"/contract/call", jsonContentType, buffer)
	defer resp.Body.Close()
	if error != nil {
		fmt.Println(error)
		return []byte(""), "wrong!", 0, ""
	} else {
		//fmt.Println("合约调用正常")
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(body)
		result := pretty.Pretty(body)
		res := gjson.GetBytes(body, "data.hash")
		resultstr := gjson.GetBytes(body, "data.result").String()

		blockNumber := gjson.GetBytes(body, "ts").Int()
		return []byte(res.String()), string(result), blockNumber, resultstr
	}
}

func (account *SingleAccount) Storeexecutecontract(codeHash, contractAddress, functionName, functionInputs string) {
	fmt.Println("生成交易请求:")
	//var Hashrandom []byte
	randomhash := make([]byte, 32)
	t1 := time.Now()
	for i := 0; i < 50000; i++ {

		rand.Read(randomhash)
		contractcallrequest := &ContractCallRequest{
			Data: struct {
				CodeHash        string `json:"codeHash"`
				ContractAddress string `json:"contractAddress"`
				From            string `json:"from"`
				FunctionName    string `json:"functionName"`

				FuntionInputs string `json:"functionInputs"`
				Ts            int64  `json:"ts"`
			}{
				ContractAddress: utils.EncodeBytesToHexStringWith0x(randomhash),
				From:            utils.EncodeBytesToHexStringWith0x(account.Keypair.PubK),
				FunctionName:    functionName,
				Ts:              time.Now().UnixMilli(),
			},
			Signature: "",
		}
		fmt.Println("智能合约执行请求序号：", i+1, contractcallrequest)
	}
	t2 := time.Now()
	du := t2.Sub(t1)
	fmt.Println("耗时：", du)
	fmt.Println("性能：", 50000/du.Seconds())
}
