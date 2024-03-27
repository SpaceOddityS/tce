package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"tce/sdkInit"
	"time"
)

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Product
}

type Product struct {
	ObjectType string `json:"ObjectType"`
	Id         string `json:"Id"`         //溯源码
	Name       string `json:"Name"`       // 奶牛场名字
	ProdutName string `json:"ProdutName"` //产品名称
	Date       string `json:"Date"`       // 生产日期
	Quality    string `json:"Quality"`    // 质量等级

	State string `json:"State"`
}

type User struct {
	Username string
	Password string
}

//嵌套获取历史信息

type HistoryItem struct {
	TxId      string
	Education Product
}

type TraceId struct {
	TxId string
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client

	return handler, nil
}
