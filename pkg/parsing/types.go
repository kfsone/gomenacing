package parsing

// EntityPacket is a small wrapper for passing opaque entities around, and
// specifically for storage.
type EntityPacket struct {
	ObjectId uint32 // unique identifier for the entry
	Data     []byte // the data of the message
}
