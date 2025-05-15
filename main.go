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

    sdk "chainmaker.org/chainmaker/sdk-go/v2"
    pb "test/proto"
    "google.golang.org/grpc"
)

const (
    sdkConfigPath = "./config/sdk_config.yml"
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

    // 解析事件数据（保留原有逻辑）
    params, err := parseEventData(eventData)
    if err != nil {
        return &pb.EventResponse{
            Success: false,
            Message: fmt.Sprintf("解析事件数据失败: %v", err),
        }, nil
    }

    // 直接使用接收消息的时间作为结束时间
    timeHash := time.Now().Format("2006-01-02 15:04:05.000")

    // 写入文件（保留原有逻辑）
    f, err := os.OpenFile("data.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err == nil {
        defer f.Close()
        content := fmt.Sprintf("time_send:%s|time_confirm:%s|lis:%s|rsend:%s|time_end:%s\n",
            params["sendts"], params["confirmtx"], params["lis"], params["rsend"], timeHash)
        f.WriteString(content)
        // 另起一行打印pre_node字段，如果存在的话
        if preNode, exists := params["pre_node"]; exists && preNode != "" {
            preNodeContent := fmt.Sprintf("pre_node:%s\n", preNode)
            f.WriteString(preNodeContent)
        }
    } else {
        log.Printf("写入 data.out 失败: %v", err)
    }

    // 简单增加交易计数并打印（保留原有逻辑）
    transactionCounter++
    if preNode, exists := params["pre_node"]; exists && preNode != "" {
        fmt.Printf("交易 #%d 的 pre_node: %s\n", transactionCounter, preNode)
    }
    return &pb.EventResponse{
        Success: true,
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

    // 获取pre_nodes字段
    if preNodes, ok := jsonData["pre_nodes"]; ok {
        params["pre_node"] = preNodes // 保持内部名称不变，便于其他代码使用
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

func main() {
    client, err := sdk.NewChainClient(
        sdk.WithConfPath(sdkConfigPath),
    )
    if err != nil {
        log.Fatalf("创建 Chainmaker 客户端失败: %v", err)
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