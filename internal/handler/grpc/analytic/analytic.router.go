package analytic

import (
	guard "antrein/bc-dashboard/application/middleware"
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"context"
	"log"
	"net/http"

	pb "github.com/antrein/proto-repository/pb/bc"
	"github.com/gorilla/mux"

	"google.golang.org/grpc"
)

type Client struct {
	cfg        *config.Config
	grpcClient *grpc.ClientConn
}

func New(cfg *config.Config, gc *grpc.ClientConn) *Client {
	return &Client{
		cfg:        cfg,
		grpcClient: gc,
	}
}

func (c *Client) RegisterRoute(app *mux.Router) {
	app.HandleFunc("/bc/dashboard/analytic", guard.DefaultGuard(c.StreamAnalyticData))
	app.HandleFunc("/bc/dashboard/analytic/{id}", guard.AuthGuard(c.cfg, c.GetProjectAnalytic))

}

func (c *Client) StreamAnalyticData(g *guard.GuardContext) error {
	// Set headers for SSE
	ctx := context.Background()
	g.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	g.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	g.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	g.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	g.ResponseWriter.Header().Set("Connection", "keep-alive")

	projectID := g.Request.URL.Query().Get("project_id")

	client := pb.NewAnalyticServiceClient(c.grpcClient)

	stream, err := client.StreamRealtimeData(ctx, &pb.AnalyticRequest{
		ProjectId: projectID,
	})
	if err != nil {
		log.Println(err)
		return g.ReturnError(http.StatusInternalServerError, "Error connecting to gRPC stream")
	}

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return g.ReturnSuccess("Stream done")
		default:
			analyticData, err := stream.Recv()
			if err != nil {
				log.Println("Error receiving data from gRPC:", err)
				return g.ReturnError(http.StatusInternalServerError, "Error sent data from stream")
			}
			data := dto.Analytic{
				ProjectID:         projectID,
				TimeStamp:         analyticData.GetTimestamp().AsTime(),
				TotalUsersInQueue: int(analyticData.TotalUsersInQueue),
				TotalUsersInRoom:  int(analyticData.TotalUsersInRoom),
				TotalUsers:        int(analyticData.TotalUsers),
			}
			err = g.ReturnEvent(data)
			if err != nil {
				log.Println("Error sending data to client:", err)
				return g.ReturnError(http.StatusInternalServerError, "Error sent data from stream")
			}
		}
	}

	return g.ReturnSuccess("Stream done")
}

func (c *Client) GetProjectAnalytic(g *guard.AuthGuardContext) error {
	ok := guard.IsMethod(g.Request, "GET")
	if !ok {
		return g.ReturnError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	projectID := guard.GetParam(g.Request, "id")
	ctx := context.Background()
	client := pb.NewAnalyticServiceClient(c.grpcClient)

	stream, err := client.StreamRealtimeData(ctx, &pb.AnalyticRequest{
		ProjectId: projectID,
	})

	if err != nil {
		log.Println(err)
		return g.ReturnError(http.StatusInternalServerError, "Error connecting to gRPC stream")
	}

	analyticData, err := stream.Recv()
	if err != nil {
		log.Println("Error receiving data from gRPC:", err)
		return g.ReturnError(http.StatusInternalServerError, "Error sent data from stream")
	}
	resp := dto.Analytic{
		ProjectID:         projectID,
		TimeStamp:         analyticData.GetTimestamp().AsTime(),
		TotalUsersInQueue: int(analyticData.TotalUsersInQueue),
		TotalUsersInRoom:  int(analyticData.TotalUsersInRoom),
		TotalUsers:        int(analyticData.TotalUsers),
	}

	return g.ReturnSuccess(resp)
}
