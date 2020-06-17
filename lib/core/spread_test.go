package spread

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	espb "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto"
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
	ts, _ := ptypes.TimestampProto(time.Now())
	_, err := s.AddEvent(context.Background(), &espb.AddEventRequest{
		Event: &espb.Event{
			Position: &espb.Position{X: 0, Y: 0},
			Timestamp: ts,
			Heuristics: []*espb.HeuristicValue{
				&espb.HeuristicValue{
					Heuristic: espb.Heuristic_HEURISTIC_MORALITY,
					Value: 0,
				},
			},
			SpreadType: espb.SpreadType_SPREAD_TYPE_SIMPLE_LINEAR,
			SpreadRate: 0,
		},
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
