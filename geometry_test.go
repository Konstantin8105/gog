package gog

import (
	"fmt"
	"math"
	"testing"
)

type TestCase struct {
	name string
	ps   []Point
	it   IntersectionType
	pi   Point
}

var tcs = []TestCase{
	{ // 0
		// *1,3,4 //
		// |      //
		// |      //
		// *2     //
		ps: []Point{
			Point{X: 0, Y: 8}, // 1
			Point{X: 0, Y: 2}, // 2
			Point{X: 0, Y: 8}, // 3
			Point{X: 0, Y: 8}, // 4
		},
		it: ZeroLengthBeam1 |
			VerticalBeam0 |
			Point0Beam0onPoint0Beam1 |
			Point0Beam0onPoint1Beam1 |
			HorizontalBeam1 | VerticalBeam1 |
			Collinear,
		pi: Point{X: 0, Y: 8},
	},
	{ // 1
		// *1  *3 //
		// |   |  //
		// |   |  //
		// *2  *4 //
		ps: []Point{
			Point{X: 0, Y: 8}, // 1
			Point{X: 0, Y: 2}, // 2
			Point{X: 0, Y: 8}, // 3
			Point{X: 0, Y: 2}, // 4
		},
		it: VerticalBeam0 | VerticalBeam1 |
			Point0Beam0onPoint0Beam1 | Point1Beam0onPoint1Beam1 |
			Collinear,
		pi: Point{X: 0, Y: 8},
	},
	{ // 2
		// *1  *3 //
		// |   |  //
		// |   |  //
		// *2  *4 //
		ps: []Point{
			Point{X: 2, Y: 8}, // 1
			Point{X: 0, Y: 2}, // 2
			Point{X: 2, Y: 8}, // 3
			Point{X: 0, Y: 2}, // 4
		},
		it: Point0Beam0onPoint0Beam1 | Point1Beam0onPoint1Beam1 |
			Collinear,
		pi: Point{X: 2, Y: 8},
	},
	{ // 3
		// *1  *3 //
		// |   |  //
		// |   |  //
		// *2  *4 //
		ps: []Point{
			Point{X: 0, Y: 8}, // 1
			Point{X: 0, Y: 2}, // 2
			Point{X: 4, Y: 8}, // 3
			Point{X: 4, Y: 2}, // 4
		},
		it: VerticalBeam0 | VerticalBeam1 | Parallel,
	},
	{ // 4
		// *1  //
		// |   //
		// *2  //
		//
		// *3  //
		// |   //
		// *4  //
		ps: []Point{
			Point{X: 2, Y: 8}, // 1
			Point{X: 2, Y: 7}, // 2
			Point{X: 2, Y: 6}, // 3
			Point{X: 2, Y: 5}, // 4
		},
		it: VerticalBeam0 | VerticalBeam1 | Collinear,
	},
	{ // 5
		// *1  //
		// |   //
		// *2  //
		//
		// *3  //
		// |   //
		// *4  //
		ps: []Point{
			Point{X: 5, Y: 5}, // 1
			Point{X: 4, Y: 4}, // 2
			Point{X: 3, Y: 3}, // 3
			Point{X: 2, Y: 2}, // 4
		},
		it: Collinear,
	},
	{ // 6
		// *1  //
		// |   //
		// *2,3//
		// |   //
		// *4  //
		ps: []Point{
			Point{X: 2, Y: 8}, // 1
			Point{X: 2, Y: 6}, // 2
			Point{X: 2, Y: 6}, // 3
			Point{X: 2, Y: 5}, // 4
		},
		it: VerticalBeam0 | VerticalBeam1 |
			Point1Beam0onPoint0Beam1 |
			Collinear,
		pi: Point{X: 2, Y: 6},
	},
	{ // 7
		// *1  //
		// |   //
		// *2,3//
		// |   //
		// *4  //
		ps: []Point{
			Point{X: 5, Y: 5}, // 1
			Point{X: 4, Y: 4}, // 2
			Point{X: 4, Y: 4}, // 3
			Point{X: 2, Y: 2}, // 4
		},
		it: Point1Beam0onPoint0Beam1 |
			Collinear,
		pi: Point{X: 4, Y: 4},
	},
	{ // 8
		// *1,2,3,4  //
		ps: []Point{
			Point{X: 5, Y: 5}, // 1
			Point{X: 5, Y: 5}, // 2
			Point{X: 5, Y: 5}, // 3
			Point{X: 5, Y: 5}, // 4
		},
		it: VerticalBeam0 |
			VerticalBeam1 |
			HorizontalBeam0 |
			HorizontalBeam1 |
			ZeroLengthBeam0 |
			ZeroLengthBeam1 |
			Point0Beam0onPoint0Beam1 |
			Point1Beam0onPoint0Beam1 |
			Point0Beam0onPoint1Beam1 |
			Point1Beam0onPoint1Beam1 |
			Collinear,
		pi: Point{X: 5, Y: 5},
	},
	{ // 9
		//     *2  //
		//    /    //
		//   *3    //
		//  /  \   //
		// *1   *4 //
		ps: []Point{
			Point{X: 1, Y: 1}, // 1
			Point{X: 4, Y: 4}, // 2
			Point{X: 2, Y: 2}, // 3
			Point{X: 5, Y: 0}, // 4
		},
		it: Point0Beam1inBeam0,
		pi: Point{X: 2, Y: 2},
	},
	{ // 10
		// *3   *1  //
		//  \  /    //
		//    X     //
		//  /  \    //
		// *2   *4  //
		ps: []Point{
			Point{X: 1, Y: 1}, // 1
			Point{X: 4, Y: 4}, // 2
			Point{X: 0, Y: 5}, // 3
			Point{X: 5, Y: 0}, // 4
		},
		it: IntersectOnBeam0 | IntersectOnBeam1,
		pi: Point{X: 2.5, Y: 2.5},
	},
	{ // 11
		//     *1  //
		//    /    //
		//   *4    //
		//  /  \   //
		// *2   *3 //
		ps: []Point{
			Point{X: 1, Y: 1}, // 1
			Point{X: 4, Y: 4}, // 2
			Point{X: 5, Y: 0}, // 3
			Point{X: 2, Y: 2}, // 4
		},
		it: Point1Beam1inBeam0,
		pi: Point{X: 2, Y: 2},
	},
	{ // 12
		//     *4  //
		//    /    //
		//   *1    //
		//  /  \   //
		// *3   *2 //
		ps: []Point{
			Point{X: 2, Y: 2}, // 1
			Point{X: 5, Y: 0}, // 2
			Point{X: 4, Y: 4}, // 3
			Point{X: 1, Y: 1}, // 4
		},
		it: Point0Beam0inBeam1,
		pi: Point{X: 2, Y: 2},
	},
	{ // 13
		//     *4  //
		//    /    //
		//   *2    //
		//  /  \   //
		// *3   *1 //
		ps: []Point{
			Point{X: 5, Y: 0}, // 1
			Point{X: 2, Y: 2}, // 2
			Point{X: 4, Y: 4}, // 3
			Point{X: 1, Y: 1}, // 4
		},
		it: Point1Beam0inBeam1,
		pi: Point{X: 2, Y: 2},
	},
	{ // 14
		//      *4 //
		//      |  //
		//      |  //
		//   *2 |  //
		//  /   |  //
		// *1   *3 //
		ps: []Point{
			Point{X: 1, Y: 1}, // 1
			Point{X: 2, Y: 2}, // 2
			Point{X: 5, Y: 0}, // 3
			Point{X: 5, Y: 9}, // 4
		},
		it: VerticalBeam1 | IntersectBeam0Ray11 | IntersectOnBeam1,
		pi: Point{X: 5, Y: 5},
	},
	{ // 15
		//      *4 //
		//      |  //
		//      |  //
		//   *1 |  //
		//  /   |  //
		// *2   *3 //
		ps: []Point{
			Point{X: 2, Y: 2}, // 1
			Point{X: 1, Y: 1}, // 2
			Point{X: 5, Y: 0}, // 3
			Point{X: 5, Y: 9}, // 4
		},
		it: VerticalBeam1 | IntersectBeam0Ray00 | IntersectOnBeam1,
		pi: Point{X: 5, Y: 5},
	},
	{ // 16
		//      *2 //
		//      |  //
		//      |  //
		//   *4 |  //
		//  /   |  //
		// *3   *1 //
		ps: []Point{
			Point{X: 5, Y: 0}, // 1
			Point{X: 5, Y: 9}, // 2
			Point{X: 1, Y: 1}, // 3
			Point{X: 2, Y: 2}, // 4
		},
		it: VerticalBeam0 | IntersectBeam1Ray11 | IntersectOnBeam0,
		pi: Point{X: 5, Y: 5},
	},
	{ // 17
		//      *2 //
		//      |  //
		//      |  //
		//   *3 |  //
		//  /   |  //
		// *4   *1 //
		ps: []Point{
			Point{X: 5, Y: 0}, // 1
			Point{X: 5, Y: 9}, // 2
			Point{X: 2, Y: 2}, // 3
			Point{X: 1, Y: 1}, // 4
		},
		it: VerticalBeam0 | IntersectBeam1Ray00 | IntersectOnBeam0,
		pi: Point{X: 5, Y: 5},
	},
	{ // 18 : Test data - no intersection
		ps: []Point{
			Point{X: 1.098, Y: 0},
			Point{X: -1.5449, Y: 12.53},
			Point{X: 1.2, Y: 2},
			Point{X: 5, Y: 5},
		},
		it: IntersectBeam1Ray00 | IntersectOnBeam0,
		pi: Point{X: 0.7509280607532581, Y: 1.6454695216473094},
	},
	{ // 19 : Test data - no intersection
		ps: []Point{
			Point{X: 5.108, Y: 0},
			Point{X: 8.339, Y: 16.17},
			Point{X: 9, Y: 2},
			Point{X: 5, Y: 5},
		},
		it: IntersectOnBeam0 | IntersectOnBeam1,
		pi: Point{X: 5.9627881085877945, Y: 4.277908918559155},
	},
}

func init() {
	// create copy
	copy := func(t TestCase) TestCase {
		var ts TestCase
		ts.pi = t.pi
		ts.ps = append(ts.ps, t.ps...)
		ts.name = t.name
		ts.it = t.it
		return ts
	}

	// add names
	for i := range tcs {
		tcs[i].name = fmt.Sprintf("%2d", i)
	}

	var size int

	// add test with moving
	size = len(tcs)
	for i := 0; i < size; i++ {
		for _, mv := range []float64{0.0, -100.0, +100.0, math.Pi, -math.Pi} {
			tc := copy(tcs[i])
			// move
			for i := range tc.ps {
				tc.ps[i].X += mv
				tc.ps[i].Y += mv
			}
			tc.pi.X += mv
			tc.pi.Y += mv
			tc.name += fmt.Sprintf("-%5.3f", mv)
			tcs = append(tcs, tc)
		}
	}

	// add test with rotating
	size = len(tcs)
	for i := 0; i < size; i++ {
		tc := copy(tcs[i])
		// move
		for i := range tc.ps {
			tc.ps[i].X *= -1.0
			tc.ps[i].Y *= -1.0
		}
		tc.pi.X *= -1.0
		tc.pi.Y *= -1.0
		tc.name += "-rotate"
		tcs = append(tcs, tc)
	}
}

func Test(t *testing.T) {
	var types [64]int
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			pi, it := Intersection(Beam{P0: 0, P1: 1}, Beam{P0: 2, P1: 3}, &tc.ps)
			if it != tc.it {
				t.Error("Not same types")
				t.Logf("Expected : %30b", tc.it)
				t.Logf("Value    : %30b", it)
				t.Logf("Diff1    : %30b", tc.it&^it)
				t.Logf("%s", tc.it&^it)
				t.Logf("Diff2    : %30b", it&^tc.it)
				t.Logf("%s", it&^tc.it)
			} else {
				t.Logf("Value    : %30b", it)
			}
			// store
			for i := 0; i < len(types); i++ {
				if it&IntersectionType(1<<i) != 0 {
					types[i]++
				}
			}
			if eps := 1e-6; Distance(pi, tc.pi) > eps &&
				math.Abs(pi.X) > eps &&
				math.Abs(pi.Y) > eps {
				t.Errorf("2 point is too far: %#v != %#v\ndistance = %.6e",
					pi,
					tc.pi,
					Distance(pi, tc.pi),
				)
			}
		})
	}

	sum := 0
	amountFail := 0
	for i := range types {
		if i == 0 {
			continue
		}
		if 1<<i >= int(endType) {
			break
		}
		if types[i] > 0 {
			t.Logf("%2d : %40b : %3d", i, IntersectionType(1<<i), types[i])
		} else {
			t.Errorf("need checking for type: %2d", i)
			amountFail++
			t.Logf("%2d : %40b : %3d fail", i, IntersectionType(1<<i), types[i])
		}
		sum += types[i]
	}
	t.Logf("full amount = %d", sum)
	t.Logf("amount fail = %d", amountFail)
}
