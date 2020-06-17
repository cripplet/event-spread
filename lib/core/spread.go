package spread

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	espb "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto"
)

type EventSpreadService struct {
	eventsMux sync.Mutex
	events []*espb.Event
}

func (s *EventSpreadService) AddEvent(ctx context.Context, req *espb.AddEventRequest) (*espb.AddEventResponse, error) {
	s.eventsMux.Lock()
	defer s.eventsMux.Unlock()

	if req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "cannot specify an empty Event to add to event queue")
	}

	s.events = append(s.events, proto.Clone(req.GetEvent()).(*espb.Event))
	return &espb.AddEventResponse{}, nil
}

func (s *EventSpreadService) GetEventSpread(ctx context.Context, req *espb.GetEventSpreadRequest) (*espb.GetEventSpreadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetEventSpread has not been implemented")
}

func NewEventSpreadService() (*EventSpreadService, error) {
	return &EventSpreadService{
	}, nil
}

func NewServer() (*grpc.Server, error) {
	service, _ := NewEventSpreadService()
	s := grpc.NewServer()
	espb.RegisterEventSpreadServiceServer(s, service)
	return s, nil
}
