package spread

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"

	espb "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto"
)

type EventSpreadService struct {}

func (s *EventSpreadService) AddEvent(ctx context.Context, req *espb.AddEventRequest) (*espb.AddEventResponse, error) {
	return nil, status.Error(codes.Unimplemented, "AddEvent has not been implemented")
}

func (s *EventSpreadService) GetEventSpread(ctx context.Context, req *espb.GetEventSpreadRequest) (*espb.GetEventSpreadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetEventSpread has not been implemented")
}

func NewServer() (*grpc.Server, error) {
	s := grpc.NewServer()
	espb.RegisterEventSpreadServiceServer(s, &EventSpreadService{})
	return s, nil
}
