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
const (
	Boundary  = -1
	Removed   = -2
	Undefined = -3
	Fixed     = 100
	Movable   = 200
)

VARIABLES

var (
	// FindRayIntersection is global variable for switch off finding
	// intersection point on segments ray
	FindRayIntersection bool = true

	// Eps is epsilon - precision of intersection
	Eps float64 = 1e-12
)
var Debug = false

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

func Distance128(p0, p1 Point) float64
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
func PointInCircle(point Point, circle [3]Point) bool
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
func TriangleSplitByPoint(
	pt Point,
	tr0, tr1, tr2 Point,
) (
	res [][3]Point,
	lineIntersect int,
	err error,
)

TYPES

type Mesh struct {
	Points    []int // tags for points
	Triangles []Triangle
	// Has unexported fields.
}

func New(model Model) (mesh *Mesh, err error)
    New triangulation created by model

func (mesh *Mesh) AddLine(p1, p2 Point, tag int) (err error)
    AddLine is add line in triangulation with tag

func (mesh *Mesh) AddPoint(p Point, tag int) (err error)
    AddPoint is add points with tag

func (mesh Mesh) Check() (err error)
    Check triangulation on point, line, triangle rules

func (mesh *Mesh) Clockwise()
    Clockwise change all triangles to clockwise orientation

func (mesh *Mesh) Delanay() (err error)
    TODO delanay only for some triangles, if list empty then for all triangles

func (mesh *Mesh) GetMaterials(ps ...Point) (materials []int, err error)
    GetMaterials return materials for each point

func (mesh *Mesh) Materials() (err error)
    Materials indentify all triangles splitted by lines, only if points sliceis
    empty. If points slice is not empty, then return material mark number for
    each point

func (mesh *Mesh) Smooth()
    Smooth move all movable point by average distance

func (mesh *Mesh) Split(d float64) (err error)
    Split all triangles edge on distance `d`

type Model struct {
	Points    []Point  // Points is slice of points
	Lines     [][3]int // Lines store 2 index of Points and last for tag
	Arcs      [][4]int // Arcs store 3 index of Points and last for tag
	Triangles [][4]int // Triangles store 3 index of Points and last for tag/material
}
    Model of points, lines, arcs for prepare of triangulation

func (m *Model) AddArc(start, middle, end Point, tag int)
    AddArc add arc into model with specific tag

func (m *Model) AddCircle(xc, yc, r float64, tag int)
    AddCircle add arcs based on circle geometry into model with specific tag

func (m *Model) AddLine(start, end Point, tag int)
    AddLine add line into model with specific tag

func (m *Model) AddPoint(p Point) (index int)
    AddPoint return index in model slice point

func (m *Model) AddTriangle(start, middle, end Point, tag int)
    AddTriangle add triangle into model with specific tag/material

func (m *Model) ArcsToLines()
    ArcsToLines convert arc to lines

func (m *Model) ConvexHullTriangles()
    ConvexHullTriangles add triangles of model convex hull

func (m Model) Dxf() string
    Dxf return string in dxf drawing format
    https://images.autodesk.com/adsk/files/autocad_2012_pdf_dxf-reference_enu.pdf

func (model *Model) Get(mesh *Mesh, lines bool)

func (m *Model) Intersection()
    Intersection change model with finding all model intersections

func (m *Model) Merge()

func (m Model) MinPointDistance() (distance float64)
    MinPointDistance return minimal between 2 points

func (m *Model) Move(dx, dy float64)

func (m *Model) RemoveEmptyPoints()

func (m *Model) RemovePoint()

func (m *Model) Rotate(xc, yc, angle float64)

func (m *Model) Split(d float64)

func (m Model) String() string
    String return a stantard model view

type OrientationPoints int8

func Orientation(p1, p2, p3 Point) OrientationPoints

func Orientation128(p1, p2, p3 Point) OrientationPoints

type Point struct {
	X, Y float64
}
    Point is store of point coordinates

func ConvexHull(points []Point) (chain []Point)
    ConvexHull return chain of convex points

func MiddlePoint(p0, p1 Point) Point

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

type Triangle struct {
	// Has unexported fields.
}
    Triangle is data structure "Nodes, ribs и triangles" created by book
    "Algoritm building and analyse triangulation", A.B.Skvorcov

        Scketch:
        +------------------------------------+
        |              tr[0]                 |
        |  nodes[0]    ribs[0]      nodes[1] |
        | o------------------------o         |
        |  \                      /          |
        |   \                    /           |
        |    \                  /            |
        |     \                /             |
        |      \              /              |
        |       \            /  ribs[1]      |
        |        \          /   tr[1]        |
        |  ribs[2]\        /                 |
        |  tr[2]   \      /                  |
        |           \    /                   |
        |            \  /                    |
        |             \/                     |
        |              o  nodes[2]           |
        +------------------------------------+


```
