package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/myutils"
	"log"
	"net/http"
	"os"
	"tce/middlewares"
	"tce/sdkInit"
	"tce/service"
	"time"
)

//设置常数    链码名称  版本
const (
	cc_name    = "TCE_cc"
	cc_version = "1.0.0"
)

func main() {
	r := gin.Default()
	r.Use(middlewares.Cors())
	// 初始化组织信息，可以写多个组织信息
	orgs := []*sdkInit.OrgInfo{ //SDK 中封装好的组织结构体
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/tce/fixtures/channel-artifacts/Org1MSPanchors.tx",
		}, //  os.Getenv  作用是返回参数内的环境变量
	}
	// sdk 环境信息
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/tce/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH") + "/src/tce/chaincode/",
		ChaincodeVersion: cc_version,
	}
	// sdk setup   调用sdkinit 内的 Setup 函数，将config.yaml  和上面 建立好的sdk环境信息传入 ，返回一个完整的SDK
	sdk, err := sdkInit.Setup("config.yaml", &info)
	fmt.Println("-------------------")
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}
	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}
	fmt.Println("-------------------")

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}
	// invoke chaincode set statuspp
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	//----------------------------------------------------------------------------------
	pro := service.Product{
		Id:         "DF-100",
		ProdutName: "大白菜",
		Name:       "南昌百菜园",
		Date:       time.Now().Format("2006-01-02 15:04:05"),
		Quality:    "优",
		State:      "1",
	}
	//初始化服务
	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	msg, err := serviceSetup.SavePro(pro)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	pro.Id = pro.Id + "_" + "1"

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte("DF-100_1"))
	result, err := serviceSetup.FindProInfoByEntityID(bodyBytes)
	//var data []map[string]interface{}
	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(bytes.NewBuffer(result.Payload).Bytes(), &data)
	if err != nil {
		//fmt.Errorf("%s", err)
		fmt.Println(err)
		//fmt.Println("错误 结束")
		return
	} else {
		fmt.Println("根据溯源码查询信息成功：")
		fmt.Println(data)
	}
	fmt.Println("打印区块链高度为1的区块信息：")
	//-----------------------------------------------------------
	//app := controller.Application{
	//	Setup: serviceSetup,
	//}
	//
	//web.WebStart(app)
	myutils.Test(2)

	r = CollectRoute(r, serviceSetup)
	r.Run(":8000")
}
func CollectRoute(r *gin.Engine, setup *service.ServiceSetup) *gin.Engine {

	//路由
	r.POST("/Add1", func(ctx *gin.Context) {

		var requestUser = service.Product{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		} else {
			ctx.String(http.StatusOK, "绑定form成功")
		}
		requestUser.State = "1"
		fmt.Println("产品上链开始......")
		fmt.Println(requestUser)
		log.Println(requestUser.Id, requestUser.Name, requestUser.ProdutName)
		//返回结果
		Txid, err := setup.SavePro(requestUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": requestUser,
			"msg":  "上链成功！",
		})
		fmt.Println(Txid)
		ctx.JSON(200, gin.H{
			"交易ID": Txid,
		})
	})
	r.POST("/Add2", func(ctx *gin.Context) {

		var requestUser = service.Product{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		} else {
			ctx.String(http.StatusOK, "绑定form成功")
		}
		requestUser.State = "2"
		fmt.Println("产品上链开始......")
		fmt.Println(requestUser)
		log.Println(requestUser.Id, requestUser.Name, requestUser.ProdutName)
		//返回结果
		Txid, err := setup.SavePro(requestUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": requestUser,
			"msg":  "上链成功！",
		})
		fmt.Println(Txid)
		ctx.JSON(200, gin.H{
			"交易ID": Txid,
		})
	})
	r.POST("/Add3", func(ctx *gin.Context) {

		var requestUser = service.Product{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		} else {
			ctx.String(http.StatusOK, "绑定form成功")
		}
		requestUser.State = "3"
		fmt.Println("产品上链开始......")
		fmt.Println(requestUser)
		log.Println(requestUser.Id, requestUser.Name, requestUser.ProdutName)
		//返回结果
		Txid, err := setup.SavePro(requestUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": requestUser,
			"msg":  "上链成功!",
		})
		fmt.Println(Txid)
		ctx.JSON(200, gin.H{
			"交易ID": Txid,
		})
	})
	r.POST("/QyById1", func(ctx *gin.Context) {

		var requestUser = service.TraceId{}
		err := ctx.ShouldBind(&requestUser)
		requestUser.TxId = requestUser.TxId + "_" + "1"
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(requestUser.TxId))
		resp, err := setup.FindProInfoByEntityID(bodyBytes)

		var data interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "",
				"错误码": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"状态":   "已生产",
				"DATE": data,
			})
		}

	})
	r.POST("/QyById2", func(ctx *gin.Context) {

		var requestUser = service.TraceId{}
		err := ctx.ShouldBind(&requestUser)
		requestUser.TxId = requestUser.TxId + "_" + "2"
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(requestUser.TxId))
		resp, err := setup.FindProInfoByEntityID(bodyBytes)

		var data interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "",
				"错误码": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"状态":   "已在中间商",
				"DATE": data,
			})
		}

	})
	r.POST("/QyById3", func(ctx *gin.Context) {

		var requestUser = service.TraceId{}
		err := ctx.ShouldBind(&requestUser)

		var id string
		var ic string
		id = requestUser.TxId
		ic = requestUser.TxId
		requestUser.TxId = requestUser.TxId + "_" + "1"
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(requestUser.TxId))
		resp, err := setup.FindProInfoByEntityID(bodyBytes)

		var data interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "未查询到产品",
				"错误码": err.Error(),
			})
			return
		} else {
			//ctx.JSON(http.StatusOK, gin.H{
			//	"状态":   "已在终端",
			//	"DATE": data,
			//})
		}

		id = id + "_" + "2"
		var bodyBytes1 [][]byte
		bodyBytes1 = append(bodyBytes1, []byte(id))
		resp1, _ := setup.FindProInfoByEntityID(bodyBytes1)

		var data1 interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp1.Payload).Bytes(), &data1); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "未查询到产品",
				"错误码": err.Error(),
			})
			return
		} else {
			//ctx.JSON(http.StatusOK, gin.H{
			//	"状态":   "已在终端",
			//	"DATE": data1,
			//})
		}

		ic = ic + "_" + "3"
		var bodyBytes2 [][]byte
		bodyBytes2 = append(bodyBytes2, []byte(ic))
		resp2, _ := setup.FindProInfoByEntityID(bodyBytes2)

		var data2 interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp2.Payload).Bytes(), &data2); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "未查询到产品",
				"错误码": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"状态":    "已在终端",
				"DATE":  data,
				"DATE1": data1,
				"DATE2": data2,
			})
		}

	})
	r.POST("/login", func(ctx *gin.Context) {

		// 输入溯源码然后
		var requestUser = service.User{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		} else {
			ctx.String(http.StatusOK, "绑定form成功")
		}

		fmt.Println("产品上链开始......")
		fmt.Println(requestUser)
		//返回结果
		if requestUser.Username == "org1" && requestUser.Password == "123456" {

			ctx.JSON(http.StatusOK, gin.H{
				"type": "org1",
				"code": 200,
				"data": requestUser,
				"msg":  "登录成功！欢迎组织一用户",
			})
			return
		}
		if requestUser.Username == "org2" && requestUser.Password == "123456" {

			ctx.JSON(http.StatusOK, gin.H{
				"type": "org2",
				"code": 200,
				"msg":  "登录成功！ 欢迎组织二用户",
			})
			return
		}
		if requestUser.Username == "org3" && requestUser.Password == "123456" {

			ctx.JSON(http.StatusOK, gin.H{
				"type": "org3",
				"code": 200,
				"msg":  "登录成功！欢迎组织三用户",
			})
			return
		}

	})
	r.GET("/Getall", func(ctx *gin.Context) {

		resp, _ := setup.QueryAllPro()

		var data []map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":   "",
				"错误码":  err.Error(),
				"date": data,
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"DATE": data,
			})
		}

	})

	//r.GET("/Getall", func(ctx *gin.Context) {
	//
	//	resp, _ := setup.QueryAllPro()
	//
	//	var data []map[string]interface{}
	//	if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
	//		ctx.JSON(http.StatusExpectationFailed, gin.H{
	//			"失败":   "",
	//			"错误码":  err.Error(),
	//			"date": data,
	//		})
	//		return
	//	} else {
	//		ctx.JSON(http.StatusOK, gin.H{
	//			"DATE": data,
	//		})
	//	}
	//
	//})

	return r
}
