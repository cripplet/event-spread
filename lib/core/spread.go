package spread

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

        // TODO(cripplet): Use the new-style Timestamp constructors once the release picks up the syntax.
        // "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cripplet/event-spread/core/handlers"
	espb "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto"
)

// NewEventSpreadService constructs a new implementation object.
func NewEventSpreadService(dispatcher map[espb.SpreadType]handlers.EventSpreadHandler) (*EventSpreadService, error) {
	return &EventSpreadService{
		spreadTypeDispatcher: dispatcher,
	}, nil
}

// NewServer creates and returns an RPC service. This will be used in binaries
// to actually run the server against a port.
func NewServer(s *EventSpreadService) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	espb.RegisterEventSpreadServiceServer(grpcServer, s)
	return grpcServer, nil
}

// EventSpreadService implements the espb.EventSpreadService RPC.
type EventSpreadService struct {
	spreadTypeDispatcher map[espb.SpreadType]handlers.EventSpreadHandler
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

// GetEventSpread gets total influence of Heuristic types specified in the request.
func (s *EventSpreadService) GetEventSpread(ctx context.Context, req *espb.GetEventSpreadRequest) (*espb.GetEventSpreadResponse, error) {
	// Make a copy of events -- some level of inconsistency between ~concurrent requests is expected.
	events := func() []*espb.Event {
		s.eventsMux.Lock()
		defer s.eventsMux.Unlock()
		es := []*espb.Event{}
		for _, e := range s.events {
			es = append(es, proto.Clone(e).(*espb.Event))
		}
		return es
	}()

	var l []*espb.HeuristicValue
	for _, h := range req.GetHeuristics() {
		l = append(l, &espb.HeuristicValue{ Heuristic: h })
	}

	// TODO(cripplet): Make this concurrent.
	for _, e := range events {
		ch, err := handlers.EventSpread(s.spreadTypeDispatcher, e, req)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get influence, got error %v", err)
		}
		for hv := range ch {
			l = append(l, hv)
		}
	}

	return &espb.GetEventSpreadResponse{
		Values: handlers.MapToList(handlers.ListToMap(l)),
	}, nil
}
