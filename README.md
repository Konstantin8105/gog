# gog
golang geometry library between point and segments
```


package gog // import "github.com/Konstantin8105/gog"


FUNCTIONS

func Check(pps *[]Point) error
    Check - check input data

func Distance(p0, p1 Point) float64
    Distance between two points

func LinePointDistance(
	ip0, ip1 int,
	ipc int,
	pps *[]Point,
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
	ipa0, ipa1 int,
	ipb0, ipb1 int,
	pps *[]Point,
) (
	pi Point,
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

type Point struct {
	X, Y float64
}
    Point is store of point coordinates

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

)
func (s State) Has(si State) bool
    Has is mean s-State has si-State

func (s State) Not(si State) bool
    Not mean s-State have not si-State

func (s State) String() string
    String is implementation of Stringer implementation for formating output


```
