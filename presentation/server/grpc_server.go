package grpc_server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	grpcmsg "webserver/presentation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *grpcmsg.HelloRequest) (*grpcmsg.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &grpcmsg.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	// 取設定檔的port
	// configPort := viper.GetString("GRPC_PORT")
	// if configPort != nil {
	// 	port = configPort
	// }

	// 監聽指定port，這樣服務才能在該port執行。
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("端口(Port)監聽失敗: %v", err)
	}

	// 建立新 gRPC 伺服器並註冊服務。
	s := grpc.NewServer()
	// grpcmsg.RegisterGreeterServer(s, &server{})

	// 在 gRPC 伺服器上註冊反射服務。
	reflection.Register(s)

	log.Printf("服務器正在監聽 %v", lis.Addr())

	// 開始在指定port中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("斷開服務: %v", err)
	}
}

// 服務器用於實現helloworld.GreeterServer。
type server struct{}
