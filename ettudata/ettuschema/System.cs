// <auto-generated>
//  automatically generated by the FlatBuffers compiler, do not modify
// </auto-generated>

namespace ettuschema
{

using global::System;
using global::System.Collections.Generic;
using global::FlatBuffers;

/// System corresponds to an individual Elite-Dangerous star system, akin to a map.
public struct System : IFlatbufferObject
{
  private Table __p;
  public ByteBuffer ByteBuffer { get { return __p.bb; } }
  public static void ValidateVersion() { FlatBufferConstants.FLATBUFFERS_1_12_0(); }
  public static System GetRootAsSystem(ByteBuffer _bb) { return GetRootAsSystem(_bb, new System()); }
  public static System GetRootAsSystem(ByteBuffer _bb, System obj) { return (obj.__assign(_bb.GetInt(_bb.Position) + _bb.Position, _bb)); }
  public void __init(int _i, ByteBuffer _bb) { __p = new Table(_i, _bb); }
  public System __assign(int _i, ByteBuffer _bb) { __init(_i, _bb); return this; }

  /// System ID is it's upper-cased name hashed via fnv1a.
  public uint SystemId { get { int o = __p.__offset(4); return o != 0 ? __p.bb.GetUint(o + __p.bb_pos) : (uint)0; } }
  public bool MutateSystemId(uint system_id) { int o = __p.__offset(4); if (o != 0) { __p.bb.PutUint(o + __p.bb_pos, system_id); return true; } else { return false; } }
  /// Unique name of the system.
  public string Name { get { int o = __p.__offset(6); return o != 0 ? __p.__string(o + __p.bb_pos) : null; } }
#if ENABLE_SPAN_T
  public Span<byte> GetNameBytes() { return __p.__vector_as_span<byte>(6, 1); }
#else
  public ArraySegment<byte>? GetNameBytes() { return __p.__vector_as_arraysegment(6); }
#endif
  public byte[] GetNameArray() { return __p.__vector_as_array<byte>(6); }
  /// Position in the galaxy.
  public ettuschema.Coordinate? Position { get { int o = __p.__offset(8); return o != 0 ? (ettuschema.Coordinate?)(new ettuschema.Coordinate()).__assign(o + __p.bb_pos, __p.bb) : null; } }
  /// Timestamp of the last update to this entry UTC.
  public ulong TimestampUtc { get { int o = __p.__offset(10); return o != 0 ? __p.bb.GetUlong(o + __p.bb_pos) : (ulong)0; } }
  public bool MutateTimestampUtc(ulong timestamp_utc) { int o = __p.__offset(10); if (o != 0) { __p.bb.PutUlong(o + __p.bb_pos, timestamp_utc); return true; } else { return false; } }
  public string Power { get { int o = __p.__offset(12); return o != 0 ? __p.__string(o + __p.bb_pos) : null; } }
#if ENABLE_SPAN_T
  public Span<byte> GetPowerBytes() { return __p.__vector_as_span<byte>(12, 1); }
#else
  public ArraySegment<byte>? GetPowerBytes() { return __p.__vector_as_arraysegment(12); }
#endif
  public byte[] GetPowerArray() { return __p.__vector_as_array<byte>(12); }
  /// Whether anyone lives here.
  public bool Populated { get { int o = __p.__offset(14); return o != 0 ? 0!=__p.bb.Get(o + __p.bb_pos) : (bool)true; } }
  public bool MutatePopulated(bool populated) { int o = __p.__offset(14); if (o != 0) { __p.bb.Put(o + __p.bb_pos, (byte)(populated ? 1 : 0)); return true; } else { return false; } }
  /// Whether a permit is required to enter the systme.
  public bool NeedsPermit { get { int o = __p.__offset(16); return o != 0 ? 0!=__p.bb.Get(o + __p.bb_pos) : (bool)false; } }
  public bool MutateNeedsPermit(bool needs_permit) { int o = __p.__offset(16); if (o != 0) { __p.bb.Put(o + __p.bb_pos, (byte)(needs_permit ? 1 : 0)); return true; } else { return false; } }
  /// Law-Enforcement level of the system.
  public ettuschema.SecurityLevel Security { get { int o = __p.__offset(18); return o != 0 ? (ettuschema.SecurityLevel)__p.bb.GetSbyte(o + __p.bb_pos) : ettuschema.SecurityLevel.Medium; } }
  public bool MutateSecurity(ettuschema.SecurityLevel security) { int o = __p.__offset(18); if (o != 0) { __p.bb.PutSbyte(o + __p.bb_pos, (sbyte)security); return true; } else { return false; } }
  /// What is the government for the system.
  public ettuschema.Government GovernmentId { get { int o = __p.__offset(20); return o != 0 ? (ettuschema.Government)__p.bb.GetSbyte(o + __p.bb_pos) : ettuschema.Government.Corporate; } }
  public bool MutateGovernmentId(ettuschema.Government government_id) { int o = __p.__offset(20); if (o != 0) { __p.bb.PutSbyte(o + __p.bb_pos, (sbyte)government_id); return true; } else { return false; } }
  /// Which faction is the system allied to.
  public ettuschema.Allegiance AllegianceId { get { int o = __p.__offset(22); return o != 0 ? (ettuschema.Allegiance)__p.bb.GetSbyte(o + __p.bb_pos) : ettuschema.Allegiance.Independent; } }
  public bool MutateAllegianceId(ettuschema.Allegiance allegiance_id) { int o = __p.__offset(22); if (o != 0) { __p.bb.PutSbyte(o + __p.bb_pos, (sbyte)allegiance_id); return true; } else { return false; } }
  /// Guess: Elite Dangerous Internal ID
  public ulong EdAddress { get { int o = __p.__offset(24); return o != 0 ? __p.bb.GetUlong(o + __p.bb_pos) : (ulong)0; } }
  public bool MutateEdAddress(ulong ed_address) { int o = __p.__offset(24); if (o != 0) { __p.bb.PutUlong(o + __p.bb_pos, ed_address); return true; } else { return false; } }
  /// Facilities in this system.
  public ettuschema.Facility? Facilities(int j) { int o = __p.__offset(26); return o != 0 ? (ettuschema.Facility?)(new ettuschema.Facility()).__assign(__p.__indirect(__p.__vector(o) + j * 4), __p.bb) : null; }
  public int FacilitiesLength { get { int o = __p.__offset(26); return o != 0 ? __p.__vector_len(o) : 0; } }
  public ettuschema.Facility? FacilitiesByKey(uint key) { int o = __p.__offset(26); return o != 0 ? ettuschema.Facility.__lookup_by_key(__p.__vector(o), key, __p.bb) : null; }

  public static void StartSystem(FlatBufferBuilder builder) { builder.StartTable(12); }
  public static void AddSystemId(FlatBufferBuilder builder, uint systemId) { builder.AddUint(0, systemId, 0); }
  public static void AddName(FlatBufferBuilder builder, StringOffset nameOffset) { builder.AddOffset(1, nameOffset.Value, 0); }
  public static void AddPosition(FlatBufferBuilder builder, Offset<ettuschema.Coordinate> positionOffset) { builder.AddStruct(2, positionOffset.Value, 0); }
  public static void AddTimestampUtc(FlatBufferBuilder builder, ulong timestampUtc) { builder.AddUlong(3, timestampUtc, 0); }
  public static void AddPower(FlatBufferBuilder builder, StringOffset powerOffset) { builder.AddOffset(4, powerOffset.Value, 0); }
  public static void AddPopulated(FlatBufferBuilder builder, bool populated) { builder.AddBool(5, populated, true); }
  public static void AddNeedsPermit(FlatBufferBuilder builder, bool needsPermit) { builder.AddBool(6, needsPermit, false); }
  public static void AddSecurity(FlatBufferBuilder builder, ettuschema.SecurityLevel security) { builder.AddSbyte(7, (sbyte)security, 3); }
  public static void AddGovernmentId(FlatBufferBuilder builder, ettuschema.Government governmentId) { builder.AddSbyte(8, (sbyte)governmentId, 5); }
  public static void AddAllegianceId(FlatBufferBuilder builder, ettuschema.Allegiance allegianceId) { builder.AddSbyte(9, (sbyte)allegianceId, 4); }
  public static void AddEdAddress(FlatBufferBuilder builder, ulong edAddress) { builder.AddUlong(10, edAddress, 0); }
  public static void AddFacilities(FlatBufferBuilder builder, VectorOffset facilitiesOffset) { builder.AddOffset(11, facilitiesOffset.Value, 0); }
  public static VectorOffset CreateFacilitiesVector(FlatBufferBuilder builder, Offset<ettuschema.Facility>[] data) { builder.StartVector(4, data.Length, 4); for (int i = data.Length - 1; i >= 0; i--) builder.AddOffset(data[i].Value); return builder.EndVector(); }
  public static VectorOffset CreateFacilitiesVectorBlock(FlatBufferBuilder builder, Offset<ettuschema.Facility>[] data) { builder.StartVector(4, data.Length, 4); builder.Add(data); return builder.EndVector(); }
  public static void StartFacilitiesVector(FlatBufferBuilder builder, int numElems) { builder.StartVector(4, numElems, 4); }
  public static Offset<ettuschema.System> EndSystem(FlatBufferBuilder builder) {
    int o = builder.EndTable();
    return new Offset<ettuschema.System>(o);
  }

  public static VectorOffset CreateSortedVectorOfSystem(FlatBufferBuilder builder, Offset<System>[] offsets) {
    Array.Sort(offsets, (Offset<System> o1, Offset<System> o2) => builder.DataBuffer.GetUint(Table.__offset(4, o1.Value, builder.DataBuffer)).CompareTo(builder.DataBuffer.GetUint(Table.__offset(4, o2.Value, builder.DataBuffer))));
    return builder.CreateVectorOfTables(offsets);
  }

  public static System? __lookup_by_key(int vectorLocation, uint key, ByteBuffer bb) {
    int span = bb.GetInt(vectorLocation - 4);
    int start = 0;
    while (span != 0) {
      int middle = span / 2;
      int tableOffset = Table.__indirect(vectorLocation + 4 * (start + middle), bb);
      int comp = bb.GetUint(Table.__offset(4, bb.Length - tableOffset, bb)).CompareTo(key);
      if (comp > 0) {
        span = middle;
      } else if (comp < 0) {
        middle++;
        start += middle;
        span -= middle;
      } else {
        return new System().__assign(tableOffset, bb);
      }
    }
    return null;
  }
};


}
