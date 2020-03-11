package gog

import (
	"fmt"
	"math"
)

//go:generate echo "# gog"
//go:generate echo "golang geometry library between point and segments"
//go:generate echo "```\n"
//go:generate go doc -all .
//go:generate echo "\n```"


type Point struct {
	X, Y float64
}

// Segment is part of line
//
// Design of segment:
//
//	-- P00 -- P0*==========*P1 -- P11 --
//	{  ray  }   {  segment }   {  ray  }
type Segment struct {
	P0, P1 int // indexes of point
}

// State is result of intersection
type State int64

const (
	empty State = 1 << iota
	// property of single segment
	VerticalSegment0
	VerticalSegment1
	HorizontalSegment0
	HorizontalSegment1
	ZeroLengthSegment0
	ZeroLengthSegment1

	// property of both segments
	Parallel

	// intersection types
	Point0Segment0onPoint0Segment1
	Point1Segment0onPoint0Segment1
	Point0Segment0onPoint1Segment1
	Point1Segment0onPoint1Segment1

	Point0Segment0inSegment1 // 12
	Point1Segment0inSegment1
	Point0Segment1inSegment0
	Point1Segment1inSegment0

	IntersectOnSegment0 // 16
	IntersectOnSegment1

	IntersectSegment0Ray00
	IntersectSegment0Ray11
	IntersectSegment1Ray00
	IntersectSegment1Ray11

	// overlapping
	Collinear

	// TODO: Overlapping

	// last unused type
	endType
)

// Has is mean s-State has si-State
func (s State) Has(si State) bool {
	return s&si != 0
}

// Not mean s-State have not si-State
func (s State) Not(si State) bool {
	return s&si == 0
}

// String is implementation of Stringer implementation for formating output
func (s State) String() string {
	var out string
	var size int
	for i := 0; i < 64; i++ {
		if endType == 1<<i {
			size = i
			break
		}
	}
	for i := 1; i < size; i++ {
		si := State(1 << i)
		out += fmt.Sprintf("%2d\t%30b\t", i, int(si))
		if s.Has(si) {
			out += "found"
		} else {
			out += "not found"
		}
		out += "\n"
	}
	return out
}

// eps is epsilon - precision of intersection
const eps float64 = 1e-6

func Intersection(b0, b1 Segment, ps *[]Point) (
	p Point,
	t State,
) {

	// TODO: check inout data

	// TODO: check output intersection point

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
		ti     State
	}{
		{isTrue: math.Abs(x1-x3) < eps && math.Abs(y1-y3) < eps, ti: Point0Segment0onPoint0Segment1},
		{isTrue: math.Abs(x1-x4) < eps && math.Abs(y1-y4) < eps, ti: Point0Segment0onPoint1Segment1},
		{isTrue: math.Abs(x2-x3) < eps && math.Abs(y2-y3) < eps, ti: Point1Segment0onPoint0Segment1},
		{isTrue: math.Abs(x2-x4) < eps && math.Abs(y2-y4) < eps, ti: Point1Segment0onPoint1Segment1},
		{isTrue: math.Abs(x1-x2) < eps, ti: VerticalSegment0},
		{isTrue: math.Abs(x3-x4) < eps, ti: VerticalSegment1},
		{isTrue: math.Abs(y1-y2) < eps, ti: HorizontalSegment0},
		{isTrue: math.Abs(y3-y4) < eps, ti: HorizontalSegment1},
		{isTrue: math.Abs(x1-x2) < eps && math.Abs(y1-y2) < eps, ti: ZeroLengthSegment0},
		{isTrue: math.Abs(x3-x4) < eps && math.Abs(y3-y4) < eps, ti: ZeroLengthSegment1},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	switch {
	case t.Has(Point0Segment0onPoint0Segment1) || t.Has(Point0Segment0onPoint1Segment1):
		p = (*ps)[b0.P0]
	case t.Has(Point1Segment0onPoint0Segment1) || t.Has(Point1Segment0onPoint1Segment1):
		p = (*ps)[b0.P1]
	}

	// if zero, then vertical/horizontal
	B := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(B) < eps || t.Has(ZeroLengthSegment0) || t.Has(ZeroLengthSegment1) {
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
		ti     State
	}{
		{isTrue: math.Abs(x1-p.X) < eps && math.Abs(y1-p.Y) < eps, ti: Point0Segment0inSegment1},
		{isTrue: math.Abs(x2-p.X) < eps && math.Abs(y2-p.Y) < eps, ti: Point1Segment0inSegment1},
		{isTrue: math.Abs(x3-p.X) < eps && math.Abs(y3-p.Y) < eps, ti: Point0Segment1inSegment0},
		{isTrue: math.Abs(x4-p.X) < eps && math.Abs(y4-p.Y) < eps, ti: Point1Segment1inSegment0},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{
			isTrue: math.Min(x1, x2)-eps <= p.X && p.X <= math.Max(x1, x2)+eps &&
				math.Min(y1, y2)-eps <= p.Y && p.Y <= math.Max(y1, y2)+eps &&
				t.Not(Point0Segment0inSegment1) && t.Not(Point1Segment0inSegment1) &&
				t.Not(Point0Segment1inSegment0) && t.Not(Point1Segment1inSegment0),
			ti: IntersectOnSegment0,
		},
		{
			isTrue: math.Min(x3, x4)-eps <= p.X && p.X <= math.Max(x3, x4)+eps &&
				math.Min(y3, y4)-eps <= p.Y && p.Y <= math.Max(y3, y4)+eps &&
				t.Not(Point0Segment0inSegment1) && t.Not(Point1Segment0inSegment1) &&
				t.Not(Point0Segment1inSegment0) && t.Not(Point1Segment1inSegment0),
			ti: IntersectOnSegment1,
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
		ti     State
	}{
		{
			isTrue: disB0P0p < disB0P1p && disB0 < disB0P1p,
			ti:     IntersectSegment0Ray00,
		},
		{
			isTrue: disB0P1p < disB0P0p && disB0 < disB0P0p,
			ti:     IntersectSegment0Ray11,
		},
		{
			isTrue: disB1P0p < disB1P1p && disB1 < disB1P1p,
			ti:     IntersectSegment1Ray00,
		},
		{
			isTrue: disB1P1p < disB1P0p && disB1 < disB1P0p,
			ti:     IntersectSegment1Ray11},
	} {
		if c.isTrue {
			t |= c.ti
		}
	}

	return
}

// Distance between two points
func Distance(p0, p1 Point) float64 {
	return math.Hypot(p0.X-p1.X, p0.Y-p1.Y)
}
