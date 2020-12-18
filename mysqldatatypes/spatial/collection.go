package spatial

import (
	"bytes"
	"encoding/binary"
)

type GeometryCollectionData []GenericGeometry

func (cd *GeometryCollectionData) decodeFrom(data *bytes.Reader, srid Srid) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*cd = make([]GenericGeometry, length)
	for i := uint32(0); i < length; i++ {
		g := &((*cd)[i])
		g.srid = srid
		if err := g.decodeFrom(data, false); nil != err {
			return err
		}
	}
	return nil
}

func (cd *GeometryCollectionData) encodeTo(data *bytes.Buffer) {
	length := uint32(len(*cd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*cd)[i].encodeTo(data, false)
	}
}

type GeometryCollection struct {
	baseGeometry
	Data GeometryCollectionData
}

func NewGeometryCollection(srid Srid) *GeometryCollection {
	return &GeometryCollection{baseGeometry: baseGeometry{srid: srid}}
}

func (c *GeometryCollection) Decode(data []byte) error {
	return c.decodeFrom(bytes.NewReader(data), true)
}

func (c *GeometryCollection) Encode() []byte {
	data := newEncodeBuffer()
	c.encodeTo(data, true)
	return data.Bytes()
}

func (c *GeometryCollection) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	if _, err := c.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_COLLECTION); nil != err {
		return err
	}
	return c.Data.decodeFrom(data, c.srid)
}

func (c *GeometryCollection) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	c.encodeHeaderTo(data, encodeSrid, GEOMETRY_TYPE_COLLECTION)
	c.Data.encodeTo(data)
}
