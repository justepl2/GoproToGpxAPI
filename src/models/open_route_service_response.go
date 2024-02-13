package models

import "encoding/xml"

// Structs for Unmarshal OpenRouteSevice GPX
type Extensions struct {
	Duration float64 `xml:"duration"`
	Distance float64 `xml:"distance"`
	Type     int     `xml:"type"`
	Step     int     `xml:"step"`
}

type Rtept struct {
	Lat        string       `xml:"lat,attr"`
	Lon        string       `xml:"lon,attr"`
	Extensions []Extensions `xml:"extensions"`
}

type Rte struct {
	Rtept []Rtept `xml:"rtept"`
}

type OpenRouteServiceResponse struct {
	XMLName xml.Name `xml:"gpx"`
	Rte     Rte      `xml:"rte"`
}
