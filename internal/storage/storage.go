package storage

type Storage interface {
	Read() ([]byte, error)
	Write([]byte) error
}
