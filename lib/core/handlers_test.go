package handlers

import (
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	espb "github.com/cripplet/event-spread/lib/proto/event_spread_go_proto"

)

const (
	heuristicEnum = espb.Heuristic_HEURISTIC_MORALITY
	spreadTypeEnum = espb.SpreadType_SPREAD_TYPE_INSTANT_GLOBAL
)

var (
	handlerDispatcher = map[espb.SpreadType]EventSpreadHandler{
		spreadTypeEnum: &InstantGlobalEventSpreadHandler{},
	}
)

func TestListToMapEmptyList(t *testing.T) {
	m := ListToMap(nil)
	if len(m) != 0 {
		t.Errorf("unexpected non-nil map %v generated from empty list", m)
	}
}

func TestMapToListEmptyMap(t *testing.T) {
	l := MapToList(nil)
	if len(l) != 0 {
		t.Errorf("unexpected non-nil list %v generated from empty map", l)
	}
}

func TestListToMapSimple(t *testing.T) {
	e := &espb.HeuristicValue{
		Heuristic: heuristicEnum,
		Value: 1,
	}

	l := []*espb.HeuristicValue{ proto.Clone(e).(*espb.HeuristicValue) }
	m := ListToMap(l)
	if !proto.Equal(m[heuristicEnum], e) {
		t.Errorf("unexpected diff of HeuristicValue map entry, expected %v but got %v", e, m[heuristicEnum])
	}
}

func TestListToMapMerge(t *testing.T) {
	e := &espb.HeuristicValue{
		Heuristic: heuristicEnum,
		Value: 3,
	}

	l := []*espb.HeuristicValue{
		&espb.HeuristicValue{ Heuristic: heuristicEnum, Value: 1, },
		&espb.HeuristicValue{ Heuristic: heuristicEnum, Value: 2, },
	}
	m := ListToMap(l)
	if !proto.Equal(m[heuristicEnum], e) {
		t.Errorf("unexpected diff of HeuristicValue map entry, expected %v but got %v", e, m[heuristicEnum])
	}
}

func TestMapToListSimple(t *testing.T) {
	e := &espb.HeuristicValue{
		Heuristic: heuristicEnum,
		Value: 1,
	}

	m := map[espb.Heuristic]*espb.HeuristicValue{
		heuristicEnum: proto.Clone(e).(*espb.HeuristicValue),
	}
	l := MapToList(m)
	if len(l) != 1 || !proto.Equal(l[0], e) {
		t.Errorf("unexpected length %v of HeuristicValue list entry %v, expected %v", len(l), l, e)
	}
}

func EventSpreadRequestHelper(t *testing.T, initTime time.Time, queryTime time.Time) (*espb.Event, *espb.GetEventSpreadRequest, *espb.HeuristicValue, []*espb.HeuristicValue, error) {
	i, _ := ptypes.TimestampProto(initTime)
	q, _ := ptypes.TimestampProto(queryTime)

	hv := &espb.HeuristicValue{
		Heuristic: heuristicEnum,
		Value: 1,
	}

	e := &espb.Event{
		SpreadType: spreadTypeEnum,
		Heuristics: []*espb.HeuristicValue{ hv },
		Timestamp: i,
	}

	req := &espb.GetEventSpreadRequest{
		Heuristics: []espb.Heuristic{ heuristicEnum },
		Timestamp: q,
	}

	l, err := (&InstantGlobalEventSpreadHandler{}).EventSpread(e, req)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return e, req, hv, l, nil
}

func TestInstantGlobalBefore(t *testing.T) {
	_, _, _, l, err := EventSpreadRequestHelper(
		t,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1999, time.December, 31, 23, 59, 59, 999999999, time.UTC),
	)
	expectedHeuristicValue := &espb.HeuristicValue{
		Heuristic: heuristicEnum,
		Value: 0,
	}

	if err != nil {
		t.Errorf("unexpected error when calculating EventSpread %v, expected nil", err)
	}
	m := ListToMap(l)
	hv, found := m[heuristicEnum]
	if !found {
		t.Errorf("specified input %v was not found in EventSpread return value %v", heuristicEnum, hv)
	}
	if !proto.Equal(expectedHeuristicValue, hv) {
		t.Errorf("unespected diff of HeuristicValue in EventSpread return value, expected %v but got %v", expectedHeuristicValue, hv)
	}
}

func TestInstantGlobalDuring(t *testing.T) {
	_, _, inputHeuristicValue, l, err := EventSpreadRequestHelper(
		t,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	)

	if err != nil {
		t.Errorf("unexpected error when calculating EventSpread %v, expected nil", err)
	}
	m := ListToMap(l)
	hv, found := m[heuristicEnum]
	if !found {
		t.Errorf("specified input %v was not found in EventSpread return value %v", heuristicEnum, hv)
	}
	if !proto.Equal(inputHeuristicValue, hv) {
		t.Errorf("unespected diff of HeuristicValue in EventSpread return value, expected %v but got %v", inputHeuristicValue, hv)
	}
}

func TestInstantGlobalAfter(t *testing.T) {
	_, _, inputHeuristicValue, l, err := EventSpreadRequestHelper(
		t,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, time.January, 1, 0, 0, 0, 1, time.UTC),
	)

	if err != nil {
		t.Errorf("unexpected error when calculating EventSpread %v, expected nil", err)
	}
	m := ListToMap(l)
	hv, found := m[heuristicEnum]
	if !found {
		t.Errorf("specified input %v was not found in EventSpread return value %v", heuristicEnum, hv)
	}
	if !proto.Equal(inputHeuristicValue, hv) {
		t.Errorf("unespected diff of HeuristicValue in EventSpread return value, expected %v but got %v", inputHeuristicValue, hv)
	}
}

func TestEventSpreadAsync(t *testing.T) {
	event, req, _, l, err := EventSpreadRequestHelper(
		t,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, time.January, 1, 0, 0, 0, 1, time.UTC),
	)
	if err != nil {
		t.Errorf("unexpected error when calculating EventSpread %v, expected nil", err)
	}

	ch, err := EventSpread(handlerDispatcher, event, req)
	if err != nil {
		t.Errorf("unexpected error when calculating EventSpread %v, expected nil", err)
	}
	lAsync := []*espb.HeuristicValue{}
	for v := range ch {
		lAsync = append(lAsync, v)
	}

	expected := ListToMap(l)[heuristicEnum]
	actual := ListToMap(lAsync)[heuristicEnum]
	if !proto.Equal(expected, actual) {
		t.Errorf("asynchronous EventSpread did not generate the same values as the synchronous method, expected %v but got %v", expected, actual)
	}
}
