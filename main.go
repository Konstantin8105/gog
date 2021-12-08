// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Konstantin8105/gog"
)

func main() {
	var m gog.Model
	var state int
	// view result in dxf format
	view := func() {
		if err := ioutil.WriteFile(fmt.Sprintf("stage%02d.dxf", state), []byte(m.Dxf()), 0644); err != nil {
			panic(err)
		}
		state++
	}
	// create model
	m.AddCircle(0, 0, 1, 1)
	m.AddCircle(0, 0, 0.5, 1)
	m.AddCircle(0, 0, 0.75, 1)
	m.AddLine(gog.Point{-1, 0},gog.Point{1, 0}, 2)
	m.AddLine(gog.Point{0, -1},gog.Point{0, 1}, 3)
	//fmt.Fprintf(os.Stdout, "Only structural lines:\n%s", m)
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
	mesh, err := gog.New(m)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		return
	}
	view() // 8
	err = mesh.Delanay()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		// return
	}
	mesh.Split(0.2)
	m.Get(mesh)
	view() // 9
	// fmt.Fprintf(os.Stdout, "After intersection:\n%s", m)
	// fmt.Fprintf(os.Stdout, "Minimal distance between points:\n%.4f", m.MinPointDistance())
}