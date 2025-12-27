package s2x

import (
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	geo "github.com/ndx-technologies/mm-geo"
)

const (
	EarthRadius = 6_371_000
)

type S2RegionCovererConfig struct {
	MinLevel int `json:"min_level"`
	MaxLevel int `json:"max_level"`
	LevelMod int `json:"level_mod"`
	MaxCells int `json:"max_cells"`
}

func (s S2RegionCovererConfig) IsZero() bool {
	return s.MinLevel == 0 && s.MaxLevel == 0 && s.LevelMod == 0 && s.MaxCells == 0
}

func (s *S2RegionCovererConfig) ToS2RegionCoverer() *s2.RegionCoverer {
	if s == nil {
		return nil
	}
	return &s2.RegionCoverer{
		MinLevel: s.MinLevel,
		MaxLevel: s.MaxLevel,
		MaxCells: s.MaxCells,
		LevelMod: s.LevelMod,
	}
}

// chordArea computes area in meters to chord area in radians
func chordArea(s float64) float64 { return s / (EarthRadius * EarthRadius) }

func CoveringCellIDs(location geo.LocationQuery, config S2RegionCovererConfig) s2.CellUnion {
	ltl := s2.LatLngFromDegrees(location.Location.Lat, location.Location.Lon)

	rc := config.ToS2RegionCoverer()

	chordAngle := s1.ChordAngleFromSquaredLength(chordArea(float64(location.Distance * location.Distance)))

	// not using FastCovering, since it returns ~4 cells of the same level, resulting in very inaccurate area
	// not using Covering, since it is outer covering resulting in non-circle area and too large
	return rc.InteriorCovering(s2.CapFromCenterChordAngle(s2.PointFromLatLng(ltl), chordAngle))
}

func CoveringCellIDsApproxExterior(location geo.LocationQuery, config S2RegionCovererConfig) s2.CellUnion {
	ltl := s2.LatLngFromDegrees(location.Location.Lat, location.Location.Lon)

	rc := config.ToS2RegionCoverer()

	chordAngle := s1.ChordAngleFromSquaredLength(chordArea(float64(location.Distance * location.Distance)))

	return rc.FastCovering(s2.CapFromCenterChordAngle(s2.PointFromLatLng(ltl), chordAngle))
}

func CellIDsToAddFromCell(cell s2.CellID, config S2RegionCovererConfig) (cells []s2.CellID) {
	for cell := cell; cell.Level() >= config.MinLevel; cell = cell.Parent(cell.Level() - 1) {
		if cell.Level() <= config.MaxLevel {
			cells = append(cells, cell)
		}
	}
	return cells
}

func CellIDsToAdd(location geo.Location, config S2RegionCovererConfig) (cells []s2.CellID) {
	ltl := s2.LatLngFromDegrees(location.Lat, location.Lon)

	for cell := s2.CellIDFromLatLng(ltl); cell.Level() >= config.MinLevel; cell = cell.Parent(cell.Level() - 1) {
		if cell.Level() <= config.MaxLevel {
			cells = append(cells, cell)
		}
	}

	return cells
}

func CellIDsByLevel(cells []s2.CellID) map[int][]s2.CellID {
	if len(cells) == 0 {
		return nil
	}
	m := make(map[int][]s2.CellID)
	for _, cell := range cells {
		level := cell.Level()
		m[level] = append(m[level], cell)
	}
	return m
}

func ApproxAreaRadians(r geo.Span) float64 {
	// delaying constants as much as possible for compile time constant arithmetics
	return math.Pi / 180 * math.Pi / 180 * r.DeltaLat * r.DeltaLon
}

// ApproxCell is smallest s2 cell whose area is greater or equal to area
func ApproxCell(s geo.Region) s2.CellID {
	ltl := s2.LatLngFromDegrees(s.Center.Lat, s.Center.Lon)
	area := ApproxAreaRadians(s.Span)

	for cell := s2.CellIDFromLatLng(ltl); cell.Level() >= 1; cell = cell.Parent(cell.Level() - 1) {
		if s2.CellFromCellID(cell).ApproxArea() >= area {
			return cell
		}
	}

	return s2.CellIDFromLatLng(ltl)
}

func CoveringCellIDsUniformForCell(cell s2.CellID, level int) []s2.CellID {
	// cannot use RectBound and Covering, since it will produce gaps and overlaps
	// cannot use Children() since it gives wrong cells in chessboard pattern
	// instead, using similar approach to Denormalize, that is used in Covering s2 methods
	if level <= cell.Level() {
		return []s2.CellID{cell}
	}

	if cell.IsLeaf() {
		return []s2.CellID{cell}
	}

	var ids []s2.CellID

	end := cell.ChildEndAtLevel(level)
	for ci := cell.ChildBeginAtLevel(level); ci != end; ci = ci.Next() {
		ids = append(ids, ci)
	}

	return ids
}

func CoveringCellIDsExteriorUniform(s geo.Region, level int) s2.CellUnion {
	rc := s2.RegionCoverer{
		MinLevel: level,
		MaxLevel: level,
		MaxCells: 32,
	}

	center := s2.LatLngFromDegrees(s.Center.Lat, s.Center.Lon)
	size := s2.LatLngFromDegrees(s.Span.DeltaLat, s.Span.DeltaLon)

	return rc.Covering(s2.RectFromCenterSize(center, size))
}

func ExteriorBoundForCells(cells []s2.CellID) s2.Rect {
	g := s2.CellUnion(cells)
	return g.RectBound()
}

func CellPolygon(id s2.CellID) geo.Polygon {
	cell := s2.CellFromCellID(id)

	vertices := make([]geo.Location, 4)

	for i := range 4 {
		ll := s2.LatLngFromPoint(cell.Vertex(i))
		vertices[i].Lon = ll.Lng.Degrees()
		vertices[i].Lat = ll.Lat.Degrees()
	}

	return geo.Polygon{Vertices: vertices}
}

func RectToGeoRegion(r s2.Rect) geo.Region {
	center := r.Center()
	size := r.Size()

	return geo.Region{
		Center: geo.Location{Lat: center.Lat.Degrees(), Lon: center.Lng.Degrees()},
		Span:   geo.Span{DeltaLat: size.Lat.Degrees(), DeltaLon: size.Lng.Degrees()},
	}
}
