package spatial

import (
	"bytes"
)

type GenericGeometry struct {
	baseGeometry
	typ                    GeometryType
	PointData              *PointData
	LineStringData         *LineStringData
	PolygonData            *PolygonData
	MultiPointData         *MultiPointData
	MultiLineStringData    *MultiLineStringData
	MultiPolygonData       *MultiPolygonData
	GeometryCollectionData *GeometryCollectionData
}

func NewGenericGeometry(srid Srid, typ GeometryType) *GenericGeometry {
	return &GenericGeometry{
		baseGeometry: baseGeometry{srid: srid},
		typ:          typ,
	}
}

func (g *GenericGeometry) Decode(data []byte) error {
	return g.decodeFrom(bytes.NewReader(data), true)
}

func (g *GenericGeometry) Encode() []byte {
	data := newEncodeBuffer()
	g.encodeTo(data, true)
	return data.Bytes()
}

func (g *GenericGeometry) Type() GeometryType {
	return g.typ
}

func (g *GenericGeometry) decodeFrom(data *bytes.Reader, decodeSrid bool) error {
	var err error
	if g.typ, err = g.decodeHeaderFrom(data, decodeSrid, GEOMETRY_TYPE_GENERIC); nil != err {
		return err
	}

	switch g.typ {
	case GEOMETRY_TYPE_POINT:
		gd := &PointData{}
		g.PointData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_LINE_STRING:
		gd := &LineStringData{}
		g.LineStringData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_POLYGON:
		gd := &PolygonData{}
		g.PolygonData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_MULTI_POINT:
		gd := &MultiPointData{}
		g.MultiPointData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_MULTI_LINE_STRING:
		gd := &MultiLineStringData{}
		g.MultiLineStringData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_MULTI_POLYGON:
		gd := &MultiPolygonData{}
		g.MultiPolygonData, err = gd, gd.decodeFrom(data)
	case GEOMETRY_TYPE_COLLECTION:
		gd := &GeometryCollectionData{}
		g.GeometryCollectionData, err = gd, gd.decodeFrom(data, g.srid)
	}
	return err
}

func (g *GenericGeometry) encodeTo(data *bytes.Buffer, encodeSrid bool) {
	g.encodeHeaderTo(data, encodeSrid, g.typ)
	switch g.typ {
	case GEOMETRY_TYPE_POINT:
		g.PointData.encodeTo(data)
	case GEOMETRY_TYPE_LINE_STRING:
		g.LineStringData.encodeTo(data)
	case GEOMETRY_TYPE_POLYGON:
		g.PolygonData.encodeTo(data)
	case GEOMETRY_TYPE_MULTI_POINT:
		g.MultiPointData.encodeTo(data)
	case GEOMETRY_TYPE_MULTI_LINE_STRING:
		g.MultiLineStringData.encodeTo(data)
	case GEOMETRY_TYPE_MULTI_POLYGON:
		g.MultiPolygonData.encodeTo(data)
	case GEOMETRY_TYPE_COLLECTION:
		g.GeometryCollectionData.encodeTo(data)
	}
}
