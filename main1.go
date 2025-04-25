// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"

// 	pb "test/proto" // 导入生成的 proto 包，根据实际路径调整
// 	"google.golang.org/grpc"
// )

// // 服务器结构体
// type server struct {
// 	pb.UnimplementedEndpointChainEventListenerServer // 嵌入未实现的服务接口
// }

// // 实现 ListenForEvents 方法
// func (s *server) ListenForEvents(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
// 	eventData := req.GetEventData()
// 	fmt.Printf("Received event data: %s\n", eventData)

// 	// 构造响应
// 	response := &pb.EventResponse{
// 		Success: true,
// 		Message: "Accepted",
// 	}

// 	fmt.Printf("Sending response: %+v\n", response)
// 	return response, nil
// }

// func main() {
// 	// 监听所有网络接口
// 	lis, err := net.Listen("tcp", "0.0.0.0:50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	// 创建 gRPC 服务器
// 	s := grpc.NewServer()

// 	// 注册服务
// 	pb.RegisterEndpointChainEventListenerServer(s, &server{})

// 	fmt.Println("Server running at 0.0.0.0:50051")

// 	// 启动服务器
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }


// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"
// 	"sync"
// 	"time"

// 	"chainmaker.org/chainmaker/pb-go/v2/common"
// 	sdk "chainmaker.org/chainmaker/sdk-go/v2"
// 	pb "test/proto" // 替换为你的 proto 生成包的实际路径
// 	"google.golang.org/grpc"
// )

// const (
// 	sdkConfigPath = "./config/sdk_config.yml"
// 	contractName  = "Ethcontract"
// 	invokeMethod  = "transferContent"
// 	maxWorkers    = 10 // 最大并发工作协程数
// 	queueSize     = 1000 // 请求队列缓冲区大小
// )

// type Job struct {
// 	ctx      context.Context
// 	req      *pb.EventRequest
// 	respChan chan *pb.EventResponse
// }

// // gRPC 服务端结构体
// type server struct {
// 	pb.UnimplementedEndpointChainEventListenerServer
// 	client       *sdk.ChainClient
// 	eventChan    chan *common.ContractEventInfo // 用于传递链上事件
// 	pendingReqs  sync.Map                   // 存储等待事件的请求
// 	jobChan      chan Job                   // 请求队列
// }

// // 实现 ListenForEvents 方法
// func (s *server) ListenForEvents(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
// 	respChan := make(chan *pb.EventResponse, 1)
// 	job := Job{ctx: ctx, req: req, respChan: respChan}

// 	// 将请求放入队列
// 	select {
// 	case s.jobChan <- job:
// 		// 等待响应
// 		select {
// 		case resp := <-respChan:
// 			return resp, nil
// 		case <-ctx.Done():
// 			return &pb.EventResponse{Success: false, Message: "请求超时或取消"}, nil
// 		}
// 	default:
// 		return &pb.EventResponse{Success: false, Message: "服务器忙，队列已满"}, nil
// 	}
// }

// // worker 处理请求
// func (s *server) startWorkers() {
// 	for i := 0; i < maxWorkers; i++ {
// 		go func() {
// 			for job := range s.jobChan {
// 				resp, err := s.processRequest(job.ctx, job.req)
// 				if err != nil {
// 					log.Printf("处理请求失败: %v", err)
// 				}
// 				job.respChan <- resp
// 			}
// 		}()
// 	}
// }

// // 处理单个请求
// func (s *server) processRequest(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
// 	startTime := time.Now()
// 	defer func() {
// 		fmt.Printf("请求处理耗时: %v\n", time.Since(startTime))
// 	}()

// 	eventData := req.GetEventData()
// 	fmt.Printf("Received event data: %s\n", eventData)

// 	// 解析事件数据
// 	params, err := parseEventData(eventData)
// 	if err != nil {
// 		return &pb.EventResponse{Success: false, Message: fmt.Sprintf("解析事件数据失败: %v", err)}, nil
// 	}

// 	// 调用合约（带重试）
// 	txId, err := s.invokeContractWithRetry(params, 3)
// 	if err != nil {
// 		return &pb.EventResponse{Success: false, Message: fmt.Sprintf("调用合约失败: %v", err)}, nil
// 	}

// 	// 注册等待事件
// 	respChan := make(chan *pb.EventResponse, 1)
// 	s.pendingReqs.Store(txId, respChan)

// 	// 等待链上事件或超时
// 	select {
// 	case resp := <-respChan:
// 		s.pendingReqs.Delete(txId)
// 		return resp, nil
// 	case <-time.After(100 * time.Second):
// 		s.pendingReqs.Delete(txId)
// 		return &pb.EventResponse{Success: false, Message: "等待链上事件超时"}, nil
// 	case <-ctx.Done():
// 		s.pendingReqs.Delete(txId)
// 		return &pb.EventResponse{Success: false, Message: "请求被取消"}, nil
// 	}
// }

// // 解析事件数据
// func parseEventData(eventData string) (map[string]string, error) {
// 	params := make(map[string]string)
// 	parts := strings.Split(eventData, ",")
// 	for _, part := range parts {
// 		kv := strings.SplitN(part, ":", 2)
// 		if len(kv) == 2 {
// 			key := strings.Trim(kv[0], "\"{")
// 			value := strings.Trim(kv[1], "\"}")
// 			switch key {
// 			case "from_address":
// 				params["from"] = value
// 			case "to_address":
// 				params["to"] = value
// 			case "content_type":
// 				params["content_type"] = value
// 			case "content":
// 				params["content_address"] = value
// 			}
// 		}
// 	}
// 	if len(params) < 4 {
// 		return nil, fmt.Errorf("事件数据缺少必要字段")
// 	}
// 	return params, nil
// }

// // 检查合约是否存在
// func checkContract(client *sdk.ChainClient) error {
// 	contract, err := client.GetContractInfo(contractName)
// 	if err != nil {
// 		return fmt.Errorf("获取合约信息失败: %v", err)
// 	}
// 	if contract == nil || contract.Name == "" {
// 		return fmt.Errorf("合约 %s 不存在", contractName)
// 	}
// 	fmt.Printf("合约已存在: %+v\n", contract)
// 	return nil
// }

// // 调用合约并返回交易ID（带重试）
// func (s *server) invokeContractWithRetry(params map[string]string, retries int) (string, error) {
// 	for i := 0; i < retries; i++ {
// 		fmt.Println("====================== 调用 Ethcontract 合约 ======================")
// 		kvs := []*common.KeyValuePair{
// 			{Key: "from", Value: []byte(params["from"])},
// 			{Key: "to", Value: []byte(params["to"])},
// 			{Key: "content_type", Value: []byte(params["content_type"])},
// 			{Key: "content_address", Value: []byte(params["content_address"])},
// 		}

// 		resp, err := s.client.InvokeContract(contractName, invokeMethod, "", kvs, -1, true)
// 		if err == nil && resp.Code == common.TxStatusCode_SUCCESS {
// 			fmt.Printf("交易成功 - TxId: %s, BlockHeight: %d, Message: %s\n",
// 				resp.TxId, resp.TxBlockHeight, resp.ContractResult.Message)
// 			return resp.TxId, nil
// 		}
// 		fmt.Printf("调用合约失败，第 %d 次重试: %v\n", i+1, err)
// 		time.Sleep(2 * time.Second) // 重试前等待
// 	}
// 	return "", fmt.Errorf("调用合约失败，已达最大重试次数")
// }

// // 处理链上事件
// func (s *server) handleEvent(event *common.ContractEventInfo) {
// 	if respChan, ok := s.pendingReqs.Load(event.TxId); ok {
// 		respChan.(chan *pb.EventResponse) <- &pb.EventResponse{
// 			Success: true,
// 			Message: "accept",
// 		}
// 		s.pendingReqs.Delete(event.TxId)
// 	}
// }

// // 订阅并监听链上事件
// func subscribeToEvents(client *sdk.ChainClient, s *server) {
// 	fmt.Println("====================== 开始订阅链上事件 ======================")
// 	ctx := context.Background()
// 	eventCh, err := client.SubscribeContractEvent(ctx, 0, -1, contractName, "TransferEvent")
// 	if err != nil {
// 		log.Fatalf("订阅事件失败: %v", err)
// 	}

// 	go func() {
// 		for event := range eventCh {
// 			contractEvent, ok := event.(*common.ContractEventInfo)
// 			if !ok {
// 				continue
// 			}
// 			fmt.Printf("接收到链上事件:\n")
// 			fmt.Printf("TxId: %s\n", contractEvent.TxId)
// 			fmt.Printf("Topic: %s\n", contractEvent.Topic)
// 			fmt.Printf("ContractName: %s\n", contractEvent.ContractName)
// 			for i, data := range contractEvent.EventData {
// 				fmt.Printf("EventData[%d]: %s\n", i, string(data))
// 			}
// 			fmt.Println("-------------------")
// 			s.handleEvent(contractEvent)
// 		}
// 	}()
// }

// func main() {
// 	// 初始化 Chainmaker 客户端
// 	client, err := sdk.NewChainClient(
// 		sdk.WithConfPath(sdkConfigPath),
// 	)
// 	if err != nil {
// 		log.Fatalf("创建 Chainmaker 客户端失败: %v", err)
// 	}

// 	// 检查合约是否存在
// 	if err := checkContract(client); err != nil {
// 		log.Fatalf("合约检查失败: %v", err)
// 	}

// 	// 创建事件通道和请求队列
// 	eventChan := make(chan *common.ContractEventInfo, 10)
// 	jobChan := make(chan Job, queueSize)

// 	// 初始化服务
// 	srv := &server{
// 		client:      client,
// 		eventChan:   eventChan,
// 		pendingReqs: sync.Map{},
// 		jobChan:     jobChan,
// 	}

// 	// 启动 worker pool
// 	go srv.startWorkers()

// 	// 订阅事件
// 	subscribeToEvents(client, srv)

// 	// 启动 gRPC 服务
// 	lis, err := net.Listen("tcp", "0.0.0.0:50051")
// 	if err != nil {
// 		log.Fatalf("监听端口失败: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	pb.RegisterEndpointChainEventListenerServer(grpcServer, srv)
// 	fmt.Println("gRPC 服务运行在 0.0.0.0:50051")

// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("gRPC 服务启动失败: %v", err)
// 	}
// }