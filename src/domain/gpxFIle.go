package domain

import "encoding/xml"

type Trkpt struct {
	Lat string  `xml:"lat,attr"`
	Lon string  `xml:"lon,attr"`
	Ele float64 `xml:"ele"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trk struct {
	Name   string   `xml:"name"`
	Trkseg []Trkseg `xml:"trkseg"`
}

type GpxStruct struct {
	FileName string
	XMLName  xml.Name `xml:"gpx"`
	Trk      Trk      `xml:"trk"`
}
