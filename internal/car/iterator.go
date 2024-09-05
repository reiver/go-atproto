package car

import (
	"bytes"
	"io"

	ipld_car "github.com/ipld/go-car"
	"github.com/reiver/go-erorr"

	"github.com/reiver/go-atproto/iter"
)

// Iterator turns a *car.CarReader from packag "github.com/ipld/go-car" into an iterator.
//
// CAR = "Content Addressable aRchives"
type Iterator struct {
	CarReader *ipld_car.CarReader
	closed bool
	err error
	rawData []byte
}

var _ iter.Iterator = &Iterator{}

func NewIteratorFromBytes(value []byte) (*Iterator, error) {
	if nil == value {
		return nil, errNilBytes
	}

	var carReader *ipld_car.CarReader
	{
		var reader io.Reader = bytes.NewReader(value)
		if nil == reader {
			return nil, errNilReader
		}

		var err error

		carReader, err = ipld_car.NewCarReader(reader)
		if nil != err {
			return nil, erorr.Errorf("atproto: problem creating CAR (Content Addressable aRchives) reader: %w", err)
		}
	}

	return NewIteratorFromCarReader(carReader)
}

func NewIteratorFromCarReader(carReader *ipld_car.CarReader) (*Iterator, error) {
	if nil == carReader {
		return nil, errNilCarReader
	}

	var iterator Iterator
	iterator.CarReader = carReader

	return &iterator, nil
}

func (receiver *Iterator) Close() error {
	if nil == receiver {
		return errNilReceiver
	}

	if receiver.closed {
		return nil
	}

	receiver.CarReader = nil
	receiver.closed = true
	return nil
}

func (receiver *Iterator) Decode(dst any) error {
	if nil == receiver {
		return errNilReceiver
	}

	switch casted := dst.(type) {
	case *[]byte:
		*casted = receiver.rawData
	case *string:
		*casted = string(receiver.rawData)
	default:
		return erorr.Errorf("atproto: cannot decode CAR (Content Addressable aRchives) raw-data into something of type %T", dst)
	}

	return nil
}

func (receiver *Iterator) Err() error {
	if nil == receiver {
		return errNilReceiver
	}

	return receiver.err
}

func (receiver *Iterator) Next() bool {
	if nil == receiver {
		return false
	}

	if receiver.closed {
		return false
	}
	if nil != receiver.err {
		return false
	}

	var carReader *ipld_car.CarReader = receiver.CarReader
	if nil == carReader {
		receiver.err = errNilCarReader
		return false
	}

	block, err := carReader.Next()
	if io.EOF == err || erorr.Is(err, io.EOF) {
		return false
	}
	if nil != err {
		receiver.err = err
		return false
	}
	if nil == block {
		receiver.err = errNilBlock
		return false
	}

	var rawData []byte = block.RawData()
	if nil == rawData {
		receiver.err = errNilBlockRawData
		return false
	}

	receiver.rawData = rawData
	return true
}
