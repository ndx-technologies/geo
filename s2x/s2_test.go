package s2x_test

import (
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/golang/geo/s2"

	geo "github.com/ndx-technologies/geo"
	"github.com/ndx-technologies/geo/s2x"
)

func edgeFromArea(s float64) float64 {
	return math.Round(math.Sqrt(s * s2x.EarthRadius * s2x.EarthRadius))
}

func TestS2_interior(t *testing.T) {
	c := s2x.S2RegionCovererConfig{
		MinLevel: 1,
		MaxLevel: 19,
		LevelMod: 2,
		MaxCells: 32,
	}

	location := geo.Location{
		Lat: 22.27711532908625,
		Lon: 114.15855363568056,
	}

	t.Run("when very small location query, then it is empty", func(t *testing.T) {
		// smaller than cell size at level 19
		locationQuery := geo.LocationQuery{Location: location, Distance: 5}

		cover := s2x.CoveringCellIDs(locationQuery, c)
		if len(cover) != 0 {
			t.Error(len(cover))
		}
		if edge := edgeFromArea(cover.ApproxArea()); edge != 0 {
			t.Error(edge)
		}
	})

	t.Run("when small location query, then it is small", func(t *testing.T) {
		locationQuery := geo.LocationQuery{Location: location, Distance: 30}

		cover := s2x.CoveringCellIDs(locationQuery, c)
		if len(cover) == 0 {
			t.Error("cover should not be empty")
		}
		if len(cover) > c.MaxCells {
			t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
		}
		if edge := edgeFromArea(cover.ApproxArea()); edge != 33.0 {
			t.Error(33.0, edge)
		}
	})

	t.Run("when big location query, then it is expressed in few cells", func(t *testing.T) {
		locationQuery := geo.LocationQuery{Location: location, Distance: 10_000}

		cover := s2x.CoveringCellIDs(locationQuery, c)
		if len(cover) == 0 {
			t.Error("cover should not be empty")
		}
		if len(cover) > c.MaxCells {
			t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
		}
		if edge := edgeFromArea(cover.ApproxArea()); edge != 15192.0 {
			t.Error(15192.0, edge)
		}
	})

	t.Run("when very large location query, then it uses small number of cells", func(t *testing.T) {
		locationQuery := geo.LocationQuery{Location: location, Distance: 100_000_000}

		cover := s2x.CoveringCellIDs(locationQuery, c)
		if len(cover) > c.MaxCells {
			t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
		}
	})
}

func TestS2_exterior(t *testing.T) {
	t.Run("when all levels used", func(t *testing.T) {
		c := s2x.S2RegionCovererConfig{
			MinLevel: 1,
			MaxLevel: 19,
			LevelMod: 2,
			MaxCells: 32,
		}

		location := geo.Location{
			Lat: 22.27711532908625,
			Lon: 114.15855363568056,
		}

		t.Run("when very small location query, then it is not empty", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 5}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) == 0 {
				t.Error("cover should not be empty")
			}
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
			if edge := edgeFromArea(cover.ApproxArea()); edge != 39 {
				t.Error(39, edge)
			}
		})

		t.Run("when small location query, then it uses small number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 30}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) == 0 {
				t.Error("cover should not be empty")
			}
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
			if edge := edgeFromArea(cover.ApproxArea()); edge != 616.0 {
				t.Error(616.0, edge)
			}
		})

		t.Run("when big location query, then it uses small number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 10_000}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) == 0 {
				t.Error("cover should not be empty")
			}
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
			if edge := edgeFromArea(cover.ApproxArea()); edge != 157738.0 {
				t.Error(157738.0, edge)
			}
		})

		t.Run("when very large location query, then it uses small number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 100_000_000}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
		})
	})

	t.Run("when only deep levels used", func(t *testing.T) {
		c := s2x.S2RegionCovererConfig{
			MinLevel: 8,  // 27km
			MaxLevel: 10, // 7km
			LevelMod: 1,
			MaxCells: 4,
		}

		location := geo.Location{
			Lat: 22.27711532908625,
			Lon: 114.15855363568056,
		}

		t.Run("when very small location query, then it uses small number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 5}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) == 0 {
				t.Error("cover should not be empty")
			}
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
			if edge := edgeFromArea(cover.ApproxArea()); edge != 9_857.0 {
				t.Error(9_857.0, edge)
			}
		})

		t.Run("when big location query, then it uses small number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 10_000}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) == 0 {
				t.Error("cover should not be empty")
			}
			if len(cover) > c.MaxCells {
				t.Errorf("cover length %d > max cells %d", len(cover), c.MaxCells)
			}
			if edge := edgeFromArea(cover.ApproxArea()); edge != 78_876.0 {
				t.Error(78_876.0, edge)
			}
		})

		t.Run("when very large location query, then it has large number of cells", func(t *testing.T) {
			locationQuery := geo.LocationQuery{Location: location, Distance: 100_000_000}

			cover := s2x.CoveringCellIDsApproxExterior(locationQuery, c)
			if len(cover) <= c.MaxCells {
				t.Errorf("cover length %d <= max cells %d", len(cover), c.MaxCells)
			}
		})
	})
}

func TestCellIDsToAdd(t *testing.T) {
	c := s2x.S2RegionCovererConfig{
		MinLevel: 7,
		MaxLevel: 10,
		LevelMod: 2,
		MaxCells: 16,
	}

	t.Run("hashes matches levels", func(t *testing.T) {
		location := geo.Location{
			Lat: 22.27711532908625,
			Lon: 114.15855363568056,
		}

		cells := s2x.CellIDsToAdd(location, c)
		if len(cells) != c.MaxLevel-c.MinLevel+1 {
			t.Error(c.MaxLevel-c.MinLevel+1, len(cells))
		}

		exp := []s2.CellID{
			s2.CellIDFromToken("340401"),
			s2.CellIDFromToken("340404"),
			s2.CellIDFromToken("34041"),
			s2.CellIDFromToken("34044"),
		}
		if !slices.Equal(exp, cells) {
			t.Error(exp, cells)
		}
	})
}

func TestCellIDsToAddFromCell(t *testing.T) {
	c := s2x.S2RegionCovererConfig{
		MinLevel: 7,
		MaxLevel: 10,
		LevelMod: 2,
		MaxCells: 16,
	}

	t.Run("hashes matches levels", func(t *testing.T) {
		cells := s2x.CellIDsToAddFromCell(s2.CellIDFromToken("340401"), c)
		if len(cells) != c.MaxLevel-c.MinLevel+1 {
			t.Error(c.MaxLevel-c.MinLevel+1, len(cells))
		}

		exp := []s2.CellID{
			s2.CellIDFromToken("340401"),
			s2.CellIDFromToken("340404"),
			s2.CellIDFromToken("34041"),
			s2.CellIDFromToken("34044"),
		}
		if !slices.Equal(exp, cells) {
			t.Error(exp, cells)
		}
	})
}

func TestCellIDsByLevel(t *testing.T) {
	tests := []struct {
		cells []s2.CellID
		m     map[int][]s2.CellID
	}{
		{
			cells: nil,
			m:     nil,
		},
		{
			cells: []s2.CellID{
				s2.CellIDFromToken("47e66d"),
				s2.CellIDFromToken("47e673"),
				s2.CellIDFromToken("340404"),
			},
			m: map[int][]s2.CellID{
				9: {
					s2.CellIDFromToken("340404"),
				},
				10: {
					s2.CellIDFromToken("47e66d"),
					s2.CellIDFromToken("47e673"),
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := s2x.CellIDsByLevel(tt.cells)
			if len(m) != len(tt.m) {
				t.Fatalf("expected len %d, got %d", len(tt.m), len(m))
			}
			for k, v1 := range tt.m {
				if !slices.Equal(v1, m[k]) {
					t.Error(k, v1, m[k])
				}
			}
		})
	}
}

func TestApproxAreaRadians(t *testing.T) {
	tests := []struct {
		r    geo.Span
		area float64
	}{
		{
			r: geo.Span{
				DeltaLat: 0.0001,
				DeltaLon: 0.0001,
			},
			area: 123,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			area := s2x.ApproxAreaRadians(tt.r)
			expAreaRadians := tt.area / s2x.EarthRadius / s2x.EarthRadius
			if math.Abs(expAreaRadians-area) >= 1 {
				t.Error(expAreaRadians, area)
			}
		})
	}
}

func TestApproxCell(t *testing.T) {
	tests := []struct {
		name  string
		r     geo.Region
		cell  s2.CellID
		level int
	}{
		{
			name: "10m x 10m",
			r: geo.Region{
				Center: geo.Location{
					Lat: 22.276623143719352,
					Lon: 114.15826014588411,
				},
				Span: geo.Span{
					DeltaLat: 0.0001,
					DeltaLon: 0.0001,
				},
			},
			cell:  s2.CellIDFromToken("3404006f84c"),
			level: 19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := s2x.ApproxCell(tt.r)
			if cell.Level() != tt.level {
				t.Error(tt.level, cell.Level())
			}
			if cell != tt.cell {
				t.Error(tt.cell, cell, cell.ToToken())
			}
		})
	}
}

func TestCoveringCellIDsExteriorUniform(t *testing.T) {
	tests := []struct {
		name     string
		r        geo.Region
		cell     geo.Span
		covering s2.CellUnion
		level    int
	}{
		{
			name: "city",
			r: geo.Region{
				Center: geo.Location{
					Lat: 22.276623143719352,
					Lon: 114.15826014588411,
				},
				Span: geo.Span{
					DeltaLat: 0.01,
					DeltaLon: 0.01,
				},
			},
			level: 13,
			covering: s2.CellUnion{
				s2.CellIDFromToken("34040064"),
				s2.CellIDFromToken("3404006c"),
				s2.CellIDFromToken("34040074"),
				s2.CellIDFromToken("3404007c"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			covering := s2x.CoveringCellIDsExteriorUniform(tt.r, tt.level)

			if len(covering) >= 32 {
				t.Error(len(covering))
			}

			var tokens []string

			for _, cell := range covering {
				if cell.Level() != tt.level {
					t.Error(tt.level, cell.Level())
				}
				tokens = append(tokens, cell.ToToken())
			}

			if !slices.Equal(tt.covering, covering) {
				t.Error(tt.covering, covering, tokens)
			}
		})
	}
}

func TestCoveringCellIDsUniformForCell(t *testing.T) {
	tests := []struct {
		cell     s2.CellID
		level    int
		count    int
		covering []s2.CellID
	}{
		{
			cell:     s2.CellIDFromToken("3404003"),
			level:    12,
			count:    1,
			covering: []s2.CellID{s2.CellIDFromToken("3404003")},
		},
		{
			cell:     s2.CellIDFromToken("3404003"),
			level:    13,
			count:    4,
			covering: []s2.CellID{0x3404002400000000, 0x3404002c00000000, 0x3404003400000000, 0x3404003c00000000},
		},
		{
			cell:     s2.CellIDFromToken("3404003"),
			level:    14,
			count:    16,
			covering: []s2.CellID{0x3404002100000000, 0x3404002300000000, 0x3404002500000000, 0x3404002700000000, 0x3404002900000000, 0x3404002b00000000, 0x3404002d00000000, 0x3404002f00000000, 0x3404003100000000, 0x3404003300000000, 0x3404003500000000, 0x3404003700000000, 0x3404003900000000, 0x3404003b00000000, 0x3404003d00000000, 0x3404003f00000000},
		},
		{
			cell:     s2.CellIDFromToken("3404003"),
			level:    15,
			count:    64,
			covering: []s2.CellID{0x3404002040000000, 0x34040020c0000000, 0x3404002140000000, 0x34040021c0000000, 0x3404002240000000, 0x34040022c0000000, 0x3404002340000000, 0x34040023c0000000, 0x3404002440000000, 0x34040024c0000000, 0x3404002540000000, 0x34040025c0000000, 0x3404002640000000, 0x34040026c0000000, 0x3404002740000000, 0x34040027c0000000, 0x3404002840000000, 0x34040028c0000000, 0x3404002940000000, 0x34040029c0000000, 0x3404002a40000000, 0x3404002ac0000000, 0x3404002b40000000, 0x3404002bc0000000, 0x3404002c40000000, 0x3404002cc0000000, 0x3404002d40000000, 0x3404002dc0000000, 0x3404002e40000000, 0x3404002ec0000000, 0x3404002f40000000, 0x3404002fc0000000, 0x3404003040000000, 0x34040030c0000000, 0x3404003140000000, 0x34040031c0000000, 0x3404003240000000, 0x34040032c0000000, 0x3404003340000000, 0x34040033c0000000, 0x3404003440000000, 0x34040034c0000000, 0x3404003540000000, 0x34040035c0000000, 0x3404003640000000, 0x34040036c0000000, 0x3404003740000000, 0x34040037c0000000, 0x3404003840000000, 0x34040038c0000000, 0x3404003940000000, 0x34040039c0000000, 0x3404003a40000000, 0x3404003ac0000000, 0x3404003b40000000, 0x3404003bc0000000, 0x3404003c40000000, 0x3404003cc0000000, 0x3404003d40000000, 0x3404003dc0000000, 0x3404003e40000000, 0x3404003ec0000000, 0x3404003f40000000, 0x3404003fc0000000},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			covering := s2x.CoveringCellIDsUniformForCell(tt.cell, tt.level)

			if len(covering) >= 10_000 {
				t.Error(len(covering))
			}

			var tokens []string

			for _, cell := range covering {
				if cell.Level() != tt.level {
					t.Error(tt.level, cell.Level())
				}
				tokens = append(tokens, cell.ToToken())
			}

			if !slices.Equal(tt.covering, covering) {
				t.Error(tt.covering, covering, tokens)
			}
		})
	}
}

func angleDegreesFromDist(dist float64) float64 { return dist / s2x.EarthRadius * 180 / math.Pi }

func TestAvgEdgeMetric(t *testing.T) {
	const epsilon = 0.1

	tests := []struct {
		level       int
		edgeDegrees float64
	}{
		{level: 10, edgeDegrees: angleDegreesFromDist(10_000)},
		{level: 13, edgeDegrees: angleDegreesFromDist(1100)},
		{level: 16, edgeDegrees: angleDegreesFromDist(150)},
		{level: 20, edgeDegrees: angleDegreesFromDist(9)},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.level), func(t *testing.T) {
			edge := s2.AvgEdgeMetric.Value(int(tt.level)) * 180 / math.Pi
			if math.Abs(tt.edgeDegrees-edge)/tt.edgeDegrees >= epsilon {
				t.Error(tt.edgeDegrees, edge)
			}
		})
	}
}

func TestCellPolygon(t *testing.T) {
	tests := []struct {
		cell s2.CellID
		p    geo.Polygon
	}{
		{
			cell: s2.CellIDFromToken("3404003"),
			p: geo.Polygon{
				Vertices: []geo.Location{
					{Lat: 22.242312674946938, Lon: 114.16771589830563}, // low left
					{Lat: 22.23857226523964, Lon: 114.19149514780166},  // low right
					{Lat: 22.260892063882842, Lon: 114.19149514780166}, // top right
					{Lat: 22.264635440125442, Lon: 114.16771589830563}, // top left
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.cell.ToToken(), func(t *testing.T) {
			p := s2x.CellPolygon(tt.cell)
			if len(p.Vertices) != len(tt.p.Vertices) {
				t.Fatal(len(tt.p.Vertices), len(p.Vertices))
			}
			for i, v := range tt.p.Vertices {
				if p.Vertices[i] != v {
					t.Error(v, p.Vertices[i], i)
				}
			}
		})
	}
}

func TestCellCenter(t *testing.T) {
	epsilon := 0.000001 // 10m precision

	tests := []struct {
		cell s2.CellID
		lat  float64
		lon  float64
	}{
		{
			cell: s2.CellIDFromToken("3404003"),
			lat:  (22.23857226523964 + 22.264635440125442) / 2,  // from cell vertices
			lon:  (114.19149514780166 + 114.16771589830563) / 2, // from cell vertices
		},
	}
	for _, tt := range tests {
		t.Run(tt.cell.ToToken(), func(t *testing.T) {
			latlon := tt.cell.LatLng()
			lat := latlon.Lat.Degrees()
			lon := latlon.Lng.Degrees()
			if math.Abs(tt.lat-lat) >= epsilon {
				t.Error(tt.lat, lat)
			}
			if math.Abs(tt.lon-lon) >= epsilon {
				t.Error(tt.lon, lon)
			}
		})
	}
}
