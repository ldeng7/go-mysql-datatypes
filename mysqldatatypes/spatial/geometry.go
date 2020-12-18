package spatial

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type GeometryType uint32
type Srid uint32

const (
	GEOMETRY_TYPE_GENERIC GeometryType = iota
	GEOMETRY_TYPE_POINT
	GEOMETRY_TYPE_LINE_STRING
	GEOMETRY_TYPE_POLYGON
	GEOMETRY_TYPE_MULTI_POINT
	GEOMETRY_TYPE_MULTI_LINE_STRING
	GEOMETRY_TYPE_MULTI_POLYGON
	GEOMETRY_TYPE_COLLECTION
)

var (
	byteOrder                   = binary.LittleEndian
	geometryInstantiableMinType = GEOMETRY_TYPE_POINT
	geometryInstantiableMaxType = GEOMETRY_TYPE_COLLECTION
)

type baseGeometry struct {
	srid Srid
}

func (g baseGeometry) Srid() Srid {
	return g.srid
}

func (g *baseGeometry) decodeHeaderFrom(data *bytes.Reader,
	decodeSrid bool, expectedType GeometryType) (GeometryType, error) {
	if decodeSrid {
		if err := binary.Read(data, byteOrder, &g.srid); nil != err {
			return 0, err
		}
	}

	if _, err := data.ReadByte(); nil != err {
		return 0, err
	}

	var typ GeometryType
	if err := binary.Read(data, byteOrder, &typ); nil != err {
		return 0, err
	} else if (typ < geometryInstantiableMinType && typ > geometryInstantiableMaxType) ||
		(GEOMETRY_TYPE_GENERIC != expectedType && typ != expectedType) {
		return 0, errors.New("unexpected geometry type")
	}
	return typ, nil
}

func (g *baseGeometry) encodeHeaderTo(data *bytes.Buffer, encodeSrid bool, typ GeometryType) {
	if encodeSrid {
		binary.Write(data, byteOrder, g.Srid())
	}
	data.WriteByte(0x01)
	binary.Write(data, byteOrder, typ)
}

func newEncodeBuffer() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 25))
}
