package gog

import (
	"fmt"
	"math"

	"github.com/Konstantin8105/errors"
	"github.com/Konstantin8105/pow"
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

	VerticalSegmentA // vertical segment A
	VerticalSegmentB // vertical segment B

	HorizontalSegmentA // horizontal segment A
	HorizontalSegmentB // horizontal segment B

	ZeroLengthSegmentA // zero length segment A
	ZeroLengthSegmentB // zero length segment B

	// Segment A and segment B are parallel.
	// Intersection point data is not valid.
	Parallel

	// Segment A and segment B are collinear.
	// Intersection point data is not valid.
	Collinear

	OnSegmentA // intersection point on segment A
	OnSegmentB // intersection point on segment B

	OnRay00SegmentA // intersection point on ray 00 segment A
	OnRay11SegmentA // intersection point on ray 11 segment A
	OnRay00SegmentB // intersection point on ray 00 segment B
	OnRay11SegmentB // intersection point on ray 11 segment B

	OnPoint0SegmentA // intersection point on point 0 segment A
	OnPoint1SegmentA // intersection point on point 1 segment A
	OnPoint0SegmentB // intersection point on point 0 segment B
	OnPoint1SegmentB // intersection point on point 1 segment B

	OverlapP0AP0B // overlapping point 0 segment A and point 0 segment B
	OverlapP0AP1B // overlapping point 0 segment A and point 1 segment B
	OverlapP1AP0B // overlapping point 1 segment A and point 0 segment B
	OverlapP1AP1B // overlapping point 1 segment A and point 1 segment B

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
	pa0, pa1 Point,
	pb0, pb1 Point,
) (
	pi Point,
	st State,
) {
	// check input data of points is outside of that function

	var (
		x1 = pa0.X
		y1 = pa0.Y

		x2 = pa1.X
		y2 = pa1.Y

		x3 = pb0.X
		y3 = pb0.Y

		x4 = pb1.X
		y4 = pb1.Y
	)

	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{isTrue: math.Abs(x1-x3) < eps && math.Abs(y1-y3) < eps, ti: OverlapP0AP0B},
		{isTrue: math.Abs(x1-x4) < eps && math.Abs(y1-y4) < eps, ti: OverlapP0AP1B},
		{isTrue: math.Abs(x2-x3) < eps && math.Abs(y2-y3) < eps, ti: OverlapP1AP0B},
		{isTrue: math.Abs(x2-x4) < eps && math.Abs(y2-y4) < eps, ti: OverlapP1AP1B},
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
	case st.Has(OverlapP0AP0B) || st.Has(OverlapP0AP1B):
		pi = pa0
	case st.Has(OverlapP1AP0B) || st.Has(OverlapP1AP1B):
		pi = pa1
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
		{isTrue: math.Abs(x1-pi.X) < eps && math.Abs(y1-pi.Y) < eps, ti: OnPoint0SegmentA},
		{isTrue: math.Abs(x2-pi.X) < eps && math.Abs(y2-pi.Y) < eps, ti: OnPoint1SegmentA},
		{isTrue: math.Abs(x3-pi.X) < eps && math.Abs(y3-pi.Y) < eps, ti: OnPoint0SegmentB},
		{isTrue: math.Abs(x4-pi.X) < eps && math.Abs(y4-pi.Y) < eps, ti: OnPoint1SegmentB},
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
			isTrue: st.Not(OnPoint0SegmentA) && st.Not(OnPoint1SegmentA) &&
				math.Min(x1, x2)-eps <= pi.X && pi.X <= math.Max(x1, x2)+eps &&
				math.Min(y1, y2)-eps <= pi.Y && pi.Y <= math.Max(y1, y2)+eps,
			ti: OnSegmentA,
		},
		{
			isTrue: st.Not(OnPoint0SegmentB) && st.Not(OnPoint1SegmentB) &&
				math.Min(x3, x4)-eps <= pi.X && pi.X <= math.Max(x3, x4)+eps &&
				math.Min(y3, y4)-eps <= pi.Y && pi.Y <= math.Max(y3, y4)+eps,
			ti: OnSegmentB,
		},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	// is intersect point on ray?
	if st.Not(OnPoint0SegmentA) && st.Not(OnPoint1SegmentA) && st.Not(OnSegmentA) {
		disB0P0p := Distance(pa0, pi)
		disB0P1p := Distance(pa1, pi)
		if disB0P0p < disB0P1p {
			st |= OnRay00SegmentA
		} else {
			st |= OnRay11SegmentA
		}
	}
	if st.Not(OnPoint0SegmentB) && st.Not(OnPoint1SegmentB) && st.Not(OnSegmentB) {
		disB1P0p := Distance(pb0, pi)
		disB1P1p := Distance(pb1, pi)
		if disB1P0p < disB1P1p {
			st |= OnRay00SegmentB
		} else {
			st |= OnRay11SegmentB
		}
	}

	return
}

// LinePointDistance return distance between line and point
//
// Equation of line:
//	(y2-y1)*(x-x1) = (x2-x1)(y-y1)
//	dy*(x-x1) = dx*(y-y1)
//	dy*x-dy*x1-dx*y+dx*y1 = 0
//	Ax+By+C = 0
//	A = dy
//	B = -dx
//	C = -dy*x1+dx*y1
//
// Distance from point (xm,ym) to line:
//	d = |(A*xm+B*ym+C)/sqrt(A^2+B^2)|
func LinePointDistance(
	ip0, ip1 int,
	ipc int,
	pps *[]Point,
) (distance float64) {
	var (
		dy = (*pps)[ip1].Y - (*pps)[ip0].Y
		dx = (*pps)[ip1].X - (*pps)[ip0].X
		x1 = (*pps)[ip0].X
		y1 = (*pps)[ip0].Y
		// parameters of line
		A = dy
		B = -dx
		C = -dy*x1 + dx*y1
		// coordinates of point
		xm = (*pps)[ipc].X
		ym = (*pps)[ipc].Y
	)
	distance = math.Abs((A*xm + B*ym + C) / math.Sqrt(pow.E2(A)+pow.E2(B)))
	return
}

// Distance between two points
func Distance(p0, p1 Point) float64 {
	return math.Hypot(p0.X-p1.X, p0.Y-p1.Y)
}

// Rotate point about (0,0) on angle
func Rotate(angle float64, point Point) (p Point) {
	p.X = math.Cos(angle)*point.X - math.Sin(angle)*point.Y
	p.Y = math.Sin(angle)*point.X + math.Cos(angle)*point.Y
	return
}
