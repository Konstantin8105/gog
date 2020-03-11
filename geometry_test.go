package gog

import (
	"fmt"
	"math"
	"os"
	"testing"
)

type TestCase struct {
	name string
	ps   []Point
	it   State
	pi   Point
}

func Example() {
	// *2   *0  //
	//  \  /    //
	//    X     //
	//  /  \    //
	// *1   *3  //
	pps := []Point{
		Point{X: 1, Y: 1}, // 0
		Point{X: 4, Y: 4}, // 1
		Point{X: 0, Y: 5}, // 2
		Point{X: 5, Y: 0}, // 3
	}

	if err := Check(&pps); err != nil {
		panic(err)
	}
	pi, st := SegmentAnalisys(
		0, 1,
		2, 3,
		&pps,
	)
	fmt.Fprintf(os.Stdout, "Intersection point: %s\n", pi)
	fmt.Fprintf(os.Stdout, "Intersection state:\n%s\n", st)
	// Output:
	// Intersection point: [2.50000e+00,2.50000e+00]
	// Intersection state:
	//  1	                            10	not found
	//  2	                           100	not found
	//  3	                          1000	not found
	//  4	                         10000	not found
	//  5	                        100000	not found
	//  6	                       1000000	not found
	//  7	                      10000000	not found
	//  8	                     100000000	not found
	//  9	                    1000000000	found
	// 10	                   10000000000	found
	// 11	                  100000000000	not found
	// 12	                 1000000000000	not found
	// 13	                10000000000000	not found
	// 14	               100000000000000	not found
	// 15	              1000000000000000	not found
	// 16	             10000000000000000	not found
	// 17	            100000000000000000	not found
	// 18	           1000000000000000000	not found
	// 19	          10000000000000000000	not found
	// 20	         100000000000000000000	not found
	// 21	        1000000000000000000000	not found
	// 22	       10000000000000000000000	not found
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
		it: ZeroLengthSegmentB |
			VerticalSegmentA |
			Point0SegmentAonPoint0SegmentB |
			Point0SegmentAonPoint1SegmentB |
			HorizontalSegmentB | VerticalSegmentB |
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
		it: VerticalSegmentA | VerticalSegmentB |
			Point0SegmentAonPoint0SegmentB | Point1SegmentAonPoint1SegmentB |
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
		it: Point0SegmentAonPoint0SegmentB | Point1SegmentAonPoint1SegmentB |
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
		it: VerticalSegmentA | VerticalSegmentB | Parallel,
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
		it: VerticalSegmentA | VerticalSegmentB | Collinear,
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
		it: VerticalSegmentA | VerticalSegmentB |
			Point1SegmentAonPoint0SegmentB |
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
		it: Point1SegmentAonPoint0SegmentB |
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
		it: VerticalSegmentA |
			VerticalSegmentB |
			HorizontalSegmentA |
			HorizontalSegmentB |
			ZeroLengthSegmentA |
			ZeroLengthSegmentB |
			Point0SegmentAonPoint0SegmentB |
			Point1SegmentAonPoint0SegmentB |
			Point0SegmentAonPoint1SegmentB |
			Point1SegmentAonPoint1SegmentB |
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
		it: Point0SegmentBinSegmentA,
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
		it: OnSegmentA | OnSegmentB,
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
		it: Point1SegmentBinSegmentA,
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
		it: Point0SegmentAinSegmentB,
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
		it: Point1SegmentAinSegmentB,
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
		it: VerticalSegmentB | OnRay11SegmentA | OnSegmentB,
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
		it: VerticalSegmentB | OnRay00SegmentA | OnSegmentB,
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
		it: VerticalSegmentA | OnRay11SegmentB | OnSegmentA,
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
		it: VerticalSegmentA | OnRay00SegmentB | OnSegmentA,
		pi: Point{X: 5, Y: 5},
	},
	{ // 18 : Test data - no intersection
		ps: []Point{
			Point{X: 1.098, Y: 0},
			Point{X: -1.5449, Y: 12.53},
			Point{X: 1.2, Y: 2},
			Point{X: 5, Y: 5},
		},
		it: OnRay00SegmentB | OnSegmentA,
		pi: Point{X: 0.7509280607532581, Y: 1.6454695216473094},
	},
	{ // 19 : Test data - no intersection
		ps: []Point{
			Point{X: 5.108, Y: 0},
			Point{X: 8.339, Y: 16.17},
			Point{X: 9, Y: 2},
			Point{X: 5, Y: 5},
		},
		it: OnSegmentA | OnSegmentB,
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
			if err := Check(&tc.ps); err != nil {
				t.Fatal(err)
			}
			pi, it := SegmentAnalisys(0, 1, 2, 3, &tc.ps)
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
				if it&State(1<<i) != 0 {
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
			t.Logf("%2d : %40b : %3d", i, State(1<<i), types[i])
		} else {
			t.Errorf("need checking for type: %2d", i)
			amountFail++
			t.Logf("%2d : %40b : %3d fail", i, State(1<<i), types[i])
		}
		sum += types[i]
	}
	t.Logf("full amount = %d", sum)
	t.Logf("amount fail = %d", amountFail)
}

func TestCheckError(t *testing.T) {
	tcs := [][]Point{
		[]Point{Point{X: 1.0, Y: math.NaN()}},
		[]Point{Point{X: math.NaN(), Y: 1.0}},
		[]Point{Point{X: math.Inf(1), Y: 1.0}},
		[]Point{Point{X: 1.0, Y: math.Inf(1)}},
	}
	for i := range tcs {
		if err := Check(&tcs[i]); err == nil {
			t.Errorf("Not valid error in case %d", i)
		}
	}
}
