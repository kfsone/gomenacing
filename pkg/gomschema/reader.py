# How to read a .gom file

import logging

from . import gomschema_pb2 as gom

HEADER_CLASS_MAP = {
        'CCommodity':       gom.Commodity,
        'CSystem':          gom.System,
        'CFacility':        gom.Facility,
        'CListing':         gom.FacilityListing,
        }

def read_gom_file(fullpath):
    """
        Attempts to load and deserialize a .GOM file.

        If successful, returns the header and a dictionary of {id: entity}.
    """

    with open(fullpath, "rb") as fh:
        # Verify the first four bytes are 'GOMD'
        if fh.read(4) != b'GOMD':
            raise ValueError(f"{fullpath}: does not appear to be a GOMD file")

        # The next 8 bytes are a hex representation of the length of the header.
        header_len_bytes = fh.read(8).decode("ascii")
        header_len = int(header_len_bytes, 16)

        if header_len <= 0:
            raise ValueError(f"{fullpath}: does not contain valid header data")

        header = gom.Header()
        header.ParseFromString(fh.read(header_len))

        header_type = header.header_type
        try:
            message_name = gom.Header.Type.Name(header_type)
        except ValueError:
            raise ValueError("Header type %d is not recognized by this reader.", header_type)

        try:
            constructor = HEADER_CLASS_MAP[message_name]
        except KeyError:
            raise ValueError(f"{message_name} is not recognized by this reader.")

        logging.debug(f"{fullpath}: header: len {header_len}, sizes {len(header.sizes)}, type {message_name}")

        messages = {}
        for size in header.sizes:
            message = constructor()
            message.ParseFromString(fh.read(size))
            messages[message.id] = message

        logging.debug(f"{fullpath} read {len(header.sizes)} {message_name}s")

        return header, messages

