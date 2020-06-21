package spread

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	// "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cripplet/event-spread/lib/core/handlers"
	espb "github.com/cripplet/event-spread/lib/proto/event_spread_go_proto"
)

// TODO(cripplet): Decide if we need to do full e2e testing via
//   https://godoc.org/google.golang.org/grpc/test/bufconn
//   https://godoc.org/google.golang.org/grpc/test/grpc_testing

var (
	initialTime = time.Now()

	// TODO(cripplet): Use the new-style Timestamp constructors once the release picks up the syntax.
	// ts = timestamppb.New(initialTime)
	ts, _ = ptypes.TimestampProto(initialTime)

	trivialHeuristicValue = &espb.HeuristicValue{
		Heuristic: espb.Heuristic_HEURISTIC_MORALITY,
		Value:     100,
	}

	trivialEvent = &espb.Event{
		Position:   &espb.Position{X: 0, Y: 0},
		Timestamp:  ts,
		Heuristics: []*espb.HeuristicValue{trivialHeuristicValue},
		SpreadType: espb.SpreadType_SPREAD_TYPE_INSTANT_GLOBAL,
	}

	dispatcher = map[espb.SpreadType]handlers.EventSpreadHandler{
		espb.SpreadType_SPREAD_TYPE_INSTANT_GLOBAL: &handlers.InstantGlobalEventSpreadHandler{},
	}
)

func TestAddNullEvent(t *testing.T) {
	s := NewEventSpreadService(nil)
	_, err := s.AddEvent(context.Background(), &espb.AddEventRequest{})
	if err == nil {
		t.Error("unexpectedly succeeded in adding empty Event to event queue")
	}
}

func TestAddEvent(t *testing.T) {
	s := NewEventSpreadService(nil)

	_, err := s.AddEvent(context.Background(), &espb.AddEventRequest{Event: trivialEvent})
	if err != nil {
		t.Errorf("unexpectedly received error when adding Event: %v", err)
	}

	s.eventsMux.Lock()
	defer s.eventsMux.Unlock()
	if len(s.events) != 1 {
		t.Errorf("unexpected event queue length: expected %v, but got %v", 1, len(s.events))
	}
}

func EventSpreadHelper(
	t *testing.T,
	s *EventSpreadService,
	events []*espb.Event,
	queryTime time.Time,
	hs []espb.Heuristic) (*espb.GetEventSpreadResponse, error) {
	for _, e := range events {
		s.AddEvent(context.Background(), &espb.AddEventRequest{Event: e})
	}

	q, _ := ptypes.TimestampProto(queryTime)
	req := &espb.GetEventSpreadRequest{
		Heuristics: hs,
		Timestamp:  q,
	}

	return s.GetEventSpread(context.Background(), req)
}

func TestGetEventSpreadNullEvents(t *testing.T) {
	h := espb.Heuristic_HEURISTIC_MORALITY
	s := NewEventSpreadService(dispatcher)

	resp, err := EventSpreadHelper(
		t,
		s,
		[]*espb.Event{},
		initialTime.Add(time.Second),
		[]espb.Heuristic{h},
	)
	if err != nil {
		t.Errorf("unexpectedly received error when querying for event value: %v", err)
	}

	if len(resp.GetValues()) != 1 {
		t.Errorf("unexpected response length when querying for event value: expected 1 != %v", len(resp.GetValues()))
	}
	if resp.GetValues()[0].GetHeuristic() != h {
		t.Errorf("unexpected Heuristic enum value: expected %v != %v", h, resp.GetValues()[0].GetHeuristic())
	}
	if resp.GetValues()[0].GetValue() != 0 {
		t.Errorf("unexpected heuristic value: expected 0 != %v", trivialEvent.GetHeuristics()[0].GetValue())
	}
}

func TestGetEventSpread(t *testing.T) {
	h := espb.Heuristic_HEURISTIC_MORALITY
	s := NewEventSpreadService(dispatcher)

	resp, err := EventSpreadHelper(
		t,
		s,
		[]*espb.Event{trivialEvent},
		initialTime.Add(time.Second),
		[]espb.Heuristic{h},
	)
	if err != nil {
		t.Errorf("unexpectedly received error when querying for event value: %v", err)
	}

	if len(resp.GetValues()) != 1 {
		t.Errorf("unexpected response length when querying for event value: expected 1 != %v", len(resp.GetValues()))
	}
	if resp.GetValues()[0].GetHeuristic() != h {
		t.Errorf("unexpected Heuristic enum value: expected %v != %v", h, resp.GetValues()[0].GetHeuristic())
	}
	if resp.GetValues()[0].GetValue() != trivialHeuristicValue.GetValue() {
		t.Errorf("unexpected heuristic value: expected %v != %v", trivialHeuristicValue.GetValue(), trivialEvent.GetHeuristics()[0].GetValue())
	}
}
