package main

import (
	"antrein/bc-dashboard/application/common/repository"
	"antrein/bc-dashboard/application/common/resource"
	"antrein/bc-dashboard/application/common/usecase"
	"antrein/bc-dashboard/application/grpc"
	"antrein/bc-dashboard/application/rest"
	"antrein/bc-dashboard/model/config"
	"context"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	resource, err := resource.NewCommonResource(cfg, ctx)
	if err != nil {
		log.Fatal(err)
	}
	repository, err := repository.NewCommonRepository(cfg, resource)
	if err != nil {
		log.Fatal(err)
	}
	uc, err := usecase.NewCommonUsecase(cfg, repository)
	if err != nil {
		log.Fatal(err)
	}

	rest_app, err := rest.ApplicationDelegate(cfg, uc, resource)
	if err != nil {
		log.Fatal(err)
	}

	// Start gRPC server concurrently
	go func() {
		grpc_app, err := grpc.ApplicationDelegate(cfg, uc, resource)
		if err != nil {
			log.Fatal(err)
		}
		if err := grpc.StartServer(cfg, grpc_app); err != nil {
			log.Fatal(err)
		}
	}()

	if err = rest.StartServer(cfg, rest_app); err != nil {
		log.Fatal(err)
	}
}
