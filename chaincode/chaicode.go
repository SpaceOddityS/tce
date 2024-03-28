package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Product struct {
	ObjectType string `json:"ObjectType"`
	Id         string `json:"Id"`   //溯源码
	Name       string `json:"Name"` // 奶牛场名字
	ProdutName string `json:"ProdutName"`
	Date       string `json:"Date"`    // 生产日期
	Quality    string `json:"Quality"` // 质量等级
	State      string `json:"State"`
}
type Chaincode struct {
}
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Product
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()
	if fun == "addPro" {
		return t.addPro(stub, args) // 添加信息
	} else if fun == "queryProInfoByEntityID" {
		return t.queryProInfoByEntityID(stub, args) // 根据身份证号码及姓名查询详情
	} else if fun == "QueryAllPro" {
		return t.QueryAllPro(stub, args)
	}
	return shim.Error("指定的函数名称错误")

}

const DOC_TYPE = "proObj"

func PutPro(stub shim.ChaincodeStubInterface, edu Product, state string) ([]byte, bool) {

	edu.ObjectType = DOC_TYPE
	b, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}

	id := edu.Id + "_" + state
	// 保存edu状态
	err = stub.PutState(id, b)
	if err != nil {
		return nil, false
	}

	return b, true
}
func (t *Chaincode) addPro(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var pro Product
	err := json.Unmarshal([]byte(args[0]), &pro)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	//// 查重: 身份证号码必须唯一
	//_, exist := GetEduInfo(stub, pro.Id)
	//if exist {
	//	return shim.Error("要添加的身份证号码已存在")
	//}

	_, bl := PutPro(stub, pro, pro.State)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息添加成功"))
}
func (t *Chaincode) queryProInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	b, _ := stub.GetState(args[0])
	return shim.Success(b)
}
func (t *Chaincode) QueryAllPro(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	startKey := ""
	endKey := ""

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error("错了")
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		//CoughDB  是一个键值对数据库，支持富查询   值是JSON形式
		if err != nil {
			return shim.Error("错了")
		}

		product := new(Product)
		_ = json.Unmarshal(queryResponse.Value, product)

		queryResult := QueryResult{Key: queryResponse.Key, Record: product}
		results = append(results, queryResult)
	}
	jsontex, _ := json.Marshal(results)
	return shim.Success(jsontex)
}
func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
