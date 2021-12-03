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

func AngleBetween(center, from, mid, to, a Point) (res bool)
    AngleBetween return true for angle case from <= a <= to

func Arc(Arc0, Arc1, Arc2 Point) (xc, yc, r float64)
func ArcSplitByPoint(Arc0, Arc1, Arc2 Point, pi ...Point) (res [][3]Point, err error)
    ArcSplit return points of arcs with middle point if pi is empty or slice of
    arcs.

        DO NOT CHECKED POINT ON ARC

func Check(pps ...Point) error
    Check - check input data

func Distance(p0, p1 Point) float64
    Distance between two points

func Line(p0, p1 Point) (A, B, C float64)
    line parameters

        Ax+By+C = 0

func LineArc(Line0, Line1 Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
)
func LineLine(
	pa0, pa1 Point,
	pb0, pb1 Point,
) (
	pi []Point,
	stA, stB State,
)
    LineLine return analisys of two segments

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

func Linear(
	a11, a12, b1 float64,
	a21, a22, b2 float64,
) (x, y float64)
    linear equations solving:

        a11*x + a12*y = b1
        a21*x + a22*y = b2

func PointArc(pt Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
)
func PointLine(
	pt Point,
	pb0, pb1 Point,
) (
	pi []Point,
	stA, stB State,
)
func PointLineDistance(
	pc Point,
	p0, p1 Point,
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

func PointPoint(
	pt0, pt1 Point,
) (
	pi []Point,
	stA, stB State,
)

TYPES

type Model struct {
	Points []Point  // Points is slice of points
	Lines  [][3]int // Lines store 2 index of Points and last for tag
	Arcs   [][4]int // Arcs store 3 index of Points and last for tag
}
    Model of points, lines, arcs for prepare of triangulation

func (m *Model) AddArc(start, middle, end Point, tag int)
    AddArc add arc into model with specific tag

func (m *Model) AddCircle(xc, yc, r float64, tag int, isHole bool)
    AddCircle add arcs based on circle geometry into model with specific tag

func (m *Model) AddLine(start, end Point, tag int)
    AddLine add line into model with specific tag

func (m *Model) AddPoint(p Point) (index int)
    AddPoint return index in model slice point

func (m *Model) Intersection()
    Intersection change model with finding all model intersections

func (m Model) MinPointDistance() (distance float64)
    MinPointDistance return minimal between 2 points

func (m *Model) RemoveEmptyPoints()

func (m *Model) RemovePoint()

func (m *Model) Split()

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

func Rotate(xc, yc, angle float64, point Point) (p Point)
    Rotate point about (xc,yc) on angle

func (p Point) String() string
    String is implementation of Stringer implementation for formating output

type State int64
    State is result of intersection

const (
	VerticalSegment State // vertical segment

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

)
func (s State) Has(si State) bool
    Has is mean s-State has si-State

func (s State) Not(si State) bool
    Not mean s-State have not si-State

func (s State) String() string
    String is implementation of Stringer implementation for formating output


```
