# automatically generated by the FlatBuffers compiler, do not modify

# namespace: gomschema

import flatbuffers

class FacilityListing(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAsFacilityListing(cls, buf, offset):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = FacilityListing()
        x.Init(buf, n + offset)
        return x

    # FacilityListing
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

# /// Commodities this facility sells.
    # FacilityListing
    def Supply(self, j):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            x = self._tab.Vector(o)
            x += flatbuffers.number_types.UOffsetTFlags.py_type(j) * 24
            from .Trade import Trade
            obj = Trade()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

    # FacilityListing
    def SupplyLength(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            return self._tab.VectorLen(o)
        return 0

# /// Commodities this facility buys.
    # FacilityListing
    def Demand(self, j):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            x = self._tab.Vector(o)
            x += flatbuffers.number_types.UOffsetTFlags.py_type(j) * 24
            from .Trade import Trade
            obj = Trade()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

    # FacilityListing
    def DemandLength(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            return self._tab.VectorLen(o)
        return 0

def FacilityListingStart(builder): builder.StartObject(2)
def FacilityListingAddSupply(builder, supply): builder.PrependUOffsetTRelativeSlot(0, flatbuffers.number_types.UOffsetTFlags.py_type(supply), 0)
def FacilityListingStartSupplyVector(builder, numElems): return builder.StartVector(24, numElems, 8)
def FacilityListingAddDemand(builder, demand): builder.PrependUOffsetTRelativeSlot(1, flatbuffers.number_types.UOffsetTFlags.py_type(demand), 0)
def FacilityListingStartDemandVector(builder, numElems): return builder.StartVector(24, numElems, 8)
def FacilityListingEnd(builder): return builder.EndObject()
