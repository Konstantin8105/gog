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
	// Boundary edge
	Boundary = -1

	// Removed element
	Removed = -2

	// Undefined state only for not valid algorithm
	Undefined = -3
	Fixed     = 100
	Movable   = 200
)
const Eps3D = 1e-5
    Eps3D is default epsilon for 3D operations


VARIABLES

var (
	// ErrorDivZero is typical result with not acceptable solving
	ErrorDivZero = fmt.Errorf("div value is too small")

	// ErrorNotValidSystem is typical return only if system of
	// linear equation have not valid data
	ErrorNotValidSystem = fmt.Errorf("not valid system")
)
var (
	// Debug only for debugging
	Debug = false
	// Log only for minimal logging
	Log = false
)
var (
	// Eps is epsilon - precision of intersection
	Eps = 1e-10
)

FUNCTIONS

func AngleBetween(center, from, mid, to, a Point) (res bool)
    AngleBetween return true for angle case from <= a <= to

func Arc(Arc0, Arc1, Arc2 Point) (xc, yc, r float64)
    Arc return parameters of circle

func ArcSplitByPoint(Arc0, Arc1, Arc2 Point, pi ...Point) (res [][3]Point, err error)
    ArcSplitByPoint return points of arcs with middle point if pi is empty or
    slice of arcs.

        DO NOT CHECKED POINT ON ARC

func Area(
	tr0, tr1, tr2 Point,
) float64
    Area return area of triangle

func BorderIntersection(ps1, ps2 []Point3d) (intersect bool)
    BorderIntersection return true only if Borders are intersect

func Check(pps ...Point) error
    Check - check input data

func Distance(p0, p1 Point) float64
    Distance between two points

func Distance128(p0, p1 Point) float64
    Distance128 is distance between 2 points with 128-bit precisions

func Distance3d(p0, p1 Point3d) float64
    Distance3d is distance between 2 points in 3D

func IsParallelLine3d(
	a0, a1 Point3d,
	b0, b1 Point3d,
) (
	parallel bool,
)
    IsParallelLines3d return true, if lines are parallel

func Line(p0, p1 Point) (A, B, C float64)
    Line parameters by formula: Ax+By+C = 0

func LineArc(Line0, Line1 Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
)
    LineArc return state and intersections points between line and arc

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

func LineLine3d(
	a0, a1 Point3d,
	b0, b1 Point3d,
) (
	ratioA, ratioB float64,
	intersect bool,
)
    LineLine3d return intersection of two points. Point on line corner ignored

func Linear(
	a11, a12, b1 float64,
	a21, a22, b2 float64,
) (x, y float64, err error)
    Linear equations solving:

        a11*x + a12*y = b1
        a21*x + a22*y = b2

func Plane(
	p1, p2, p3 Point3d,
) (
	A, B, C, D float64,
)
    Plane equation `A*x+B*y+C*z+D=0`

func PlaneAverage(
	ps []Point3d,
) (
	A, B, C, D float64,
)
    PlaneAverage return parameters of average plane for points

func PointArc(pt Point, Arc0, Arc1, Arc2 Point) (
	pi []Point,
	stA, stB State,
)
    PointArc return state and intersections points between point and arc

func PointInCircle(point Point, circle [3]Point) bool
    PointInCircle return true only if point inside circle based on 3 circles
    points

func PointLine(
	pt Point,
	pb0, pb1 Point,
) (
	pi []Point,
	stA, stB State,
)
    PointLine return states between point and line.

func PointLine3d(
	p Point3d,
	l0, l1 Point3d,
) (
	intersect bool,
)
    PointLine3d return true only if point located on line segment

func PointLineDistance(
	pc Point,
	p0, p1 Point,
) (distance float64)
    PointLineDistance return distance between line and point.

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

func PointOnPlane3d(
	A, B, C, D float64,
	p Point3d,
) (
	on bool,
)
    PointOnPlane3d return true if all points on plane

func PointPoint(
	pt0, pt1 Point,
) (
	pi []Point,
	stA, stB State,
)
    PointPoint return states between two points.

func PointPoint3d(
	p0 Point3d,
	p1 Point3d,
) (
	intersect bool,
)
    PointPoint3d return true only if points have same coordinate

func PointTriangle3d(
	p Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
)
    PointTriangle3d return true only if point located inside triangle but do not
    check point on triangle edge

func SamePoints(p0, p1 Point) bool
    SamePoints return true only if point on very distance or with same
    coordinates

func SamePoints3d(p0, p1 Point3d) bool
    SamePoints3d return true only if point on very distance or with same
    coordinates

func TriangleSplitByPoint(
	pt Point,
	tr0, tr1, tr2 Point,
) (
	res [][3]Point,
	lineIntersect int,
	err error,
)
    TriangleSplitByPoint split triangle on triangles only if point inside
    triangle or on triangle edge

func ZeroLine3d(
	l0, l1 Point3d,
) (
	zero bool,
)
    ZeroLine3d return true only if lenght of line segment is zero

func ZeroTriangle3d(
	t0, t1, t2 Point3d,
) (
	zero bool,
)
    ZeroTriangle3d return true only if triangle have zero area


TYPES

type Mesh struct {
	Points    []int    // tags for points
	Triangles [][3]int // indexes of near triangles

	// Has unexported fields.
}
    Mesh is based structure of triangulation. Triangle is data structure
    "Nodes, ribs и triangles" created by book "Algoritm building and analyse
    triangulation", A.B.Skvorcov

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

func New(model Model) (mesh *Mesh, err error)
    New triangulation created by model

func (mesh *Mesh) AddLine(inp1, inp2 Point) (err error)
    AddLine is add line in triangulation with tag

func (mesh *Mesh) AddPoint(p Point, tag int, triIndexes ...int) (idp int, err error)
    AddPoint is add points with tag

func (mesh Mesh) Check() (err error)
    Check triangulation on point, line, triangle rules

func (mesh *Mesh) Clockwise()
    Clockwise change all triangles to clockwise orientation

func (mesh *Mesh) Delanay(triIndexes ...int) (err error)
    TODO delanay only for some triangles, if list empty then for all triangles

func (mesh *Mesh) GetMaterials(ps ...Point) (materials []int, err error)
    GetMaterials return materials for each point

func (mesh *Mesh) Materials() (err error)
    Materials indentify all triangles splitted by lines, only if points sliceis
    empty. If points slice is not empty, then return material mark number for
    each point

func (mesh *Mesh) RemoveMaterials(ps ...Point) (err error)
    RemoveMaterials remove material by specific points

func (mesh *Mesh) Smooth(pts ...int) (err error)
    Smooth move all movable point by average distance

func (mesh *Mesh) Split(d float64) (err error)
    Split all triangles edge on distance `d`

type Model struct {
	Points    []Point  // Points is slice of points
	Lines     [][3]int // Lines store 2 index of Points and last for tag
	Arcs      [][4]int // Arcs store 3 index of Points and last for tag
	Triangles [][4]int // Triangles store 3 index of Points and last for tag/material
	Quadrs    [][5]int // Rectanges store 4 index of Points and last for tag/material
}
    Model of points, lines, arcs for prepare of triangulation

func (m *Model) AddArc(start, middle, end Point, tag int)
    AddArc add arc into model with specific tag

func (m *Model) AddCircle(xc, yc, r float64, tag int)
    AddCircle add arcs based on circle geometry into model with specific tag

func (m *Model) AddLine(start, end Point, tag int)
    AddLine add line into model with specific tag

func (m *Model) AddModel(from Model)
    AddModel inject model into model

func (m *Model) AddMultiline(tag int, ps ...Point)
    AddMultiline add many lines with specific tag

func (m *Model) AddPoint(p Point) (index int)
    AddPoint return index in model slice point

func (m *Model) AddTriangle(start, middle, end Point, tag int)
    AddTriangle add triangle into model with specific tag/material

func (m *Model) ArcsToLines()
    ArcsToLines convert arc to lines

func (m *Model) Combine(factorOneLine float64) (err error)
    Combine triangles to quadr with same tag

        factorOneLine from 1 to 2/sqrt(2) = 1.41

    Recommendation value is 1.05

func (m *Model) ConvexHullTriangles()
    ConvexHullTriangles add triangles of model convex hull

func (src Model) Copy() (dst Model)
    Copy return copy of Model

func (m Model) Dxf() string
    Dxf return string in dxf drawing format
    https://images.autodesk.com/adsk/files/autocad_2012_pdf_dxf-reference_enu.pdf

func (model *Model) Get(mesh *Mesh)
    Get add into Model all triangles from Mesh Recommendation after `Get` :
    model.Intersection()

func (m *Model) Intersection()
    Intersection change model with finding all model intersections

func (m Model) JSON() (_ string, err error)
    JSON convert model in JSON format

func (to *Model) Merge(from Model)
    Merge `from` model to `to` model

func (m Model) MinPointDistance() (distance float64)
    MinPointDistance return minimal between 2 points

func (m *Model) Move(dx, dy float64)
    Move all points of model

func (m *Model) Read(filename string) (err error)
    Read model from file with filename in JSON format

func (m *Model) RemoveEmptyPoints()
    RemoveEmptyPoints removed point not connected to line, arcs, triangles

func (m *Model) RemovePoint(remove func(p Point) bool)
    RemovePoint removed point in accoding to function `filter`

func (m *Model) Rotate(xc, yc, angle float64)
    Rotate all points of model around point {xc,yc}

func (m *Model) Split(d float64)
    Split all model lines, arcs by distance `d`

func (m Model) String() string
    String return a stantard model view

func (m Model) TagProperty() (length []float64, area []float64)
    TagProperty return length of lines, area of triangles for each tag. Arcs are
    ignored

func (m Model) Write(filename string) (err error)
    Write model into file with filename in JSON format

type OrientationPoints int8
    OrientationPoints is orientation points state

func Orientation(p1, p2, p3 Point) OrientationPoints

func Orientation128(p1, p2, p3 Point) OrientationPoints

type Point struct {
	X, Y float64
}
    Point is store of point coordinates

func ConvexHull(points []Point, withoutCollinearPoints bool) (chain []int, res []Point)
    ConvexHull return chain of convex points

func MiddlePoint(p0, p1 Point) Point
    MiddlePoint calculate middle point precisionally.

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

type Point3d [3]float64
    Point3d is point coordinate in 3D decart system

func BorderPoints(ps ...Point3d) (min, max Point3d)
    BorderPoints return (min..max) points coordinates

func LineTriangle3dI1(
	l0, l1 Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
	pi []Point3d,
)
    LineTriangle3dI1 return intersection points for case if line and triangle is
    not on one plane. line intersect triangle in one point

func LineTriangle3dI2(
	l0, l1 Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
	pi []Point3d,
)
    LineTriangle3dI2 return intersection points if line and triangle located on
    one plane. Line on triangle plane Line is not zero ignore triangle point on
    line

func Mirror3d(plane [3]Point3d, points ...Point3d) (mir []Point3d)
    Mirror3d return mirror points by mirror plane

func PointLineRatio3d(
	l0, l1 Point3d,
	ratio float64,
) (
	p Point3d,
)
    PointLineRatio3d return point in accroding to line ratio

func TriangleTriangle3d(
	a0, a1, a2 Point3d,
	b0, b1, b2 Point3d,
) (
	intersect bool,
	pi []Point3d,
)
    TriangleTriangle3d return intersection points between two triangles.
    do not intersect with egdes

type State int64
    State is result of intersection

const (

	// VerticalSegment return if segment is vertical
	VerticalSegment State

	// HorizontalSegment return if segment is horizontal
	HorizontalSegment

	// ZeroLengthSegment return for zero length segment
	ZeroLengthSegment

	// Parallel is segment A and segment B.
	// Intersection point data is not valid.
	Parallel

	// Collinear return if:
	// Segment A and segment B are collinear.
	// Intersection point data is not valid.
	Collinear

	// OnSegment is intersection point on segment
	OnSegment

	// OnPoint0Segment intersection point on point 0 segment
	OnPoint0Segment

	// OnPoint1Segment intersection point on point 1 segment
	OnPoint1Segment

	// ArcIsLine return only if wrong arc is line
	ArcIsLine

	// ArcIsPoint return only if wrong arc is point
	ArcIsPoint
)
func (s State) Has(si State) bool
    Has is mean s-State has si-State

func (s State) Not(si State) bool
    Not mean s-State have not si-State

func (s State) String() string
    String is implementation of Stringer implementation for formating output

```
