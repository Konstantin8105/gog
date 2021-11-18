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
	it   State
	pi   Point

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
	pi, st := SegmentAnalisys(
		pps[0], pps[1],
		pps[2], pps[3],
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
	// 23	      100000000000000000000000	not found
	// 24	     1000000000000000000000000	not found
	// 25	    10000000000000000000000000	not found
	// 26	   100000000000000000000000000	not found
	// 27	  1000000000000000000000000000	not found
	// 28	 10000000000000000000000000000	not found
// 29	100000000000000000000000000000	not found
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
			OverlapP0AP0B |
			OverlapP0AP1B |
			HorizontalSegmentB | VerticalSegmentB |
			Collinear,
		pi:  Point{X: 0, Y: 8},
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
		it: VerticalSegmentA | VerticalSegmentB |
			OverlapP0AP0B |
			OverlapP1AP1B |
			Collinear,
		pi:  Point{X: 0, Y: 8},
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
		it: OverlapP0AP0B |
			OverlapP1AP1B |
			Collinear,
		pi:  Point{X: 2, Y: 8},
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
		it:  VerticalSegmentA | VerticalSegmentB | Parallel,
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
		it:  VerticalSegmentA | VerticalSegmentB | Collinear,
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
		it:  Collinear,
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
		it: VerticalSegmentA | VerticalSegmentB |
			OverlapP1AP0B |
			Collinear,
		pi:  Point{X: 2, Y: 6},
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
		it: OverlapP1AP0B |
			Collinear,
		pi:  Point{X: 4, Y: 4},
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
		it: VerticalSegmentA |
			VerticalSegmentB |
			HorizontalSegmentA |
			HorizontalSegmentB |
			ZeroLengthSegmentA |
			ZeroLengthSegmentB |
			OverlapP0AP0B |
			OverlapP0AP1B |
			OverlapP1AP0B |
			OverlapP1AP1B |
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
		it:  OnPoint0SegmentB | OnSegmentA,
		pi:  Point{X: 2, Y: 2},
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
		it:  OnSegmentA | OnSegmentB,
		pi:  Point{X: 2.5, Y: 2.5},
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
		it:  OnPoint1SegmentB | OnSegmentA,
		pi:  Point{X: 2, Y: 2},
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
		it:  OnPoint0SegmentA | OnSegmentB,
		pi:  Point{X: 2, Y: 2},
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
		it:  OnPoint1SegmentA | OnSegmentB,
		pi:  Point{X: 2, Y: 2},
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
		it:  VerticalSegmentB | OnRay11SegmentA | OnSegmentB,
		pi:  Point{X: 5, Y: 5},
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
		it:  VerticalSegmentB | OnRay00SegmentA | OnSegmentB,
		pi:  Point{X: 5, Y: 5},
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
		it:  VerticalSegmentA | OnRay11SegmentB | OnSegmentA,
		pi:  Point{X: 5, Y: 5},
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
		it:  VerticalSegmentA | OnRay00SegmentB | OnSegmentA,
		pi:  Point{X: 5, Y: 5},
		dbp: 6,
	},
	{ // 18 : Test data - no intersection
		ps: []Point{
			Point{X: 1.098, Y: 0},
			Point{X: -1.5449, Y: 12.53},
			Point{X: 1.2, Y: 2},
			Point{X: 5, Y: 5},
		},
		it:  OnRay00SegmentB | OnSegmentA,
		pi:  Point{X: 0.7509280607532581, Y: 1.6454695216473094},
		dbp: 2.4656014677e+00,
	},
	{ // 19 : Test data - no intersection
		ps: []Point{
			Point{X: 5.108, Y: 0},
			Point{X: 8.339, Y: 16.17},
			Point{X: 9, Y: 2},
			Point{X: 5, Y: 5},
		},
		it:  OnSegmentA | OnSegmentB,
		pi:  Point{X: 5.9627881085877945, Y: 4.277908918559155},
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
		it:  Collinear,
		dbp: 7.0710678119e-01,
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
		ts.bp = t.bp
		ts.dbp = t.dbp
		return ts
	}

	// add names
	for i := range tcs {
		tcs[i].name = fmt.Sprintf("%d", i)
		tcs[i].bp = Point{X: -1, Y: -2}
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
			tc.bp.X += mv
			tc.bp.Y += mv
			tc.name += fmt.Sprintf(":%+5.3f", mv)
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
		tc.bp.X *= -1.0
		tc.bp.Y *= -1.0
		tc.name += ":rotate"
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
			pi, it := SegmentAnalisys(tc.ps[0], tc.ps[1], tc.ps[2], tc.ps[3])
			if it != tc.it {
				t.Error("Not same types")
				t.Logf("Expected : %30b", tc.it)
				t.Logf("Value    : %30b", it)
				t.Logf("Diff1    : %30b", tc.it&^it)
				t.Logf("%s", view(tc.it&^it, names))
				t.Logf("Diff2    : %30b", it&^tc.it)
				t.Logf("%s", view(it&^tc.it, names))
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

	tcs := []struct {
		Line [2]Point
		Arc  [3]Point

		pi []Point
		it State
	}{
		{ // 0
			Line: [2]Point{{-2, 0}, {2, 0}},
			Arc:  [3]Point{{0, -1}, {1, 0}, {0, 1}},

			pi: []Point{{1, 0}},
			it: HorizontalSegmentA | OnSegmentA | OnSegmentB | LineFromArcCenter,
		},
		{ // 1
			Line: [2]Point{{-2, 4}, {2, 4}},
			Arc:  [3]Point{{0, -1}, {1, 0}, {0, 1}},

			pi: []Point{},
			it: HorizontalSegmentA | LineOutside,
		},
		{ // 2
			Line: [2]Point{{-2, 4}, {2, 4}},
			Arc:  [3]Point{{0, -1}, {0, 0}, {0, 1}},

			pi: []Point{},
 			it: HorizontalSegmentA | LineOutside,
			// ArcOneLine | OnSegmentA | OnRay11SegmentB | VerticalSegmentB,
		},
		{ // 3
			Line: [2]Point{{-2, 4}, {2, 4}},
			Arc:  [3]Point{{1, 0}, {1, 0}, {0, 1}},

			pi: []Point{{-3,4}},
			it: OnRay00SegmentA | OnRay11SegmentB | HorizontalSegmentA |
			 ArcIsLine | Arc01indentical,
		},
		{ // 4
			Line: [2]Point{{-2, 4}, {2, 4}},
			Arc:  [3]Point{{1, 0}, {0, 1}, {0, 1}},

			pi: []Point{{-3,4}},
			it: OnRay00SegmentA | OnRay11SegmentB | HorizontalSegmentA |
			 ArcIsLine | Arc12indentical,
		},
		{ // 5
			Line: [2]Point{{-2, 4}, {2, 4}},
			Arc:  [3]Point{{0, 1}, {0, 1}, {0, 1}},

			pi: []Point{},
			it:  HorizontalSegmentA | ArcIsPoint |
				Arc01indentical | Arc02indentical | Arc12indentical,
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("ARC%02d", i), func(t *testing.T) {
			pi, it := ArcLineAnalisys(tcs[i].Line, tcs[i].Arc)
			if it != tc.it {
				t.Error("Not same types")
				t.Logf("Expected : %30b", tc.it)
				t.Logf("Value    : %30b", it)
				t.Logf("Diff1    : %30b", tc.it&^it)
				t.Logf("%s", view(tc.it&^it, names))
				t.Logf("Diff2    : %30b", it&^tc.it)
				t.Logf("%s", view(it&^tc.it, names))
			}
			// store
			for i := 0; i < len(types); i++ {
				if it&State(1<<i) != 0 {
					types[i]++
				}
			}
			if len(pi) != len(tcs[i].pi) {
				t.Errorf("not valid sizes %d != %d", len(pi), len(tcs[i].pi))
			} else {
				for i := range pi {
					if eps := 1e-6; eps < Distance(pi[i], tc.pi[i]) &&
						math.Abs(pi[i].X) > eps &&
						math.Abs(pi[i].Y) > eps {
						t.Errorf("2 point is too far: %#v != %#v\ndistance = %.6e",
							pi[i],
							tc.pi[i],
							Distance(pi[i], tc.pi[i]),
						)
					}
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
		SegmentAnalisys(
			pps[0], pps[1],
			pps[2], pps[3],
		)
	}
}

func TestLinePointDistance(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.ps = append(tc.ps, tc.bp)
			if err := Check(tc.ps...); err != nil {
				t.Fatal(err)
			}
			d := LinePointDistance(tc.ps[0], tc.ps[1], tc.ps[4])
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
		for _, angle := range angles {
			t.Run(fmt.Sprintf("%s:%+.2f", tc.name, angle), func(t *testing.T) {
				for index, p := range tc.ps {
					pw := Rotate(angle, p)
					if Distance(p, pw) < eps {
						t.Errorf("No change %d: %v %v", index, p, pw)
					}
					pw = Rotate(-angle, pw)
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
			expect:  []Point{Point{0, 0}, Point{X: -4, Y: 4}},
		},
		{
			segment: [2]Point{Point{X: 4, Y: 5}, Point{X: -4, Y: -3}},
			mirror:  [2]Point{Point{X: -1, Y: 1}, Point{X: 5, Y: 1}},
			expect:  []Point{Point{0, 1}, Point{X: -4, Y: 5}},
		},
		{
			segment: [2]Point{Point{X: 4, Y: 10}, Point{X: 4, Y: 0}},
			mirror:  [2]Point{Point{X: 0, Y: 0}, Point{X: 1, Y: 1}},
			expect:  []Point{Point{4, 4}, Point{X: 0, Y: 4}},
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
			if Distance(tc.expect[0], ml0) > eps || Distance(tc.expect[1], ml1) > eps {
				t.Errorf("Not valid points: %v %v", ml0, ml1)
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
