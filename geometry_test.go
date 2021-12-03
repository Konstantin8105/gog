package gog

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Konstantin8105/cs"
)

func getStates() (names []string) {
	out, err := exec.Command("go", "doc", "-all", "State").CombinedOutput()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(out), "\n")
	for i := range lines {
		if strings.Contains(lines[i], "const (") {
			lines = lines[i+1:]
			break
		}
	}
	for i := range lines {
		if strings.Contains(lines[i], ")") {
			lines = lines[:i]
			break
		}
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
		index := strings.Index(lines[i], "//")
		if index < 0 {
			if 0 < len(lines[i]) {
				names = append(names, lines[i])
			}
			continue
		}
		lines[i] = lines[i][:index]
		if 0 < len(lines[i]) {
			names = append(names, lines[i])
		}
	}
	return
}

func view(s State, names []string) string {
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
		out += fmt.Sprintf("%2d\t%25s\t", i, names[i-1])
		if s.Has(si) {
			out += "found"
		} else {
			out += "not found"
		}
		out += "\n"
	}
	return out
}

type TestCase struct {
	name string
	ps   []Point
	itA  State
	itB  State
	pi   []Point

	bp  Point   // base point
	dbp float64 // distance between base point and line between 0 and 1
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

	if err := Check(pps...); err != nil {
		panic(err)
	}
	pi, stA, stB := LineLine(
		pps[0], pps[1],
		pps[2], pps[3],
	)
	fmt.Fprintf(os.Stdout, "Intersection point: %s\n", pi)
	fmt.Fprintf(os.Stdout, "Intersection state A:\n%s\n", stA)
	fmt.Fprintf(os.Stdout, "Intersection state B:\n%s\n", stB)
	// Output:
	// Intersection point: [[2.50000e+00,2.50000e+00]]
	// Intersection state A:
	//  1	                            10	not found
	//  2	                           100	not found
	//  3	                          1000	not found
	//  4	                         10000	not found
	//  5	                        100000	not found
	//  6	                       1000000	found
	//  7	                      10000000	not found
	//  8	                     100000000	not found
	//  9	                    1000000000	not found
	// 10	                   10000000000	not found
	// 11	                  100000000000	not found
	// 12	                 1000000000000	not found
	//
	// Intersection state B:
	//  1	                            10	not found
	//  2	                           100	not found
	//  3	                          1000	not found
	//  4	                         10000	not found
	//  5	                        100000	not found
	//  6	                       1000000	found
	//  7	                      10000000	not found
	//  8	                     100000000	not found
	//  9	                    1000000000	not found
	// 10	                   10000000000	not found
	// 11	                  100000000000	not found
	// 12	                 1000000000000	not found
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
		itA: VerticalSegment | OnPoint0Segment,
		itB: ZeroLengthSegment | VerticalSegment | HorizontalSegment |
			OnPoint0Segment | OnPoint1Segment,
		pi:  []Point{},
		dbp: 1,
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
		itA: VerticalSegment | OnPoint0Segment | OnPoint1Segment | Collinear,
		itB: VerticalSegment | OnPoint0Segment | OnPoint1Segment | Collinear,
		pi:  []Point{},
		dbp: 1,
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
		itA: OnPoint0Segment | OnPoint1Segment | Collinear,
		itB: OnPoint0Segment | OnPoint1Segment | Collinear,
		pi:  []Point{},
		dbp: 3.1622776602e-01,
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
		itA: VerticalSegment | Parallel,
		itB: VerticalSegment | Parallel,
		dbp: 1,
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
		itA: VerticalSegment | Collinear,
		itB: VerticalSegment | Collinear,
		dbp: 3,
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
		itA: Collinear,
		itB: Collinear,
		dbp: 7.0710678119e-01,
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
		itA: VerticalSegment | OnPoint1Segment | Collinear,
		itB: VerticalSegment | OnPoint0Segment | Collinear,
		pi:  []Point{},
		dbp: 3,
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
		itA: OnPoint1Segment | Collinear,
		itB: OnPoint0Segment | Collinear,
		pi:  []Point{},
		dbp: 7.0710678119e-01,
	},
	{ // 8
		// *1,2,3,4  //
		ps: []Point{
			Point{X: 5, Y: 5}, // 1
			Point{X: 5, Y: 5}, // 2
			Point{X: 5, Y: 5}, // 3
			Point{X: 5, Y: 5}, // 4
		},
		itA: ZeroLengthSegment | VerticalSegment | HorizontalSegment | OnPoint0Segment | OnPoint1Segment,
		itB: ZeroLengthSegment | VerticalSegment | HorizontalSegment | OnPoint0Segment | OnPoint1Segment,
		pi:  []Point{},
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
		itA: OnSegment,
		itB: OnPoint0Segment,
		pi:  []Point{{X: 2, Y: 2}},
		dbp: 7.0710678119e-01,
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
		itA: OnSegment,
		itB: OnSegment,
		pi:  []Point{{X: 2.5, Y: 2.5}},
		dbp: 7.0710678119e-01,
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
		itA: OnSegment,
		itB: OnPoint1Segment,
		pi:  []Point{{X: 2, Y: 2}},
		dbp: 7.0710678119e-01,
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
		itA: OnPoint0Segment,
		itB: OnSegment,
		pi:  []Point{{X: 2, Y: 2}},
		dbp: 4.9923017660e+00,
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
		itA: OnPoint1Segment,
		itB: OnSegment,
		pi:  []Point{{X: 2, Y: 2}},
		dbp: 4.9923017660e+00,
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
		itA: OnRay11Segment,
		itB: VerticalSegment,
		pi:  []Point{},
		dbp: 7.0710678119e-01,
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
		itA: OnRay00Segment,
		itB: VerticalSegment,
		pi:  []Point{},
		dbp: 7.0710678119e-01,
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
		itA: VerticalSegment,
		itB: OnRay11Segment,
		pi:  []Point{},
		dbp: 6,
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
		itA: VerticalSegment,
		itB: OnRay00Segment,
		pi:  []Point{},
		dbp: 6,
	},
	{ // 18 : Test data - no intersection
		ps: []Point{
			Point{X: 1.098, Y: 0},
			Point{X: -1.5449, Y: 12.53},
			Point{X: 1.2, Y: 2},
			Point{X: 5, Y: 5},
		},
		itB: OnRay00Segment,
		pi:  []Point{},
		dbp: 2.4656014677e+00,
	},
	{ // 19 : Test data - no intersection
		ps: []Point{
			Point{X: 5.108, Y: 0},
			Point{X: 8.339, Y: 16.17},
			Point{X: 9, Y: 2},
			Point{X: 5, Y: 5},
		},
		itA: OnSegment,
		itB: OnSegment,
		pi:  []Point{{X: 5.9627881085877945, Y: 4.277908918559155}},
		dbp: 5.5977179784e+00,
	},
	{ // 20
		// *1  //
		// |   //
		// *3  //
		// |   //
		// *2  //
		// |   //
		// *4  //
		ps: []Point{
			Point{X: 5, Y: 5}, // 1
			Point{X: 3, Y: 3}, // 2
			Point{X: 4, Y: 4}, // 3
			Point{X: 2, Y: 2}, // 4
		},
		itA: Collinear,
		itB: Collinear,
		dbp: 7.0710678119e-01,
	},
	{ // 21
		ps: []Point{
			{-2, 0}, {2, 0},
			{0, -1}, {1, 0}, {0, 1},
		},
		itA: HorizontalSegment | OnSegment,
		itB: OnSegment,
		pi:  []Point{{1, 0}},
	},
	{ // 22
		ps: []Point{
			{-2, 4}, {2, 4},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{},
		itA: HorizontalSegment,
	},
	{ // 23
		ps: []Point{
			{-2, 4}, {2, 4},
			{0, -1}, {0, 0}, {0, 1}},

		pi:  []Point{},
		itA: HorizontalSegment,
		itB: VerticalSegment | OnRay11Segment | ArcIsLine,
	},
	{ // 24
		ps: []Point{
			{-2, 4}, {2, 4},
			{1, 0}, {1, 0}, {0, 1}},

		pi:  []Point{},
		itA: OnRay00Segment | HorizontalSegment,
		itB: OnRay11Segment | ArcIsLine,
	},
	{ // 25
		ps: []Point{
			{-2, 4}, {2, 4},
			{1, 0}, {0, 1}, {0, 1}},

		pi:  []Point{},
		itA: OnRay00Segment | HorizontalSegment,
		itB: OnRay11Segment | ArcIsLine,
	},
	{ // 26
		ps: []Point{
			{-2, 4}, {2, 4},
			{0, 1}, {0, 1}, {0, 1}},

		pi:  []Point{},
		itA: HorizontalSegment | OnRay11Segment,
		itB: ArcIsPoint | VerticalSegment | HorizontalSegment | ZeroLengthSegment,
	},
	{ // 27
		ps: []Point{
			{-2, 1}, {2, 1},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{Point{0, 1}},
		itA: HorizontalSegment | OnSegment,
		itB: OnPoint1Segment,
	},
	{ // 28
		ps: []Point{
			{0, 1}, {2, 1},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{},
		itA: HorizontalSegment | OnPoint0Segment,
		itB: OnPoint1Segment,
	},
	{ // 29
		ps: []Point{
			{0, 1}, {0, -1},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{},
		itA: VerticalSegment | OnPoint0Segment | OnPoint1Segment,
		itB: OnPoint0Segment | OnPoint1Segment,
	},
	{ // 30
		ps: []Point{
			{1, 1}, {1, -1},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{{1, 0}},
		itA: VerticalSegment | OnSegment,
		itB: OnSegment,
	},
	{ // 31
		ps: []Point{
			{2, 1}, {2, -1},
			{0, -1}, {1, 0}, {0, 1}},

		pi:  []Point{},
		itA: VerticalSegment,
	},
	{ // 32
		ps: []Point{
			{1, 1}, {1, -1},
			{0, 0}, {1, 1}, {0, 2}},

		pi:  []Point{Point{1, 1}},
		itA: VerticalSegment | OnPoint0Segment,
		itB: OnSegment,
	},
	{ // 33
		ps: []Point{
			{-1, 1}, {3, 1},
			{1, 0}, {2, 1}, {1, 2}},

		pi:  []Point{{2, 1}},
		itA: HorizontalSegment | OnSegment,
		itB: OnSegment,
	},
	{ // 34
		ps: []Point{
			{2, 4}, {4, 2},
			{2, 3}, {3, 2}, {2, 1}},

		pi: []Point{},
	},
	{ // 35
		ps: []Point{
			{2, 3}, {3, 2},
			{2, 3}, {3, 2}, {2, 1}},

		pi:  []Point{{3, 2}},
		itA: OnPoint0Segment | OnPoint1Segment,
		itB: OnPoint0Segment | OnSegment,
	},
	{ // 36
		ps: []Point{
			{2, 2 + 0.70710678*2}, {2 + 0.70710678*2, 2},
			{2, 3}, {3, 2}, {2, 1}},

		pi:  []Point{{2 + 0.70710678, 2 + 0.70710678}},
		itA: OnSegment,
		itB: OnSegment,
	},
	{ // 37
		ps: []Point{
			{-2, 0}, {2, 0},
			{0, 1}, {1, 0}, {0, -1},
		},
		itA: HorizontalSegment | OnSegment,
		itB: OnSegment,
		pi:  []Point{{1, 0}},
	},
	{ // 38
		ps: []Point{
			{1, 1}, {1, 1},
			{1, 1}, {1, 1}, {1, 1},
		},
		itA: HorizontalSegment | VerticalSegment | ZeroLengthSegment |
			OnPoint0Segment | OnPoint1Segment,
		itB: HorizontalSegment | VerticalSegment | ZeroLengthSegment |
			OnPoint0Segment | OnPoint1Segment |
			ArcIsPoint,
		pi: []Point{},
	},
	{ // 39
		ps: []Point{
			{0, -1}, {0, 0},
			{0, -1}, {1, 0}, {0, 1},
		},
		itA: VerticalSegment | OnPoint0Segment,
		itB: OnPoint0Segment,
		pi:  []Point{},
	},
	{ // 40
		ps: []Point{
			{-0.0, -2.2}, {-0.0, -0.20},
			{+1.0, -1.2}, {-0.0, -2.20}, {1.0, -3.2},
		},
		pi:  []Point{Point{0, -2.2}},
		itA: VerticalSegment | OnPoint0Segment,
		itB: OnSegment,
	},
}

func init() {
	// create copy
	copy := func(t TestCase) TestCase {
		var ts TestCase
		ts.pi = make([]Point, len(t.pi))
		copy(ts.pi, t.pi)
		ts.ps = make([]Point, len(t.ps))
		copy(ts.ps, t.ps)
		ts.name = t.name
		ts.itA = t.itA
		ts.itB = t.itB
		ts.bp = t.bp
		ts.dbp = t.dbp
		return ts
	}

	// add names
	for i := range tcs {
		tcs[i].name = fmt.Sprintf("t%02d", i)
		tcs[i].bp = Point{X: -1, Y: -2}
	}

	var size int

	// TODO line to arc

	// add test with moving
	size = len(tcs)
	for i := 0; i < size; i++ {
		for _, mx := range []float64{-1.0, -1.0, +1.0} {
			for _, my := range []float64{1.2, -2.0, +3.0} {
				tc := copy(tcs[i])
				// move
				for i := range tc.ps {
					tc.ps[i].X += mx
					tc.ps[i].Y += my
				}
				for i := range tc.pi {
					tc.pi[i].X += mx
					tc.pi[i].Y += my
				}
				tc.bp.X += mx
				tc.bp.Y += my
				tc.name += fmt.Sprintf("_Move%+5.3f%+5.3f_", mx, my)
				tcs = append(tcs, tc)
			}
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
		for i := range tc.pi {
			tc.pi[i].X *= -1.0
			tc.pi[i].Y *= -1.0
		}
		tc.bp.X *= -1.0
		tc.bp.Y *= -1.0
		tc.name += "_Rotate"
		tcs = append(tcs, tc)
	}
}

func Test(t *testing.T) {
	names := getStates()
	var types [64]int
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if err := Check(tc.ps...); err != nil {
				t.Fatal(err)
			}
			var (
				pi       []Point
				itA, itB State
			)
			switch len(tc.ps) {
			case 4:
				pi, itA, itB = LineLine(
					tc.ps[0], tc.ps[1],
					tc.ps[2], tc.ps[3],
				)
			case 5:
				pi, itA, itB = LineArc(
					tc.ps[0], tc.ps[1],
					tc.ps[2], tc.ps[3], tc.ps[4],
				)
			default:
				t.Fatal("not valid data")
			}
			for index, s := range [][2]State{
				[2]State{itA, tc.itA},
				[2]State{itB, tc.itB},
			} {
				if s[0] != s[1] {
					t.Errorf("Not same types: %d", index)
					t.Logf("Points   : %v", tc.ps)
					t.Logf("Expected : %30b", s[0])
					t.Logf("Actual   : %30b", s[1])
					t.Logf("Diff1    : %30b", s[1]&^s[0])
					t.Logf("View:\n%s", view(s[1]&^s[0], names))
					t.Logf("Diff2    : %30b", s[0]&^s[1])
					t.Logf("View:\n%s", view(s[0]&^s[1], names))
				}
				// store
				for i := 0; i < len(types); i++ {
					if itA&State(1<<i) != 0 {
						types[i]++
					}
					if itB&State(1<<i) != 0 {
						types[i]++
					}
				}
			}
			if len(pi) != len(tc.pi) {
				t.Errorf("not valid sizes %d != %d", len(pi), len(tc.pi))
				t.Errorf("Points  : %v", tc.ps)
				t.Errorf("Actual  : %v", pi)
				t.Errorf("Expected: %v", tc.pi)
			} else {
				bs := make([]bool, len(pi))
				for i := range pi {
					for j := range tc.pi {
						if eps := 1e-6; Distance(pi[i], tc.pi[j]) < eps {
							bs[j] = true
						}
					}
				}
				for i := range pi {
					if bs[i] {
						continue
					}
					t.Errorf("2 point is too far: %#v != %#v\ndistance = %.6e",
						pi[i],
						tc.pi[i],
						Distance(pi[i], tc.pi[i]),
					)
				}
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
			t.Logf("%2d : %30s : %3d", i, names[i-1], types[i])
		} else {
			t.Errorf("need checking for type: %2d", i)
			amountFail++
			t.Logf("%2d : %30s : %3d fail", i, names[i-1], types[i])
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
		if err := Check(tcs[i]...); err == nil {
			t.Errorf("Not valid error in case %d", i)
		}
	}
}

func TestCodeStyle(t *testing.T) {
	cs.All(t)
}

func Benchmark(b *testing.B) {
	pps := []Point{
		Point{X: 1, Y: 1}, // 0
		Point{X: 4, Y: 4}, // 1
		Point{X: 0, Y: 5}, // 2
		Point{X: 5, Y: 0}, // 3
	}

	if err := Check(pps...); err != nil {
		panic(err)
	}
	for n := 0; n < b.N; n++ {
		LineLine(
			pps[0], pps[1],
			pps[2], pps[3],
		)
	}
}

func TestLinePointDistance(t *testing.T) {
	for _, tc := range tcs {
		if len(tc.pi) != 4 {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			tc.ps = append(tc.ps, tc.bp)
			if err := Check(tc.ps...); err != nil {
				t.Fatal(err)
			}
			d := PointLineDistance(tc.ps[0], tc.ps[1], tc.ps[4])
			if eps := 1e-6; math.Abs(d-tc.dbp) > eps {
				t.Errorf("Not valid distance: %.10e != expected %.2e",
					d, tc.dbp,
				)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	angles := []float64{0.01, math.Pi / 4, math.Pi / 2, 2, 5, -0.1}
	eps := 0.001
	for _, tc := range tcs {
		if len(tc.pi) != 4 {
			continue
		}
		for _, angle := range angles {
			t.Run(fmt.Sprintf("%s:%+.2f", tc.name, angle), func(t *testing.T) {
				for index, p := range tc.ps {
					pw := Rotate(0, 0, angle, p)
					if Distance(p, pw) < eps {
						t.Errorf("No change %d: %v %v", index, p, pw)
					}
					pw = Rotate(0, 0, -angle, pw)
					if Distance(p, pw) > eps {
						t.Errorf("Some change %d: %v %v", index, p, pw)
					}
				}
			})
		}
	}
}

func TestMirrorLine(t *testing.T) {
	tcs := []struct {
		segment [2]Point
		mirror  [2]Point
		expect  []Point
	}{
		{
			segment: [2]Point{Point{X: 4, Y: 4}, Point{X: -4, Y: -4}},
			mirror:  [2]Point{Point{X: -1, Y: 0}, Point{X: 5, Y: 0}},
			expect:  []Point{Point{4, -4}, Point{X: -4, Y: 4}},
		},
		{
			segment: [2]Point{Point{X: 4, Y: 5}, Point{X: -4, Y: -3}},
			mirror:  [2]Point{Point{X: -1, Y: 1}, Point{X: 5, Y: 1}},
			expect:  []Point{Point{4, -3}, Point{X: -4, Y: 5}},
		},
		{
			segment: [2]Point{Point{X: 4, Y: 10}, Point{X: 4, Y: 0}},
			mirror:  [2]Point{Point{X: 0, Y: 0}, Point{X: 1, Y: 1}},
			expect:  []Point{Point{10, 4}, Point{X: 0, Y: 4}},
		},
	}
	eps := 0.001
	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			ml0, ml1, err := MirrorLine(
				tc.segment[0],
				tc.segment[1],
				tc.mirror[0],
				tc.mirror[1],
			)
			if err != nil {
				t.Fatal(err)
			}
			if eps < Distance(tc.expect[0], ml0) || eps < Distance(tc.expect[1], ml1) {
				t.Errorf("Not valid points: %v %v", ml0, ml1)
				t.Logf("Points : %v", tc.segment )
				t.Logf("Mirror : %v", tc.mirror )
				t.Logf("Distance: %v %v", Distance(tc.expect[0], ml0), Distance(tc.expect[1], ml1) )
			}
		})
	}
}

func TestOrientation(t *testing.T) {
	tcs := []struct {
		ps [3]Point
		or OrientationPoints
	}{
		{
			ps: [3]Point{Point{0, 0}, Point{1, 1}, Point{2, 2}},
			or: CollinearPoints,
		},
		{
			ps: [3]Point{Point{0, 0}, Point{1, 1}, Point{2, 20}},
			or: CounterClockwisePoints,
		},
		{
			ps: [3]Point{Point{0, 0}, Point{1, 1}, Point{2, -2}},
			or: ClockwisePoints,
		},
	}
	for _, tc := range tcs {
		p0 := tc.ps[0]
		p1 := tc.ps[1]
		p2 := tc.ps[2]
		or := Orientation(p0, p1, p2)
		if or != tc.or {
			t.Errorf("Not valid: %v", tc)
		}
	}
}

func TestAngleBetween(t *testing.T) {
	tcs := []struct {
		name          string
		xc, yc        float64
		from, mid, to Point
		a             Point
		expect        bool
	}{
		{ // 0
			from:   Point{+1, +0},
			mid:    Point{+0, +1},
			to:     Point{-1, +0},
			a:      Point{+0, +1},
			expect: true,
		}, { // 1
			from:   Point{+1, +0},
			mid:    Point{+0, +1},
			to:     Point{-1, +0},
			a:      Point{+0, -1},
			expect: false,
		}, { // 2
			from:   Point{-1, +0},
			mid:    Point{+0, +1},
			to:     Point{+1, +0},
			a:      Point{+0, +1},
			expect: true,
		}, { // 3
			from:   Point{-1, +0},
			mid:    Point{+0, +1},
			to:     Point{+1, +0},
			a:      Point{+0, -1},
			expect: false,
		}, { // 4
			from:   Point{+0, +1},
			mid:    Point{-1, +0},
			to:     Point{+0, -1},
			a:      Point{-1, +0},
			expect: true,
		}, { // 5
			from:   Point{+0, +1},
			mid:    Point{-1, +0},
			to:     Point{+0, -1},
			a:      Point{+1, +0},
			expect: false,
		}, { // 6
			from:   Point{+0, -1},
			mid:    Point{-1, +0},
			to:     Point{+0, +1},
			a:      Point{-1, +0},
			expect: true,
		}, { // 7
			from:   Point{+0, -1},
			mid:    Point{-1, +0},
			to:     Point{+0, +1},
			a:      Point{+1, +0},
			expect: false,
		},
	}
	for i := range tcs {
		tcs[i].name = fmt.Sprintf("t%02d", i)
	}
	// move
	size := len(tcs)
	for _, xc := range []float64{2.1, 0.0, -2.0} {
		for _, yc := range []float64{2.0, 0.0, -3.0} {
			for i := 0; i < size; i++ {
				c := tcs[i]
				c.name += fmt.Sprintf("move%.0f%.0f", xc, yc)
				c.from.X += xc
				c.from.Y += yc
				c.mid.X += xc
				c.mid.Y += yc
				c.to.X += xc
				c.to.Y += yc
				c.a.X += xc
				c.a.Y += yc
				c.xc = xc
				c.yc = yc
				tcs = append(tcs, c)
			}
		}
	}

	for i := range tcs {
		t.Run(fmt.Sprintf("%s", tcs[i].name), func(t *testing.T) {
			res := AngleBetween(Point{tcs[i].xc, tcs[i].yc}, tcs[i].from, tcs[i].mid, tcs[i].to, tcs[i].a)
			if res != tcs[i].expect {
				t.Errorf("not valid: %v", tcs[i])
			}
		})
	}

}

func ExampleArcSplitByPoint() {
	tcs := [][]Point{
		[]Point{ // 0
			Point{X: -2, Y: 0},
			Point{X: 0, Y: +2},
			Point{X: +2, Y: 0},
		},
		[]Point{ // 1
			Point{X: +2, Y: 0},
			Point{X: 0, Y: +2},
			Point{X: -2, Y: 0},
		},
		[]Point{ // 2
			Point{X: -2, Y: 0},
			Point{X: 0, Y: +2},
			Point{X: +2, Y: 0},
			Point{X: 0, Y: +2},
		},
		[]Point{ // 3
			Point{X: +2, Y: 0},
			Point{X: 0, Y: +2},
			Point{X: -2, Y: 0},
			Point{X: 0, Y: +2},
		},
		[]Point{ // 4
			Point{X: -2, Y: 0},
			Point{X: 0, Y: +2},
			Point{X: +2, Y: 0},
			Point{+1.41421, +1.41421},
		},
		[]Point{ // 5
			Point{0, 1}, Point{-1, 0}, Point{0, -1},
			Point{-1, 0},
		},
	}
	for index, tc := range tcs {
		fmt.Fprintf(os.Stdout, "case %d:\n", index)
		res, err := ArcSplitByPoint(tc[0], tc[1], tc[2], tc[3:]...)
		if err != nil {
			panic(fmt.Errorf("index %d: %v", index, err))
		}
		for i := range res {
			for j := range res[i] {
				fmt.Fprintf(os.Stdout, "[%02d,%02d] = %+.5f\n", i, j, res[i][j])
			}
		}
	}
	// Output:
	// case 0:
	// [00,00] = {-2.00000 +0.00000}
	// [00,01] = {-1.41421 +1.41421}
	// [00,02] = {+0.00000 +2.00000}
	// [01,00] = {+0.00000 +2.00000}
	// [01,01] = {+1.41421 +1.41421}
	// [01,02] = {+2.00000 +0.00000}
	// case 1:
	// [00,00] = {+2.00000 +0.00000}
	// [00,01] = {+1.41421 +1.41421}
	// [00,02] = {+0.00000 +2.00000}
	// [01,00] = {+0.00000 +2.00000}
	// [01,01] = {-1.41421 +1.41421}
	// [01,02] = {-2.00000 +0.00000}
	// case 2:
	// [00,00] = {-2.00000 +0.00000}
	// [00,01] = {-1.41421 +1.41421}
	// [00,02] = {-0.00000 +2.00000}
	// [01,00] = {-0.00000 +2.00000}
	// [01,01] = {+1.41421 +1.41421}
	// [01,02] = {+2.00000 +0.00000}
	// case 3:
	// [00,00] = {+2.00000 +0.00000}
	// [00,01] = {+1.41421 +1.41421}
	// [00,02] = {-0.00000 +2.00000}
	// [01,00] = {-0.00000 +2.00000}
	// [01,01] = {-1.41421 +1.41421}
	// [01,02] = {-2.00000 +0.00000}
	// case 4:
	// [00,00] = {-2.00000 +0.00000}
	// [00,01] = {-0.76537 +1.84776}
	// [00,02] = {+1.41421 +1.41421}
	// [01,00] = {+1.41421 +1.41421}
	// [01,01] = {+1.84776 +0.76537}
	// [01,02] = {+2.00000 +0.00000}
	// case 5:
	// [00,00] = {+0.00000 +1.00000}
	// [00,01] = {-0.70711 +0.70711}
	// [00,02] = {-1.00000 +0.00000}
	// [01,00] = {-1.00000 +0.00000}
	// [01,01] = {-0.70711 -0.70711}
	// [01,02] = {-0.00000 -1.00000}
}
