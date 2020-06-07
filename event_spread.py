from typing import Any, Dict, List, Tuple, NamedTuple

import collections
import dataclasses
import datetime
import enum


class Heuristic(enum.Enum):
  HEURISTIC_UNDEFINED = 0
  HEURISTIC_MORALITY = 1


class Pos(NamedTuple):
  x: float
  y: float


class _Event(object):
  pos: Pos
  timestamp: datetime.datetime
  influence: Dict[Heuristic, float]
  spread_rate: float


Event = dataclasses.dataclass(_Event)


class _Global(object):
  dim: Pos
  events: List[Event]
  influence: Dict[Heuristic, float]


Global = dataclasses.dataclass(_Global)


def GetEventInfluence(
    e: Event,
    pos: Pos,
    timestamp: datetime.datetime,
    heuristic: Heuristic) -> float:
  d = ((pos.x - e.pos.x) ** 2 + (pos.y - e.pos.y) ** 2) ** .5
  return e.influence[heuristic] if (
      (timestamp - e.timestamp).total_seconds() * e.spread_rate >= d
  ) else 0


def GetInfluence(
    g: Global,
    pos: Pos,
    timestamp: datetime.datetime,
    heuristic: Heuristic) -> float:
  return g.influence[heuristic] + sum(
      [GetEventInfluence(e, pos, timestamp, heuristic) for e in g.events]
  )


if __name__ == '__main__':
  early_event = Event(
      pos=Pos(x=10, y=10),
      timestamp=datetime.datetime.fromtimestamp(0) + datetime.timedelta(seconds=10),
      influence=collections.defaultdict(
          float,
          {Heuristic.HEURISTIC_MORALITY: 50},
      ),
      spread_rate=1,
  )
  late_event = Event(
      pos=Pos(x=10, y=10),
      timestamp=datetime.datetime.fromtimestamp(0) + datetime.timedelta(seconds=20),
      influence=collections.defaultdict(
          float,
          {Heuristic.HEURISTIC_MORALITY: 50},
      ),
      spread_rate=1,
  )

  g = Global(
    dim=Pos(x=10, y=10),
    events=[early_event, late_event],
    influence=collections.defaultdict(float),
  )

  assert GetEventInfluence(
      early_event, early_event.pos, early_event.timestamp, Heuristic.HEURISTIC_MORALITY
  ) == early_event.influence[Heuristic.HEURISTIC_MORALITY]

  assert GetEventInfluence(
      early_event, Pos(0, 0), early_event.timestamp, Heuristic.HEURISTIC_MORALITY) == 0

  assert GetEventInfluence(
      early_event,
      Pos(g.dim.x * 2, g.dim.y * 2),
      early_event.timestamp + datetime.timedelta(
          seconds=(max(g.dim.x, g.dim.y) * 2 ** .5 / early_event.spread_rate)
      ),
      Heuristic.HEURISTIC_MORALITY,
  ) == early_event.influence[Heuristic.HEURISTIC_MORALITY]
