package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	pb "test/proto" // 替换为你的 proto 生成包的实际路径
	"google.golang.org/grpc"
)

const (
	sdkConfigPath = "./config/sdk_config.yml"
	contractName  = "Ethcontract"
	invokeMethod  = "transferContent"
)


type server struct {
    pb.UnimplementedEndpointChainEventListenerServer
    client *sdk.ChainClient
}


func (s *server) ListenForEvents(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
    eventData := req.GetEventData()
    timeRec := time.Now().Format("2006-01-02-15:04:05:00")

    params, err := parseEventData(eventData)
    if err != nil {
        return &pb.EventResponse{
            Success: false,
            Message: fmt.Sprintf("解析事件数据失败: %v", err),
        }, nil
    }

    addrEth := params["content_address"]
    txId, err := s.invokeContract(params)
    if err != nil {
        timeSend := time.Now().Format("2006-01-02-15:04:05:00")
        responseMsg := fmt.Sprintf("Addr_eth:%s|Hash:%s|Time_rec:%s|Time_hash:%s|Time_send:%s|Stat:0",
            addrEth, "", timeRec, "", timeSend)
        return &pb.EventResponse{
            Success: false,
            Message: responseMsg,
        }, nil
    }

    // 为当前请求创建专用通道
    txIdChan := make(chan string, 1)
    eventChan := make(chan *common.ContractEventInfo, 1)

    // 启动订阅协程
    go subscribeToEvents(s.client, eventChan, txIdChan)

    // 发送 txId 到订阅
    txIdChan <- txId
    fmt.Printf("已发送 TxId 到订阅: %s\n", txId)

    // 等待链上事件
    fmt.Println("Now waiting for event...")
    var timeHash string
    select {
    case event := <-eventChan:
        fmt.Println("End Waiting for event...")
        fmt.Printf("event.TxId:%s\n", event.TxId)
        fmt.Printf("txID:%s\n", txId)
        timeHash = time.Now().Format("2006-01-02-15:04:05:00")
        if event.TxId == txId {
            fmt.Println("Accepted")
            timeSend := time.Now().Format("2006-01-02-15:04:05:00")
            responseMsg := fmt.Sprintf("Addr_eth:%s|Hash:%s|Time_rec:%s|Time_hash:%s|Time_send:%s|Stat:1",
                addrEth, event.TxId, timeRec, timeHash, timeSend)
            fmt.Printf("responseMsg:%s\n", responseMsg)
            return &pb.EventResponse{
                Success: true,
                Message: responseMsg,
            }, nil
        } else {
            fmt.Println("Wrong Answer")
        }
    case <-time.After(100 * time.Second):
        fmt.Println("End Waiting for event...")
        timeSend := time.Now().Format("2006-01-02-15:04:05:00")
        responseMsg := fmt.Sprintf("Addr_eth:%s|Hash:%s|Time_rec:%s|Time_hash:%s|Time_send:%s|Stat:0",
            addrEth, txId, timeRec, "", timeSend)
        return &pb.EventResponse{
            Success: false,
            Message: responseMsg,
        }, nil
    }

    timeSend := time.Now().Format("2006-01-02-15:04:05:00")
    responseMsg := fmt.Sprintf("Addr_eth:%s|Hash:%s|Time_rec:%s|Time_hash:%s|Time_send:%s|Stat:1",
        addrEth, txId, timeRec, timeHash, timeSend)
    return &pb.EventResponse{
        Success: false,
        Message: responseMsg,
    }, nil
}

// 解析事件数据
func parseEventData(eventData string) (map[string]string, error) {
	fmt.Printf("Received event data: %s\n", eventData)
	params := make(map[string]string)
	parts := strings.Split(eventData, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.Trim(kv[0], "\"{")
			value := strings.Trim(kv[1], "\"}")
			switch key {
			case "from_address":
				params["from"] = value
			case "to_address":
				params["to"] = value
			case "content_type":
				params["content_type"] = value
			case "content":
				params["content_address"] = value
			}
		}
	}
	if len(params) < 4 {
		return nil, fmt.Errorf("事件数据缺少必要字段")
	}
	return params, nil
}

// 检查合约是否存在
func checkContract(client *sdk.ChainClient) error {
	contract, err := client.GetContractInfo(contractName)
	//fmt.Printf(contract.Name)
	if err != nil {
		return fmt.Errorf("获取合约信息失败: %v", err)
	}
	if contract == nil || contract.Name == "" {
		return fmt.Errorf("合约 %s 不存在", contractName)
	}
	fmt.Printf("合约已存在: %+v\n", contract)
	return nil
}

// 调用合约并返回交易ID
func (s *server) invokeContract(params map[string]string) (string, error) {
	fmt.Println("====================== 调用 Ethcontract 合约 ======================")

	kvs := []*common.KeyValuePair{
		{Key: "from", Value: []byte(params["from"])},
		{Key: "to", Value: []byte(params["to"])},
		{Key: "content_type", Value: []byte(params["content_type"])},
		{Key: "content_address", Value: []byte(params["content_address"])},
	}

	resp, err := s.client.InvokeContract(contractName, invokeMethod, "", kvs, -1, true)
	if err != nil {
		return "", fmt.Errorf("调用合约失败: %v", err)
	}
	if resp.Code != common.TxStatusCode_SUCCESS {
		return "", fmt.Errorf("合约执行失败: [code:%d]/[msg:%s]", resp.Code, resp.Message)
	}

	fmt.Printf("交易成功 - TxId: %s, BlockHeight: %d, Message: %s\n",
		resp.TxId, resp.TxBlockHeight, resp.ContractResult.Message)
	return resp.TxId, nil
}

func main() {
    client, err := sdk.NewChainClient(
        sdk.WithConfPath(sdkConfigPath),
    )
    if err != nil {
        log.Fatalf("创建 Chainmaker 客户端失败: %v", err)
    }

    if err := checkContract(client); err != nil {
        log.Fatalf("合约检查失败: %v", err)
    }

    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("监听端口失败: %v", err)
    }

    grpcServer := grpc.NewServer()
    srv := &server{
        client: client,
    }
    pb.RegisterEndpointChainEventListenerServer(grpcServer, srv)
    fmt.Println("gRPC 服务运行在 0.0.0.0:50051")

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("gRPC 服务启动失败: %v", err)
    }
}