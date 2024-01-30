package grpctry

import (
	context "context"
	"fmt"
	"log"
	logsgenerated "logger-service/cmd/api/logs"
	data "logger-service/cmd/data"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type LogGrpcServer struct {
	logsgenerated.UnimplementedLogServeiceServer
	Models data.Models
}

func NewLogGrpcServer(mongo *mongo.Client) *LogGrpcServer {
	return &LogGrpcServer{
		Models: data.New(mongo),
	}
}

func (l *LogGrpcServer) WriteLog(ctx context.Context, req *logsgenerated.LogRequest) (*logsgenerated.LogResponse, error) {
	input := req.GetLogEntry()
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}

	err := logEntry.Create()
	if err != nil {
		res := &logsgenerated.LogResponse{Message: "Error while creating log entry"}
		return res, err
	}

	// Loggging success

	res := &logsgenerated.LogResponse{Message: "Log entry created successfully using gRPC"}

	return res, nil
}

func (logGrpcServer *LogGrpcServer) ListenTOGRPC(GrpcPort string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	logsgenerated.RegisterLogServeiceServer(s, logGrpcServer)

	log.Printf("Starting gRPC server on port %v", GrpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
