package handlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	// TODO(cripplet): Use the new-style Timestamp constructors once the release picks up the syntax.
	// "google.golang.org/protobuf/types/known/timestamppb"

	espb "github.com/cripplet/event-spread/lib/proto/event_spread_go_proto"
)

// EventSpread will calculate the global influence contribution by a single Event object and return a
// channel containing the influence broken down by Heuristic. The channel is closed after all values
// have been pushed into the channel.
func EventSpread(dispatcher map[espb.SpreadType]EventSpreadHandler, e *espb.Event, req *espb.GetEventSpreadRequest) (<-chan *espb.HeuristicValue, error) {
	vc := make(chan *espb.HeuristicValue)

	h, found := dispatcher[e.GetSpreadType()]
	if !found {
		return nil, status.Errorf(codes.Unimplemented, "specified event SpreadType %v has not been implemented", e.GetSpreadType())
	}

	hv, err := h.EventSpread(e, req)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(vc)
		for _, h := range hv {
			vc <- h
		}
	}()
	return vc, nil
}

// EventSpreadHandler defines an interface by which different SpreadType enums may define their own
// influence contribution, e.g. SPREAD_TYPE_INSTANT_GLOBAL will return the full HeuristicValue
// influence immediately after the Event is added to the queue.
//
// Handlers may need to query some global map object. In this case, the implemented handlers
// should connect to a map provider service whenever possible.
type EventSpreadHandler interface {
	// IsPropagated returns true if the given Event may be considered propagated globally
	// at the given timestamp. This is useful for marking tombstoned objects when advancing
	// the global state for performance optimization.
	IsPropagated(e *espb.Event, t *timestamp.Timestamp) (bool, error)

	// EventSpread does the heavy lifting in calculating the values of the influence propagation.
	EventSpread(e *espb.Event, req *espb.GetEventSpreadRequest) ([]*espb.HeuristicValue, error)
}

// ListToMap is a convenience function to help implement contains-syntax testing and comparison.
// Multiple HeuristicValue instances with the same Heuristic are merged together as the sum of
// individual values.
func ListToMap(heuristicValues []*espb.HeuristicValue) map[espb.Heuristic]*espb.HeuristicValue {
	heuristicMap := map[espb.Heuristic]*espb.HeuristicValue{}
	for _, hv := range heuristicValues {
		h, found := heuristicMap[hv.GetHeuristic()]
		if found {
			h.Value += hv.GetValue()
		} else {
			heuristicMap[hv.GetHeuristic()] = proto.Clone(hv).(*espb.HeuristicValue)
		}
	}
	return heuristicMap
}

// MapToList is a convenience function for restoring data into a format useful for data storage,
// e.g. in relevant protos.
func MapToList(heuristicMap map[espb.Heuristic]*espb.HeuristicValue) []*espb.HeuristicValue {
	var heuristicList []*espb.HeuristicValue
	for _, hv := range heuristicMap {
		heuristicList = append(heuristicList, proto.Clone(hv).(*espb.HeuristicValue))
	}
	return heuristicList
}

// InstantGlobalEventSpreadHandler will handle influence propagation calculations for SpreadType
// SPREAD_TYPE_INSTANT_GLOBAL.
//
// TODO(cripplet): Move to separate file.
type InstantGlobalEventSpreadHandler struct{}

// IsPropagated will return True if the input timestamp is at or after the event timstamp, since
// influence propagation is instantaneous for this implementation.
func (h *InstantGlobalEventSpreadHandler) IsPropagated(e *espb.Event, t *timestamp.Timestamp) (bool, error) {
	i, err := ptypes.Timestamp(e.GetTimestamp())
	if err != nil {
		return false, err
	}
	q, err := ptypes.Timestamp(t)
	if err != nil {
		return false, err
	}
	return !q.Before(i), nil
}

// EventSpread implements the EventSpreadHandler.EventSpread function for SpreadType SPREAD_TYPE_INSTANT_GLOBAL.
func (h *InstantGlobalEventSpreadHandler) EventSpread(e *espb.Event, req *espb.GetEventSpreadRequest) ([]*espb.HeuristicValue, error) {
	isPropagated, err := h.IsPropagated(e, req.GetTimestamp())
	if err != nil {
		return nil, err
	}

	var l []*espb.HeuristicValue
	for _, h := range req.GetHeuristics() {
		l = append(l, &espb.HeuristicValue{Heuristic: h})
	}

	if !isPropagated {
		return MapToList(ListToMap(l)), nil
	}

	m := ListToMap(e.GetHeuristics())
	for _, h := range req.GetHeuristics() {
		if v, found := m[h]; found {
			l = append(l, v)
		}
	}
	return MapToList(ListToMap(l)), nil
}
