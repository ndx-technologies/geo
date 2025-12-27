package geo

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (s Location) IsZero() bool { return s.Lat == 0 && s.Lon == 0 }

type LocationQuery struct {
	Location Location `json:"location"`
	Distance int      `json:"distance"` // meters
}

func (s LocationQuery) IsZero() bool { return s.Location.IsZero() || s.Distance == 0 }

type Span struct {
	DeltaLat float64 `json:"delta_lat"`
	DeltaLon float64 `json:"delta_lon"`
}

type Region struct {
	Center Location `json:"center"`
	Span   Span     `json:"span"` // size of whole region
}

func (s Region) Vertices() []Location {
	return []Location{
		{Lat: s.Center.Lat + s.Span.DeltaLat/2, Lon: s.Center.Lon + s.Span.DeltaLon/2},
		{Lat: s.Center.Lat + s.Span.DeltaLat/2, Lon: s.Center.Lon - s.Span.DeltaLon/2},
		{Lat: s.Center.Lat - s.Span.DeltaLat/2, Lon: s.Center.Lon - s.Span.DeltaLon/2},
		{Lat: s.Center.Lat - s.Span.DeltaLat/2, Lon: s.Center.Lon + s.Span.DeltaLon/2},
	}
}

type Polygon struct {
	Vertices []Location `json:"vertices"`
}
