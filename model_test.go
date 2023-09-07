package gog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Konstantin8105/compare"
)

func TestCopy(t *testing.T) {
	m := Model{
		Points:    []Point{{1, 2}},
		Lines:     [][3]int{{3, 4, 5}, {6, 7, 8}},
		Arcs:      [][4]int{{11, 12, 13, 14}, {15, 16, 17, 18}},
		Triangles: [][4]int{{21, 22, 23, 34}},
	}
	equal := func(m0, m1 *Model) bool {
		return m0.String() == m1.String()
	}
	var c Model
	if equal(&m, &c) {
		t.Fatalf("strange 1")
	}
	c = m.Copy()
	if !equal(&m, &c) {
		t.Fatalf("not full copy\n%s\n%s", m.String(), c.String())
	}
	c.Lines[1][1] = -4
	if equal(&m, &c) {
		t.Fatalf("strange 2")
	}
}

func TestModel(t *testing.T) {
	if Log == false {
		Log = true
		defer func() {
			Log = false
		}()
	}
	var (
		m     Model
		state int
		buf   bytes.Buffer
	)
	// view result in dxf format
	view := func() {
		if err := ioutil.WriteFile(
			fmt.Sprintf("ExampleModelOnState%02d.dxf", state),
			[]byte(m.Dxf()),
			0644,
		); err != nil {
			panic(err)
		}
		state++
	}
	// create model
	m.AddCircle(0, 0, 1, 1)
	m.AddLine(Point{-1, 0}, Point{1, 0}, 2)
	m.AddLine(Point{0, -1}, Point{0, 1}, 3)
	fmt.Fprintf(&buf, "Only structural lines:\n%s", m)
	view() // 0
	m.Intersection()
	view() // 1
	m.Split(0.2)
	view() // 2
	m.ArcsToLines()
	view() // 3
	m.RemoveEmptyPoints()
	view() // 4
	m.ConvexHullTriangles()
	view() // 5
	m.Intersection()
	view() // 6
	m.RemoveEmptyPoints()
	view() // 7
	m.Triangles = nil
	mesh, err := New(m)
	if err != nil {
		t.Fatal(err)
	}
	err = mesh.Materials()
	if err != nil {
		t.Fatal(err)
	}
	//
	for _, p := range []Point{
		Point{+0.5, +0.5},
		Point{+0.5, -0.5},
		Point{-0.5, -0.5},
		Point{-0.5, +0.5},
	} {
		mat, err := mesh.GetMaterials(p)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&buf, "%v %v\n", p, mat)
	}
	m.Get(mesh)
	view() // 8
	fmt.Fprintf(&buf, "After intersection:\n%s", m)
	fmt.Fprintf(&buf, "Minimal distance between points:\n%.4f\n", m.MinPointDistance())
	if err := m.Combine(1.05); err != nil {
		t.Fatal(err)
	}
	view() // 9
	fmt.Fprintf(&buf, "After combine:\n%s", m)

	compare.Test(t, ".example", buf.Bytes())
}
