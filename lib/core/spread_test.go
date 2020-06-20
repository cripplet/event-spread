package spread

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	// timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	espb "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto"
)

// TODO(cripplet): Decide if we need to do full e2e testing via
//   https://godoc.org/google.golang.org/grpc/test/bufconn
//   https://godoc.org/google.golang.org/grpc/test/grpc_testing

var (
	initialTime = time.Now()

	// TODO(cripplet): Use the new-style Timestamp constructors once the release picks up the syntax.
	// ts = timestamppb.New(initialTime)
	ts, _ = ptypes.TimestampProto(initialTime)

	trivialEvent = &espb.Event{
		Position: &espb.Position{X: 0, Y: 0},
		Timestamp: ts,
		Heuristics: []*espb.HeuristicValue{
			&espb.HeuristicValue{
				Heuristic: espb.Heuristic_HEURISTIC_MORALITY,
				Value: 100,
			},
		},
		SpreadType: espb.SpreadType_SPREAD_TYPE_SIMPLE_LINEAR,
		SpreadRate: 1,
	}
)

func TestAddNullEvent(t *testing.T) {
	s, _ := NewEventSpreadService()
	_, err := s.AddEvent(context.Background(), &espb.AddEventRequest{})
	if err == nil {
		t.Error("unexpectedly succeeded in adding empty Event to event queue")
	}
}

func TestAddEvent(t *testing.T) {
	s, _ := NewEventSpreadService()

	_, err := s.AddEvent(context.Background(), &espb.AddEventRequest{
		Event: trivialEvent,
	})
	if err != nil {
		t.Errorf("unexpectedly received error when adding Event: %v", err)
	}

	s.eventsMux.Lock()
	defer s.eventsMux.Unlock()
	if len(s.events) != 1 {
		t.Errorf("unexpected event queue length: expected %v, but got %v", 1, len(s.events))
	}
}

func TestGetEventSpreadNullMatch(t *testing.T) {
	h := espb.Heuristic_HEURISTIC_MORALITY
	propagateTime, _ := ptypes.TimestampProto(initialTime.Add(10 * time.Second))

	s, _ := NewEventSpreadService()

	resp, err := s.GetEventSpread(context.Background(), &espb.GetEventSpreadRequest{
		Heuristics: []espb.Heuristic{h},
		Timestamp: propagateTime,
	})
	if err != nil {
		t.Errorf("unexpectedly received error when querying for event value: %v", err)
	}

	if len(resp.GetValues()) != 1 {
		t.Errorf("unexpected response length when querying for event value: expected 1 != %v", len(resp.GetValues()))
	}
	if resp.GetValues()[0].GetHeuristic() != h {
		t.Errorf("unexpected Heuristic enum value: expected %v != %v", h, resp.GetValues()[0].GetHeuristic())
	}
	if resp.GetValues()[0].GetValue() != trivialEvent.GetHeuristics()[0].GetValue() {
		t.Errorf("unexpected heuristic value: expected %v != %v", trivialEvent.GetHeuristics()[0].GetValue(), resp.GetValues()[0].GetValue())
	}
}
