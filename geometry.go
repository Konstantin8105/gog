package gog

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

type BeamType int

const (
	Segment BeamType = iota
	Ray
	// Line
)

// Design:
//
//	-- P00 -- P0*==========*P1 -- P11 --
//	{  ray  }   {   beam   }   {  ray  }
//
type Beam struct {
	P0, P1 int // index of point
	Type   BeamType
}

type IntersectionType int64

const (
	empty IntersectionType = 1 << iota
	// property of single beam
	VerticalBeam0
	VerticalBeam1
	HorizontalBeam0
	HorizontalBeam1
	ZeroLengthBeam0
	ZeroLengthBeam1

	// property of both beams
	Parallel

	// intersection types
	Point0Beam0onPoint0Beam1
	Point1Beam0onPoint0Beam1
	Point0Beam0onPoint1Beam1
	Point1Beam0onPoint1Beam1

	Point0Beam0inBeam1 // 12
	Point1Beam0inBeam1
	Point0Beam1inBeam0
	Point1Beam1inBeam0

	IntersectOnBeam0 // 16
	IntersectOnBeam1

	IntersectBeam0Ray00
	IntersectBeam0Ray11
	IntersectBeam1Ray00
	IntersectBeam1Ray11

	// overlapping
	Collinear

	// TODO: Overlapping

	// last unused type
	endType
)

func is(t, ti IntersectionType) bool {
	return t&ti != 0
}
func not(t, ti IntersectionType) bool {
	return t&ti == 0
}

func (t IntersectionType) String() string {
	var out string
	var size int
	for i := 0; i < 64; i++ {
		if endType == 1<<i {
			size = i
			break
		}
	}
	for i := 1; i < size; i++ {
		ti := IntersectionType(1 << i)
		out += fmt.Sprintf("%2d\t%30b\t", i, int(ti))
		if is(t, ti) {
			out += "found"
		} else {
			out += "not found"
		}
		out += "\n"
	}
	return out
}

var eps float64 = 1e-6

func Intersection(b0, b1 Beam, ps *[]Point) (
	p Point,
	t IntersectionType,
) {

	// TODO: check inout data

	// TODO: check output intersection poinrt

	// see https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	var (
		x1 = (*ps)[b0.P0].X
		y1 = (*ps)[b0.P0].Y

		x2 = (*ps)[b0.P1].X
		y2 = (*ps)[b0.P1].Y

		x3 = (*ps)[b1.P0].X
		y3 = (*ps)[b1.P0].Y

		x4 = (*ps)[b1.P1].X
		y4 = (*ps)[b1.P1].Y
	)

	for _, c := range [...]struct {
		isTrue bool
		ti     IntersectionType
	}{
		{isTrue: math.Abs(x1-x3) < eps && math.Abs(y1-y3) < eps, ti: Point0Beam0onPoint0Beam1},
		{isTrue: math.Abs(x1-x4) < eps && math.Abs(y1-y4) < eps, ti: Point0Beam0onPoint1Beam1},
		{isTrue: math.Abs(x2-x3) < eps && math.Abs(y2-y3) < eps, ti: Point1Beam0onPoint0Beam1},
		{isTrue: math.Abs(x2-x4) < eps && math.Abs(y2-y4) < eps, ti: Point1Beam0onPoint1Beam1},
		{isTrue: math.Abs(x1-x2) < eps, ti: VerticalBeam0},
		{isTrue: math.Abs(x3-x4) < eps, ti: VerticalBeam1},
		{isTrue: math.Abs(y1-y2) < eps, ti: HorizontalBeam0},
		{isTrue: math.Abs(y3-y4) < eps, ti: HorizontalBeam1},
		{isTrue: math.Abs(x1-x2) < eps && math.Abs(y1-y2) < eps, ti: ZeroLengthBeam0},
		{isTrue: math.Abs(x3-x4) < eps && math.Abs(y3-y4) < eps, ti: ZeroLengthBeam1},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	switch {
	case is(t, Point0Beam0onPoint0Beam1) || is(t, Point0Beam0onPoint1Beam1):
		p = (*ps)[b0.P0]
	case is(t, Point1Beam0onPoint0Beam1) || is(t, Point1Beam0onPoint1Beam1):
		p = (*ps)[b0.P1]
	}

	// if zero, then vertical/horizontal
	B := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(B) < eps || is(t, ZeroLengthBeam0) || is(t, ZeroLengthBeam1) {
		if math.Abs((x3-x1)*(y2-y1)-(x2-x1)*(y3-y1)) < eps {
			t |= Collinear
		} else {
			t |= Parallel
		}
		return
	}

	// intersection point
	A12 := x1*y2 - y1*x2
	A34 := x3*y4 - y3*x4
	p.X = (A12*(x3-x4) - (x1-x2)*A34) / B
	p.Y = (A12*(y3-y4) - (y1-y2)*A34) / B

	// is intersect point on line?
	for _, c := range [...]struct {
		isTrue bool
		ti     IntersectionType
	}{
		{isTrue: math.Abs(x1-p.X) < eps && math.Abs(y1-p.Y) < eps, ti: Point0Beam0inBeam1},
		{isTrue: math.Abs(x2-p.X) < eps && math.Abs(y2-p.Y) < eps, ti: Point1Beam0inBeam1},
		{isTrue: math.Abs(x3-p.X) < eps && math.Abs(y3-p.Y) < eps, ti: Point0Beam1inBeam0},
		{isTrue: math.Abs(x4-p.X) < eps && math.Abs(y4-p.Y) < eps, ti: Point1Beam1inBeam0},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	for _, c := range [...]struct {
		isTrue bool
		ti     IntersectionType
	}{
		{
			isTrue: math.Min(x1, x2)-eps <= p.X && p.X <= math.Max(x1, x2)+eps &&
				math.Min(y1, y2)-eps <= p.Y && p.Y <= math.Max(y1, y2)+eps &&
				not(t, Point0Beam0inBeam1) && not(t, Point1Beam0inBeam1) &&
				not(t, Point0Beam1inBeam0) && not(t, Point1Beam1inBeam0),
			ti: IntersectOnBeam0,
		},
		{
			isTrue: math.Min(x3, x4)-eps <= p.X && p.X <= math.Max(x3, x4)+eps &&
				math.Min(y3, y4)-eps <= p.Y && p.Y <= math.Max(y3, y4)+eps &&
				not(t, Point0Beam0inBeam1) && not(t, Point1Beam0inBeam1) &&
				not(t, Point0Beam1inBeam0) && not(t, Point1Beam1inBeam0),
			ti: IntersectOnBeam1,
		},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	// is intersect point on ray?
	var (
		disB0    = Distance((*ps)[b0.P0], (*ps)[b0.P1])
		disB0P0p = Distance((*ps)[b0.P0], p)
		disB0P1p = Distance((*ps)[b0.P1], p)

		disB1    = Distance((*ps)[b1.P0], (*ps)[b1.P1])
		disB1P0p = Distance((*ps)[b1.P0], p)
		disB1P1p = Distance((*ps)[b1.P1], p)
	)
	for _, c := range [...]struct {
		isTrue bool
		ti     IntersectionType
	}{
		{
			isTrue: disB0P0p < disB0P1p && disB0 < disB0P1p,
			ti:     IntersectBeam0Ray00,
		},
		{
			isTrue: disB0P1p < disB0P0p && disB0 < disB0P0p,
			ti:     IntersectBeam0Ray11,
		},
		{
			isTrue: disB1P0p < disB1P1p && disB1 < disB1P1p,
			ti:     IntersectBeam1Ray00,
		},
		{
			isTrue: disB1P1p < disB1P0p && disB1 < disB1P0p,
			ti:     IntersectBeam1Ray11},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	return
}

func Distance(p0, p1 Point) float64 {
	return math.Hypot(p0.X-p1.X, p0.Y-p1.Y)
}
