# automatically generated by the FlatBuffers compiler, do not modify

# namespace: ettuschema

import flatbuffers
from flatbuffers.compat import import_numpy
np = import_numpy()

# Trade is a discrete entry for a Commodity that is or can be traded at with
# a number of units and a value.
class Trade(object):
    __slots__ = ['_tab']

    # Trade
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

    # Which commodity this descrbes.
    # Trade
    def CommodityId(self): return self._tab.Get(flatbuffers.number_types.Uint64Flags, self._tab.Pos + flatbuffers.number_types.UOffsetTFlags.py_type(0))
    # How many units
    # Trade
    def Units(self): return self._tab.Get(flatbuffers.number_types.Uint32Flags, self._tab.Pos + flatbuffers.number_types.UOffsetTFlags.py_type(8))
    # How many credits
    # Trade
    def Credits(self): return self._tab.Get(flatbuffers.number_types.Uint16Flags, self._tab.Pos + flatbuffers.number_types.UOffsetTFlags.py_type(12))
    # Unix timestamp UTC.
    # Trade
    def TimestampUtc(self): return self._tab.Get(flatbuffers.number_types.Uint64Flags, self._tab.Pos + flatbuffers.number_types.UOffsetTFlags.py_type(16))

def CreateTrade(builder, commodityId, units, credits, timestampUtc):
    builder.Prep(8, 24)
    builder.PrependUint64(timestampUtc)
    builder.Pad(2)
    builder.PrependUint16(credits)
    builder.PrependUint32(units)
    builder.PrependUint64(commodityId)
    return builder.Offset()
