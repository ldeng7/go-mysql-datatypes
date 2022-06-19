package spatial

import (
	"bytes"
	"encoding/binary"
)

type LineStringData []PointData

func (ld *LineStringData) decodeFrom(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*ld = make([]PointData, length)
	for i := uint32(0); i < length; i++ {
		if err := (*ld)[i].decodeFrom(data); nil != err {
			return err
		}
	}
	return nil
}

func (ld *LineStringData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*ld))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*ld)[i].encodeTo(data)
	}
}

type LineString struct {
	baseGeometry
	Data LineStringData
}

func NewLineString(srid Srid) *LineString {
	return &LineString{baseGeometry: baseGeometry{srid: srid}}
}

func (l *LineString) Decode(data []byte) error {
	return l.decodeFrom(bytes.NewReader(data), true)
}

func (l *LineString) Encode() []byte {
	data := newEncodeBuffer()
	l.encodeTo(data, true)
	return data.Bytes()
}

func (l *LineString) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := l.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_LINE_STRING); nil != err {
		return err
	}
	return l.Data.decodeFrom(data)
}

func (l *LineString) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	l.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_LINE_STRING)
	l.Data.encodeTo(data)
}
