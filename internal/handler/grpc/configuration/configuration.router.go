package configuration

import (
	"antrein/bc-dashboard/internal/usecase/configuration"
	"context"
	"errors"

	pb "github.com/antrein/proto-repository/pb/bc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedProjectConfigServiceServer
	usecase *configuration.Usecase
}

func New(usecase *configuration.Usecase) *Server {
	return &Server{
		usecase: usecase,
	}
}

func (s *Server) GetProjectConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ProjectConfigResponse, error) {
	projectID := in.GetProjectId()
	resp, err := s.usecase.GetProjectConfigByID(ctx, projectID)
	if err != nil {
		return nil, errors.New(err.Error)
	}
	return &pb.ProjectConfigResponse{
		ProjectId:       projectID,
		Threshold:       int32(resp.Threshold),
		SessionTime:     int32(resp.SessionTime),
		Host:            resp.Host,
		BaseUrl:         resp.BaseURL,
		MaxUsersInQueue: int32(resp.MaxUsersInQueue),
		QueueStart:      timestamppb.New(resp.QueueStart),
		QueueEnd:        timestamppb.New(resp.QueueEnd),
		QueuePageStyle:  resp.QueuePageStyle,
		QueueHtmlPage:   resp.QueueHTMLPage,
		QueuePageTitle:  resp.QueuePageTitle,
		QueuePageLogo:   resp.QueuePageLogo,
		IsConfigure:     resp.IsConfigure,
	}, nil
}
