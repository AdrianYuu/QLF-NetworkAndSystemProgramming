package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType       // Jadi 2 karena 1 + 1 dari yang sebelumnya.
)

type Payload interface {
	io.WriterTo
	io.ReaderFrom
	Bytes() []byte
}

type Binary[]byte

func(m Binary) Bytes() []byte{
	return m
}

func(m Binary) String() string{
	return string(m)
}

func(m Binary) WriteTo(w io.Writer) (int64, error){
	err := binary.Write(w, binary.BigEndian, BinaryType) // Order placing nya dari kecil 
	 											  	     // ke besar dan besar ke kecil
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))

	if err != nil {
		return 0, err
	}

	n, err := w.Write(m)

	return int64(n + 5), err
}

func(m *Binary) ReadFrom(r io.Reader) (int64, error){
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ)

	if err != nil {
		return 0, err
	}

	var size int32

	err = binary.Read(r, binary.BigEndian, &size)

	if err != nil {
		return 0, err
	}

	*m = make([]byte, size)

	s, err := r.Read(*m)
	
	return int64(5 + s), err
}

func Decode(r io.Reader)(Payload, error){
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ)

	if err != nil {
		return nil, err
	}

	payload := new(Binary)

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))

	if err != nil {
		return nil, err
	}

	return payload, nil
}