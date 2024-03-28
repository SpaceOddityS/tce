package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) SavePro(pro Product) (string, error) {

	eventID := "eventAddPro1"
	//注册事件
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)

	//事件defer
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(pro)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}
	//rep 是执行调用链码需要的参数
	req := channel.Request{
		ChaincodeID: t.ChaincodeID,                //通道名字
		Fcn:         "addPro",                     //函数名
		Args:        [][]byte{b, []byte(eventID)}, //函数参数
	}

	//t.Client.Execute(req) 是在后端执行该函数
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	//事件结果
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) FindProInfoByEntityID(entityID [][]byte) (channel.Response, error) {

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "queryProInfoByEntityID",
		Args:        entityID,
	}
	resp, _ := t.Client.Query(req)

	return resp, nil
}

func (t *ServiceSetup) QueryAllPro() (channel.Response, error) {

	eventID := "QueryAllPro"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "QueryAllPro",
		Args:        [][]byte{[]byte(eventID)},
	}
	respone, _ := t.Client.Query(req)

	_ = eventResult(notifier, eventID)

	return respone, nil
}
