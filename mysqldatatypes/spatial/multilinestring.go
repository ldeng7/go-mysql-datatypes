package spatial

import (
	"bytes"
	"encoding/binary"
)

type MultiLineStringData []LineStringData

func (mld *MultiLineStringData) decodeFrom(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mld = make([]LineStringData, length)
	var l LineString
	for i := uint32(0); i < length; i++ {
		if err := l.decodeFrom(data, false); nil != err {
			return err
		}
		(*mld)[i] = l.Data
	}
	return nil
}

func (mld *MultiLineStringData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*mld))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*baseGeometry).encodeHeaderTo(nil, data, false, GEOMETRY_TYPE_LINE_STRING)
		(*mld)[i].encodeTo(data)
	}
}

type MultiLineString struct {
	baseGeometry
	Data MultiLineStringData
}

func NewMultiLineString(srid Srid) *MultiLineString {
	return &MultiLineString{baseGeometry: baseGeometry{srid: srid}}
}

func (ml *MultiLineString) Decode(data []byte) error {
	return ml.decodeFrom(bytes.NewReader(data), true)
}

func (ml *MultiLineString) Encode() []byte {
	data := newEncodeBuffer()
	ml.encodeTo(data, true)
	return data.Bytes()
}

func (ml *MultiLineString) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := ml.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_MULTI_LINE_STRING); nil != err {
		return err
	}
	return ml.Data.decodeFrom(data)
}

func (ml *MultiLineString) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	ml.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_MULTI_LINE_STRING)
	ml.Data.encodeTo(data)
}
