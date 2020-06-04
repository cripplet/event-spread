from typing import Any, List, Tuple, NamedTuple

import dataclasses
import datetime


class Pos(NamedTuple):
  x: float
  y: float


class _Event(object):
  pos: Pos
  timestamp: datetime.datetime
  influence: float
  spread_rate: float


Event = dataclasses.dataclass(_Event)


class _Global(object):
  dim: Pos
  events: List[Event]
  influence: int


Global = dataclasses.dataclass(_Global)


def GetEventInfluence(
    e: Event,
    pos: Pos,
    timestamp: datetime.datetime) -> float:
  d = ((pos.x - e.pos.x) ** 2 + (pos.y - e.pos.y) ** 2) ** .5
  return e.influence if (
      (timestamp - e.timestamp).total_seconds() * e.spread_rate >= d
  ) else 0


def GetInfluence(
    g: Global,
    pos: Pos,
    timestamp: datetime.datetime) -> float:
  return g.influence + sum(
      [GetEventInfluence(e, pos, timestamp) for e in g.events]
  )


def UpdateGlobalEvents(
    g: Global,
    timestamp: datetime.datetime) -> None:
  all_influence = {
      e: GetEventInfluence(
          e,
          Pos(x=(g.dim.x * 2), y=(g.dim.y * 2)),
          timestamp
      ) for e in g.events
  }
  g.events = [e for e in all_influence if all_influence[e]]
  g.influence += sum(all_influence.values())


if __name__ == '__main__':
  early_event = Event(
      pos=Pos(x=10, y=10),
      timestamp=datetime.datetime.fromtimestamp(0) + datetime.timedelta(seconds=10),
      influence=50,
      spread_rate=1,
  )
  late_event = Event(
      pos=Pos(x=10, y=10),
      timestamp=datetime.datetime.fromtimestamp(0) + datetime.timedelta(seconds=20),
      influence=50,
      spread_rate=1,
  )

  g = Global(
    dim=Pos(x=10, y=10),
    events=[early_event, late_event],
    influence=0,
  )

  assert GetEventInfluence(
      early_event, early_event.pos, early_event.timestamp
  ) == early_event.influence

  assert GetEventInfluence(
      early_event, Pos(0, 0), early_event.timestamp) == 0

  assert GetEventInfluence(
      early_event,
      Pos(g.dim.x * 2, g.dim.y * 2),
      early_event.timestamp + datetime.timedelta(
          seconds=(max(g.dim.x, g.dim.y) * 2 ** .5 / early_event.spread_rate)
      )
  ) == early_event.influence
