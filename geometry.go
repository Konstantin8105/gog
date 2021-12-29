package gog

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/Konstantin8105/errors"
	"github.com/Konstantin8105/pow"
)

//go:generate rm -f *.dxf
//go:generate rm -f *.out
//go:generate rm -f *.test
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

	VerticalSegment // vertical segment

	HorizontalSegment // horizontal segment

	ZeroLengthSegment // zero length segment

	// Segment A and segment B are parallel.
	// Intersection point data is not valid.
	Parallel

	// Segment A and segment B are collinear.
	// Intersection point data is not valid.
	Collinear

	OnSegment // intersection point on segment

	OnRay00Segment // intersection point on ray 00 segment
	OnRay11Segment // intersection point on ray 11 segment

	OnPoint0Segment // intersection point on point 0 segment
	OnPoint1Segment // intersection point on point 1 segment

	ArcIsLine  // wrong arc is line
	ArcIsPoint // wrong arc is point

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
	Eps float64 = 1e-12
)

func PointPoint(
	pt0, pt1 Point,
) (
	pi []Point,
	stA, stB State,
) {
	stA |= ZeroLengthSegment | VerticalSegment | HorizontalSegment
	if Distance(pt0, pt1) < Eps {
		stA |= OnPoint0Segment | OnPoint1Segment
	}
	stB = stA
	return
}

func PointLine(
	pt Point,
	pb0, pb1 Point,
) (
	pi []Point,
	stA, stB State,
) {
	// Point - Point
	if Distance(pb0, pb1) < Eps {
		return PointPoint(pt, pb0)
	}
	// Point - Line

	stA |= ZeroLengthSegment | VerticalSegment | HorizontalSegment

	for _, c := range [...]struct {
		isTrue   bool
		tiA, tiB State
	}{
		{isTrue: Distance(pt, pb0) < Eps, tiA: OnPoint0Segment | OnPoint1Segment, tiB: OnPoint0Segment},
		{isTrue: Distance(pt, pb1) < Eps, tiA: OnPoint0Segment | OnPoint1Segment, tiB: OnPoint1Segment},
		{isTrue: math.Abs(pb0.X-pb1.X) < Eps, tiB: VerticalSegment},
		{isTrue: math.Abs(pb0.Y-pb1.Y) < Eps, tiB: HorizontalSegment},
	} {
		if c.isTrue {
			stA |= c.tiA
			stB |= c.tiB
		}
	}

	if math.Abs(Distance(pt, pb0)+Distance(pt, pb1)-Distance(pb0, pb1)) < Eps &&
		stB.Not(OnPoint0Segment) && stB.Not(OnPoint1Segment) {
		stA |= OnPoint0Segment | OnPoint1Segment
		stB |= OnSegment
	}

	if stB.Has(OnSegment) && stB.Not(OnPoint0Segment) && stB.Not(OnPoint1Segment) {
		pi = []Point{pt}
	}

	// is point on ray
	if FindRayIntersection &&
		stB.Not(OnPoint0Segment) &&
		stB.Not(OnPoint1Segment) &&
		stB.Not(OnSegment) {
		if Distance(pb0, pt) < Distance(pb1, pt) {
			stB |= OnRay00Segment
		} else {
			stB |= OnRay11Segment
		}
	}

	return
}

// LineLine return analisys of two segments
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
func LineLine(
	pa0, pa1 Point,
	pb0, pb1 Point,
) (
	pi []Point,
	stA, stB State,
) {
	// Point - Point
	if Distance(pa0, pa1) < Eps && Distance(pb0, pb1) < Eps {
		return PointPoint(pa0, pb0)
	}
	// Point - Line
	if Distance(pa0, pa1) < Eps {
		return PointLine(pa0, pb0, pb1)
	}
	if Distance(pb0, pb1) < Eps {
		pi, stA, stB = PointLine(pb0, pa0, pa1)
		stA, stB = stB, stA
		return
	}
	// Line - Line

	for _, c := range [...]struct {
		isTrue   bool
		tiA, tiB State
	}{
		{isTrue: Distance(pa0, pb0) < Eps, tiA: OnPoint0Segment, tiB: OnPoint0Segment},
		{isTrue: Distance(pa0, pb1) < Eps, tiA: OnPoint0Segment, tiB: OnPoint1Segment},
		{isTrue: Distance(pa1, pb0) < Eps, tiA: OnPoint1Segment, tiB: OnPoint0Segment},
		{isTrue: Distance(pa1, pb1) < Eps, tiA: OnPoint1Segment, tiB: OnPoint1Segment},
		{isTrue: math.Abs(pa0.X-pa1.X) < Eps, tiA: VerticalSegment},
		{isTrue: math.Abs(pa0.Y-pa1.Y) < Eps, tiA: HorizontalSegment},
		{isTrue: math.Abs(pb0.X-pb1.X) < Eps, tiB: VerticalSegment},
		{isTrue: math.Abs(pb0.Y-pb1.Y) < Eps, tiB: HorizontalSegment},
	} {
		if c.isTrue {
			stA |= c.tiA
			stB |= c.tiB
		}
	}

	// collinear or parallel
	Aa, Ba, Ca := Line(pa0, pa1)
	Ab, Bb, Cb := Line(pb0, pb1)
	if math.Abs((pa1.Y-pa0.Y)*(pb1.X-pb0.X)-(pb1.Y-pb0.Y)*(pa1.X-pa0.X)) < Eps {
		collinear := false
		switch {
		case stA.Has(VerticalSegment) && stB.Has(VerticalSegment):
			if math.Abs(pa0.X-pb0.X) < Eps {
				collinear = true
			}
		case stA.Has(HorizontalSegment) && stB.Has(HorizontalSegment):
			if math.Abs(pa0.Y-pb0.Y) < Eps {
				collinear = true
			}
		default:
			if Eps < math.Abs(Aa) && Eps < math.Abs(Ab) && math.Abs(Ca/Aa-Cb/Ab) < Eps {
				collinear = true
			}
		}

		if collinear {
			stA |= Collinear
			stB |= Collinear
		} else {
			stA |= Parallel
			stB |= Parallel
		}
		return
	}

	// intersection point
	x, y := Linear(Aa, Ba, -Ca, Ab, Bb, -Cb)
	root := Point{X: x, Y: y}
	{
		_, _, stBa := PointLine(root, pa0, pa1)
		_, _, stBb := PointLine(root, pb0, pb1)
		if stBa.Has(OnSegment) &&
			(stBb.Has(OnSegment) || stBb.Has(OnPoint0Segment) || stBb.Has(OnPoint1Segment)) {
			stA |= OnSegment
		}
		if stBb.Has(OnSegment) &&
			(stBa.Has(OnSegment) || stBa.Has(OnPoint0Segment) || stBa.Has(OnPoint1Segment)) {
			stB |= OnSegment
		}

		if stBa.Has(OnRay00Segment) {
			stA |= OnRay00Segment
		}
		if stBa.Has(OnRay11Segment) {
			stA |= OnRay11Segment
		}

		if stBb.Has(OnRay00Segment) {
			stB |= OnRay00Segment
		}
		if stBb.Has(OnRay11Segment) {
			stB |= OnRay11Segment
		}
	}
	if stA.Has(OnSegment) || stB.Has(OnSegment) {
		pi = []Point{root}
	}

	for _, c := range [...]struct {
		isTrue   bool
		tiA, tiB State
	}{
		{isTrue: Distance(pa0, root) < Eps, tiA: OnPoint0Segment},
		{isTrue: Distance(pa1, root) < Eps, tiA: OnPoint1Segment},
		{isTrue: Distance(pb0, root) < Eps, tiB: OnPoint0Segment},
		{isTrue: Distance(pb1, root) < Eps, tiB: OnPoint1Segment},
	} {
		if c.isTrue {
			stA |= c.tiA
			stB |= c.tiB
		}
	}

	return
}

func MiddlePoint(p0, p1 Point) Point {
	const prec = 128

	var (
		half = new(big.Float).SetPrec(prec).SetFloat64(0.5)
		x0   = new(big.Float).SetPrec(prec).SetFloat64(p0.X)
		x1   = new(big.Float).SetPrec(prec).SetFloat64(p1.X)
		y0   = new(big.Float).SetPrec(prec).SetFloat64(p0.Y)
		y1   = new(big.Float).SetPrec(prec).SetFloat64(p1.Y)
	)
	x0.Mul(x0, half)
	x1.Mul(x1, half)
	y0.Mul(y0, half)
	y1.Mul(y1, half)

	x0.Add(x0, x1)
	y0.Add(y0, y1)

	x, _ := x0.Float64()
	y, _ := y0.Float64()

	mid := Point{X: x, Y: y}

	// Simple float64 algoritm:
	// mid := Point{
	// 	X: p0.X*0.5 + p1.X*0.5,
	// 	Y: p0.Y*0.5 + p1.Y*0.5,
	// }
	return mid
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
func PointLineDistance(
	pc Point,
	p0, p1 Point,
) (distance float64) {
	A, B, C := Line(p0, p1)
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
func Line(p0, p1 Point) (A, B, C float64) {
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

func Distance128(p0, p1 Point) float64 {
	const prec = 128

	var (
		x0   = new(big.Float).SetPrec(prec).SetFloat64(p0.X)
		x1   = new(big.Float).SetPrec(prec).SetFloat64(p1.X)
		y0   = new(big.Float).SetPrec(prec).SetFloat64(p0.Y)
		y1   = new(big.Float).SetPrec(prec).SetFloat64(p1.Y)
		x    = new(big.Float).SetPrec(prec).Sub(x0, x1)
		y    = new(big.Float).SetPrec(prec).Sub(y0, y1)
		xx   = new(big.Float).SetPrec(prec).Mul(x, x)
		yy   = new(big.Float).SetPrec(prec).Mul(y, y)
		summ = new(big.Float).SetPrec(prec).Add(xx, yy)
		s    = new(big.Float).SetPrec(prec).Sqrt(summ)
	)

	sf, _ := s.Float64()
	return sf
}

// Distance between two points
func Distance(p0, p1 Point) float64 {
	v := math.Hypot(p0.X-p1.X, p0.Y-p1.Y)
	if 100*Eps < v {
		return v
	}
	return Distance128(p0, p1)
}

// Rotate point about (xc,yc) on angle
func Rotate(xc, yc, angle float64, point Point) (p Point) {
	p.X = math.Cos(angle)*(point.X-xc) - math.Sin(angle)*(point.Y-yc) + xc
	p.Y = math.Sin(angle)*(point.X-xc) + math.Cos(angle)*(point.Y-yc) + yc
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
	if Distance(mp0, mp1) < Eps {
		panic("mirror line is point")
	}

	A, B, C := Line(mp0, mp1)

	mir := func(x1, y1 float64) Point {
		temp := -2 * (A*x1 + B*y1 + C) / (A*A + B*B)
		return Point{X: temp*A + x1, Y: temp*B + y1}
	}

	ml0 = mir(sp0.X, sp0.Y)
	ml1 = mir(sp1.X, sp1.Y)
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
	if 100*Eps < math.Abs(v) {
		switch {
		case math.Abs(v) < Eps:
			return CollinearPoints
		case 0 < v:
			return ClockwisePoints
		}
		return CounterClockwisePoints
	}
	return Orientation128(p1, p2, p3)
}

func Orientation128(p1, p2, p3 Point) OrientationPoints {
	const prec = 128

	var (
		x1 = new(big.Float).SetPrec(prec).SetFloat64(p1.X)
		x2 = new(big.Float).SetPrec(prec).SetFloat64(p2.X)
		x3 = new(big.Float).SetPrec(prec).SetFloat64(p3.X)

		y1 = new(big.Float).SetPrec(prec).SetFloat64(p1.Y)
		y2 = new(big.Float).SetPrec(prec).SetFloat64(p2.Y)
		y3 = new(big.Float).SetPrec(prec).SetFloat64(p3.Y)

		y21 = new(big.Float).SetPrec(prec).Sub(y2, y1)
		y32 = new(big.Float).SetPrec(prec).Sub(y3, y2)

		x21 = new(big.Float).SetPrec(prec).Sub(x2, x1)
		x32 = new(big.Float).SetPrec(prec).Sub(x3, x2)

		left  = new(big.Float).SetPrec(prec).Mul(y21, x32)
		right = new(big.Float).SetPrec(prec).Mul(x21, y32)

		s = new(big.Float).SetPrec(prec).Sub(left, right)
	)

	v, _ := s.Float64()

	switch {
	case math.Abs(v) < Eps:
		return CollinearPoints
	case 0 < v:
		return ClockwisePoints
	}
	return CounterClockwisePoints
}

func PointArc(pt Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
) {
	// Point - Point
	if Distance(Arc0, Arc1) < Eps && Distance(Arc1, Arc2) < Eps {
		pi, stA, stB = PointPoint(pt, Arc0)
		stB |= ArcIsPoint
		return
	}
	// Point - Line
	{
		A01, B01, C01 := Line(Arc0, Arc1)
		A12, B12, C12 := Line(Arc1, Arc2)
		if math.Abs(A01-A12) < Eps &&
			math.Abs(B01-B12) < Eps &&
			math.Abs(C01-C12) < Eps {
			pi, stA, stB = PointLine(pt, Arc0, Arc2)
			stB |= ArcIsLine
			return
		}
		if Distance(Arc0, Arc1) < Eps {
			pi, stA, stB = PointLine(pt, Arc0, Arc2)
			stB |= ArcIsLine
			return
		}
		if Distance(Arc1, Arc2) < Eps {
			pi, stA, stB = PointLine(pt, Arc0, Arc2)
			stB |= ArcIsLine
			return
		}
	}
	// Point - Arc

	stA |= ZeroLengthSegment | VerticalSegment | HorizontalSegment

	xc, yc, r := Arc(Arc0, Arc1, Arc2)
	radius := Distance(Point{X: xc, Y: yc}, pt)
	if radius < r-Eps || r+Eps < radius {
		// point is outside of arc
		return
	}
	// point is on arc corner ?
	if Distance(pt, Arc0) < Eps {
		stB |= OnPoint0Segment
	}
	if Distance(pt, Arc2) < Eps {
		stB |= OnPoint1Segment
	}

	// point is on arc ?
	if stB.Not(OnPoint0Segment) && stB.Not(OnPoint1Segment) && AngleBetween(
		Point{X: xc, Y: yc},
		Arc0,
		Arc1,
		Arc2,
		pt,
	) {
		stB |= OnSegment
	}

	return
}

func LineArc(Line0, Line1 Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
) {
	// Point - Arc
	if Distance(Line0, Line1) < Eps {
		return PointArc(Line0, Arc0, Arc1, Arc2)
	}
	// Line - Point
	if Distance(Arc0, Arc1) < Eps && Distance(Arc1, Arc2) < Eps {
		pi, stA, stB = PointLine(Arc0, Line0, Line1)
		stA, stB = stB, stA
		stB |= ArcIsPoint
		return
	}
	// Line - Line
	if Distance(Arc0, Arc1) < Eps {
		pi, stA, stB = LineLine(Line0, Line1, Arc0, Arc2)
		stB |= ArcIsLine
		return
	}
	if Distance(Arc1, Arc2) < Eps {
		pi, stA, stB = LineLine(Line0, Line1, Arc0, Arc2)
		stB |= ArcIsLine
		return
	}
	{
		A01, B01, C01 := Line(Arc0, Arc1)
		A12, B12, C12 := Line(Arc1, Arc2)
		if math.Abs(A01-A12) < Eps &&
			math.Abs(B01-B12) < Eps &&
			math.Abs(C01-C12) < Eps {
			pi, stA, stB = LineLine(Line0, Line1, Arc0, Arc2)
			stB |= ArcIsLine
			return
		}
	}
	// Line - Arc

	for _, c := range [...]struct {
		isTrue   bool
		tiA, tiB State
	}{
		{isTrue: math.Abs(Line0.X-Line1.X) < Eps, tiA: VerticalSegment},
		{isTrue: math.Abs(Line0.Y-Line1.Y) < Eps, tiA: HorizontalSegment},
		{isTrue: Distance(Line0, Arc0) < Eps, tiA: OnPoint0Segment, tiB: OnPoint0Segment},
		{isTrue: Distance(Line0, Arc2) < Eps, tiA: OnPoint0Segment, tiB: OnPoint1Segment},
		{isTrue: Distance(Line1, Arc0) < Eps, tiA: OnPoint1Segment, tiB: OnPoint0Segment},
		{isTrue: Distance(Line1, Arc2) < Eps, tiA: OnPoint1Segment, tiB: OnPoint1Segment},
	} {
		if c.isTrue {
			stA |= c.tiA
			stB |= c.tiB
		}
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
	xc, yc, r := Arc(Arc0, Arc1, Arc2)

	// line may be horizontal, vertical, other
	A, B, C := Line(Line0, Line1)

	// Find intersections
	//	A*x+B*y+C = 0               :   line equations
	//	(x-xc)^2 + (y-yc)^2 = r^2   : circle equations
	var roots []Point
	switch {
	case stA.Has(HorizontalSegment):
		// line is horizontal
		//	A = 0
		//
		//	B*y+C = 0
		//	(x-xc)^2 + (y-yc)^2 = r^2
		//
		//	y = -C/B
		//	(x-xc)^2 = r^2 - (-C/B-yc)^2
		//
		//	x = +/- sqrt(r^2 - (-C/B-yc)^2) + xc
		D := pow.E2(r) - pow.E2(-C/B-yc)
		switch {
		case D < -Eps:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			roots = append(roots, Point{X: +xc, Y: Line0.Y})
		default:
			// 0 < D
			roots = append(roots,
				Point{X: +math.Sqrt(D) + xc, Y: Line0.Y},
				Point{X: -math.Sqrt(D) + xc, Y: Line0.Y},
			)
		}

	case stA.Has(VerticalSegment):
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
		case D < -Eps:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			roots = append(roots, Point{X: Line0.X, Y: +yc})
		default:
			// 0 < D
			roots = append(roots,
				Point{X: Line0.X, Y: +math.Sqrt(D) + yc},
				Point{X: Line0.X, Y: -math.Sqrt(D) + yc},
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
		case D < -Eps:
			// no intersection
		case D < Eps:
			// D == 0
			// have one root
			y := -b / (2.0 * a)
			x := -(B*y + C) * 1 / A
			roots = append(roots, Point{X: x, Y: y})
		default:
			// 0 < D
			{
				y := (-b + math.Sqrt(D)) / (2.0 * a)
				x := -(B*y + C) * 1 / A
				roots = append(roots, Point{X: x, Y: y})
			}
			{
				y := (-b - math.Sqrt(D)) / (2.0 * a)
				x := -(B*y + C) * 1 / A
				roots = append(roots, Point{X: x, Y: y})
			}
		}
	}

	for _, root := range roots {
		_, _, stBa := PointLine(root, Line0, Line1)
		_, _, stBb := PointArc(root, Arc0, Arc1, Arc2)

		added := false

		if stBa.Has(OnSegment) &&
			(stBb.Has(OnSegment) || stBb.Has(OnPoint0Segment) || stBb.Has(OnPoint1Segment)) {
			added = true
			stA |= OnSegment
		}

		if stBb.Has(OnSegment) &&
			(stBa.Has(OnSegment) || stBa.Has(OnPoint0Segment) || stBa.Has(OnPoint1Segment)) {
			added = true
			stB |= OnSegment
		}

		if !added {
			continue
		}

		pi = append(pi, root)

		for _, c := range [...]struct {
			isTrue   bool
			tiA, tiB State
		}{
			{
				isTrue: Distance(Line0, root) < Eps &&
					(stBa.Has(OnSegment) || stBa.Has(OnPoint0Segment) || stBa.Has(OnPoint1Segment)),
				tiA: OnPoint0Segment,
			},
			{
				isTrue: Distance(Line1, root) < Eps &&
					(stBa.Has(OnSegment) || stBa.Has(OnPoint0Segment) || stBa.Has(OnPoint1Segment)),
				tiA: OnPoint1Segment,
			},
			{
				isTrue: Distance(Arc0, root) < Eps &&
					(stBb.Has(OnSegment) || stBb.Has(OnPoint0Segment) || stBb.Has(OnPoint1Segment)),
				tiB: OnPoint0Segment,
			},
			{
				isTrue: Distance(Arc2, root) < Eps &&
					(stBb.Has(OnSegment) || stBb.Has(OnPoint0Segment) || stBb.Has(OnPoint1Segment)),
				tiB: OnPoint1Segment,
			},
		} {
			if c.isTrue {
				stA |= c.tiA
				stB |= c.tiB
			}
		}
	}

	return
}

// ArcSplit return points of arcs with middle point if pi is empty or
// slice of arcs.
//	DO NOT CHECKED POINT ON ARC
func ArcSplitByPoint(Arc0, Arc1, Arc2 Point, pi ...Point) (res [][3]Point, err error) {
	switch Orientation(Arc0, Arc1, Arc2) {
	case CollinearPoints:
		panic("collinear")
	case ClockwisePoints:
		res, err = ArcSplitByPoint(Arc2, Arc1, Arc0, pi...)
		for i := range res {
			res[i][0], res[i][2] = res[i][2], res[i][0]
		}
		return
	}
	// CounterClockwisePoints

	for _, c := range [...]struct {
		isTrue bool
	}{
		{isTrue: Distance(Arc0, Arc1) < Eps},
		{isTrue: Distance(Arc1, Arc2) < Eps},
		{isTrue: Distance(Arc0, Arc2) < Eps},
	} {
		if c.isTrue {
			err = fmt.Errorf("invalid points of arc")
			return
		}
	}

	// remove points on corners or same points
againRemove:
	for i, p := range pi {
		for _, c := range [...]struct {
			isTrue bool
		}{
			{isTrue: Distance(Arc0, p) < Eps},
			{isTrue: Distance(Arc2, p) < Eps},
		} {
			if c.isTrue {
				pi = append(pi[:i], pi[i+1:]...)
				goto againRemove
			}
		}
		for j := range pi {
			if i < j && Distance(pi[i], pi[j]) < Eps {
				pi = append(pi[:i], pi[i+1:]...)
				goto againRemove
			}
		}
	}

	// parameter of arc
	xc, yc, r := Arc(Arc0, Arc1, Arc2)

	// angle for rotate
	angle0 := -math.Atan2(Arc0.Y-yc, Arc0.X-xc) - math.Pi + 0.01

	// rotate
	ps := []Point{
		Rotate(xc, yc, +angle0, Arc0),
		Rotate(xc, yc, +angle0, Arc2),
	}
	for i := range pi {
		ps = append(ps, Rotate(xc, yc, +angle0, pi[i]))
	}

	// points angles
	var b []float64
	for i := range ps {
		b = append(b, math.Atan2(ps[i].Y-yc, ps[i].X-xc))
	}
	sort.Float64s(b)

	// remove same angles
again:
	for i := 1; i < len(b); i++ {
		if math.Abs(b[i]-b[i-1]) < Eps {
			b = append(b[:i-1], b[i:]...)
			goto again
		}
	}

	// add middle angles
	if len(pi) == 0 {
		for _, f := range []float64{0.25, 0.5, 0.75} {
			b = append(b, b[0]+f*(b[1]-b[0]))
		}
	} else {
		for i, size := 0, len(b)-1; i < size; i++ {
			b = append(b, b[i]+0.5*(b[i+1]-b[i]))
		}
	}
	sort.Float64s(b)

	if b[0] < -math.Pi {
		panic(fmt.Errorf("debug: %v", b))
	}

	ps = []Point{}
	for _, angle := range b {
		p := Point{
			X: xc + r*math.Cos(angle-angle0),
			Y: yc + r*math.Sin(angle-angle0),
		}
		ps = append(ps, p)
	}

	// prepare arcs
	// 0-1-2=3=4-5-6
	// len=7 arcs=3
	// len=5 arcs=2
	// len=3 arcs=1
	for i := 0; i <= (len(ps)-1)/2; i += 2 {
		res = append(res, [3]Point{ps[i], ps[i+1], ps[i+2]})
	}

	return
}

// linear equations solving:
//	a11*x + a12*y = b1
//	a21*x + a22*y = b2
func Linear(
	a11, a12, b1 float64,
	a21, a22, b2 float64,
) (x, y float64) {
	if math.Abs(a11) < Eps {
		if math.Abs(a12) < Eps {
			panic("cannot solve linear equations")
		}
		// swap parameters
		a11, a12 = a12, a11
		a21, a22 = a22, a21
		defer func() {
			x, y = y, x
		}()
	}
	// float64 algoritm
	y = (b2 - a21/a11*b1) / (a22 - a21/a11*a12)
	x = (b1 - a12*y) * 1 / a11
	return
}

func Arc(Arc0, Arc1, Arc2 Point) (xc, yc, r float64) {
	var (
		x1, x2, x3 = Arc0.X, Arc1.X, Arc2.X
		y1, y2, y3 = Arc0.Y, Arc1.Y, Arc2.Y
		a11        = 2 * (x1 - x2)
		a12        = 2 * (y1 - y2)
		a21        = 2 * (x1 - x3)
		a22        = 2 * (y1 - y3)
		b1         = (pow.E2(x1) - pow.E2(x2)) + (pow.E2(y1) - pow.E2(y2))
		b2         = (pow.E2(x1) - pow.E2(x3)) + (pow.E2(y1) - pow.E2(y3))
	)
	xc, yc = Linear(a11, a12, b1, a21, a22, b2)

	//	(xi-xc)^2+(yi-yc)^2 = R^2
	r1 := math.Hypot (x1-xc, y1-yc)
	r2 := math.Hypot (x2-xc, y2-yc)
	r3 := math.Hypot (x3-xc, y3-yc)
	r = (r1 + r2 + r3) / 3.0
	// find angles
	return
}

// AngleBetween return true for angle case from <= a <= to
func AngleBetween(center, from, mid, to, a Point) (res bool) {
	switch Orientation(from, mid, to) {
	case CollinearPoints:
		panic("collinear")
	case ClockwisePoints:
		return AngleBetween(center, to, mid, from, a)
	}
	// CounterClockwisePoints

	ps := []Point{from, mid, to, a}
	for i := range ps {
		ps[i] = Point{X: ps[i].X - center.X, Y: ps[i].Y - center.Y}
	}

	// angle for rotate
	angle0 := -math.Atan2(ps[0].Y, ps[0].X) - math.Pi + 0.01

	// rotate
	for i := range ps {
		ps[i] = Rotate(0, 0, +angle0, ps[i])
	}

	// points angles
	var b []float64
	for i := range ps {
		b = append(b, math.Atan2(ps[i].Y, ps[i].X))
	}

	if b[0] < -math.Pi {
		panic(fmt.Errorf("debug : %v", b))
	}

	if b[0] < b[3] && b[3] < b[2] {
		return true
	}

	return false
}

func TriangleSplitByPoint(
	pt Point,
	tr0, tr1, tr2 Point,
) (
	res [][3]Point,
	lineIntersect int,
	err error,
) {
	// check valid triangle
	for _, c := range [...]struct {
		isTrue bool
	}{
		{isTrue: Distance(tr0, tr1) < Eps},
		{isTrue: Distance(tr1, tr2) < Eps},
		{isTrue: Distance(tr0, tr2) < Eps},
	} {
		if c.isTrue {
			err = fmt.Errorf("invalid points of triangle")
			return
		}
	}
	// point in triangle box ?
	{
		var (
			xmax = -math.MaxFloat64
			ymax = -math.MaxFloat64
			xmin = +math.MaxFloat64
			ymin = +math.MaxFloat64
		)
		for _, tr := range []Point{tr0, tr1, tr2} {
			xmax = math.Max(xmax, tr.X)
			ymax = math.Max(ymax, tr.Y)
			xmin = math.Min(xmin, tr.X)
			ymin = math.Min(ymin, tr.Y)
		}
		if pt.X < xmin || xmax < pt.X || pt.Y < ymin || ymax < pt.Y {
			// point outside triangle
			return
		}
	}
	// point on corner ?
	for _, c := range [...]struct {
		isTrue bool
	}{
		{isTrue: Distance(tr0, pt) < Eps},
		{isTrue: Distance(tr1, pt) < Eps},
		{isTrue: Distance(tr2, pt) < Eps},
	} {
		if c.isTrue {
			// point on corner
			// no need a split
			return
		}
	}
	// point on the side ?
	for _, line := range []struct {
		Line  [2]Point
		Free  Point
		state int
	}{
		{
			// tr0 --- pt --- tr1  //
			//   \           /     //
			//    \         /      //
			//     \       /       //
			//      \ tr2 /        //
			Line:  [2]Point{tr0, tr1},
			Free:  tr2,
			state: 0,
		},
		{
			// tr0 ---------- tr1  //
			//   \           /     //
			//    \         pt     //
			//     \       /       //
			//      \ tr2 /        //
			Line:  [2]Point{tr1, tr2},
			Free:  tr0,
			state: 1,
		},
		{
			// tr0 ---------- tr1  //
			//   \           /     //
			//    pt        /      //
			//     \       /       //
			//      \ tr2 /        //
			Line:  [2]Point{tr2, tr0},
			Free:  tr1,
			state: 2,
		},
	} {
		_, _, stBl := PointLine(pt, line.Line[0], line.Line[1])
		if !stBl.Has(OnSegment) {
			// point is outside side
			continue
		}
		// point on side
		switch Orientation(tr0, tr1, tr2) {
		case ClockwisePoints:
			res = [][3]Point{
				{line.Line[0], pt, line.Free},
				{pt, line.Line[1], line.Free},
			}
		case CounterClockwisePoints:
			res = [][3]Point{
				{line.Free, pt, line.Line[0]},
				{line.Free, line.Line[1], pt},
			}
		default:
			panic("strange situation")
		}
		lineIntersect = line.state
		return
	}

	// point in body ?
	orient := [3]OrientationPoints{
		Orientation(tr0, pt, tr1),
		Orientation(tr1, pt, tr2),
		Orientation(tr2, pt, tr0),
	}
	if orient[0] != orient[1] ||
		orient[1] != orient[2] ||
		orient[0] != orient[2] {
		// point is outside triangle
		return
	}
	// point inside triangle
	res = [][3]Point{
		{tr0, tr1, pt},
		{tr1, tr2, pt},
		{tr2, tr0, pt},
	}
	return
}

func PointInCircle(point Point, circle [3]Point) bool {
	xc, yc, r := Arc(circle[0], circle[1], circle[2])
	return Distance(Point{xc, yc}, point)+Eps < r
}

// ConvexHull return chain of convex points
func ConvexHull(points []Point) (chain []Point) {
	if len(points) < 3 {
		// points slice is small
		return
	}

	// copy of points
	{
		c := make([]Point, len(points))
		copy(c, points)
		points = c
	}
	// sorting
	sort.Slice(points, func(i, j int) bool {
		if points[i].Y == points[j].Y {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	// lower hull
	var hull []Point
	for _, point := range points {
		for 2 <= len(hull) && (Orientation(hull[len(hull)-2], hull[len(hull)-1], point) == CollinearPoints ||
			Orientation(hull[len(hull)-2], hull[len(hull)-1], point) == ClockwisePoints) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}
	chain = append(chain, hull...)

	// upper hull
	hull = []Point{}
	for i := len(points) - 1; 0 <= i; i-- {
		point := points[i]
		for 2 <= len(hull) && (Orientation(hull[len(hull)-2], hull[len(hull)-1], point) == CollinearPoints ||
			Orientation(hull[len(hull)-2], hull[len(hull)-1], point) == ClockwisePoints) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, point)
	}

	// merge hulls
	if 0 < len(chain) && 0 < len(hull) && Distance(chain[len(chain)-1], hull[0]) < Eps {
		hull = hull[1:]
	}
	if 0 < len(chain) && 0 < len(hull) && Distance(chain[0], hull[len(hull)-1]) < Eps {
		hull = hull[:len(hull)-1]
	}
	chain = append(chain, hull...)

	return
}
