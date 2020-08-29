# automatically generated by the FlatBuffers compiler, do not modify

# namespace: ettuschema

import flatbuffers
from flatbuffers.compat import import_numpy
np = import_numpy()

# Encapsulation of all the data.
class Galaxy(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAsGalaxy(cls, buf, offset):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = Galaxy()
        x.Init(buf, n + offset)
        return x

    @classmethod
    def GalaxyBufferHasIdentifier(cls, buf, offset, size_prefixed=False):
        return flatbuffers.util.BufferHasIdentifier(buf, offset, b"\x67\x6F\x6D\x64", size_prefixed=size_prefixed)

    # Galaxy
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

    # Semantically-versioned schema id.
    # Galaxy
    def SchemaVersion(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # Human-friendly description of what is enclosed, e.g "import from source X" or
    # "complete local database". Entirely descriptive.
    # Galaxy
    def Description(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # Human-friendly attributition, if relevant.
    # Galaxy
    def Attribution(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(8))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # UTC Unix time of generation.
    # Galaxy
    def TimestampUtc(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(10))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint64Flags, o + self._tab.Pos)
        return 0

    # Items recognized by this data.
    # Galaxy
    def Commodities(self, j):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        if o != 0:
            x = self._tab.Vector(o)
            x += flatbuffers.number_types.UOffsetTFlags.py_type(j) * 4
            x = self._tab.Indirect(x)
            from ettuschema.Commodity import Commodity
            obj = Commodity()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

    # Galaxy
    def CommoditiesLength(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        if o != 0:
            return self._tab.VectorLen(o)
        return 0

    # Galaxy
    def CommoditiesIsNone(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        return o == 0

    # Systems recognized by this data (presence of facilities optional).
    # Galaxy
    def Systems(self, j):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(14))
        if o != 0:
            x = self._tab.Vector(o)
            x += flatbuffers.number_types.UOffsetTFlags.py_type(j) * 4
            x = self._tab.Indirect(x)
            from ettuschema.System import System
            obj = System()
            obj.Init(self._tab.Bytes, x)
            return obj
        return None

    # Galaxy
    def SystemsLength(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(14))
        if o != 0:
            return self._tab.VectorLen(o)
        return 0

    # Galaxy
    def SystemsIsNone(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(14))
        return o == 0

    # Fields reserved for any user-specific notes.
    # Galaxy
    def UserData(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(16))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # 3rd-party application values that prefer .ini format.
    # Galaxy
    def IniData(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(18))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # 3rd-party application values that prefer .json format.
    # Galaxy
    def JsonData(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(20))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # 3rd-party application values that prefer .yaml because they hate life.
    # Galaxy
    def YamlData(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(22))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

def GalaxyStart(builder): builder.StartObject(10)
def GalaxyAddSchemaVersion(builder, schemaVersion): builder.PrependUOffsetTRelativeSlot(0, flatbuffers.number_types.UOffsetTFlags.py_type(schemaVersion), 0)
def GalaxyAddDescription(builder, description): builder.PrependUOffsetTRelativeSlot(1, flatbuffers.number_types.UOffsetTFlags.py_type(description), 0)
def GalaxyAddAttribution(builder, attribution): builder.PrependUOffsetTRelativeSlot(2, flatbuffers.number_types.UOffsetTFlags.py_type(attribution), 0)
def GalaxyAddTimestampUtc(builder, timestampUtc): builder.PrependUint64Slot(3, timestampUtc, 0)
def GalaxyAddCommodities(builder, commodities): builder.PrependUOffsetTRelativeSlot(4, flatbuffers.number_types.UOffsetTFlags.py_type(commodities), 0)
def GalaxyStartCommoditiesVector(builder, numElems): return builder.StartVector(4, numElems, 4)
def GalaxyAddSystems(builder, systems): builder.PrependUOffsetTRelativeSlot(5, flatbuffers.number_types.UOffsetTFlags.py_type(systems), 0)
def GalaxyStartSystemsVector(builder, numElems): return builder.StartVector(4, numElems, 4)
def GalaxyAddUserData(builder, userData): builder.PrependUOffsetTRelativeSlot(6, flatbuffers.number_types.UOffsetTFlags.py_type(userData), 0)
def GalaxyAddIniData(builder, iniData): builder.PrependUOffsetTRelativeSlot(7, flatbuffers.number_types.UOffsetTFlags.py_type(iniData), 0)
def GalaxyAddJsonData(builder, jsonData): builder.PrependUOffsetTRelativeSlot(8, flatbuffers.number_types.UOffsetTFlags.py_type(jsonData), 0)
def GalaxyAddYamlData(builder, yamlData): builder.PrependUOffsetTRelativeSlot(9, flatbuffers.number_types.UOffsetTFlags.py_type(yamlData), 0)
def GalaxyEnd(builder): return builder.EndObject()
