# gog
golang geometry library between point and segments
```

package gog // import "github.com/Konstantin8105/gog"


CONSTANTS

const (
	CollinearPoints        OrientationPoints = -1
	ClockwisePoints                          = 0
	CounterClockwisePoints                   = 1
)

VARIABLES

var (
	// FindRayIntersection is global variable for switch off finding
	// intersection point on segments ray
	FindRayIntersection bool = true

	// Eps is epsilon - precision of intersection
	Eps float64 = 1e-6
)

FUNCTIONS

func AngleBetween(from, a, to float64) (res bool)
    AngleBetween return true for angle case from <= a <= to

func ArcLineAnalisys(Line0, Line1 Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	st State,
)
func ArcSplit(Arc0, Arc1, Arc2 Point) (res [2][3]Point, err error)
    ArcSplit return points of 2 arcs

func Check(pps ...Point) error
    Check - check input data

func Distance(p0, p1 Point) float64
    Distance between two points

func LinePointDistance(
	p0, p1 Point,
	pc Point,
) (distance float64)
    LinePointDistance return distance between line and point

    Equation of line:

        (y2-y1)*(x-x1) = (x2-x1)(y-y1)
        dy*(x-x1) = dx*(y-y1)
        dy*x-dy*x1-dx*y+dx*y1 = 0
        Ax+By+C = 0
        A = dy
        B = -dx
        C = -dy*x1+dx*y1

    Distance from point (xm,ym) to line:

        d = |(A*xm+B*ym+C)/sqrt(A^2+B^2)|

func SegmentAnalisys(
	pa0, pa1 Point,
	pb0, pb1 Point,
) (
	pi []Point,
	st State,
)
    SegmentAnalisys return analisys of two segments

    Design of segments:

                                                    //
        <-- rb00 -- pb0*==========*pb1 -- rb11 -->  // Segment B
                                                    //
        <-- ra00 -- pa0*==========*pa1 -- ra11 -->  // Segment A
        {   ray   }{      segment     }{   ray   }  //
                                                    //

    Input data:

        ipa0, ipa1 - point indexes of segment A
        ipb0, ipb1 - point indexes of segment B
        pps      - pointer of point slice

    Output data:

        pi - intersection point
        st - states of analisys

    Reference:

        [1]  https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection


TYPES

type OrientationPoints int8

func Orientation(p1, p2, p3 Point) OrientationPoints

type Point struct {
	X, Y float64
}
    Point is store of point coordinates

func MirrorLine(
	sp0, sp1 Point,
	mp0, mp1 Point,
) (
	ml0, ml1 Point,
	err error,
)
    MirrorLine return intersection point and second mirrored point from mirror
    line (mp0-mp1) and ray (sp0-sp1)

func Rotate(angle float64, point Point) (p Point)
    Rotate point about (0,0) on angle

func (p Point) String() string
    String is implementation of Stringer implementation for formating output

type State int64
    State is result of intersection

const (
	VerticalSegmentA State // vertical segment A
	VerticalSegmentB       // vertical segment B

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

	Arc01indentical   // arc points 0 and 1 zero length
	Arc12indentical   // arc points 1 and 2 zero length
	Arc02indentical   // arc points 0 and 2 zero length
	ArcIsLine         // wrong arc is line
	ArcIsPoint        // wrong arc is point
	LineFromArcCenter // line intersect center of arc
	LineOutside       // line is outside of arc

)
func (s State) Has(si State) bool
    Has is mean s-State has si-State

func (s State) Not(si State) bool
    Not mean s-State have not si-State

func (s State) String() string
    String is implementation of Stringer implementation for formating output


```
