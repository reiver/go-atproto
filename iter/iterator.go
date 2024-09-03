package iter

type Iterator interface {
	Close() error
	Decode(any) error
	Err() error
	Next() bool
}
