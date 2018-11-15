package planar

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestCentroidArea(t *testing.T) {
	for _, g := range orb.AllGeometries {
		CentroidArea(g)
	}
}

func TestCentroidArea_MultiPoint(t *testing.T) {
	mp := orb.MultiPoint{{0, 0}, {1, 1.5}, {2, 0}}

	centroid, area := CentroidArea(mp)
	expected := orb.Point{1, 0.5}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 0 {
		t.Errorf("area should be 0: %f", area)
	}
}

func TestCentroidArea_LineString(t *testing.T) {
	cases := []struct {
		name   string
		ls     orb.LineString
		result orb.Point
	}{
		{
			name:   "simple",
			ls:     orb.LineString{{0, 0}, {3, 4}},
			result: orb.Point{1.5, 2},
		},
		{
			name:   "empty line",
			ls:     orb.LineString{},
			result: orb.Point{0, 0},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ls); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestCentroidArea_MultiLineString(t *testing.T) {
	cases := []struct {
		name   string
		ls     orb.MultiLineString
		result orb.Point
	}{
		{
			name:   "simple",
			ls:     orb.MultiLineString{{{0, 0}, {3, 4}}},
			result: orb.Point{1.5, 2},
		},
		{
			name:   "two lines",
			ls:     orb.MultiLineString{{{0, 0}, {0, 1}}, {{1, 0}, {1, 1}}},
			result: orb.Point{0.5, 0.5},
		},
		{
			name:   "multiple empty lines",
			ls:     orb.MultiLineString{{{1, 0}}, {{2, 1}}},
			result: orb.Point{1.5, 0.5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ls); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestCentroid_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   orb.Ring
		result orb.Point
	}{
		{
			name:   "triangle, cw",
			ring:   orb.Ring{{0, 0}, {1, 3}, {2, 0}, {0, 0}},
			result: orb.Point{1, 1},
		},
		{
			name:   "triangle, ccw",
			ring:   orb.Ring{{0, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: orb.Point{1, 1},
		},
		{
			name:   "square, cw",
			ring:   orb.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: orb.Point{0.5, 0.5},
		},
		{
			name:   "non-closed square, cw",
			ring:   orb.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}},
			result: orb.Point{0.5, 0.5},
		},
		{
			name:   "triangle, ccw",
			ring:   orb.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: orb.Point{0.5, 0.5},
		},
		{
			name:   "redudent points",
			ring:   orb.Ring{{0, 0}, {1, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: orb.Point{1, 1},
		},
		{
			name: "3 points",
			ring: orb.Ring{{0, 0}, {1, 0}, {0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e8
				tc.ring[i][1] -= 1e8
			}

			tc.result[0] += 1e8
			tc.result[1] -= 1e8

			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestArea_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   orb.Ring
		result float64
	}{
		{
			name:   "simple box, ccw",
			ring:   orb.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			ring:   orb.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: -1,
		},
		{
			name:   "even number of points",
			ring:   orb.Ring{{0, 0}, {1, 0}, {1, 1}, {0.4, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "3 points",
			ring:   orb.Ring{{0, 0}, {1, 0}, {0, 0}},
			result: 0.0,
		},
		{
			name:   "4 points",
			ring:   orb.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			result: 0.5,
		},
		{
			name:   "6 points",
			ring:   orb.Ring{{1, 1}, {2, 1}, {2, 1.5}, {2, 2}, {1, 2}, {1, 1}},
			result: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, val := CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e15
				tc.ring[i][1] -= 1e15
			}

			_, val = CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that are rendant last point is implicit
			tc.ring = tc.ring[:len(tc.ring)-1]
			_, val = CentroidArea(tc.ring)
			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

func TestCentroid_RingAdv(t *testing.T) {
	ring := orb.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5}, {2, 1}, {3, 1}, {3, 0}, {0, 0}}

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	expected := orb.Point{1.5, 0.45}
	if c, _ := CentroidArea(ring); !c.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", c, expected)
	}
}

func TestCentroidArea_Polygon(t *testing.T) {
	r1 := orb.Ring{{0, 0}, {4, 0}, {4, 3}, {0, 3}, {0, 0}}
	r1.Reverse()

	r2 := orb.Ring{{2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 1}}
	poly := orb.Polygon{r1, r2}

	centroid, area := CentroidArea(poly)
	if !centroid.Equal(orb.Point{21.5 / 11.0, 1.5}) {
		t.Errorf("%v", 21.5/11.0)
		t.Errorf("incorrect centroid: %v", centroid)
	}

	if area != 11 {
		t.Errorf("incorrect area: %v != 11", area)
	}

	// empty polygon
	e := orb.Point{0.5, 1}
	c, _ := CentroidArea(orb.Polygon{{{0, 1}, {1, 1}, {0, 1}}})
	if !c.Equal(e) {
		t.Errorf("incorrect point: %v != %v", c, e)
	}
}

func TestCentroidArea_Bound(t *testing.T) {
	b := orb.Bound{Min: orb.Point{0, 2}, Max: orb.Point{1, 3}}
	centroid, area := CentroidArea(b)

	expected := orb.Point{0.5, 2.5}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 1 {
		t.Errorf("incorrect area: %f != 1", area)
	}

	b = orb.Bound{Min: orb.Point{0, 2}, Max: orb.Point{0, 2}}
	centroid, area = CentroidArea(b)

	expected = orb.Point{0, 2}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 0 {
		t.Errorf("area should be zero: %f", area)
	}
}

func TestCentroidArea_Hole(t *testing.T) {
	geom := orb.Polygon{
		{
			{
				-102.2690493,
				40.9939916,
			},
			{
				-102.268951,
				40.9940453,
			},
			{
				-102.2690464,
				40.9941449,
			},
			{
				-102.2689287,
				40.9942091,
			},
			{
				-102.2688078,
				40.9940829,
			},
			{
				-102.2684138,
				40.9942979,
			},
			{
				-102.2683195,
				40.9941995,
			},
			{
				-102.2684496,
				40.9941286,
			},
			{
				-102.2683239,
				40.9939973,
			},
			{
				-102.2682875,
				40.9940171,
			},
			{
				-102.2682667,
				40.9939954,
			},
			{
				-102.2680722,
				40.9941015,
			},
			{
				-102.2679264,
				40.9939493,
			},
			{
				-102.2680433,
				40.9938855,
			},
			{
				-102.268059,
				40.9939019,
			},
			{
				-102.2682353,
				40.9938057,
			},
			{
				-102.2682124,
				40.9937818,
			},
			{
				-102.2682594,
				40.9937561,
			},
			{
				-102.2682808,
				40.9937785,
			},
			{
				-102.2683036,
				40.993766,
			},
			{
				-102.2683272,
				40.9937907,
			},
			{
				-102.2685813,
				40.9936521,
			},
			{
				-102.2688361,
				40.9939181,
			},
			{
				-102.2689632,
				40.9938488,
			},
			{
				-102.2690246,
				40.993913,
			},
			{
				-102.2689913,
				40.9939311,
			},
			{
				-102.2690493,
				40.9939916,
			},
		},
		{
			{
				-102.2687189,
				40.9939914,
			},
			{
				-102.2685563,
				40.9938227,
			},
			{
				-102.2684289,
				40.9938926,
			},
			{
				-102.2684666,
				40.9939317,
			},
			{
				-102.2684333,
				40.99395,
			},
			{
				-102.2685582,
				40.9940796,
			},
			{
				-102.2687189,
				40.9939914,
			},
		},
	}
	centroid, _ := CentroidArea(geom)

	bounds := geom.Bound()

	if !bounds.Contains(centroid) {
		t.Errorf("centroid not within the bounds of holey polygon: %v", centroid)
	}
}
