//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/Konstantin8105/gog"
)

func main() {
	var m gog.Model
	var state int
	// view result in dxf format
	view := func() {
		if err := os.WriteFile(
			fmt.Sprintf("stage%02d.dxf", state),
			[]byte(m.Dxf()),
			0644,
		); err != nil {
			panic(err)
		}
		state++
	}
	// create model
	m.AddCircle(0, 0, 1, 1)
	m.AddCircle(0, 0, 0.5, 1)
	m.AddCircle(0, 0, 0.75, 1)
	m.AddLine(gog.Point{X: -1, Y: 0}, gog.Point{X: 1, Y: 0}, 2)
	m.AddLine(gog.Point{X: 0, Y: -1}, gog.Point{X: 0, Y: 1}, 3)
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
	m.Get(mesh)
	view() // 9
	m.Triangles = nil
	err = mesh.Split(0.2)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		// return
	}
	mesh.Smooth()
	err = mesh.Materials()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		// return
	}
	err = mesh.Check()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v\n", err)
		// return
	}
	m.Get(mesh)
	view() // 10
}
