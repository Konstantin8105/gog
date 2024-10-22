package gog

import (
	"fmt"
	"math"
	"testing"
)

// cpu: Intel(R) Xeon(R) CPU E3-1240 V2 @ 3.40GHz
// BenchmarkArcSplitByPoint/WithoutPoint-4         	  702073	      1678 ns/op	     592 B/op	      13 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	  656035	      1786 ns/op	     616 B/op	      13 allocs/op
//
// BenchmarkArcSplitByPoint/WithoutPoint-4         	  662924	      1577 ns/op	     496 B/op	      10 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	  783086	      1518 ns/op	     456 B/op	       9 allocs/op
//
// BenchmarkArcSplitByPoint/WithoutPoint-4         	  877332	      1434 ns/op	     552 B/op	      12 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	  828496	      1412 ns/op	     552 B/op	      12 allocs/op
//
// BenchmarkArcSplitByPoint/WithoutPoint-4         	 1245926	       959.5 ns/op	     320 B/op	       6 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	 1356279	       865.0 ns/op	     272 B/op	       5 allocs/op
//
// BenchmarkArcSplitByPoint/WithoutPoint-4         	 1297148	       939.4 ns/op	     240 B/op	       5 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	 1397499	       815.6 ns/op	     192 B/op	       4 allocs/op
//
// BenchmarkArcSplitByPoint/WithoutPoint-4         	 1328588	       909.7 ns/op	     240 B/op	       5 allocs/op
// BenchmarkArcSplitByPoint/WithPoint-4            	 1495153	       850.5 ns/op	     192 B/op	       4 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/LineLine3d-8         	 6218899	       192.9 ns/op	       0 B/op	       0 allocs/op
// Benchmark/LineLine-8           	  803428	      1535 ns/op	      48 B/op	       3 allocs/op
// Benchmark/ArcSplitNoPoint-8    	  646896	      1777 ns/op	     240 B/op	       5 allocs/op
// Benchmark/ArcSplitPoint-8      	  724832	      1562 ns/op	     192 B/op	       4 allocs/op
// Benchmark/Split-8              	       1	1632291808 ns/op	45196712 B/op	  257346 allocs/op
// Benchmark/New-8                	       1	2294259064 ns/op	29537544 B/op	  851695 allocs/op
// Benchmark/Triangulation-8      	      88	  15954597 ns/op	  527464 B/op	   10823 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU E3-1240 V2 @ 3.40GHz
// Benchmark/LineLine3d-4         	 7504720	       150.1 ns/op	       0 B/op	       0 allocs/op
// Benchmark/LineLine-4           	 1259445	       882.8 ns/op	      48 B/op	       3 allocs/op
// Benchmark/ArcSplitNoPoint-4    	  884901	      1201 ns/op	     240 B/op	       5 allocs/op
// Benchmark/ArcSplitPoint-4      	 1228023	      1091 ns/op	     192 B/op	       4 allocs/op
// Benchmark/Split-4              	       2	 877110258 ns/op	12689864 B/op	  233365 allocs/op
// Benchmark/New-4                	       1	1513065216 ns/op	21567288 B/op	  849719 allocs/op
// Benchmark/Triangulation-4      	     122	   9730151 ns/op	  470788 B/op	   10084 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU E3-1240 V2 @ 3.40GHz
// Benchmark/LineLine3d-4         	 8530584	       137.4 ns/op	       0 B/op	       0 allocs/op
// Benchmark/LineLine-4           	 1496764	       797.0 ns/op	      48 B/op	       3 allocs/op
// Benchmark/ArcSplitNoPoint-4    	 1651136	       723.0 ns/op	     192 B/op	       3 allocs/op
// Benchmark/ArcSplitPoint-4      	 1828788	       624.1 ns/op	     144 B/op	       2 allocs/op
// Benchmark/Split-4              	       2	 811401014 ns/op	11925824 B/op	  222620 allocs/op
// Benchmark/New-4                	       1	1363595467 ns/op	21564008 B/op	  849643 allocs/op
// Benchmark/Triangulation-4      	     142	   8284052 ns/op	  446659 B/op	    9645 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU E3-1240 V2 @ 3.40GHz
// Benchmark/LineLine3d-4         	 8361111	       146.6 ns/op	       0 B/op	       0 allocs/op
// Benchmark/LineLine-4           	 1391710	       837.9 ns/op	      48 B/op	       3 allocs/op
// Benchmark/ArcSplitNoPoint-4    	 1486761	       757.0 ns/op	     192 B/op	       3 allocs/op
// Benchmark/ArcSplitPoint-4      	 1942293	       623.4 ns/op	     144 B/op	       2 allocs/op
// Benchmark/Split-4              	       2	 813958700 ns/op	 8520072 B/op	   70906 allocs/op
// Benchmark/New-4                	       1	1291804152 ns/op	 4608328 B/op	   51161 allocs/op
// Benchmark/Triangulation-4      	     147	   7875483 ns/op	  311667 B/op	    3905 allocs/op
func Benchmark(b *testing.B) {
	b.Run("LineLine3d", func(b *testing.B) {
		pps := []Point3d{
			Point3d{1, 1, 0}, // 0
			Point3d{4, 4, 0}, // 1
			Point3d{0, 5, 0}, // 2
			Point3d{5, 0, 0}, // 3
		}
		ra, rb, tint := LineLine3d(
			pps[0], pps[1],
			pps[2], pps[3],
		)
		if !tint {
			panic(fmt.Errorf("%v %v %v", ra, rb, tint))
		}
		for n := 0; n < b.N; n++ {
			LineLine3d(
				pps[0], pps[1],
				pps[2], pps[3],
			)
		}
	})
	b.Run("LineLine", func(b *testing.B) {
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
	})
	b.Run("ArcSplitNoPoint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := ArcSplitByPoint(Point{-2, 0}, Point{0, +2}, Point{+2, 0})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("ArcSplitPoint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := ArcSplitByPoint(Point{0, -1}, Point{1, 0}, Point{0, 1}, Point{1, 0})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Split", func(b *testing.B) {
		pps := []Point{
			Point{X: 1, Y: 1}, // 0
			Point{X: 4, Y: 4}, // 1
			Point{X: 0, Y: 5}, // 2
			Point{X: 5, Y: 0}, // 3
		}
		var dist float64 = 0.1
		for n := 0; n < b.N; n++ {
			var model Model
			for i := range pps {
				if i == 0 {
					continue
				}
				model.AddLine(pps[i-1], pps[i], 10)
			}
			mesh, err := New(model)
			if err != nil {
				b.Fatal(err)
			}

			// distance
			err = mesh.Split(dist)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("New", func(b *testing.B) {
		var pps []Point
		size := 200
		for i := 0; i < size; i++ {
			pps = append(pps, Point{X: float64(i), Y: float64(i)})
		}
		for i := 0; i < size; i++ {
			pps = append(pps, Point{X: float64(size - i), Y: float64(i)})
		}
		for i := 0; i < size; i++ {
			pps = append(pps, Point{X: float64(size - i), Y: float64(size - i)})
		}
		for i := 0; i < size; i++ {
			pps = append(pps, Point{X: float64(i), Y: float64(size - i)})
		}
		for i := 0; i < size; i++ {
			pps = append(pps, Point{X: float64(size/2 + i + 1), Y: float64(size/2 + 1 - i)})
		}
		var model Model
		for i := range pps {
			if i == 0 {
				continue
			}
			model.AddLine(pps[i-1], pps[i], 10)
		}

		// distance
		var dist float64 = 1.0
		model.Split(dist)
		model.Intersection()

		new := func() {
			mesh, err := New(model)
			if err != nil {
				panic(err)
			}
			_ = mesh
		}
		new()

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			new()
		}
	})
	b.Run("Triangulation", func(b *testing.B) {
		pps := []Point{
			Point{X: 1, Y: 1}, // 0
			Point{X: 4, Y: 4}, // 1
			Point{X: 0, Y: 5}, // 2
			Point{X: 5, Y: 0}, // 3
		}

		for n := 0; n < b.N; n++ {
			var model Model
			for i := range pps {
				if i == 0 {
					continue
				}
				model.AddLine(pps[i-1], pps[i], 10)
			}

			// distance
			var dist float64
			dist = model.MinPointDistance()
			{
				xmax := -math.MaxFloat64
				xmin := +math.MaxFloat64
				for i := range model.Points {
					xmax = math.Max(xmax, model.Points[i].X)
					xmin = math.Min(xmin, model.Points[i].X)
				}
				dist = math.Min(dist, math.Abs(xmax-xmin)/10.0)
			}
			model.Intersection()
			model.Split(dist)
			model.ArcsToLines()
			mesh, err := New(model)
			if err != nil {
				b.Fatal(err)
			}
			err = mesh.Delanay()
			if err != nil {
				b.Fatal(err)
			}
			err = mesh.Split(dist)
			if err != nil {
				b.Fatal(err)
			}
			err = mesh.Smooth()
			if err != nil {
				b.Fatal(err)
			}
			err = mesh.Check()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
