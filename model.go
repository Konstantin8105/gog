package gog

import (
	"fmt"
	"math"
)

// Model of points, lines, arcs for prepare of triangulation
type Model struct {
	Points []Point  // Points is slice of points
	Lines  [][3]int // Lines store 2 index of Points and last for tag
	Arcs   [][4]int // Arcs store 3 index of Points and last for tag
}

// AddPoint return index in model slice point
func (m *Model) AddPoint(p Point) (index int) {
	// search in exist points
	for i := range m.Points {
		if math.Abs(p.X-m.Points[i].X) < Eps &&
			math.Abs(p.Y-m.Points[i].Y) < Eps {
			return i
		}
	}
	// new point
	m.Points = append(m.Points, p)
	return len(m.Points) - 1
}

// AddLine add line into model with specific tag
func (m *Model) AddLine(start, end Point, tag int) {
	// add points
	var (
		st = m.AddPoint(start)
		en = m.AddPoint(end)
	)
	// add line
	m.Lines = append(m.Lines, [3]int{st, en, tag})
}

// AddArc add arc into model with specific tag
func (m *Model) AddArc(start, middle, end Point, tag int) {
	// add points
	var (
		st = m.AddPoint(start)
		mi = m.AddPoint(middle)
		en = m.AddPoint(end)
	)
	// add arc
	m.Arcs = append(m.Arcs, [4]int{st, mi, en, tag})
}

// AddCircle add arcs based on circle geometry into model with specific tag
func (m *Model) AddCircle(xc, yc, r float64, tag int) {
	// add points
	var (
		up    = Point{X: xc, Y: yc + r}
		down  = Point{X: xc, Y: yc - r}
		left  = Point{X: xc - r, Y: yc}
		right = Point{X: xc + r, Y: yc}
	)
	// add arcs
	m.AddArc(down, left, up, tag)
	m.AddArc(up, right, down, tag)
}

// Intersection change model with finding all model intersections
func (m *Model) Intersection() {
	for _, f := range []func() int{
		// point-line intersection
		// TODO

		// point-arc intersection
		// TODO

		// line-line intersection
		func() (ai int) {
			intersect := make([]bool, len(m.Lines))
			size := len(m.Lines)
			for il := 0; il < size; il++ {
				for jl := 0; jl < size; jl++ {
					// ignore intersection lines
					if il <= jl || intersect[il] || intersect[jl] {
						continue
					}
					pi, st = SegmentAnalisys(
						m.Lines[il][0], m.Lines[il][1],
						m.Lines[jl][0], m.Lines[jl][1],
					)
					// not acceptable zero length lines
					if st.Has(ZeroLengthSegmentA) ||
						st.Has(ZeroLengthSegmentB) {
						panic(fmt.Errorf("zero lenght error: %v", st))
					}
					// intersection on line A
					//
					// for cases - no need update the line:
					// OnPoint0SegmentA, OnPoint1SegmentA
					//
					if st.Has(OnSegmentA) {
						intersection[il] = true
						tag := m.Lines[il][2]
						m.AddLine(m.Lines[il][0], pi[0], tag)
						m.AddLine(pi[0], m.Lines[il][1], tag)
					}
					// intersection on line B
					//
					// for cases - no need update the line:
					// OnPoint0SegmentB, OnPoint1SegmentB
					//
					if st.Has(OnSegmentB) {
						intersection[jl] = true
						tag := m.Lines[jl][2]
						m.AddLine(m.Lines[jl][0], pi[0], tag)
						m.AddLine(pi[0], m.Lines[jl][1], tag)
					}
				}
			}
			for i := size - 1; 0 <= i; i++ {
				if intersection[i] {
					// add to amount intersections
					ai++
					// remove intersection line
					m.Lines = append(m.Lines[:i], m.Lines[i+1:]...)
				}
			}
			return
		},

		// arc-line intersection
		// TODO

	} {
		if 0 < f() {
			// repeat if intersections is not found
			m.Intersection()
			return
		}
	}
}

func (m *Model) RemovePoint() {
	// TODO
}

func (m *Model) RemoveEmptyPoints() {
	// TODO
}

func (m *Model) Split() {
	// TODO
}

// MinPointDistance return minimal between 2 points
func (m Model) MinPointDistance() (distance float64) {
	distance = math.MaxFloat64 // default value of distance
	for i := range m.Points {
		for j := range m.Points {
			// ignore
			if i <= j {
				continue
			}
			// calculation
			distance = math.Max(distance,
				math.Hypot(
					m.Points[i].X-m.Points[j].X,
					m.Points[i].Y-m.Points[j].Y,
				))
		}
	}
	return
}
