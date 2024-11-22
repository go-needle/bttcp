package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode encodes the message with it's length
func Encode(message []byte) ([]byte, error) {
	var length = int64(len(message))
	var pkg = new(bytes.Buffer)
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode decodes message
func Decode(reader *bufio.Reader) ([]byte, error) {
	lengthByte, _ := reader.Peek(8)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int64
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	if int64(reader.Buffered()) < length+8 {
		return nil, err
	}

	pack := make([]byte, int(8+length))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}
	return pack[8:], nil
}
