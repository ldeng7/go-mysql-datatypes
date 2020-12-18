package spatial

import (
	"bytes"
	"encoding/binary"
)

type PointData struct {
	X, Y float64
}

func (pd *PointData) decodeFrom(data *bytes.Reader) error {
	if err := binary.Read(data, byteOrder, &pd.X); nil != err {
		return err
	}
	return binary.Read(data, byteOrder, &pd.Y)
}

func (pd *PointData) encodeTo(data *bytes.Buffer) {
	binary.Write(data, byteOrder, pd.X)
	binary.Write(data, byteOrder, pd.Y)
}

type Point struct {
	baseGeometry
	Data PointData
}

func NewPoint(srid Srid) *Point {
	return &Point{baseGeometry: baseGeometry{srid: srid}}
}

func (p *Point) Decode(data []byte) error {
	return p.decodeFrom(bytes.NewReader(data), true)
}

func (p *Point) Encode() []byte {
	data := newEncodeBuffer()
	p.encodeTo(data, true)
	return data.Bytes()
}

func (p *Point) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := p.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_POINT); nil != err {
		return err
	}
	return p.Data.decodeFrom(data)
}

func (p *Point) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	p.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_POINT)
	p.Data.encodeTo(data)
}
