package spatial

import (
	"bytes"
	"encoding/binary"
)

type PolygonData []LineStringData

func (pd *PolygonData) decodeFrom(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*pd = make([]LineStringData, length)
	for i := uint32(0); i < length; i++ {
		if err := (*pd)[i].decodeFrom(data); nil != err {
			return err
		}
	}
	return nil
}

func (pd *PolygonData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*pd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*pd)[i].encodeTo(data)
	}
}

type Polygon struct {
	baseGeometry
	Data PolygonData
}

func NewPolygon(srid Srid) *Polygon {
	return &Polygon{baseGeometry: baseGeometry{srid: srid}}
}

func (p *Polygon) Decode(data []byte) error {
	return p.decodeFrom(bytes.NewReader(data), true)
}

func (p *Polygon) Encode() []byte {
	data := newEncodeBuffer()
	p.encodeTo(data, true)
	return data.Bytes()
}

func (p *Polygon) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := p.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_POLYGON); nil != err {
		return err
	}
	return p.Data.decodeFrom(data)
}
func (p *Polygon) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	p.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_POLYGON)
	p.Data.encodeTo(data)
}
