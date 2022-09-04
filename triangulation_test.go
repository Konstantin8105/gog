package gog

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"path/filepath"
	"runtime/debug"
	"testing"
)

func BenchmarkTriangulation(b *testing.B) {
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
			dist = math.Max(dist, math.Abs(xmax-xmin)/10.0)
		}
		model.Intersection()
		model.Split(dist)
		model.ArcsToLines()
		mesh, err := New(model)
		if err != nil {
			panic(err)
		}
		err = mesh.Delanay()
		if err != nil {
			panic(err)
		}
		err = mesh.Check()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSplit(b *testing.B) {
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
			panic(err)
		}

		// distance
		mesh.Split(dist)
	}
}

func TestTriangulation(t *testing.T) {
	pnts := [][]Point{
		{ // 0
			Point{48, 47},
			Point{39, 0},
			Point{0, 1},
			Point{10, 7},
			Point{19, 45},
		}, { // 1
			Point{35, 22},
			Point{20, 40},
			Point{40, 16},
			Point{24, 29},
			Point{33, 34},
		}, { // 2
			Point{33, 11},
			Point{25, 25},
			Point{37, 4},
			Point{29, 22},
			Point{35, 36},
		}, { // 3
			Point{8, 34},
			Point{11, 35},
			Point{44, 46},
			Point{4, 42},
			Point{20, 11},
		}, { // 4
			Point{31, 24},
			Point{0, 0},
			Point{33, 4},
			Point{28, 36},
			Point{43, 23},
		}, { // 5
			Point{11, 24},
			Point{16, 16},
			Point{1, 1},
			Point{44, 44},
			Point{8, 21},
		}, { // 6
			Point{24, 0},
			Point{22, 9},
			Point{19, 44},
			Point{1, 12},
			Point{6, 3},
		}, { // 7
			Point{0, 34},
			Point{10, 48},
			Point{48, 46},
			Point{0, 46},
			Point{11, 32},
		}, { // 8
			Point{34, 13},
			Point{37, 45},
			Point{23, 49},
			Point{24, 16},
		}, { // 9
			Point{-18.45, -18.45},
			Point{-20.5, 72.05},
			Point{65.5, -20.5},
			Point{11, 24},
			Point{16, 16},
			Point{1, 1},
			Point{44, 44},
			Point{8, 21},
		}, { // 10
			Point{20, 6},
			Point{6, 11},
			Point{24, 7},
			Point{22, 9},
			Point{9, 22},
		}, { // 11
			Point{1, 37},
			Point{31, 20},
			Point{31, 24},
			Point{31, 5},
			Point{14, 16},
		}, { // 12
			Point{36, 8},
			Point{47, 43},
			Point{36, 37},
			Point{27, 17},
			Point{35, 26},
		}, { // 13
			Point{0, 34},
			Point{10, 48},
			Point{48, 46},
		}, { // 14
			Point{28, 17},
			Point{39, 24},
			Point{45, 25},
			Point{7, 38},
			Point{16, 0},
		}, { // 15
			Point{27, 32},
			Point{27, 1},
			Point{3, 37},
			Point{27, 28},
			Point{44, 14},
		}, { // 16
			Point{20, 38},
			Point{34, 17},
			Point{16, 8},
			Point{43, 2},
			Point{25, 47},
		}, { // 17
			Point{0, 37},
			Point{38, 13},
			Point{12, 35},
			Point{8, 33},
			Point{32, 37},
		}, { // 18
			Point{33, 16},
			Point{9, 24},
			Point{23, 37},
			Point{18, 2},
			Point{26, 28},
		}, { // 19
			Point{34, 45},
			Point{17, 25},
			Point{0, 31},
			Point{25, 0},
			Point{15, 24},
		}, { // 20
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
		}, { // 21
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
			Point{2, 0},
		}, { // 22
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
			Point{2, 0},
			Point{1, 0.5},
		}, { // 23
			Point{10, 40},
			Point{36, 27},
			Point{1, 12},
			Point{6, 42},
			Point{41, 24},
		}, { // 24
			Point{0, 0},
			Point{1, 0},
			Point{1, 0.1},
		},
	}

	if Debug == false {
		Debug = true
		defer func() {
			Debug = false
		}()
	}

	if Log == false {
		Log = true
		defer func() {
			Log = false
		}()
	}

	type testcase struct {
		name  string
		model Model
	}
	tcs := []testcase{}

	convert := func(name string, pnts []Point) (t []testcase) {
		for _, lines := range []bool{false, true} {
			var model Model
			for i := range pnts {
				model.AddPoint(pnts[i])
			}
			if lines {
				for i := 1; i < len(pnts); i += 2 {
					model.AddLine(pnts[i-1], pnts[i], 8)
				}
			}
			t = append(t, testcase{
				name:  fmt.Sprintf("%s.%v", name, lines),
				model: model,
			})
		}
		return
	}

	// other models
	for _, size := range []int{5, 10} {
		for _, f := range []struct {
			name   string
			f      func(size int) []Point
			noLine bool
		}{
			{"random", getRandomPoints, false},
			{"circle", getCirclePoints, false},
			{"lineonline", getLineOnLine, true},
			{"intriangle", getInTriangles, true},
		} {
			t := convert(
				fmt.Sprintf("%s%02d", f.name, size),
				f.f(size))
			if !f.noLine {
				tcs = append(tcs, t...)
			} else {
				tcs = append(tcs, t[0])
			}
		}
	}

	// model with/without lines
	for index, pnt := range pnts {
		tcs = append(tcs, convert(
			fmt.Sprintf("t%02d", index),
			pnt)...)
	}

	// read JSON models
	{
		files, err := filepath.Glob("*.json")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			var model Model
			if err = model.Read(file); err != nil {
				panic(err)
			}
			tcs = append(tcs, testcase{
				name:  fmt.Sprintf("file%s", file),
				model: model,
			})
		}
	}

	// small square
	for i := 1; i < 5; i++ { // TODO : need more 5
		var model Model
		var (
			v = math.Pow10(-i)
			a = Point{X: 0, Y: 0}
			b = Point{X: v, Y: 0}
			c = Point{X: v, Y: v}
			d = Point{X: 0, Y: v}
		)
		model.AddLine(a, b, 10)
		model.AddLine(b, c, 10)
		model.AddLine(c, d, 10)
		model.AddLine(d, a, 10)
		model.AddLine(a, c, 10)
		tcs = append(tcs, testcase{
			name:  fmt.Sprintf("small%d", i),
			model: model,
		})
	}

	for _, ts := range tcs {
		t.Run(ts.name, func(t *testing.T) {
			// t.Logf("%#v", ts.model)
			defer func() {
				if r := recover(); r != nil {
					t.Logf("stacktrace from panic: %s", string(debug.Stack()))
					t.Fatal(r)
				}
			}()
			// distance
			var dist float64
			dist = ts.model.MinPointDistance()
			{
				xmax := -math.MaxFloat64
				xmin := +math.MaxFloat64
				for i := range ts.model.Points {
					xmax = math.Max(xmax, ts.model.Points[i].X)
					xmin = math.Min(xmin, ts.model.Points[i].X)
				}
				dist = math.Max(dist, math.Abs(xmax-xmin)/10.0)
			}
			ts.model.Intersection()
			ts.model.Split(dist)
			ts.model.ArcsToLines()
			if err := ioutil.WriteFile(
				ts.name+".model.dxf",
				[]byte(ts.model.Dxf()),
				0644,
			); err != nil {
				t.Error(err)
			}
			mesh, err := New(ts.model)
			if err != nil {
				t.Fatal(err)
			}
			err = mesh.Delanay()
			if err != nil {
				t.Fatal(err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 1: %v", err)
			}
			defer func() {
				// write dxf file
				ts.model.Get(mesh)
				if err := ioutil.WriteFile(
					ts.name+".dxf",
					[]byte(ts.model.Dxf()),
					0644,
				); err != nil {
					t.Error(err)
				}
			}()
			err = mesh.Split(dist)
			if err != nil {
				t.Fatalf("check 1a: %v", err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 2: %v", err)
			}
			err = mesh.Smooth()
			if err != nil {
				t.Fatalf("check 2a: %v", err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 3: %v", err)
			}
			err = mesh.Materials()
			if err != nil {
				t.Fatalf("check 4: %v", err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 5: %v", err)
			}
			for i := range mesh.model.Triangles {
				mesh.model.Triangles[i][3] = 42
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 6: %v", err)
			}
		})
	}
}

const width = 600.0

func getRandomPoints(size int) []Point {
	coords := make([]Point, size)
	for j := 0; j < size; j++ {
		coords[j] = Point{
			X: rand.Float64() * width,
			Y: rand.Float64() * width,
		}
	}
	return coords
}

func getCirclePoints(size int) []Point {
	coords := make([]Point, size+1)
	for j := 0; j < size; j++ {
		coords[j] = Point{
			X: width/2.*math.Sin(2.*math.Pi/float64(size)*float64(j)) + width/2,
			Y: width/2.*math.Cos(2.*math.Pi/float64(size)*float64(j)) + width/2,
		}
	}
	coords[len(coords)-1] = Point{X: width / 2.0, Y: width / 2.0}
	return coords
}

func getLineOnLine(size int) []Point {
	coords := make([]Point, size+3)
	coords[0] = Point{X: 0, Y: 0}
	coords[1] = Point{X: width, Y: 0}
	coords[2] = Point{X: width, Y: width}
	for j := 3; j < len(coords); j++ {
		coords[j] = Point{
			X: 10.0 + (width-20.0)*float64(j)/float64(len(coords)),
			Y: 10.0 + (width-20.0)*float64(j)/float64(len(coords)),
		}
	}
	return coords
}

func getInTriangles(size int) []Point {
	coords := make([]Point, size+4)
	coords[0] = Point{X: 0, Y: 0}
	coords[1] = Point{X: width, Y: 0}
	coords[2] = Point{X: 0, Y: width}
	coords[3] = Point{X: width, Y: width}
	for j := 4; j < len(coords); j++ {
		coords[j] = Point{
			X: 10.0 + (width-20.0)*float64(j)/float64(len(coords)),
			Y: width / 2.0,
		}
	}
	return coords
}
