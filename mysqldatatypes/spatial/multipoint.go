package spatial

import (
	"bytes"
	"encoding/binary"
)

type MultiPointData []PointData

func (mpd *MultiPointData) decodeFrom(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mpd = make([]PointData, length)
	var p Point
	for i := uint32(0); i < length; i++ {
		if err := p.decodeFrom(data, false); nil != err {
			return err
		}
		(*mpd)[i] = p.Data
	}
	return nil
}

func (mpd *MultiPointData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*mpd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*baseGeometry).encodeHeaderTo(nil, data, false, GEOMETRY_TYPE_POINT)
		(*mpd)[i].encodeTo(data)
	}
}

type MultiPoint struct {
	baseGeometry
	Data MultiPointData
}

func NewMultiPoint(srid Srid) *MultiPoint {
	return &MultiPoint{baseGeometry: baseGeometry{srid: srid}}
}

func (mp *MultiPoint) Decode(data []byte) error {
	return mp.decodeFrom(bytes.NewReader(data), true)
}

func (mp *MultiPoint) Encode() []byte {
	data := newEncodeBuffer()
	mp.encodeTo(data, true)
	return data.Bytes()
}

func (mp *MultiPoint) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := mp.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_MULTI_POINT); nil != err {
		return err
	}
	return mp.Data.decodeFrom(data)
}

func (mp *MultiPoint) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	mp.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_MULTI_POINT)
	mp.Data.encodeTo(data)
}
