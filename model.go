package gog

import "math"

type Model struct {
	Points []Point  // Points is slice of points
	Lines  [][2]int // Lines store 2 index of Points
	Arcs   [][3]int // Arcs store 3 index of Points
}

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

func (m *Model) AddLine(start, end Point) {
	// add points
	var (
		st = m.AddPoint(start)
		en = m.AddPoint(end)
	)
	// add line
	m.Lines = append(m.Lines, [2]int{st, en})
}

func (m *Model) AddArc(start, middle, end Point) {
	// add points
	var (
		st = m.AddPoint(start)
		mi = m.AddPoint(middle)
		en = m.AddPoint(end)
	)
	// add arc
	m.Arcs = append(m.Arcs, [3]int{st, mi, en})
}

func (m *Model) AddCircle(xc, yc, r float64) {
	// add points
	var (
		up    = Point{X: xc, Y: yc + r}
		down  = Point{X: xc, Y: yc - r}
		left  = Point{X: xc - r, Y: yc}
		right = Point{X: xc + r, Y: yc}
	)
	// add arcs
	m.AddArc(down, left, up)
	m.AddArc(up, right, down)
}

func (m *Model) Split() {
}
