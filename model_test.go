package gog

import (
	"fmt"
	"os"
)

func ExampleModel() {
	var m Model
	m.AddCircle(0, 0, 1, 1)
	m.AddLine(Point{-1, 0}, Point{1, 0}, 2)
	m.AddLine(Point{0, -1}, Point{0, 1}, 2)
	m.Intersection()
	fmt.Fprintf(os.Stdout, "%v", m)
	// Output:
	// {[[0.00000e+00,-1.00000e+00] [1.00000e+00,0.00000e+00] [0.00000e+00,1.00000e+00] [-1.00000e+00,0.00000e+00] [0.00000e+00,0.00000e+00] [7.07107e-01,-7.07107e-01] [-7.07107e-01,-7.07107e-01] [-7.07107e-01,7.07107e-01]] [[0 4 2] [2 4 2] [1 4 2] [3 4 2]] [[0 5 1 1] [0 6 3 1] [2 7 3 1]]}
}
