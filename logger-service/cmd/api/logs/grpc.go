package logs

import (
	context "context"
	"fmt"
	"log"
	main "logger-service/cmd/api"
	"logger-service/cmd/api/handellers"
	"logger-service/cmd/data"
	"net"
)

type LogGrpcServer struct {
	UnimplementedLogServeiceServer
	Models data.Models
}

func (l *LogGrpcServer) WriteLog(ctx context.Context, req *LogRequest) (*LogResponse, error) {
	input := req.GetLogEntry()
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}

	err := logEntry.Create()
	if err != nil {
		res := &LogResponse{Message: "Error while creating log entry"}
		return res, err
	}

	// Loggging success

	res := &LogResponse{Message: "Log entry created successfully using gRPC"}

	return res, nil
}

func listenTOGRPC(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", main.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	
}
