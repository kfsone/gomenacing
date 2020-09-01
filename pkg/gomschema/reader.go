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

func getMessages(header *Header, source io.Reader, unmarshal func([]byte, int) error) error {
	buffer := make([]byte, 4096)

	for idx, size := range header.Sizes {
		if size > uint32(cap(buffer)) {
			buffer = make([]byte, ((size+4095)/4096)*4096)
		}

		// Adjust the *size* of the buffer to only what we expect to read.
		buffer = buffer[:size]
		read, err := io.ReadFull(source, buffer)
		if uint32(read) != size {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			log.Printf("error reading entry %d: %s", idx+1, err)
			return err
		}

		if err = unmarshal(buffer, idx); err != nil {
			return err
		}
	}

	return nil
}

func loadCommodities(header *Header, source io.Reader) ([]Commodity, error) {
	items := make([]Commodity, len(header.Sizes))
	var unmarshal = func(buffer []byte, index int) error {
		return proto.Unmarshal(buffer, &items[index])
	}
	if err := getMessages(header, source, unmarshal); err != nil {
		return nil, err
	}
	return items, nil
}

func loadSystems(header *Header, source io.Reader) ([]System, error) {
	items := make([]System, len(header.Sizes))
	var unmarshal = func(buffer []byte, index int) error {
		return proto.Unmarshal(buffer, &items[index])
	}
	if err := getMessages(header, source, unmarshal); err != nil {
		return nil, err
	}
	return items, nil
}

func loadFacilities(header *Header, source io.Reader) ([]Facility, error) {
	items := make([]Facility, len(header.Sizes))
	var unmarshal = func(buffer []byte, index int) error {
		return proto.Unmarshal(buffer, &items[index])
	}
	if err := getMessages(header, source, unmarshal); err != nil {
		return nil, err
	}
	return items, nil
}

func loadListings(header *Header, source io.Reader) ([]FacilityListing, error) {
	items := make([]FacilityListing, len(header.Sizes))
	var unmarshal = func(buffer []byte, index int) error {
		return proto.Unmarshal(buffer, &items[index])
	}
	if err := getMessages(header, source, unmarshal); err != nil {
		return nil, err
	}
	return items, nil
}

func ReadGOMFile(source io.Reader) (interface{}, error) {
	// First four bytes have to be our magic.
	magic := make([]byte, 4)
	if _, err := io.ReadFull(source, magic); err != nil {
		return nil, err
	}
	if string(magic) != MAGIC {
		return nil, errors.New("unsupported file format")
	}

	// The next 8 bytes should be the size.
	sizeBytes := make([]byte, 8)
	if _, err := io.ReadFull(source, sizeBytes); err != nil {
		return nil, err
	}
	size, err := strconv.ParseUint(string(sizeBytes), 16, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to parse length: %w", err)
	}

	// We should now be able to read the header.
	headerBytes := make([]byte, size)
	if _, err := io.ReadFull(source, headerBytes); err != nil {
		return nil, fmt.Errorf("unable to load header: %w", err)
	}
	header := &Header{}
	if err = proto.Unmarshal(headerBytes, header); err != nil {
		return nil, fmt.Errorf("unable to parse header: %w", err)
	}

	log.Printf("Loaded %s header with %d sizes.", Header_Type_name[int32(header.HeaderType)], len(header.Sizes))

	var list interface{}
	switch header.HeaderType {
	case Header_CCommodity:
		list, err = loadCommodities(header, source)

	case Header_CSystem:
		list, err = loadSystems(header, source)

	case Header_CFacility:
		list, err = loadFacilities(header, source)

	case Header_CListing:
		list, err = loadListings(header, source)

	default:
		return nil, fmt.Errorf("unable to load %s headers", Header_Type_name[int32(header.HeaderType)])
	}

	return list, err
}
