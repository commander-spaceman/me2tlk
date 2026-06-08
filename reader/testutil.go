package reader

import "encoding/binary"

func BuildTestFile() *File {
	var buf []byte

	writeI32 := func(v int32) {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		buf = append(buf, b...)
	}

	writeI32(int32(TLKMagic))
	writeI32(3)
	writeI32(2)
	writeI32(1)
	writeI32(0)
	writeI32(2)
	writeI32(4)

	writeI32(1)
	writeI32(0)

	writeI32(-66)
	writeI32(1)

	writeI32(-67)
	writeI32(-1)

	buf = append(buf, 0b00011010, 0, 0, 0)

	f, _ := Parse(buf, "test.tlk")
	return f
}
