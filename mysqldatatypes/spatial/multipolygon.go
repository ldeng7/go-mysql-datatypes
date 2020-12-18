package spatial

import (
	"bytes"
	"encoding/binary"
)

type MultiPolygonData []PolygonData

func (mpd *MultiPolygonData) decodeFrom(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mpd = make([]PolygonData, length)
	var p Polygon
	for i := uint32(0); i < length; i++ {
		if err := p.decodeFrom(data, false); nil != err {
			return err
		}
		(*mpd)[i] = p.Data
	}
	return nil
}

func (mpd *MultiPolygonData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*mpd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*baseGeometry).encodeHeaderTo(nil, data, false, GEOMETRY_TYPE_POLYGON)
		(*mpd)[i].encodeTo(data)
	}
}

type MultiPolygon struct {
	baseGeometry
	Data MultiPolygonData
}

func NewMultiPolygon(srid Srid) *MultiPolygon {
	return &MultiPolygon{baseGeometry: baseGeometry{srid: srid}}
}

func (mp *MultiPolygon) Decode(data []byte) error {
	return mp.decodeFrom(bytes.NewReader(data), true)
}

func (mp *MultiPolygon) Encode() []byte {
	data := newEncodeBuffer()
	mp.encodeTo(data, true)
	return data.Bytes()
}

func (mp *MultiPolygon) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := mp.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_MULTI_POLYGON); nil != err {
		return err
	}
	return mp.Data.decodeFrom(data)
}

func (mp *MultiPolygon) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	mp.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_MULTI_POLYGON)
	mp.Data.encodeTo(data)
}
