package gog

import (
	"fmt"
	"os"
)

func ExampleModel() {
	var m Model
	m.AddCircle(0, 0, 1, 1)
	m.AddLine(Point{-1, 0}, Point{1, 0}, 2)
	m.AddLine(Point{0, -1}, Point{0, 1}, 3)
	fmt.Fprintf(os.Stdout, "Only structural lines:\n%s", m)
	m.Intersection()
	m.Split(0.2)
	m.ArcsToLines()
	m.RemoveEmptyPoints()
	m.ConvexHullTriangles()
	m.Intersection()
	m.RemoveEmptyPoints()
	fmt.Fprintf(os.Stdout, "After intersection:\n%s", m)
	fmt.Fprintf(os.Stdout, "Minimal distance between points:\n%.4f", m.MinPointDistance())
	// Output:
	// Only structural lines:
	// Points:
	// 000	{+0.0000 -1.0000}
	// 001	{+1.0000 +0.0000}
	// 002	{+0.0000 +1.0000}
	// 003	{-1.0000 +0.0000}
	// Lines:
	// 000	[  1   3   2]
	// 001	[  0   2   3]
	// Arcs:
	// 000	[  0   1   2   1]
	// 001	[  0   3   2   1]
	// After intersection:
	// Points:
	// 000	{+0.0000 -1.0000}
	// 001	{+1.0000 +0.0000}
	// 002	{+0.0000 +1.0000}
	// 003	{-1.0000 +0.0000}
	// 004	{+0.0000 +0.0000}
	// 005	{+0.7071 -0.7071}
	// 006	{-0.7071 -0.7071}
	// 007	{-0.7071 +0.7071}
	// 008	{+0.0000 -0.8333}
	// 009	{+0.0000 -0.6667}
	// 010	{+0.0000 -0.5000}
	// 011	{+0.0000 -0.3333}
	// 012	{+0.0000 -0.1667}
	// 013	{+0.0000 +0.8333}
	// 014	{+0.0000 +0.6667}
	// 015	{+0.0000 +0.5000}
	// 016	{+0.0000 +0.3333}
	// 017	{+0.0000 +0.1667}
	// 018	{+0.8333 +0.0000}
	// 019	{+0.6667 +0.0000}
	// 020	{+0.5000 +0.0000}
	// 021	{+0.3333 +0.0000}
	// 022	{+0.1667 +0.0000}
	// 023	{-0.8333 +0.0000}
	// 024	{-0.6667 +0.0000}
	// 025	{-0.5000 +0.0000}
	// 026	{-0.3333 +0.0000}
	// 027	{-0.1667 +0.0000}
	// 028	{+0.1951 -0.9808}
	// 029	{+0.3827 -0.9239}
	// 030	{+0.5556 -0.8315}
	// 031	{+0.8315 -0.5556}
	// 032	{+0.9239 -0.3827}
	// 033	{+0.9808 -0.1951}
	// 034	{-0.1951 -0.9808}
	// 035	{-0.3827 -0.9239}
	// 036	{-0.5556 -0.8315}
	// 037	{-0.8315 -0.5556}
	// 038	{-0.9239 -0.3827}
	// 039	{-0.9808 -0.1951}
	// 040	{-0.1951 +0.9808}
	// 041	{-0.3827 +0.9239}
	// 042	{-0.5556 +0.8315}
	// 043	{-0.8315 +0.5556}
	// 044	{-0.9239 +0.3827}
	// 045	{-0.9808 +0.1951}
	// Lines:
	// 000	[  0   8   3]
	// 001	[  8   9   3]
	// 002	[  9  10   3]
	// 003	[ 10  11   3]
	// 004	[ 11  12   3]
	// 005	[  4  12   3]
	// 006	[  2  13   3]
	// 007	[ 13  14   3]
	// 008	[ 14  15   3]
	// 009	[ 15  16   3]
	// 010	[ 16  17   3]
	// 011	[  4  17   3]
	// 012	[  1  18   2]
	// 013	[ 18  19   2]
	// 014	[ 19  20   2]
	// 015	[ 20  21   2]
	// 016	[ 21  22   2]
	// 017	[  4  22   2]
	// 018	[  3  23   2]
	// 019	[ 23  24   2]
	// 020	[ 24  25   2]
	// 021	[ 25  26   2]
	// 022	[ 26  27   2]
	// 023	[  4  27   2]
	// 024	[  0  28   1]
	// 025	[ 28  29   1]
	// 026	[ 29  30   1]
	// 027	[  5  30   1]
	// 028	[  5  31   1]
	// 029	[ 31  32   1]
	// 030	[ 32  33   1]
	// 031	[  1  33   1]
	// 032	[  0  34   1]
	// 033	[ 34  35   1]
	// 034	[ 35  36   1]
	// 035	[  6  36   1]
	// 036	[  6  37   1]
	// 037	[ 37  38   1]
	// 038	[ 38  39   1]
	// 039	[  3  39   1]
	// 040	[  2  40   1]
	// 041	[ 40  41   1]
	// 042	[ 41  42   1]
	// 043	[  7  42   1]
	// 044	[  7  43   1]
	// 045	[ 43  44   1]
	// 046	[ 44  45   1]
	// 047	[  3  45   1]
	// Triangles:
	// 000	[  0  28  29  -1]
	// 001	[ 28  29  30  -1]
	// 002	[  5  30  29  -1]
	// 003	[ 30   5  31  -1]
	// 004	[  5  31  32  -1]
	// 005	[ 31  32  33  -1]
	// 006	[  1  33  32  -1]
	// 007	[  2  40  41  -1]
	// 008	[ 40  41  42  -1]
	// 009	[  7  42  41  -1]
	// 010	[ 42   7  43  -1]
	// 011	[  7  43  44  -1]
	// 012	[ 43  44  45  -1]
	// 013	[  3  45  44  -1]
	// 014	[ 39   3  45  -1]
	// 015	[  3  39  38  -1]
	// 016	[ 37  38  39  -1]
	// 017	[  6  37  38  -1]
	// 018	[ 36   6  37  -1]
	// 019	[  6  36  35  -1]
	// 020	[ 34  35  36  -1]
	// 021	[  1   2  13  -1]
	// 022	[  2  40  13  -1]
	// 023	[ 13   1  40  -1]
	// 024	[  2   1  18  -1]
	// 025	[  1  33  18  -1]
	// 026	[ 18   2  33  -1]
	// Minimal distance between points:
	// 0.1667
}
