package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	pb "test/proto"
	"google.golang.org/grpc"
)

const (
	sdkConfigPath = "./config/sdk_config.yml"
	contractName  = "Ethcontract"
	invokeMethod  = "transferContent"
)

// 简单交易计数器
var transactionCounter int

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

			// 写入文件（不使用互斥锁）
			f, err := os.OpenFile("data.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				defer f.Close()
				content := fmt.Sprintf("time_send:%s|time_confirm:%s|lis:%s|rsend:%s|time_end:%s\n",
					params["sendts"], params["confirmtx"], params["lis"], params["rsend"], timeHash)
				f.WriteString(content)
			} else {
				log.Printf("写入 data.out 失败: %v", err)
			}

			// 简单增加交易计数并打印
			transactionCounter++
			fmt.Printf("===== 成功接收交易 #%d - TxId: %s =====\n", transactionCounter, txId)

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

func parseTimestamps(line string) (sendts string, confirmtx string, lis string, rsend string) {
	// 先用 '|' 分割所有字段
	parts := strings.Split(line, "|")
	for _, part := range parts {
		if strings.HasPrefix(part, "sendts:") {
			// 拆分 "key:value" 并取 value
			sendts = strings.TrimPrefix(part, "sendts:")
		}
		if strings.HasPrefix(part, "confirmtx:") {
			confirmtx = strings.TrimPrefix(part, "confirmtx:")
		}
		if strings.HasPrefix(part, "lis:") {
			lis = strings.TrimPrefix(part, "lis:")
		}
		if strings.HasPrefix(part, "rsend:") {
			rsend = strings.TrimPrefix(part, "rsend:")
		}
	}
	return
}

// 解析事件数据
func parseEventData(eventData string) (map[string]string, error) {
	fmt.Printf("Received event data: %s\n", eventData)

	// 使用JSON解析
	var jsonData map[string]string
	if err := json.Unmarshal([]byte(eventData), &jsonData); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	params := make(map[string]string)

	// 映射字段
	if from, ok := jsonData["from_address"]; ok {
		params["from"] = from
	}
	if to, ok := jsonData["to_address"]; ok {
		params["to"] = to
	}
	if contentType, ok := jsonData["content_type"]; ok {
		params["content_type"] = contentType
	}
	if content, ok := jsonData["content"]; ok {
		params["content_address"] = content
	}

	fmt.Printf("Parsed parameters: %+v\n", params)
	if len(params) < 4 {
		fmt.Println("Bomb")
		return nil, fmt.Errorf("事件数据缺少必要字段")
	}

	// 解析时间戳
	params["sendts"], params["confirmtx"], params["lis"], params["rsend"] = parseTimestamps(params["content_type"])
	fmt.Printf("Parsed parameters with timestamps: %+v\n", params)
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

func subscribeToEvents(client *sdk.ChainClient, eventChan chan *common.ContractEventInfo, txIdChan chan string) {
	currentHeight, err := client.GetCurrentBlockHeight()
	if err != nil {
		log.Fatalf("获取当前区块高度失败: %v", err)
	}
	startBlock := int64(currentHeight)
	fmt.Printf("当前区块高度: %d\n", startBlock)

	ctx := context.Background()
	eventCh, err := client.SubscribeContractEvent(ctx, startBlock, -1, contractName, "TransferEvent")
	if err != nil {
		log.Fatalf("订阅事件失败: %v", err)
	}

	recentEvents := make(map[string]*common.ContractEventInfo)
	expectedTxId := ""

	for {
		select {
		case txId := <-txIdChan:
			expectedTxId = txId
			fmt.Printf("更新期待的 TxId: %s\n", expectedTxId)
			if cachedEvent, exists := recentEvents[expectedTxId]; exists {
				fmt.Printf("从缓存中找到匹配的事件: TxId=%s\n", expectedTxId)
				eventChan <- cachedEvent
				delete(recentEvents, expectedTxId)
			}
		case event := <-eventCh:
			contractEvent, ok := event.(*common.ContractEventInfo)
			if !ok {
				log.Printf("接收到的事件类型无效: %T", event)
				continue
			}
			fmt.Printf("接收到链上事件: TxId=%s, Topic=%s, BlockHeight=%d\n",
				contractEvent.TxId, contractEvent.Topic, contractEvent.BlockHeight)
			recentEvents[contractEvent.TxId] = contractEvent
			if expectedTxId != "" && contractEvent.TxId == expectedTxId {
				fmt.Printf("找到匹配的事件: TxId=%s\n", contractEvent.TxId)
				eventChan <- contractEvent
				delete(recentEvents, contractEvent.TxId)
			}
		}
	}
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