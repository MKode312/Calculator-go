package agent

import (
	"calculator_go/internal/utils/agent/calculation"
	"calculator_go/internal/utils/agent/validator"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	itp "calculator_go/internal/utils/agent/infix_to_postfix"
	pb "calculator_go/protos/gen"

	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	pb.CalculatorServiceServer
}

func NewServer() *Server {
	return &Server{}
}

type CalculatorServiceServer interface {
	Calculate(context.Context, *pb.ExpressionRequest) (*pb.ExpressionResponse, error)
	mustEmbedUnimplementedCalculatorServiceServer()
}

func (s *Server) Calculate(ctx context.Context, in *pb.ExpressionRequest) (*pb.ExpressionResponse, error) {

	expr := in.GetExpression()
	if !validator.IsValidExpression(expr) {
		return nil, fmt.Errorf("invalid expression: %s", in.Expression)
	}
	postfixed := itp.ToPostfix(expr)
	res, err := calculation.Evaluate(postfixed)
	if err != nil {
		return nil, err
	}
	log.Println("successful operation!")
	return &pb.ExpressionResponse{Result: res}, nil
}

func RunAgentServer() {

	host, ok := os.LookupEnv("AGENT_HOST")
	if !ok {
		log.Print("AGENT_HOST not set, using 0.0.0.0")
		host = "0.0.0.0"
	}

	port, ok := os.LookupEnv("AGENT_PORT")
	if !ok {
		log.Print("AGENT_PORT not set, using 5000")
		port = "5000"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Printf("tcp listener started at %s", addr)

	grpcServer := grpc.NewServer()

	expressionServiceServer := NewServer()

	pb.RegisterCalculatorServiceServer(grpcServer, expressionServiceServer)

	go log.Fatal(grpcServer.Serve(lis))
}
