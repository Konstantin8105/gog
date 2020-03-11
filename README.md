# gog
golang geometry library between point and segments
```


package gog // import "github.com/Konstantin8105/gog"


FUNCTIONS

func Distance(p0, p1 Point) float64
    Distance between two points

func Intersection(b0, b1 Segment, ps *[]Point) (
	p Point,
	t State,
)

TYPES

type Point struct {
	X, Y float64
}

type Segment struct {
	P0, P1 int // indexes of point
}
    Segment is part of line

    Design of segment:

        -- P00 -- P0*==========*P1 -- P11 --
        {  ray  }   {  segment }   {  ray  }

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
