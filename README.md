# gog
golang geometry library between point and segments
```


package gog // import "github.com/Konstantin8105/gog"


FUNCTIONS

func Check(pps *[]Point) error
    Check - check input data

func Distance(p0, p1 Point) float64
    Distance between two points

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

	// property of single segment
	VerticalSegment0 State
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
)
func (s State) Has(si State) bool
    Has is mean s-State has si-State

func (s State) Not(si State) bool
    Not mean s-State have not si-State

func (s State) String() string
    String is implementation of Stringer implementation for formating output


```
