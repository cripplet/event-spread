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

// EventSpreadService implements the espb.EventSpreadService RPC.
type EventSpreadService struct {
	eventsMux sync.Mutex
	events []*espb.Event
}

// AddEvent adds the specified espb.Event object into the state queue. Events
// will be used to calculate heuristic values from the global state.
func (s *EventSpreadService) AddEvent(ctx context.Context, req *espb.AddEventRequest) (*espb.AddEventResponse, error) {
	s.eventsMux.Lock()
	defer s.eventsMux.Unlock()

	if req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "cannot specify an empty Event to add to event queue")
	}

	s.events = append(s.events, proto.Clone(req.GetEvent()).(*espb.Event))
	return &espb.AddEventResponse{}, nil
}

func eventSpread(e *espb.Event, req *espb.GetEventSpreadRequest, vc chan<- *espb.HeuristicValue) error {
	// TODO(cripplet): Implement.
	return nil
}

func (s *EventSpreadService) GetEventSpread(ctx context.Context, req *espb.GetEventSpreadRequest) (*espb.GetEventSpreadResponse, error) {
	s.eventsMux.Lock()
	defer s.eventsMux.Unlock()

	vc := make(chan *espb.HeuristicValue)
	ec := make(chan error)
	var wg sync.WaitGroup

	wg.Add(len(s.events))
	for _, e := range s.events {
		go func(e *espb.Event, req *espb.GetEventSpreadRequest, vc chan<- *espb.HeuristicValue, ec chan<- error) {
			defer wg.Done()
			ec <- eventSpread(e, req, vc)
		}(e, req, vc, ec)
	}
	wg.Wait()
	close(ec)
	close(vc)

	var errors []error
	for err := range ec {
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return nil, status.Errorf(codes.Aborted, "could not calculate event spread due to error(s) %v", errors)
	}

	resp := &espb.GetEventSpreadResponse{}
	for _ = range vc {
		// TODO(cripplet): Switch over to SetValue() once setters are part of Golang protobuf API.
		// TODO(cripplet): Implement.
	}

	return resp, nil
}

// NewEventSpreadService constructs a new implementation object.
func NewEventSpreadService() (*EventSpreadService, error) {
	return &EventSpreadService{}, nil
}

// NewServer creates and returns an RPC service. This will be used in binaries
// to actually run the server against a port.
func NewServer() (*grpc.Server, error) {
	service, _ := NewEventSpreadService()
	s := grpc.NewServer()
	espb.RegisterEventSpreadServiceServer(s, service)
	return s, nil
}
