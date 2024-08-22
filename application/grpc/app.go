package grpc

import (
	"antrein/bc-dashboard/application/common/resource"
	"antrein/bc-dashboard/application/common/usecase"
	"antrein/bc-dashboard/internal/handler/grpc/configuration"
	"antrein/bc-dashboard/model/config"
	"context"

	pb "github.com/antrein/proto-repository/pb/bc"
	"google.golang.org/grpc"
)

type helloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func ApplicationDelegate(cfg *config.Config, uc *usecase.CommonUsecase, rsc *resource.CommonResource) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	// Hello service
	helloServer := &helloServer{}
	pb.RegisterGreeterServer(grpcServer, helloServer)

	// Project config service
	projectConfigServer := configuration.New(uc.ConfigUsecase)
	pb.RegisterProjectConfigServiceServer(grpcServer, projectConfigServer)

	return grpcServer, nil
}
