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

// Point is store of point coordinates
type Point struct {
	X, Y float64
}

// String is implementation of Stringer implementation for formating output
func (p Point) String() string {
	return fmt.Sprintf("[%.5e,%.5e]", p.X, p.Y)
}

// State is result of intersection
type State int64

const (
	empty State = 1 << iota
	// property of single segment
	VerticalSegmentA
	VerticalSegmentB
	HorizontalSegmentA
	HorizontalSegmentB
	ZeroLengthSegmentA
	ZeroLengthSegmentB

	// property of both segments
	Parallel

	// intersection types
	Point0SegmentAonPoint0SegmentB
	Point1SegmentAonPoint0SegmentB
	Point0SegmentAonPoint1SegmentB
	Point1SegmentAonPoint1SegmentB

	Point0SegmentAinSegmentB
	Point1SegmentAinSegmentB
	Point0SegmentBinSegmentA
	Point1SegmentBinSegmentA

	IntersectOnSegmentA
	IntersectOnSegmentB

	IntersectSegmentARay00
	IntersectSegmentARay11
	IntersectSegmentBRay00
	IntersectSegmentBRay11

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
		{isTrue: math.Abs(x1-x3) < eps && math.Abs(y1-y3) < eps, ti: Point0SegmentAonPoint0SegmentB},
		{isTrue: math.Abs(x1-x4) < eps && math.Abs(y1-y4) < eps, ti: Point0SegmentAonPoint1SegmentB},
		{isTrue: math.Abs(x2-x3) < eps && math.Abs(y2-y3) < eps, ti: Point1SegmentAonPoint0SegmentB},
		{isTrue: math.Abs(x2-x4) < eps && math.Abs(y2-y4) < eps, ti: Point1SegmentAonPoint1SegmentB},
		{isTrue: math.Abs(x1-x2) < eps, ti: VerticalSegmentA},
		{isTrue: math.Abs(x3-x4) < eps, ti: VerticalSegmentB},
		{isTrue: math.Abs(y1-y2) < eps, ti: HorizontalSegmentA},
		{isTrue: math.Abs(y3-y4) < eps, ti: HorizontalSegmentB},
		{isTrue: math.Abs(x1-x2) < eps && math.Abs(y1-y2) < eps, ti: ZeroLengthSegmentA},
		{isTrue: math.Abs(x3-x4) < eps && math.Abs(y3-y4) < eps, ti: ZeroLengthSegmentB},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	switch {
	case st.Has(Point0SegmentAonPoint0SegmentB) || st.Has(Point0SegmentAonPoint1SegmentB):
		pi = (*pps)[ipa0]
	case st.Has(Point1SegmentAonPoint0SegmentB) || st.Has(Point1SegmentAonPoint1SegmentB):
		pi = (*pps)[ipa1]
	}

	// if zero, then vertical/horizontal
	B := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(B) < eps || st.Has(ZeroLengthSegmentA) || st.Has(ZeroLengthSegmentB) {
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
		{isTrue: math.Abs(x1-pi.X) < eps && math.Abs(y1-pi.Y) < eps, ti: Point0SegmentAinSegmentB},
		{isTrue: math.Abs(x2-pi.X) < eps && math.Abs(y2-pi.Y) < eps, ti: Point1SegmentAinSegmentB},
		{isTrue: math.Abs(x3-pi.X) < eps && math.Abs(y3-pi.Y) < eps, ti: Point0SegmentBinSegmentA},
		{isTrue: math.Abs(x4-pi.X) < eps && math.Abs(y4-pi.Y) < eps, ti: Point1SegmentBinSegmentA},
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
				st.Not(Point0SegmentAinSegmentB) && st.Not(Point1SegmentAinSegmentB) &&
				st.Not(Point0SegmentBinSegmentA) && st.Not(Point1SegmentBinSegmentA),
			ti: IntersectOnSegmentA,
		},
		{
			isTrue: math.Min(x3, x4)-eps <= pi.X && pi.X <= math.Max(x3, x4)+eps &&
				math.Min(y3, y4)-eps <= pi.Y && pi.Y <= math.Max(y3, y4)+eps &&
				st.Not(Point0SegmentAinSegmentB) && st.Not(Point1SegmentAinSegmentB) &&
				st.Not(Point0SegmentBinSegmentA) && st.Not(Point1SegmentBinSegmentA),
			ti: IntersectOnSegmentB,
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
			ti:     IntersectSegmentARay00,
		},
		{
			isTrue: disB0P1p < disB0P0p && disB0 < disB0P0p,
			ti:     IntersectSegmentARay11,
		},
		{
			isTrue: disB1P0p < disB1P1p && disB1 < disB1P1p,
			ti:     IntersectSegmentBRay00,
		},
		{
			isTrue: disB1P1p < disB1P0p && disB1 < disB1P0p,
			ti:     IntersectSegmentBRay11},
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
