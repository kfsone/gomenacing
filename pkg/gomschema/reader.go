package gomschema

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"google.golang.org/protobuf/proto"
)

const MAGIC = "GOMD"

// Consumer is a callback that takes an unmarshaled proto message and its index in the current read.
type Consumer func(proto.Message, uint) error

// GOMFile is a simple class for reading
type GOMFile struct {
	source io.Reader
	header *Header
	item   proto.Message
}

func (f *GOMFile) Item() *proto.Message {
	return &f.item
}

// getMessageType will identify which type of GOM message the header represents.
func getMessageType(header *Header) proto.Message {
	switch header.HeaderType {
	case Header_CCommodity:
		return &Commodity{}

	case Header_CSystem:
		return &System{}

	case Header_CFacility:
		return &Facility{}

	case Header_CListing:
		return &FacilityListing{}

	default:
		return nil
	}
}

// readAll is a helper to apply a callback to all messages in a gomfile.
func readAll(f *GOMFile, consumer Consumer) error {
	buffer := make([]byte, 256)

	for idx, size := range f.header.Sizes {
		if size > uint32(cap(buffer)) {
			buffer = make([]byte, ((size+4095)/4096)*4096)
		}

		// Adjust the *size* of the buffer to only what we expect to read.
		buffer = buffer[:size]
		read, err := io.ReadFull(f.source, buffer)
		if uint32(read) != size {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			log.Printf("error reading entry %d: %s", idx+1, err)
			return err
		}

		if err = proto.Unmarshal(buffer, f.item); err != nil {
			return err
		}
		if err = consumer(f.item, uint(idx)); err != nil {
			return err
		}
	}

	return nil
}

// readMagic will return an error if the source does not appear to be a genuine GOM stream.
func readMagic(source io.Reader) error {
	// At the beginning of a file should be a 'magic' ident.,
	magic := make([]byte, 4)
	if _, err := io.ReadFull(source, magic); err != nil {
		return err
	}
	if string(magic) != MAGIC {
		return errors.New("unsupported file format")
	}
	return nil
}

// readSizePrefix will consume the in-stream annotation of how large the header is.
func readSizePrefix(source io.Reader) (size uint64, err error) {
	// We prefix the header with its size so we can allocate for it.
	sizeBytes := make([]byte, 8)
	if _, err := io.ReadFull(source, sizeBytes); err != nil {
		return 0, err
	}
	size, err = strconv.ParseUint(string(sizeBytes), 16, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse length: %w", err)
	}
	return size, nil
}

// readHeader will return the deserialized header from the io.Reader.
func readHeader(source io.Reader, size uint64) (header *Header, err error) {
	headerBytes := make([]byte, size)
	if _, err := io.ReadFull(source, headerBytes); err != nil {
		return nil, fmt.Errorf("unable to read header: %w", err)
	}
	header = &Header{}
	if err = proto.Unmarshal(headerBytes, header); err != nil {
		return nil, fmt.Errorf("unable to parse header: %w", err)
	}
	return
}

// OpenGOMFile will consume a .gom file header from an io.Reader and return a GOMFile
// object based on reading the header message in the source.
// See also GOMFile.Load().
func OpenGOMFile(source io.Reader) (*GOMFile, error) {
	// File layout:
	// Byte 0   1   2   3   4   5   6   7   8   9   a   b   c
	//    | G | O | M | D | n | n | n | n | n | n | n | n | n | <proto header> | <messages>

	var err error
	if err = readMagic(source); err != nil {
		return nil, err
	}
	var size uint64
	if size, err = readSizePrefix(source); err != nil {
		return nil, err
	}
	var header *Header
	if header, err = readHeader(source, size); err != nil {
		return nil, err
	}
	var item proto.Message
	if item = getMessageType(header); item == nil {
		return nil, fmt.Errorf("cannot load %s headers", Header_Type_name[int32(header.HeaderType)])
	}
	return &GOMFile{source: source, header: header, item: item}, nil
}

// Close will release resources used by a GOMFile.
func (f *GOMFile) Close() {
	f.source = nil
	f.header = nil
	f.item = nil
}

// Load will read a GOM header and return all the messages identified by it from the current source.
func (f *GOMFile) Load() (list []proto.Message, err error) {
	// Allocate space for the messages up-front.
	list = make([]proto.Message, len(f.header.Sizes))
	// Capture to populate the list.
	consumer := func(in proto.Message, idx uint) error {
		list[idx] = in
		return nil
	}

	err = readAll(f, consumer)

	return
}

// Read will consume messages from a GOMFile and pass them to your consumer.
func (f *GOMFile) Read(consumer Consumer) error {
	return readAll(f, consumer)
}
