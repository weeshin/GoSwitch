package field

type ISOField interface {
	Pack(val string, length int) ([]byte, error)
	Unpack(data []byte, length int) (val string, readLen int, err error)
}
