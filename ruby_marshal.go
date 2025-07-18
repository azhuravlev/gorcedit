package main

import (
	"bytes"
	"fmt"
	"io"
)

// https://jakegoulding.com/blog/2013/01/15/a-little-dip-into-rubys-marshal-format/
// ruby marshal is a bit complicated, but as we work only with encoded yml files,
// so we can only implement string (un)marshal. If it will be required,
// I think, that good idea will be to use https://github.com/samcday/rmarsh

// UnmarshalRubyString parses ruby marshalled string
func UnmarshalRubyString(data []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0x04 || data[1] != 0x08 {
		// data wasn't encoded or Marshal version was changed
		return nil, fmt.Errorf("unsupported Marshal version")
	}

	dataReader := bytes.NewReader(data[2:])

	tag, err := dataReader.ReadByte()
	if err != nil {
		return nil, err
	}

	// 0x22 == '"' we support only strings
	if tag != 0x22 {
		return nil, fmt.Errorf("unsupported Ruby Marshal type: %q", tag)
	}

	length, err := DecodePositiveInt(dataReader)
	if err != nil {
		return nil, err
	}
	strBuf := make([]byte, length)
	_, err = io.ReadFull(dataReader, strBuf)

	return strBuf, err
}

// MarshalRubyString encode data as ruby marshalled string
func MarshalRubyString(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	// Marshal version
	_, err := buf.Write([]byte{0x04, 0x08})
	if err != nil {
		return nil, err
	}
	// String tag
	err = buf.WriteByte('"')
	if err != nil {
		return nil, err
	}
	// Data length
	length := EncodePositiveInt(len(data))
	_, err = buf.Write(length)
	if err != nil {
		return nil, err
	}
	// Data
	_, err = buf.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodePositiveInt reads and decodes a Ruby-style positive integer from the reader.
// Used for decoding lengths (e.g. string size) from Ruby Marshal format.
func DecodePositiveInt(r *bytes.Reader) (int, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	n := int(int8(b))

	switch {
	case n == 0:
		return 0, nil
	case n >= 4:
		return n - 5, nil
	case n >= 1 && n <= 3:
		var result int
		for i := 0; i < n; i++ {
			b, err := r.ReadByte()
			if err != nil {
				return 0, err
			}
			result |= int(b) << (8 * i)
		}
		return result, nil
	default:
		return 0, fmt.Errorf("unsupported or negative integer encoding: %d", n)
	}
}

// EncodePositiveInt encodes a positive integer into Ruby Marshal compact format.
func EncodePositiveInt(n int) []byte {
	switch {
	case n == 0:
		return []byte{0x00}
	case n > 0 && n < 123:
		return []byte{byte(n + 5)}
	default:
		var payload []byte
		val := n
		for val > 0 {
			payload = append(payload, byte(val&0xff))
			val >>= 8
		}
		return append([]byte{byte(len(payload))}, payload...)
	}
}
