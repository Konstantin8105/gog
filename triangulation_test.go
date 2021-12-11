package gog

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestTriangulation(t *testing.T) {
	tcs := [][]Point{
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
		}, { // 15
			Point{20, 38},
			Point{34, 17},
			Point{16, 8},
			Point{43, 2},
			Point{25, 47},
		}, { // 16
			Point{0, 37},
			Point{38, 13},
			Point{12, 35},
			Point{8, 33},
			Point{32, 37},
		}, { // 17
			Point{33, 16},
			Point{9, 24},
			Point{23, 37},
			Point{18, 2},
			Point{26, 28},
		}, { // 18
			Point{34, 45},
			Point{17, 25},
			Point{0, 31},
			Point{25, 0},
			Point{15, 24},
		}, { // 19
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
		}, { // 20
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
			Point{2, 0},
		}, { // 21
			Point{0, 0},
			Point{1, 1},
			Point{1, 0},
			Point{2, 0},
			Point{1, 0.5},
		}, { // 22
			Point{10, 40},
			Point{36, 27},
			Point{1, 12},
			Point{6, 42},
			Point{41, 24},
		},
	}

	if Debug == false {
		Debug = true
		defer func() {
			Debug = false
		}()
	}

	for i, ts := range tcs {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatal(r)
				}
			}()
			var model Model
			for i := range ts {
				model.AddPoint(ts[i])
			}
			// TODO add lines
			dist := model.MinPointDistance()
			mesh, err := New(model)
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
			mesh.Split(dist)
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 2: %v", err)
			}
			mesh.Smooth()
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 3: %v", err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 4: %v", err)
			}
			err = mesh.Materials()
			if err != nil {
				t.Fatal(err)
			}
			err = mesh.Check()
			if err != nil {
				t.Fatalf("check 5: %v", err)
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
