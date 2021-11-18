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

	Arc01indentical
	Arc12indentical
	Arc02indentical
	ArcIsLine
	ArcOnePoint
	LineFromArcCenter
	LineOutside

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

// Check - check input data
func Check(pps ...Point) error {
	et := errors.New("Check points")
	for i := range pps {
		if x, y := pps[i].X, pps[i].Y; math.IsNaN(x) || math.IsInf(x, 0) ||
			math.IsNaN(y) || math.IsInf(y, 0) {
			et.Add(fmt.Errorf("Not valid point #%d: (%.5e,%.5e)", i, x, y))
		}
	}
	if et.IsError() {
		return et
	}
	return nil
}

var (
	// FindRayIntersection is global variable for switch off finding
	// intersection point on segments ray
	FindRayIntersection bool = true

	// Eps is epsilon - precision of intersection
	Eps float64 = 1e-6
)

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
		{isTrue: math.Abs(x1-x3) < Eps && math.Abs(y1-y3) < Eps, ti: OverlapP0AP0B},
		{isTrue: math.Abs(x1-x4) < Eps && math.Abs(y1-y4) < Eps, ti: OverlapP0AP1B},
		{isTrue: math.Abs(x2-x3) < Eps && math.Abs(y2-y3) < Eps, ti: OverlapP1AP0B},
		{isTrue: math.Abs(x2-x4) < Eps && math.Abs(y2-y4) < Eps, ti: OverlapP1AP1B},
		{isTrue: math.Abs(x1-x2) < Eps, ti: VerticalSegmentA},
		{isTrue: math.Abs(x3-x4) < Eps, ti: VerticalSegmentB},
		{isTrue: math.Abs(y1-y2) < Eps, ti: HorizontalSegmentA},
		{isTrue: math.Abs(y3-y4) < Eps, ti: HorizontalSegmentB},
		{isTrue: math.Abs(x1-x2) < Eps && math.Abs(y1-y2) < Eps, ti: ZeroLengthSegmentA},
		{isTrue: math.Abs(x3-x4) < Eps && math.Abs(y3-y4) < Eps, ti: ZeroLengthSegmentB},
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
	if math.Abs(B) < Eps || st.Has(ZeroLengthSegmentA) || st.Has(ZeroLengthSegmentB) {
		if math.Abs((x3-x1)*(y2-y1)-(x2-x1)*(y3-y1)) < Eps {
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
		{isTrue: math.Abs(x1-pi.X) < Eps && math.Abs(y1-pi.Y) < Eps, ti: OnPoint0SegmentA},
		{isTrue: math.Abs(x2-pi.X) < Eps && math.Abs(y2-pi.Y) < Eps, ti: OnPoint1SegmentA},
		{isTrue: math.Abs(x3-pi.X) < Eps && math.Abs(y3-pi.Y) < Eps, ti: OnPoint0SegmentB},
		{isTrue: math.Abs(x4-pi.X) < Eps && math.Abs(y4-pi.Y) < Eps, ti: OnPoint1SegmentB},
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
				math.Min(x1, x2)-Eps <= pi.X && pi.X <= math.Max(x1, x2)+Eps &&
				math.Min(y1, y2)-Eps <= pi.Y && pi.Y <= math.Max(y1, y2)+Eps,
			ti: OnSegmentA,
		},
		{
			isTrue: st.Not(OnPoint0SegmentB) && st.Not(OnPoint1SegmentB) &&
				math.Min(x3, x4)-Eps <= pi.X && pi.X <= math.Max(x3, x4)+Eps &&
				math.Min(y3, y4)-Eps <= pi.Y && pi.Y <= math.Max(y3, y4)+Eps,
			ti: OnSegmentB,
		},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	// is intersect point on ray?
	if FindRayIntersection {
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
	p0, p1 Point,
	pc Point,
) (distance float64) {
	A, B, C := line(p0, p1)
	var (
		// coordinates of point
		xm = pc.X
		ym = pc.Y
	)
	distance = math.Abs((A*xm + B*ym + C) / math.Sqrt(pow.E2(A)+pow.E2(B)))
	return
}

// line parameters
//	Ax+By+C = 0
func line(p0, p1 Point) (A, B, C float64) {
	var (
		dy = p1.Y - p0.Y
		dx = p1.X - p0.X
		x1 = p0.X
		y1 = p0.Y
	)
	// parameters of line
	A = dy
	B = -dx
	C = -dy*x1 + dx*y1
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

// MirrorLine return intersection point and second mirrored point from mirror
// line (mp0-mp1) and ray (sp0-sp1)
func MirrorLine(
	sp0, sp1 Point,
	mp0, mp1 Point,
) (
	ml0, ml1 Point,
	err error,
) {
	pi, ti := SegmentAnalisys(
		sp0, sp1,
		mp0, mp1,
	)
	if ti.Has(Parallel) || ti.Has(Collinear) {
		err = fmt.Errorf("Segment and mirror is not intersect")
		return
	}
	// Image of Point sp0 with respect to a line ax+bx+c=0 is line mirror mp0-mp1
	var (
		A, B, C = line(mp0, mp1)
		x1      = sp1.X
		y1      = sp1.Y
		common  = -2.0 * (A*x1 + B*y1 + C) / (pow.E2(A) + pow.E2(B))
	)
	ml0 = pi
	ml1.X = A*common + x1
	ml1.Y = B*common + y1
	if Distance(ml0, ml1) == 0.0 {
		sp1.X += (sp1.X - sp0.X) * 2.0
		sp1.Y += (sp1.Y - sp0.Y) * 2.0
		return MirrorLine(sp0, sp1, mp0, mp1)
	}
	return
}

type OrientationPoints int8

const (
	CollinearPoints        OrientationPoints = -1
	ClockwisePoints                          = 0
	CounterClockwisePoints                   = 1
)

func Orientation(p1, p2, p3 Point) OrientationPoints {
	v := (p2.Y-p1.Y)*(p3.X-p2.X) - (p2.X-p1.X)*(p3.Y-p2.Y)
	switch {
	case math.Abs(v) < 1e-6:
		return CollinearPoints
	case v > 0:
		return ClockwisePoints
	}
	return CounterClockwisePoints
}

func ArcLineAnalisys(Line [2]Point, Arc [3]Point) (
	pi []Point,
	st State,
) {
	var (
		x1 = Line[0].X
		y1 = Line[0].Y

		x2 = Line[1].X
		y2 = Line[1].Y
	)

	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{isTrue: math.Abs(x1-x2) < Eps, ti: VerticalSegmentA},
		{isTrue: math.Abs(y1-y2) < Eps, ti: HorizontalSegmentA},
		{isTrue: math.Abs(x1-x2) < Eps && math.Abs(y1-y2) < Eps, ti: ZeroLengthSegmentA},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	for _, c := range [...]struct {
		isTrue bool
		ti     State
	}{
		{isTrue: math.Abs(Arc[0].X-Arc[1].X) < Eps && math.Abs(Arc[0].Y-Arc[1].Y) < Eps, ti: Arc01indentical},
		{isTrue: math.Abs(Arc[1].X-Arc[2].X) < Eps && math.Abs(Arc[1].Y-Arc[2].Y) < Eps, ti: Arc12indentical},
		{isTrue: math.Abs(Arc[0].X-Arc[2].X) < Eps && math.Abs(Arc[0].Y-Arc[2].Y) < Eps, ti: Arc02indentical},
	} {
		if c.isTrue {
			st |= c.ti
		}
	}

	if st.Has(Arc01indentical) || st.Has(Arc12indentical) || st.Has(Arc02indentical){
		switch {
		case st.Has(Arc01indentical) && st.Has(Arc12indentical):
			st |= ArcOnePoint
			return
		case st.Has(Arc01indentical):
			pi2, st2 := SegmentAnalisys(Line[0], Line[1], Arc[0], Arc[2])
			st |= st2
			pi = append(pi, pi2)
		default:
			// for all cases and st.Has(Arc12indentical)
			pi2, st2 := SegmentAnalisys(Line[0], Line[1], Arc[0], Arc[1])
			st |= st2
			pi = append(pi, pi2)
		}
		st |= ArcIsLine
		return
	}

	// circle function
	//	(x-xc)^2 + (y-yc)^2 = R^2
	// solve circle parameters
	//	(xi-xc)^2+(yi-yc)^2 = R^2
	//	xi^2-2*xi*xc+xc^2 + yi^2-2*yi*yc+yc^2 = R^2
	//	xi^2+yi^2 -2*xi*xc-2*yi*yc = R^2-xc^2-yc^2
	// between points 1 and 2:
	//	(x1^2-x2^2) +(y1^2-y2^2) -2*xc*(x1-x2)-2*yc*(y1-y2) = 0
	//	(x1^2-x3^2) +(y1^2-y3^2) -2*xc*(x1-x3)-2*yc*(y1-y3) = 0
	//
	//	2*xc*(x1-x2)+2*yc*(y1-y2) = (x1^2-x2^2) + (y1^2-y2^2)
	//	2*xc*(x1-x3)+2*yc*(y1-y3) = (x1^2-x3^2) + (y1^2-y3^2)
	//
	//	2*(x1-x2)*xc + 2*(y1-y2)*yc = (x1^2-x2^2)+(y1^2-y2^2)
	//	2*(x1-x3)*xc + 2*(y1-y3)*yc = (x1^2-x3^2)+(y1^2-y3^2)
	//
	// solve linear equations:
	//	a11*xc + a12*yc = b1
	//	a21*xc + a22*yc = b2
	// solve:
	//	xc = (b1 - a12*yc)*1/a11
	//	a21*(b1-a12*yc)*1/a11 + a22*yc = b2
	//	yc*(a22-a21/a11*a12) = b2 - a21/a11*b1
	var xc, yc, r float64
	{
		var (
			x1, x2, x3 = Arc[0].X, Arc[1].X, Arc[2].X
			y1, y2, y3 = Arc[0].Y, Arc[1].Y, Arc[2].Y
			a11        = 2 * (x1 - x2)
			a12        = 2 * (y1 - y2)
			a21        = 2 * (x1 - x3)
			a22        = 2 * (y1 - y3)
			b1         = (pow.E2(x1) - pow.E2(x2)) + (pow.E2(y1) - pow.E2(y2))
			b2         = (pow.E2(x1) - pow.E2(x3)) + (pow.E2(y1) - pow.E2(y3))
			lin        = func(a11, a12, b1, a21, a22, b2 float64) (xc, yc float64) {
				yc = (b2 - a21/a11*b1) / (a22 - a21/a11*a12)
				xc = (b1 - a12*yc) * 1 / a11
				return
			}
		)
		if math.Abs(a11) < Eps {
			yc, xc = lin(a21, a11, b1, a22, a21, b2)
		} else {
			xc, yc = lin(a11, a12, b1, a21, a22, b2)
		}
		//	(xi-xc)^2+(yi-yc)^2 = R^2
		r1 := math.Sqrt(pow.E2(x1-xc) + pow.E2(y1-yc))
		r2 := math.Sqrt(pow.E2(x2-xc) + pow.E2(y2-yc))
		r3 := math.Sqrt(pow.E2(x3-xc) + pow.E2(y3-yc))
		r = (r1 + r2 + r3) / 3.0
	}

	if LinePointDistance(Line[0], Line[1], Point{xc,yc}) < Eps{
		st |= LineFromArcCenter
	}

	// line may be horizontal, vertical, other
	A, B, C := line(Line[0], Line[1])

	// Find intersections
	//	A*x+B*y+C = 0               :   line equations
	//	(x-xc)^2 + (y-yc)^2 = r^2   : circle equations
	var root []Point
	switch {
	case st.Has(HorizontalSegmentA):
		// line is horizontal
		//	A = 0
		//
		//	B*y+C = 0
		//	(x-xc)^2 + (y-yc)^2 = r^2
		//
		//	y = -C/B
		//	(x-xc)^2 = r^2 - (-C/B-yc)^2
		//
		//	x = +/- sqrt(r^2 - (-C/B-yc)^2) - xc
		D := pow.E2(r) - pow.E2(-C/B-yc)
		switch {
		case D < 0:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			root = append(root, Point{X: -xc, Y: Line[0].Y})
		default:
			// 0 < D
			root = append(root,
				Point{X: +math.Sqrt(D) - xc, Y: Line[0].Y},
				Point{X: -math.Sqrt(D) - xc, Y: Line[0].Y},
			)
		}

	case st.Has(VerticalSegmentA):
		// line is vertical
		// B = 0
		//
		//	A*x+C = 0
		//	(x-xc)^2 + (y-yc)^2 = r^2
		//
		//	x=-C/A
		//	(y-yc)^2 = r^2 - (-C/A-xc)^2
		//
		//	y = +/- sqrt(r^2 - (-C/A-xc)^2) - yc
		D := pow.E2(r) - pow.E2(-C/A-xc)
		switch {
		case D < 0:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			root = append(root, Point{X: Line[0].X, Y: -yc})
		default:
			// 0 < D
			root = append(root,
				Point{X: Line[0].X, Y: +math.Sqrt(D) - yc},
				Point{X: Line[0].X, Y: -math.Sqrt(D) - yc},
			)
		}

	default:
		//	A*x+B*y+C = 0               :   line function
		//	(x-xc)^2 + (y-yc)^2 = r^2   : circle function
		//
		// solve intersection:
		//	x = -(B*y+C)*1/A
		//	(-(B*y+C)*1/A-xc)^2 + (y-yc)^2 = r^2
		//	(-(B*y+C)-A*xc)^2 + (y-yc)^2*A^2 = r^2*A^2
		//	(-B*y-C-A*xc)^2 + (y-yc)^2*A^2 = r^2*A^2
		//	(-B*y-(C+A*xc))^2 + (y-yc)^2*A^2 = r^2*A^2
		//	(B*y)^2 + 2*(B*y)*(C+A*xc) + (C+A*xc)^2 + (y^2-2*y*yc+yc^2)*A^2 = r^2*A^2
		//	y^2*(B^2 + A^2) + y*(2*B*(C+A*xc) - 2*yc*A^2) + (C+A*xc)^2 + yc^2*A^2 - r^2*A^2 = 0
		//	    ==========      -------------------------   _______________________________
		var (
			a = pow.E2(B) + pow.E2(A)
			b = 2*B*(C+A*xc) - 2*yc*pow.E2(A)
			c = pow.E2(C+A*xc) + pow.E2(yc*A) - pow.E2(r*A)
			D = b*b - 4.0*a*c
		)

		// A and B of line parameters is not zero, so
		// value a is not a zero and more then zero.
		switch {
		case D < 0:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			y := -b / (2.0 * a)
			x := -(B*y + C) * 1 / A
			root = append(root, Point{X: x, Y: y})
		default:
			// 0 < D
			{
				y := (-b + math.Sqrt(D)) / (2.0 * a)
				x := -(B*y + C) * 1 / A
				root = append(root, Point{X: x, Y: y})
			}
			{
				y := (-b - math.Sqrt(D)) / (2.0 * a)
				x := -(B*y + C) * 1 / A
				root = append(root, Point{X: x, Y: y})
			}
		}
	}

	// is root on arc?
	for _, r := range root {
		a := math.Atan2(r.Y-yc, r.X-xc)
		a0 := math.Atan2(Arc[0].Y-yc, Arc[0].X-xc)
		a1 := math.Atan2(Arc[1].Y-yc, Arc[1].X-xc)
		a2 := math.Atan2(Arc[2].Y-yc, Arc[2].X-xc)
		if (a0 <= a && a <= a1) || (a1 <= a && a <= a2) {
			pi = append(pi, r)
		}
	}

	// is intersect point on line?
	for _, r := range pi {
		for _, c := range [...]struct {
			isTrue bool
			ti     State
		}{
			{isTrue: math.Abs(Line[0].X-r.X) < Eps && math.Abs(Line[0].Y-r.Y) < Eps,
				ti: OnPoint0SegmentA},
			{isTrue: math.Abs(Line[1].X-r.X) < Eps && math.Abs(Line[1].Y-r.Y) < Eps,
				ti: OnPoint1SegmentA},
		} {
			if c.isTrue {
				st |= c.ti
			}
		}
	}

	if 0 < len(pi) {
		if !st.Has(OnPoint0SegmentA) && !st.Has(OnPoint1SegmentA) {
			st |= OnSegmentA
		}
		st |= OnSegmentB
	} else {
		st |= LineOutside
	}

	return
}
