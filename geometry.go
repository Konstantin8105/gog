package gog

import (
	"fmt"
	"math"

	"github.com/Konstantin8105/errors"
)

//go:generate echo "# gog"
//go:generate echo "golang geometry library between point and segments"
//go:generate echo "```\n"
//go:generate go doc -all .
//go:generate echo "\n```"

type Point struct {
	X, Y float64
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

// Check - check input data
func Check(pps *[]Point) error {
	et := errors.New("Check points")
	for i := range *pps {
		if x, y := (*pps)[i].X, (*pps)[i].Y; math.IsNaN(x) || math.IsInf(x, 0) ||
			math.IsNaN(y) || math.IsInf(y, 0) {
			et.Add(fmt.Errorf("Not valid point #%d: (%.5e,%.5e)", i, x, y))
		}
	}
	if et.IsError() {
		return et
	}
	return nil
}

// SegmentAnalisys return analisys of two segments
//
// Design of segments:
//	                                            //
//	<-- rb00 -- pb0*==========*pb1 -- rb11 -->  // Segment B
//	                                            //
//	<-- ra00 -- pa0*==========*pa1 -- ra11 -->  // Segment A
//	{   ray   }{      segment     }{   ray   }  //
//	                                            //
//
// Input data:
//	ipa0, ipa1 - point indexes of segment A
//	ipb0, ipb1 - point indexes of segment B
//	pps      - pointer of point slice
//
// Output data:
//	pi - intersection point
//	st - states of analisys
//
// Reference:
//	[1]  https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
func SegmentAnalisys(
	ipa0, ipa1 int,
	ipb0, ipb1 int,
	pps *[]Point,
) (
	pi Point,
	st State,
) {
	// check input data of points is outside of that function

	// TODO: check output intersection point

	var (
		x1 = (*pps)[ipa0].X
		y1 = (*pps)[ipa0].Y

		x2 = (*pps)[ipa1].X
		y2 = (*pps)[ipa1].Y

		x3 = (*pps)[ipb0].X
		y3 = (*pps)[ipb0].Y

		x4 = (*pps)[ipb1].X
		y4 = (*pps)[ipb1].Y
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
			st |= c.ti
		}
	}

	switch {
	case st.Has(Point0Segment0onPoint0Segment1) || st.Has(Point0Segment0onPoint1Segment1):
		pi = (*pps)[ipa0]
	case st.Has(Point1Segment0onPoint0Segment1) || st.Has(Point1Segment0onPoint1Segment1):
		pi = (*pps)[ipa1]
	}

	// if zero, then vertical/horizontal
	B := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(B) < eps || st.Has(ZeroLengthSegment0) || st.Has(ZeroLengthSegment1) {
		if math.Abs((x3-x1)*(y2-y1)-(x2-x1)*(y3-y1)) < eps {
			st |= Collinear
		} else {
			st |= Parallel
		}
		return
	}

	// intersection point
	A12 := x1*y2 - y1*x2
	A34 := x3*y4 - y3*x4
	pi.X = (A12*(x3-x4) - (x1-x2)*A34) / B
	pi.Y = (A12*(y3-y4) - (y1-y2)*A34) / B

	// is intersect point on line?
	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{isTrue: math.Abs(x1-pi.X) < eps && math.Abs(y1-pi.Y) < eps, ti: Point0Segment0inSegment1},
		{isTrue: math.Abs(x2-pi.X) < eps && math.Abs(y2-pi.Y) < eps, ti: Point1Segment0inSegment1},
		{isTrue: math.Abs(x3-pi.X) < eps && math.Abs(y3-pi.Y) < eps, ti: Point0Segment1inSegment0},
		{isTrue: math.Abs(x4-pi.X) < eps && math.Abs(y4-pi.Y) < eps, ti: Point1Segment1inSegment0},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{
			isTrue: math.Min(x1, x2)-eps <= pi.X && pi.X <= math.Max(x1, x2)+eps &&
				math.Min(y1, y2)-eps <= pi.Y && pi.Y <= math.Max(y1, y2)+eps &&
				st.Not(Point0Segment0inSegment1) && st.Not(Point1Segment0inSegment1) &&
				st.Not(Point0Segment1inSegment0) && st.Not(Point1Segment1inSegment0),
			ti: IntersectOnSegment0,
		},
		{
			isTrue: math.Min(x3, x4)-eps <= pi.X && pi.X <= math.Max(x3, x4)+eps &&
				math.Min(y3, y4)-eps <= pi.Y && pi.Y <= math.Max(y3, y4)+eps &&
				st.Not(Point0Segment0inSegment1) && st.Not(Point1Segment0inSegment1) &&
				st.Not(Point0Segment1inSegment0) && st.Not(Point1Segment1inSegment0),
			ti: IntersectOnSegment1,
		},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	// is intersect point on ray?
	var (
		disB0    = Distance((*pps)[ipa0], (*pps)[ipa1])
		disB0P0p = Distance((*pps)[ipa0], pi)
		disB0P1p = Distance((*pps)[ipa1], pi)

		disB1    = Distance((*pps)[ipb0], (*pps)[ipb1])
		disB1P0p = Distance((*pps)[ipb0], pi)
		disB1P1p = Distance((*pps)[ipb1], pi)
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
			st |= c.ti
		}
	}

	// TODO: perpendicular

	return
}

// Distance between two points
func Distance(p0, p1 Point) float64 {
	return math.Hypot(p0.X-p1.X, p0.Y-p1.Y)
}
