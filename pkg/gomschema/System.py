# automatically generated by the FlatBuffers compiler, do not modify

# namespace: gomschema

import flatbuffers

# /// System corresponds to an individual Elite-Dangerous star system, akin to a map.
class System(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAsSystem(cls, buf, offset):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = System()
        x.Init(buf, n + offset)
        return x

    # System
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

# /// System ID is it's upper-cased name hashed via fnv1a.
    # System
    def SystemId(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint32Flags, o + self._tab.Pos)
        return 0

# /// Unique name of the system.
    # System
    def Name(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

# /// Position in the galaxy.
    # System
    def Position(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(8))
        if o != 0:
            x = o + self._tab.Pos
            from .Coordinate import Coordinate
            obj = Coordinate()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

# /// Timestamp of the last update to this entry UTC.
    # System
    def TimestampUtc(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(10))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint64Flags, o + self._tab.Pos)
        return 0

    # System
    def Power(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

# /// Whether anyone lives here.
    # System
    def Populated(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(14))
        if o != 0:
            return bool(self._tab.Get(flatbuffers.number_types.BoolFlags, o + self._tab.Pos))
        return True

# /// Whether a permit is required to enter the systme.
    # System
    def NeedsPermit(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(16))
        if o != 0:
            return bool(self._tab.Get(flatbuffers.number_types.BoolFlags, o + self._tab.Pos))
        return False

# /// Law-Enforcement level of the system.
    # System
    def Security(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(18))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int8Flags, o + self._tab.Pos)
        return 3

# /// What is the government for the system.
    # System
    def GovernmentId(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(20))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int8Flags, o + self._tab.Pos)
        return 5

# /// Which faction is the system allied to.
    # System
    def AllegianceId(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(22))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int8Flags, o + self._tab.Pos)
        return 4

# /// Guess: Elite Dangerous Internal ID
    # System
    def EdAddress(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(24))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint64Flags, o + self._tab.Pos)
        return 0

# /// Facilities in this system.
    # System
    def Facilities(self, j):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(26))
        if o != 0:
            x = self._tab.Vector(o)
            x += flatbuffers.number_types.UOffsetTFlags.py_type(j) * 4
            x = self._tab.Indirect(x)
            from .Facility import Facility
            obj = Facility()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

    # System
    def FacilitiesLength(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(26))
        if o != 0:
            return self._tab.VectorLen(o)
        return 0

def SystemStart(builder): builder.StartObject(12)
def SystemAddSystemId(builder, systemId): builder.PrependUint32Slot(0, systemId, 0)
def SystemAddName(builder, name): builder.PrependUOffsetTRelativeSlot(1, flatbuffers.number_types.UOffsetTFlags.py_type(name), 0)
def SystemAddPosition(builder, position): builder.PrependStructSlot(2, flatbuffers.number_types.UOffsetTFlags.py_type(position), 0)
def SystemAddTimestampUtc(builder, timestampUtc): builder.PrependUint64Slot(3, timestampUtc, 0)
def SystemAddPower(builder, power): builder.PrependUOffsetTRelativeSlot(4, flatbuffers.number_types.UOffsetTFlags.py_type(power), 0)
def SystemAddPopulated(builder, populated): builder.PrependBoolSlot(5, populated, 1)
def SystemAddNeedsPermit(builder, needsPermit): builder.PrependBoolSlot(6, needsPermit, 0)
def SystemAddSecurity(builder, security): builder.PrependInt8Slot(7, security, 3)
def SystemAddGovernmentId(builder, governmentId): builder.PrependInt8Slot(8, governmentId, 5)
def SystemAddAllegianceId(builder, allegianceId): builder.PrependInt8Slot(9, allegianceId, 4)
def SystemAddEdAddress(builder, edAddress): builder.PrependUint64Slot(10, edAddress, 0)
def SystemAddFacilities(builder, facilities): builder.PrependUOffsetTRelativeSlot(11, flatbuffers.number_types.UOffsetTFlags.py_type(facilities), 0)
def SystemStartFacilitiesVector(builder, numElems): return builder.StartVector(4, numElems, 4)
def SystemEnd(builder): return builder.EndObject()
