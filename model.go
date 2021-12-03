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
		if Distance(p, m.Points[i]) < Eps {
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
	if en < st {
		st, en = en, st
	}
	// do not add line with same id
	for i := range m.Lines {
		if m.Lines[i][0] == st && m.Lines[i][1] == en {
			m.Lines[i][2] = tag
			return
		}
	}
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
	if en < st {
		st, en = en, st
	}
	// do not add line with same id
	for i := range m.Arcs {
		if (m.Arcs[i][0] == st && m.Arcs[i][1] == mi && m.Arcs[i][2] == en) ||
			(m.Arcs[i][2] == st && m.Arcs[i][1] == mi && m.Arcs[i][0] == en) {
			m.Arcs[i][3] = tag
			return
		}
	}
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
	// CounterClockwisePoints
	m.AddArc(down, right, up, tag)
	m.AddArc(up, left, down, tag)
}

// Intersection change model with finding all model intersections
func (m *Model) Intersection() {
	// value `ai` is amount of intersections

	// find intersections
	fs := []func() int{
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
					// analyse
					pi, stA, stB := LineLine(
						m.Points[m.Lines[il][0]], m.Points[m.Lines[il][1]],
						m.Points[m.Lines[jl][0]], m.Points[m.Lines[jl][1]],
					)
					// no intersections
					if 0 == len(pi) {
						continue
					}
					// debug test
					if 1 < len(pi) {
						panic("not valid")
					}
					// not acceptable zero length lines
					if stA.Has(ZeroLengthSegment) ||
						stB.Has(ZeroLengthSegment) {
						panic(fmt.Errorf("zero lenght error"))
					}

					if stA.Has(OnPoint0Segment) && stA.Has(OnPoint1Segment) &&
						stB.Has(OnPoint0Segment) && stB.Has(OnPoint1Segment) {
						intersect[il] = true
						continue
					}

					// intersection on line A
					//
					// for cases - no need update the line:
					// OnPoint0SegmentA, OnPoint1SegmentA
					//
					if stA.Has(OnSegment) && !(stA.Has(OnPoint0Segment) || stA.Has(OnPoint1Segment)) {
						intersect[il] = true
						tag := m.Lines[il][2]
						m.AddLine(m.Points[m.Lines[il][0]], pi[0], tag)
						m.AddLine(pi[0], m.Points[m.Lines[il][1]], tag)
					}
					// intersection on line B
					//
					// for cases - no need update the line:
					// OnPoint0SegmentB, OnPoint1SegmentB
					//
					if stB.Has(OnSegment) && !(stB.Has(OnPoint0Segment) || stB.Has(OnPoint1Segment)) {
						intersect[jl] = true
						tag := m.Lines[jl][2]
						m.AddLine(m.Points[m.Lines[jl][0]], pi[0], tag)
						m.AddLine(pi[0], m.Points[m.Lines[jl][1]], tag)
					}
				}
			}
			for i := size - 1; 0 <= i; i-- {
				if intersect[i] {
					// add to amount intersections
					ai++
					// remove intersection line
					m.Lines = append(m.Lines[:i], m.Lines[i+1:]...)
				}
			}
			return
		},

		// arc-line intersection
		func() (ai int) {
			var (
				intersectLines = make([]bool, len(m.Lines))
				intersectArcs  = make([]bool, len(m.Arcs))
				sizeLines      = len(m.Lines)
				sizeArcs       = len(m.Arcs)
			)
			for il := 0; il < sizeLines; il++ {
				for ja := 0; ja < sizeArcs; ja++ {
					// ignore intersection lines
					if intersectLines[il] || intersectArcs[ja] {
						continue
					}
					// analyse
					pi, stA, stB := LineArc(
						// Line
						m.Points[m.Lines[il][0]], m.Points[m.Lines[il][1]],
						// Arc
						m.Points[m.Arcs[ja][0]],
						m.Points[m.Arcs[ja][1]],
						m.Points[m.Arcs[ja][2]],
					)
					// not acceptable zero length lines
					if stA.Has(ZeroLengthSegment) ||
						stB.Has(ZeroLengthSegment) {
						panic(fmt.Errorf("zero lenght error"))
					}
					// intersection on line A
					//
					// for cases - no need update the line:
					// OnPoint0Segment, OnPoint1Segment
					//
					if stA.Has(OnSegment) {
						// remove OnPoint
						roots := make([]Point, len(pi))
						copy(roots, pi)

					same1:
						for i := range roots {
							for j := 0; j < 2; j++ {
								if Distance(roots[i], m.Points[m.Lines[il][j]]) < Eps {
									roots = append(roots[:i], roots[i+1:]...)
									goto same1
								}
							}
						}

						if 0 < len(roots) {
							intersectLines[il] = true
							tag := m.Lines[il][2]
							switch len(roots) {
							case 1:
								m.AddLine(m.Points[m.Lines[il][0]], roots[0], tag)
								m.AddLine(roots[0], m.Points[m.Lines[il][1]], tag)
							case 2:
								if stA.Has(VerticalSegment) {
									if roots[1].Y < roots[0].Y {
										roots[0], roots[1] = roots[1], roots[0]
									}
									// roots[0].Y < roots[1].Y
									if m.Points[m.Lines[il][0]].Y < m.Points[m.Lines[il][1]].Y {
										// Design:
										//
										//	| Lines [1]
										//	| roots[1]
										//	| roots[0]
										//	| Lines [0]
										m.AddLine(m.Points[m.Lines[il][0]], roots[0], tag)
										m.AddLine(roots[0], roots[1], tag)
										m.AddLine(roots[1], m.Points[m.Lines[il][1]], tag)
									} else {
										// Design:
										//
										//	| Lines [0]
										//	| roots[1]
										//	| roots[0]
										//	| Lines [1]
										m.AddLine(m.Points[m.Lines[il][1]], roots[0], tag)
										m.AddLine(roots[0], roots[1], tag)
										m.AddLine(roots[1], m.Points[m.Lines[il][0]], tag)
									}
								} else {
									// Not vertical line
									if roots[1].X < roots[0].X {
										roots[0], roots[1] = roots[1], roots[0]
									}
									// roots[0].X < roots[1].X
									if m.Points[m.Lines[il][0]].X < m.Points[m.Lines[il][1]].X {
										// Design:
										//
										//	 Lines[0]    roots[0]   roots[1]   Lines[1]
										m.AddLine(m.Points[m.Lines[il][0]], roots[0], tag)
										m.AddLine(roots[0], roots[1], tag)
										m.AddLine(roots[1], m.Points[m.Lines[il][1]], tag)
									} else {
										// Design:
										//
										//	 Lines[1]    roots[0]   roots[1]   Lines[0]
										m.AddLine(m.Points[m.Lines[il][1]], roots[0], tag)
										m.AddLine(roots[0], roots[1], tag)
										m.AddLine(roots[1], m.Points[m.Lines[il][0]], tag)
									}
								}
							default:
								panic("not valid")
							}
						}
					}
					// intersection on arc B
					//
					// for cases - no need update the line:
					// OnPoint0SegmentB, OnPoint1SegmentB
					//
					if stB.Has(OnSegment) {
						tag := m.Arcs[ja][3]
						res, err := ArcSplitByPoint(
							m.Points[m.Arcs[ja][0]],
							m.Points[m.Arcs[ja][1]],
							m.Points[m.Arcs[ja][2]],
							pi...)
						if err != nil {
							// TODO	panic(err)
							err = nil
						} else {
							for i := range res {
								intersectArcs[ja] = true
								m.AddArc(res[i][0], res[i][1], res[i][2], tag)
							}
						}
					}
				}
			}
			for i := sizeLines - 1; 0 <= i; i-- {
				if intersectLines[i] {
					// add to amount intersections
					ai++
					// remove intersection line
					m.Lines = append(m.Lines[:i], m.Lines[i+1:]...)
				}
			}
			for i := sizeArcs - 1; 0 <= i; i-- {
				if intersectArcs[i] {
					// add to amount intersections
					ai++
					// remove intersection arcs
					m.Arcs = append(m.Arcs[:i], m.Arcs[i+1:]...)
				}
			}
			return
		},

		// point-arc intersection
		// TODO

		// point-line intersection
		// TODO

		// point-point intersection
		// TODO

	}
	ai := 0
	for _, f := range fs {
		ai += f()
		if 0 < ai {
			m.Intersection()
			return
		}
	}
}

func (m *Model) Merge() {
	// TODO
}

func (m *Model) Rotate() {
	// TODO
}

func (m *Model) Move() {
	// TODO
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
			distance = math.Min(distance,
				math.Hypot(
					m.Points[i].X-m.Points[j].X,
					m.Points[i].Y-m.Points[j].Y,
				))
		}
	}
	return
}
