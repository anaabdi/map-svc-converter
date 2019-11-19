package kml

import "encoding/xml"

type Kml struct {
	XMLName  xml.Name `xml:"kml"`
	Document Document `xml:"Document"`
}

type Document struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Folder      Folder `xml:"Folder"`
}

type Folder struct {
	Name      string      `xml:"name"`
	Placemark []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string  `xml:"name"`
	Description string  `xml:"description"`
	Polygon     Polygon `xml:"Polygon"`
	Point       Point   `xml:"Point"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type Polygon struct {
	OuterBoundaryIs OuterBoundaryIs `xml:"outerBoundaryIs"`
}

type OuterBoundaryIs struct {
	LinearRing LinearRing `xml:"LinearRing"`
}

type LinearRing struct {
	Coordinates string `xml:"coordinates"`
}

func (p Placemark) IsPolygon() bool {
	return p.Polygon.OuterBoundaryIs.LinearRing.Coordinates != "" &&
		p.Point.Coordinates == ""
}
