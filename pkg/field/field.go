package field

type ISOField interface {
	Pack(val string, length int) ([]byte, error)
	Unpack(data []byte, length int) (val string, readLen int, err error)
}

type BitMap interface {
	Pack(fields map[int]bool) ([]byte, error)
	Unpack(data []byte) (map[int]bool, int, error)
}
