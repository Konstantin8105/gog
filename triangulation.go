package gog

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"sort"

	eTree "github.com/Konstantin8105/errors"
)

// Mesh is based structure of triangulation.
// Triangle is data structure "Nodes, ribs и triangles" created by
// book "Algoritm building and analyse triangulation", A.B.Skvorcov
//
//	Scketch:
//	+------------------------------------+
//	|              tr[0]                 |
//	|  nodes[0]    ribs[0]      nodes[1] |
//	| o------------------------o         |
//	|  \                      /          |
//	|   \                    /           |
//	|    \                  /            |
//	|     \                /             |
//	|      \              /              |
//	|       \            /  ribs[1]      |
//	|        \          /   tr[1]        |
//	|  ribs[2]\        /                 |
//	|  tr[2]   \      /                  |
//	|           \    /                   |
//	|            \  /                    |
//	|             \/                     |
//	|              o  nodes[2]           |
//	+------------------------------------+
type Mesh struct {
	InternalModel Model    // internal model
	Points        []int    // tags for points
	Triangles     [][3]int // indexes of near triangles
	// TODO
	templorary struct {
		ignore []bool
	}
}

var (
	// Debug only for debugging
	Debug = false
	// Log only for minimal logging
	Log = false
)

const (
	// Boundary edge
	Boundary = -1

	// Removed element
	Removed = -2

	// Undefined state only for not valid algorithm
	Undefined = -3
	Fixed     = 100
	Movable   = 200
)

// New triangulation created by model
func New(model Model) (mesh *Mesh, err error) {
	if Log {
		log.Printf("New")
	}
	defer func() {
		if err != nil {
			et := eTree.New("New")
			_ = et.Add(err)
			err = et
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			et := eTree.New("panic error")
			_ = et.Add(fmt.Errorf("%v", r))
			_ = et.Add(fmt.Errorf("stacktrace from panic: %s",
				string(debug.Stack())))
			err = et
		}
	}()
	// prepare model before triangulation
	model.Intersection()
	if 0 < len(model.Arcs) {
		model.ArcsToLines()
	}
	// create a new Mesh
	mesh = new(Mesh)
	// convex
	_, cps := ConvexHull(model.Points, true) // points on convex hull
	if len(cps) < 3 {
		err = fmt.Errorf("not enought points for convex. Amount: %d", len(cps))
		return
	}
	// add last point for last triangle
	cps = append(cps, cps[0])
	// prepare mesh triangles
	for i, initialize := 3, false; i < len(cps); i++ {
		mesh.InternalModel.AddTriangle(cps[0], cps[i-2], cps[i-1], 1) // Default material
		if !initialize {
			mesh.Triangles = append(mesh.Triangles,
				[3]int{Boundary, Boundary, 1},
			)
			initialize = true
		} else {
			mesh.Triangles = append(mesh.Triangles,
				[3]int{i - 4, Boundary, i - 2},
			)
		}
	}
	// last not exist triangle and mark as boundary
	mesh.Triangles[len(mesh.Triangles)-1][2] = Boundary
	// clockwise all triangles
	mesh.Clockwise()
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("after convex: %v", err)
			return
		}
	}
	// add all points of model
	for i := range model.Points {
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("Check 0 {%d}: %v", i, err)
				return
			}
		}
		_, err = mesh.AddPoint(model.Points[i], Fixed)
		if err != nil {
			et := eTree.New("In add point")
			_ = et.Add(fmt.Errorf("point %d of %d", i, len(model.Points)))
			_ = et.Add(err)
			err = et
			return
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("Check 1 {%d}: %v", i, err)
				return
			}
		}
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("Check 2: %v", err)
			return
		}
	}
	// delanay
	err = mesh.Delanay()
	if err != nil {
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("Check 3: %v", err)
			return
		}
	}
	// add fixed tags
	if Debug {
		if len(mesh.Points) != len(mesh.InternalModel.Points) {
			err = fmt.Errorf("not equal points size: %d != %d ",
				len(mesh.Points),
				len(mesh.InternalModel.Points),
			)
			return
		}
	}
	for i := range mesh.Points {
		mesh.Points[i] = Fixed
	}

	// add fixed lines
	for i := range model.Lines {
		if err = mesh.AddLine(
			model.Points[model.Lines[i][0]],
			model.Points[model.Lines[i][1]],
		); err != nil {
			return
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("Check 5: %v", err)
				return
			}
		}
	}

	return
}

// Check triangulation on point, line, triangle rules
func (mesh Mesh) Check() (err error) {
	// if Log {
	// 	log.Printf("Check")
	// }
	et := eTree.New("check")
	defer func() {
		if et.IsError() {
			_ = et.Add(fmt.Errorf("amount of points: %5d", len(mesh.InternalModel.Points)))
			err = et
		}
	}()
	// amount of triangles
	if len(mesh.InternalModel.Triangles) != len(mesh.Triangles) {
		_ = et.Add(fmt.Errorf("sizes is not same"))
	}
	// same points
	for i := range mesh.InternalModel.Points {
		for j := range mesh.InternalModel.Points {
			if i <= j {
				continue
			}
			if Distance(mesh.InternalModel.Points[i], mesh.InternalModel.Points[j]) < Eps {
				_ = et.Add(fmt.Errorf("same points %v and %v", i, j))
			}
		}
	}
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		for _, d := range [3][2]int{{0, 1}, {1, 2}, {2, 0}} {
			id0 := mesh.InternalModel.Triangles[i][d[0]]
			id1 := mesh.InternalModel.Triangles[i][d[1]]
			if Distance(mesh.InternalModel.Points[id0], mesh.InternalModel.Points[id1]) < Eps {
				_ = et.Add(fmt.Errorf("triangle %d same points", i))
				_ = et.Add(fmt.Errorf("point %d: %.13f", id0, mesh.InternalModel.Points[id0]))
				_ = et.Add(fmt.Errorf("point %d: %.13f", id1, mesh.InternalModel.Points[id1]))
			}
		}
	}

	// undefined triangles
	for i := range mesh.InternalModel.Triangles {
		for j := 0; j < 3; j++ {
			if mesh.InternalModel.Triangles[i][j] == Undefined {
				_ = et.Add(fmt.Errorf("undefined point of triangle"))
			}
			if mesh.Triangles[i][j] == Undefined {
				_ = et.Add(fmt.Errorf("undefined triangle of triangle"))
			}
		}
		if mesh.InternalModel.Triangles[i][3] == Undefined {
			_ = et.Add(fmt.Errorf("undefined tag of triangle"))
		}
	}
	// clockwise triangles
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		or := Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][1]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][2]],
		)
		if or != ClockwisePoints {
			ew := eTree.New("Clockwise")
			_ = ew.Add(fmt.Errorf("triangle %d", i))
			_ = ew.Add(fmt.Errorf("is CounterClock : %v", or == CounterClockwisePoints))
			_ = ew.Add(fmt.Errorf("is CollinearPoints: %v", or == CollinearPoints))
			_ = et.Add(ew)
		}
	}
	// same triangles - self linked
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i][j] == i {
				_ = et.Add(fmt.Errorf("self linked triangle: %v %d", mesh.Triangles[i], i))
			}
		}
	}
	// correct remove
	for i := range mesh.Triangles {
		if mesh.Triangles[i][0] == Removed && mesh.InternalModel.Triangles[i][0] != Removed {
			_ = et.Add(fmt.Errorf("uncorrect removing"))
		}
	}
	// double triangles
	for i := range mesh.Triangles {
		if mesh.Triangles[i][0] == Removed {
			continue
		}
		tri := mesh.Triangles[i]
		if tri[0] == tri[1] && tri[0] != Boundary {
			_ = et.Add(fmt.Errorf("double triangles 0: %d %v %v", i, tri, mesh.Triangles[tri[0]]))
		}
		if tri[1] == tri[2] && tri[1] != Boundary {
			_ = et.Add(fmt.Errorf("double triangles 1: %d %v %v", i, tri, mesh.Triangles[tri[1]]))
		}
		if tri[2] == tri[0] && tri[2] != Boundary {
			_ = et.Add(fmt.Errorf("double triangles 2: %d %v %v", i, tri, mesh.Triangles[tri[2]]))
		}
	}
	// near triangle
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i][j] == Boundary {
				continue
			}
			if i != mesh.Triangles[mesh.Triangles[i][j]][0] &&
				i != mesh.Triangles[mesh.Triangles[i][j]][1] &&
				i != mesh.Triangles[mesh.Triangles[i][j]][2] {
				_ = et.Add(fmt.Errorf("near-near triangle: %d, %v, %v",
					i,
					mesh.Triangles[i],
					mesh.Triangles[mesh.Triangles[i][j]],
				))
			}
		}
	}
	// near triangles
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i][j] == Boundary {
				continue
			}
			arr := mesh.InternalModel.Triangles[mesh.Triangles[i][j]]
			found := false
			if arr[0] == mesh.InternalModel.Triangles[i][j] ||
				arr[1] == mesh.InternalModel.Triangles[i][j] ||
				arr[2] == mesh.InternalModel.Triangles[i][j] {
				found = true
			}
			if !found {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "triangle have not same side.\n")
				fmt.Fprintf(&buf, "tr   = %d\n", i)
				fmt.Fprintf(&buf, "side = %d\n", j)
				fmt.Fprintf(&buf, "triangle points      = %+2d\n",
					mesh.InternalModel.Triangles[i])
				fmt.Fprintf(&buf, "link triangles       = %+2d\n",
					mesh.Triangles[i])
				fmt.Fprintf(&buf, "near triangle points = %+2d\n",
					mesh.InternalModel.Triangles[mesh.Triangles[i][j]])
				fmt.Fprintf(&buf, "1: %v\n", mesh.InternalModel.Triangles)
				fmt.Fprintf(&buf, "2: %v\n", mesh.Triangles)
				_ = et.Add(fmt.Errorf("%s", buf.String()))
			}
		}
	}
	// undefined points
	// TODO : for i := range mesh.Points {
	// TODO : 	if mesh.Points[i] == Undefined {
	// TODO : 		_=et.Add(fmt.Errorf("undefined point: %d",i))
	// TODO : 	}
	// TODO : }

	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		em := eTree.New(fmt.Sprintf("Segment check triangle %d", i))
		for _, ind := range [][3]int{{0, 1, 2}, {2, 0, 1}, {1, 2, 0}} {
			_, _, stB := PointLine(
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[0]]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[1]]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[2]]],
			)
			if stB.Has(OnSegment) {
				_ = em.Add(fmt.Errorf("for %v", ind))
				_ = em.Add(fmt.Errorf("point %v",
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[0]]]))
				_ = em.Add(fmt.Errorf("point as line %v",
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[1]]]))
				_ = em.Add(fmt.Errorf("point as line %v",
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][ind[2]]]))
				_ = em.Add(fmt.Errorf("%v", stB))
			}
		}
		if em.IsError() {
			_ = et.Add(em)
		}
	}

	// lines inside line
	{
		em := eTree.New("line inside line")
		for i := range mesh.InternalModel.Lines {
			for j := range mesh.InternalModel.Lines {
				if i <= j {
					continue
				}
				if mesh.InternalModel.Lines[i][2] == Removed {
					continue
				}
				if mesh.InternalModel.Lines[j][2] == Removed {
					continue
				}
				i0 := mesh.InternalModel.Lines[i][0]
				i1 := mesh.InternalModel.Lines[i][1]
				j0 := mesh.InternalModel.Lines[j][0]
				j1 := mesh.InternalModel.Lines[j][1]
				if (i0 == j0 && i1 == j1) || (i1 == j0 && i0 == j1) {
					_ = em.Add(fmt.Errorf("line same points index: %d %d", i0, i1))
					_ = em.Add(fmt.Errorf("coord: %2d %.12e", i0, mesh.InternalModel.Points[i0]))
					_ = em.Add(fmt.Errorf("coord: %2d %.12e", i1, mesh.InternalModel.Points[i1]))
					_ = em.Add(fmt.Errorf("case %v", (i0 == j0 && i1 == j1)))
					_ = em.Add(fmt.Errorf("case %v", (i1 == j0 && i0 == j1)))
				}
				if i0 != j0 && i1 != j0 {
					_, _, stB := PointLine(
						mesh.InternalModel.Points[j0],
						mesh.InternalModel.Points[i0],
						mesh.InternalModel.Points[i1],
					)
					if stB.Has(OnSegment) {
						_ = em.Add(fmt.Errorf("i        %v", mesh.InternalModel.Lines[i]))
						_ = em.Add(fmt.Errorf("i0 = %d. %v", i0, mesh.InternalModel.Points[i0]))
						_ = em.Add(fmt.Errorf("i1 = %d. %v", i1, mesh.InternalModel.Points[i1]))
						_ = em.Add(fmt.Errorf("j        %v", mesh.InternalModel.Lines[j]))
						_ = em.Add(fmt.Errorf("j0 = %d. %v", j0, mesh.InternalModel.Points[j0]))
						_ = em.Add(fmt.Errorf("j1 = %d. %v", j1, mesh.InternalModel.Points[j1]))
						_ = em.Add(fmt.Errorf("line 1: %d inside %d", i, j))
						_ = em.Add(fmt.Errorf("propably intersection"))
					}
				}
				if i0 != j1 && i1 != j1 {
					_, _, stB := PointLine(
						mesh.InternalModel.Points[j1],
						mesh.InternalModel.Points[i0],
						mesh.InternalModel.Points[i1],
					)
					if stB.Has(OnSegment) {
						_ = em.Add(fmt.Errorf("i        %v", mesh.InternalModel.Lines[i]))
						_ = em.Add(fmt.Errorf("i0 = %d. %v", i0, mesh.InternalModel.Points[i0]))
						_ = em.Add(fmt.Errorf("i1 = %d. %v", i1, mesh.InternalModel.Points[i1]))
						_ = em.Add(fmt.Errorf("j        %v", mesh.InternalModel.Lines[j]))
						_ = em.Add(fmt.Errorf("j0 = %d. %v", j0, mesh.InternalModel.Points[j0]))
						_ = em.Add(fmt.Errorf("j1 = %d. %v", j1, mesh.InternalModel.Points[j1]))
						_ = em.Add(fmt.Errorf("line 2: %d inside %d", i, j))
						_ = em.Add(fmt.Errorf("propably intersection"))
					}
				}
			}
		}
		if em.IsError() {
			_ = et.Add(em)
		}
	}

	// no error
	return
}

// Get add into Model all triangles from Mesh
// Recommendation after `Get` : model.Intersection()
func (model *Model) Get(mesh *Mesh) {
	if Log {
		log.Printf("Get")
	}
	for _, tr := range mesh.InternalModel.Triangles {
		if tr[0] == Removed {
			continue
		}
		model.AddTriangle(
			mesh.InternalModel.Points[tr[0]],
			mesh.InternalModel.Points[tr[1]],
			mesh.InternalModel.Points[tr[2]],
			tr[3],
		)
	}
}

// Clockwise change all triangles to clockwise orientation
func (mesh *Mesh) Clockwise() {
	if Log {
		log.Printf("Clockwise")
	}
	for i := range mesh.InternalModel.Triangles {
		switch Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][1]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][2]],
		) {
		case CounterClockwisePoints:
			mesh.Triangles[i][0], mesh.Triangles[i][2] =
				mesh.Triangles[i][2], mesh.Triangles[i][0]
			mesh.InternalModel.Triangles[i][1], mesh.InternalModel.Triangles[i][2] =
				mesh.InternalModel.Triangles[i][2], mesh.InternalModel.Triangles[i][1]
		case CollinearPoints:
			panic(fmt.Errorf("collinear triangle: %#v", mesh.InternalModel))
		}
	}
}

// AddPoint is add points with tag
func (mesh *Mesh) AddPoint(p Point, tag int, triIndexes ...int) (idp int, err error) {
	if Log {
		log.Printf("AddPoint: %.20e. tag = %d", p, tag)
	}
	defer func() {
		if err != nil {
			et := eTree.New("AddPoint")
			_ = et.Add(fmt.Errorf("add point with tag %d, coord: %.9f", tag, p))
			_ = et.Add(err)
			err = et
		}
	}()
	if Debug {
		if err = mesh.Check(); err != nil {
			et := eTree.New("begin")
			_ = et.Add(err)
			err = et
			return
		}
	}

	add := func() (idp int) {
		// index of new point
		idp = mesh.InternalModel.AddPoint(p)
		for i := len(mesh.Points) - 1; i < idp; i++ {
			mesh.Points = append(mesh.Points, Undefined)
		}
		mesh.Points[idp] = tag
		return
	}

	// ignore points if on corner
	for _, pt := range mesh.InternalModel.Points {
		if Distance(p, pt) < Eps {
			idp = add()
			return
		}
	}

	// add points on line
	if tag != Movable {
		for i, size := 0, len(mesh.InternalModel.Lines); i < size; i++ {
			if mesh.InternalModel.Lines[i][2] == Removed {
				continue
			}
			// TODO fast box checking
			_, _, stB := PointLine(
				p,
				mesh.InternalModel.Points[mesh.InternalModel.Lines[i][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Lines[i][1]],
			)
			if !stB.Has(OnSegment) {
				continue
			}
			if tag == Movable {
				err = fmt.Errorf("movable point cannot be on line")
				return
			}
			// replace point tag
			tag = Fixed
			// index of new point
			idp = add()
			// add new lines
			tag := mesh.InternalModel.Lines[i][2]
			mesh.InternalModel.AddLine(mesh.InternalModel.Points[mesh.InternalModel.Lines[i][0]], p, tag)
			mesh.InternalModel.AddLine(p, mesh.InternalModel.Points[mesh.InternalModel.Lines[i][1]], tag)
			for p := 0; p < 3; p++ {
				mesh.InternalModel.Lines[i][p] = Removed
			}
		}
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("check 0a: %v", err)
			return
		}
	}

	addInTriangle := func(i int) (found bool, err error) {
		// ignore removed triangle
		if mesh.InternalModel.Triangles[i][0] == Removed {
			return
		}
		if i < 0 || len(mesh.Triangles)-1 < i {
			return
		}
		// split triangle
		var res [][3]Point
		var state int
		res, state, err = TriangleSplitByPoint(
			p,
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][1]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[i][2]],
		)
		if err != nil {
			panic(err)
		}
		if len(res) == 0 {
			return
		}
		// index of new point
		idp = add()
		// removed triangles
		removedTriangles := []int{i}

		// status
		var status int

		// find intersect side and near triangle if exist
		switch len(res) {
		case 2: // point on some side
			if mesh.Triangles[i][state] != Boundary {
				removedTriangles = append(removedTriangles, mesh.Triangles[i][state])
				status = 100 + state
			} else {
				status = 200 + state
			}
		case 3: // point in triangle
			status = 300
		}

		// repair near triangles
		var update []int
		update, err = mesh.repairTriangles(idp, removedTriangles, status)
		if err != nil {
			et := eTree.New("After repairTriangles")
			_ = et.Add(err)
			err = et
			return
		}

		// repair near triangles
		// list triangle indexes for Delanay update
		err = mesh.Delanay(update...)
		if err != nil {
			et := eTree.New("Delanay update")
			_ = et.Add(err)
			_ = et.Add(fmt.Errorf("len of res: %d", len(res)))
			err = et
			return
		}
		// TODO : add to delanay flip linked list
		return true, nil
	}

	if len(triIndexes) == 0 {
		triIndexes = make([]int, len(mesh.InternalModel.Triangles))
		for i := range mesh.InternalModel.Triangles {
			triIndexes[i] = i
		}
	}
	counter := 0
	for _, tri := range triIndexes {
		var added bool
		added, err = addInTriangle(tri)
		if err != nil {
			et := eTree.New("triangle indexes")
			_ = et.Add(err)
			if 20 < len(triIndexes) {
				_ = et.Add(fmt.Errorf("list triIndexes: %v ... more %d", triIndexes[:10], len(triIndexes)))
			} else {
				_ = et.Add(fmt.Errorf("list triIndexes: %v", triIndexes))
			}
			err = et
			return
		}
		if added {
			counter++
		}
	}
	if counter == 0 {
		err = fmt.Errorf("not found triangle: %v", idp)
		return
	}
	if 1 < counter {
		err = fmt.Errorf("not found trinagle: %v with counter %v", idp, counter)
		return
	}
	// outside of triangles or on corners
	if Debug {
		if errc := mesh.Check(); err != nil {
			err = eTree.New("Check at the end").Add(errc)
		}
		for itr, tps := range mesh.InternalModel.Triangles {
			if tps[0] == Removed {
				continue
			}
			same := false
			for i := 0; i < 3; i++ {
				if SamePoints(mesh.InternalModel.Points[idp], mesh.InternalModel.Points[tps[i]]) {
					same = true
				}
			}
			if same {
				continue
			}
			res, lineIntersect, err := TriangleSplitByPoint(
				mesh.InternalModel.Points[idp],
				mesh.InternalModel.Points[tps[0]],
				mesh.InternalModel.Points[tps[1]],
				mesh.InternalModel.Points[tps[2]],
			)
			if err != nil {
				continue
			}
			if len(res) == 0 {
				continue
			}
			err = errors.Join(
				fmt.Errorf("idp: %d", idp),
				fmt.Errorf("point: %.8e", mesh.InternalModel.Points[idp]),
				fmt.Errorf("# tri %v", itr),
				fmt.Errorf("tri: %v", tps),
				fmt.Errorf("line intersect: %v", lineIntersect),
				fmt.Errorf("res: %v", res),
			)
			panic(err)
		}
	}
	return
}

// shiftTriangle is shift numbers on right on one
func (mesh *Mesh) shiftTriangle(i int) {
	mesh.Triangles[i][0], mesh.Triangles[i][1], mesh.Triangles[i][2] =
		mesh.Triangles[i][2], mesh.Triangles[i][0], mesh.Triangles[i][1]
	mesh.InternalModel.Triangles[i][0], mesh.InternalModel.Triangles[i][1], mesh.InternalModel.Triangles[i][2] =
		mesh.InternalModel.Triangles[i][2], mesh.InternalModel.Triangles[i][0], mesh.InternalModel.Triangles[i][1]
}

// repairTriangles
//
// ap is index of added point
//
// rt is removed triangles
//
// state:
//
//	100 - point on line with 2 triangles
//	200 - point on line with 1 boundary triangle
//	300 - point in triangle
func (mesh *Mesh) repairTriangles(ap int, rt []int, state int) (updateTr []int, err error) {
	if Log {
		log.Printf("repairTriangles: ap=%d state = %d", ap, state)
	}
	defer func() {
		if err != nil {
			et := eTree.New("repairTriangles")
			_ = et.Add(fmt.Errorf("ap{%v} = %v", ap, mesh.InternalModel.Points[ap]))
			_ = et.Add(fmt.Errorf("state{%v}", state))
			_ = et.Add(fmt.Errorf("remove{%v}", rt))
			_ = et.Add(err)
			err = et
		}
	}()
	if Debug {
		if err = mesh.Check(); err != nil {
			err = fmt.Errorf("check 0: %v", err)
			return
		}
	}
	// create a chains
	//
	//	left|          | right
	//	    |    in    |
	//	from ---->----- to
	//	       out
	//
	//
	//	|         +-- ap  --+         |
	//	|        /    | |    \        |
	//	|  tc<--/    /   \    \-->tc  |
	//	|      /    |     |    \      |
	//	|     0-->--1-->--2-->--3     |
	//
	type chain struct {
		from, to int // point of triangle base
		// left, right int   // triangles index at left anf right
		in, out int // inside/outside index triangles
		before  int // triangle index before
	}
	var chains []chain
	tc := [2]int{Undefined, Undefined} // index of corner triangle
	if Debug {
		if tc[0] != Undefined || tc[1] != Undefined {
			panic("not set default values")
		}
	}

	// amount triangles before added
	size := len(mesh.Triangles)

	// create chain
	switch state {
	case 100:
		// point on triangle 0 line 0 with 2 triangles

		// debug checking
		if len(rt) != 2 {
			err = fmt.Errorf("removed triangles: %v", rt)
			return
		}
		// rotate second triangle to point intersect line 0
		if mesh.InternalModel.Triangles[rt[0]][0] != mesh.InternalModel.Triangles[rt[1]][1] ||
			mesh.InternalModel.Triangles[rt[0]][1] != mesh.InternalModel.Triangles[rt[1]][0] {
			mesh.shiftTriangle(rt[1])
			// repair triangles sides
			return mesh.repairTriangles(ap, rt, state)
		}
		if Debug {
			if mesh.InternalModel.Triangles[rt[0]][0] != mesh.InternalModel.Triangles[rt[1]][1] &&
				mesh.InternalModel.Triangles[rt[0]][1] != mesh.InternalModel.Triangles[rt[1]][0] {
				err = fmt.Errorf("not valid rotation")
				return
			}
			_, _, stB001 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][1]],
			)
			_, _, stB012 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][1]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][2]],
			)
			_, _, stB020 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][2]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][0]],
			)
			_, _, stB0012 := PointLine(
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][1]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][2]],
			)

			_, _, stB101 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][1]],
			)
			_, _, stB112 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][1]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][2]],
			)
			_, _, stB120 := PointLine(
				mesh.InternalModel.Points[ap],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][2]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][0]],
			)
			_, _, stB1012 := PointLine(
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][1]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][2]],
			)
			if stB001.Not(OnSegment) || stB012.Has(OnSegment) || stB020.Has(OnSegment) ||
				stB101.Not(OnSegment) || stB112.Has(OnSegment) || stB120.Has(OnSegment) {
				et := eTree.New("Orient mistake")
				_ = et.Add(fmt.Errorf("point AP : %.19f", mesh.InternalModel.Points[ap]))
				_ = et.Add(fmt.Errorf("stB0"))
				for i := 0; i < 3; i++ {
					_ = et.Add(fmt.Errorf("point %d : %.19f", i,
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[0]][i]]))
				}
				_ = et.Add(fmt.Errorf("%v\n%v", stB001.Not(OnSegment), stB001))
				_ = et.Add(fmt.Errorf("%v\n%v", stB012.Has(OnSegment), stB012))
				_ = et.Add(fmt.Errorf("%v\n%v", stB020.Has(OnSegment), stB020))
				_ = et.Add(fmt.Errorf("%v", stB0012))
				_ = et.Add(fmt.Errorf("stB1"))
				for i := 0; i < 3; i++ {
					_ = et.Add(fmt.Errorf("point %d : %.19f", i,
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[1]][i]]))
				}
				_ = et.Add(fmt.Errorf("%v\n%v", stB101.Not(OnSegment), stB101))
				_ = et.Add(fmt.Errorf("%v\n%v", stB112.Has(OnSegment), stB112))
				_ = et.Add(fmt.Errorf("%v\n%v", stB120.Has(OnSegment), stB120))
				_ = et.Add(fmt.Errorf("%v", stB1012))
				err = et
				return
			}
		}
		// debug: point in not on line
		if Debug {
			for k := 0; k < 2; k++ {
				_, _, stB := PointLine(
					mesh.InternalModel.Points[ap],
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][0]],
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][1]],
				)
				if stB.Not(OnSegment) {
					et := eTree.New("State100")
					_ = et.Add(fmt.Errorf("point is not on line %v", k))

					if _, _, stB := PointLine(
						mesh.InternalModel.Points[ap],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][0]],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][1]],
					); stB.Has(OnSegment) {
						_ = et.Add(fmt.Errorf("on segment 0 1"))
					}
					if _, _, stB := PointLine(
						mesh.InternalModel.Points[ap],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][1]],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][2]],
					); stB.Has(OnSegment) {
						_ = et.Add(fmt.Errorf("on segment 1 2"))
					}
					if _, _, stB := PointLine(
						mesh.InternalModel.Points[ap],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][2]],
						mesh.InternalModel.Points[mesh.InternalModel.Triangles[rt[k]][0]],
					); stB.Has(OnSegment) {
						_ = et.Add(fmt.Errorf("on segment 2 0"))
					}

					err = et
					return
				}
			}
		}

		//         0 1                      1 0         //
		//        /| |\                    /| |\        //
		//       / | | \                  / | | \       //
		//      /  | |  \                /  | |  \      //
		//     /   | |   \              /   | |   \     //
		//    /2   | |   1\            /1   | |   1\    //
		//   /     | |     \          /     | |     \   //
		//  /  rt1 | | rt0  \        /      | |      \  //
		// 2      0| |0      2      0-------2 2-------1 //
		// |       | |       |      1-------2 2-------0 //
		//  \      | |      /        \      | |      /  //
		//   \     | |     /          \     | |     /   //
		//    \1   | |   2/            \1   | |   1/    //
		//     \   | |   /              \   | |   /     //
		//      \  | |  /                \  | |  /      //
		//       \ | | /                  \ | | /       //
		//        \| |/                    \| |/        //
		//         1 0                      0 1         //

		// create chains
		chains = []chain{{
			from:   mesh.InternalModel.Triangles[rt[0]][1],
			to:     mesh.InternalModel.Triangles[rt[0]][2],
			in:     size,
			out:    mesh.Triangles[rt[0]][1],
			before: rt[0],
		}, {
			from:   mesh.InternalModel.Triangles[rt[0]][2],
			to:     mesh.InternalModel.Triangles[rt[0]][0],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]][2],
			before: rt[0],
		}, {
			from:   mesh.InternalModel.Triangles[rt[1]][1],
			to:     mesh.InternalModel.Triangles[rt[1]][2],
			in:     size + 2,
			out:    mesh.Triangles[rt[1]][1],
			before: rt[1],
		}, {
			from:   mesh.InternalModel.Triangles[rt[1]][2],
			to:     mesh.InternalModel.Triangles[rt[1]][0],
			in:     size + 3,
			out:    mesh.Triangles[rt[1]][2],
			before: rt[1],
		}}
		tc = [2]int{size + 3, size}

	case 101:
		// point on triangle 0 line 1 with 2 triangles
		mesh.shiftTriangle(rt[0])
		mesh.shiftTriangle(rt[0])
		return mesh.repairTriangles(ap, rt, 100)

	case 102:
		// point on triangle 0 line 2 with 2 triangles
		mesh.shiftTriangle(rt[0])
		return mesh.repairTriangles(ap, rt, 100)

	case 200:
		// debug checking
		if len(rt) != 1 {
			err = fmt.Errorf("removed triangles: %v", rt)
			return
		}
		if mesh.Triangles[rt[0]][0] != Boundary {
			err = fmt.Errorf("not valid boundary")
			return
		}
		// point on triangle boundary line 0
		chains = []chain{{
			from:   mesh.InternalModel.Triangles[rt[0]][1],
			to:     mesh.InternalModel.Triangles[rt[0]][2],
			in:     size,
			out:    mesh.Triangles[rt[0]][1],
			before: rt[0],
		}, {
			from:   mesh.InternalModel.Triangles[rt[0]][2],
			to:     mesh.InternalModel.Triangles[rt[0]][0],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]][2],
			before: rt[0],
		}}
		tc = [2]int{Boundary, Boundary}

	case 201:
		// point on triangle boundary line 1
		mesh.shiftTriangle(rt[0])
		mesh.shiftTriangle(rt[0])
		return mesh.repairTriangles(ap, rt, 200)

	case 202:
		// point on triangle boundary line 2
		mesh.shiftTriangle(rt[0])
		return mesh.repairTriangles(ap, rt, 200)

	case 300:
		// debug checking
		if len(rt) != 1 {
			err = fmt.Errorf("removed triangles: %v", rt)
			return
		}
		// point in triangle
		chains = []chain{{
			from:   mesh.InternalModel.Triangles[rt[0]][0],
			to:     mesh.InternalModel.Triangles[rt[0]][1],
			in:     size,
			out:    mesh.Triangles[rt[0]][0],
			before: rt[0],
		}, {
			from:   mesh.InternalModel.Triangles[rt[0]][1],
			to:     mesh.InternalModel.Triangles[rt[0]][2],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]][1],
			before: rt[0],
		}, {
			from:   mesh.InternalModel.Triangles[rt[0]][2],
			to:     mesh.InternalModel.Triangles[rt[0]][0],
			in:     size + 2,
			out:    mesh.Triangles[rt[0]][2],
			before: rt[0],
		}}
		tc = [2]int{size + 2, size}

	default:
		err = fmt.Errorf("not clear state %v", state)
		return
	}
	// debug checking
	if tc[0] == Undefined || tc[1] == Undefined {
		panic(fmt.Errorf("undefined corner triangle"))
	}

	// create triangles
	for i := range chains {
		mesh.InternalModel.AddTriangle(
			mesh.InternalModel.Points[chains[i].from],
			mesh.InternalModel.Points[chains[i].to],
			mesh.InternalModel.Points[ap],
			// by default used tag of first triangle
			mesh.InternalModel.Triangles[rt[0]][0],
		)
		tr := [3]int{Undefined, Undefined, Undefined}
		if chains[i].before == Undefined {
			panic("undefined")
		}

		tr[0] = chains[i].out
		mesh.swap(chains[i].out, chains[i].before, chains[i].in)
		if i == len(chains)-1 {
			tr[1] = tc[1]
		} else {
			tr[1] = chains[i+1].in
		}
		if i == 0 {
			tr[2] = tc[0]
		} else {
			tr[2] = chains[i-1].in
		}
		mesh.Triangles = append(mesh.Triangles, tr)
		updateTr = append(updateTr, tr[0], tr[1], tr[2])
	}

	// remove triangles
	for _, rem := range rt {
		mesh.InternalModel.Triangles[rem][0] = Removed
		mesh.InternalModel.Triangles[rem][1] = Removed
		mesh.InternalModel.Triangles[rem][2] = Removed
		mesh.Triangles[rem][0] = Removed
		mesh.Triangles[rem][1] = Removed
		mesh.Triangles[rem][2] = Removed
	}

	if Debug {
		if err = mesh.Check(); err != nil {
			et := eTree.New("Check 1")
			_ = et.Add(err)
			for i := range chains {
				_ = et.Add(fmt.Errorf("chains %d: %#v", i, chains[i]))
				_ = et.Add(fmt.Errorf("chains dist(%d,%d): Distance %e and {%e||%e}. Orient %v",
					chains[i].from,
					chains[i].to,
					Distance128(
						mesh.InternalModel.Points[chains[i].from],
						mesh.InternalModel.Points[chains[i].to],
					),
					Distance128(
						mesh.InternalModel.Points[chains[i].from],
						mesh.InternalModel.Points[ap],
					),
					Distance128(
						mesh.InternalModel.Points[ap],
						mesh.InternalModel.Points[chains[i].to],
					),
					Orientation(
						mesh.InternalModel.Points[chains[i].from],
						mesh.InternalModel.Points[chains[i].to],
						mesh.InternalModel.Points[ap],
					) == CollinearPoints,
				))
				_ = et.Add(fmt.Errorf("Point %d: %e",
					chains[i].from, mesh.InternalModel.Points[chains[i].from]))
				_ = et.Add(fmt.Errorf("Point %d: %e",
					chains[i].to, mesh.InternalModel.Points[chains[i].to]))
				_ = et.Add(fmt.Errorf("Point %d: %e",
					ap, mesh.InternalModel.Points[ap]))
			}
			for _, r := range rt {
				_ = et.Add(fmt.Errorf("remove triangle %d", r))
			}
			_ = et.Add(fmt.Errorf("corner triangle %d", tc))
			err = et
			return
		}
	}
	return
}

func (mesh *Mesh) swap(elem, from, to int) {
	if elem == Boundary {
		return
	}
	counter := 0
	for h := 0; h < 3; h++ {
		if from == mesh.Triangles[elem][h] {
			counter++
			mesh.Triangles[elem][h] = to
		}
	}
	if 1 < counter {
		panic("swap")
	}
}

// TODO delanay only for some triangles, if list empty then for  all triangles
func (mesh *Mesh) Delanay(triIndexes ...int) (err error) {
	if Log {
		log.Printf("Delanay: amount %d", len(triIndexes))
	}
	defer func() {
		if err != nil {
			et := eTree.New("Delanay")
			_ = et.Add(err)
			err = et
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err,
				fmt.Errorf("%v\n%s", r, string(debug.Stack())))
		}
	}()
	// triangle is success by delanay, if all points is outside of circle
	// from 3 triangle points
	delanay := func(tr, side int) (flip bool, err error) {
		if mesh.InternalModel.Triangles[tr][0] == Removed {
			return
		}
		neartr := mesh.Triangles[tr][side]
		if neartr == Boundary {
			return
		}
		if neartr == Removed {
			return
		}
		if mesh.InternalModel.Triangles[neartr][0] == Removed {
			return
		}
		// rotate near triangle
		for iter := 0; ; iter++ {
			if iter == 50 {
				err = fmt.Errorf("delanay infinite loop 1")
				return
			}
			if mesh.InternalModel.Triangles[tr][side] == mesh.InternalModel.Triangles[neartr][1] {
				break
			}
			mesh.shiftTriangle(neartr)
		}
		//       0 1       //       0 0       //       0 2       //
		//      /| |\      //      /| |\      //      /| |\      //
		//     / | | \     //     / | | \     //     / | | \     //
		//   2/  | |  \1   //   2/  | |  \0   //   2/  | |  \2   //
		//   /   | |   \   //   /   | |   \   //   /   | |   \   //
		//  /    | |    \  //  /    | |    \  //  /    | |    \  //
		// 2 near| | tr  2 // 2 near| | tr  1 // 2 near| | tr  0 //
		//  \   0| |0   /  //  \   0| |2   /  //  \   0| |1   /  //
		//   \   | |   /   //   \   | |   /   //   \   | |   /   //
		//   1\  | |  /2   //   1\  | |  /1   //   1\  | |  /0   //
		//     \ | | /     //     \ | | /     //     \ | | /     //
		//      \| |/      //      \| |/      //      \| |/      //
		//       1 0       //       1 2       //       1 1       //

		// is point in circle
		// Problem : for long triangle - possible triangle, but
		// not possible for arc
		if !PointInCircle(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][2]],
			[3]Point{
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][0]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][1]],
				mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][2]],
			},
		) {
			return
		}

		// flip only if middle side is not fixed
		{
			idp0 := mesh.InternalModel.AddPoint(mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][0]])
			idp1 := mesh.InternalModel.AddPoint(mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][1]])
			for _, line := range mesh.InternalModel.Lines {
				if line[0] != idp0 && line[0] != idp1 {
					continue
				}
				if line[1] != idp0 && line[1] != idp1 {
					continue
				}
				if line[2] == Fixed {
					return
				}
			}
		}

		// rotate triangle tr
		for iter := 0; ; iter++ {
			if iter == 50 {
				err = fmt.Errorf("delanay infinite loop 2")
				return
			}
			if mesh.InternalModel.Triangles[tr][0] == mesh.InternalModel.Triangles[neartr][1] {
				break
			}
			mesh.shiftTriangle(tr)
		}

		if Debug {
			if mesh.InternalModel.Triangles[tr][0] != mesh.InternalModel.Triangles[neartr][1] ||
				mesh.InternalModel.Triangles[tr][1] != mesh.InternalModel.Triangles[neartr][0] {
				err = fmt.Errorf("not valid input")
				return
			}
		}

		// flip
		flip = true

		//corner case:
		if ClockwisePoints != Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][1]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][2]],
		) || ClockwisePoints != Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][1]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][2]],
		) || ClockwisePoints != Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][2]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][2]],
		) || ClockwisePoints != Orientation(
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][0]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[tr][2]],
			mesh.InternalModel.Points[mesh.InternalModel.Triangles[neartr][2]],
		) {
			return false, nil
		}
		// before:         // after:        //
		//       0 1       //       0       //
		//      /| |\      //      / \      //
		//     / | | \-->  //     /   \-->  //
		//   2/  | |  \1   //   2/     \0   //
		//   /   | |   \   //   /  red  \   //
		//  /    | |    \  //  /    1    \  //
		// 2 red | | blu 2 // 2-----------1 //
		//  \   0| |0   /  // 1\    1    /2 //
		//   \   | |   /   //   \  blu  /   //
		//   1\  | |  /2   //   0\     /2   //
		//  <--\ | | /     //  <--\   /     //
		//      \| |/      //      \ /      //
		//       1 0       //       0       //
		{
			// flip points
			red := &mesh.InternalModel.Triangles[neartr]
			blu := &mesh.InternalModel.Triangles[tr]
			red[1], blu[1] =
				blu[2], red[2]
		}
		{
			// flip near triangles
			red := &mesh.Triangles[neartr]
			blu := &mesh.Triangles[tr]
			red[0], red[1], blu[0], blu[1] =
				blu[1], red[0], red[1], blu[0]
			mesh.swap(red[0], tr, neartr)
			mesh.swap(blu[0], neartr, tr)
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				et := eTree.New("after delanay")
				_ = et.Add(fmt.Errorf(
					"side %d in triangle %d", side, tr))
				for i := 0; i < 3; i++ {
					p := mesh.InternalModel.Triangles[tr][i]
					_ = et.Add(fmt.Errorf("Point %d %4d: %.8e", i, p, mesh.InternalModel.Points[p]))
				}
				_ = et.Add(fmt.Errorf(
					"near triangles of %d is %v", tr, mesh.Triangles[tr]))

				_ = et.Add(fmt.Errorf(
					"near triangle %d", neartr))
				for i := 0; i < 3; i++ {
					p := mesh.InternalModel.Triangles[neartr][i]
					_ = et.Add(fmt.Errorf("Point %d %4d: %.8e", i, p, mesh.InternalModel.Points[p]))
				}
				_ = et.Add(fmt.Errorf(
					"near triangles of %d is %v", neartr, mesh.Triangles[neartr]))

				_ = et.Add(err)
				err = et
				return
			}
		}
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("input: %v", err)
			return
		}
	}

	// initialize
	{
		size := len(mesh.InternalModel.Triangles)
		if len(mesh.templorary.ignore) < size {
			mesh.templorary.ignore = make([]bool, size*2)
		}
	}
	ignore := &mesh.templorary.ignore

	// loop of triangles
	for iter := 0; ; iter++ {
		counter := 0

		// reset values
		for i := range *ignore {
			(*ignore)[i] = false
		}

		for _, index := range triIndexes {
			if index < 0 {
				continue
			}
			(*ignore)[index] = true
		}

		for tr := range mesh.InternalModel.Triangles {
			if 0 < len(triIndexes) && !(*ignore)[tr] {
				continue
			}
			if mesh.InternalModel.Triangles[tr][0] == Removed {
				continue
			}
			var flip bool
			for side := 0; side < 3; side++ {
				flip, err = delanay(tr, side)
				if err != nil {
					return
				}
				if flip {
					counter++
					if Debug {
						err = mesh.Check()
						if err != nil {
							et := eTree.New("In loop")
							_ = et.Add(err)
							err = et
							return
						}
					}
					break
				}
			}
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("end of loop: %v", err)
				return
			}
		}
		if counter == 0 {
			break
		}
		if iter == 5000 {
			err = fmt.Errorf("global delanay infinite loop")
			return
		}
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("end: %v", err)
			return
		}
	}
	return nil
}

// RemoveMaterials remove material by specific points
func (mesh *Mesh) RemoveMaterials(ps ...Point) (err error) {
	mats, err := mesh.GetMaterials(ps...)
	if err != nil {
		return
	}
	for i := range mesh.InternalModel.Triangles {
		for _, m := range mats {
			if mesh.InternalModel.Triangles[i][3] == m {
				mesh.InternalModel.Triangles[i][0] = Removed
			}
		}
	}
	return
}

// SetMaterial change material for point
func (mesh *Mesh) SetMaterial(p Point, material int) (err error) {
	mats, err := mesh.GetMaterials(p)
	if err != nil {
		return
	}
	for i := range mesh.InternalModel.Triangles {
		for _, m := range mats {
			if mesh.InternalModel.Triangles[i][3] == m {
				mesh.InternalModel.Triangles[i][3] = material
			}
		}
	}
	return
}

// GetMaterials return materials for each point
func (mesh *Mesh) GetMaterials(ps ...Point) (materials []int, err error) {
	if Log {
		log.Printf("GetMaterials")
	}
	defer func() {
		if err != nil {
			et := eTree.New("GetMaterials")
			_ = et.Add(err)
			err = et
		}
	}()

	for _, p := range ps {
		for i := range mesh.InternalModel.Points {
			if Eps < Distance(p, mesh.InternalModel.Points[i]) {
				continue
			}
			// point on triangulation point

			// find all triangles with that point
			var near []int
			for t, tri := range mesh.InternalModel.Triangles {
				if i == tri[0] || i == tri[1] || i == tri[2] {
					near = append(near, t)
				}
			}
			if len(near) == 0 {
				err = fmt.Errorf("cannot find point in triangles")
				return
			}
			// add triangles shall have same materials
			mat := mesh.InternalModel.Triangles[near[0]][3]
			for _, n := range near {
				if mat != mesh.InternalModel.Triangles[n][3] {
					err = fmt.Errorf("not equal materials: %v", near)
					return
				}
			}
			materials = append(materials, mat)
			if Log {
				log.Printf("GetMaterials point in point: %v", materials)
			}
		}

		// Is point on triangulation triangle
		for it, tri := range mesh.InternalModel.Triangles {
			if mesh.InternalModel.Triangles[it][0] == Removed {
				continue
			}

			var res [][3]Point
			var lineIntersect int
			res, lineIntersect, err = TriangleSplitByPoint(p,
				mesh.InternalModel.Points[tri[0]],
				mesh.InternalModel.Points[tri[1]],
				mesh.InternalModel.Points[tri[2]],
			)
			if err != nil {
				return
			}
			switch len(res) {
			case 3:
				materials = append(materials, tri[3])
				if Log {
					log.Printf("GetMaterials triangle %d %v in triangle: %v",
						it, tri, materials)
				}
			case 2:
				j := lineIntersect
				// on edge
				if mesh.Triangles[it][j] == Boundary {
					continue
				}
				mat := []int{
					mesh.InternalModel.Triangles[it][3],
					mesh.InternalModel.Triangles[mesh.Triangles[it][j]][3],
				}
				if mat[1] == Boundary {
					materials = append(materials, mat[0])
					if Log {
						log.Printf("GetMaterials triangle %d %v on edge with boundary: %v",
							it, tri, materials)
					}
					continue
				}
				if mat[0] != mat[1] {
					err = fmt.Errorf("CollinearPoints: not equal materials on edge: %v", mat)
					return
				}
				materials = append(materials, mat[0])
				if Log {
					log.Printf("GetMaterials triangle %d %v and %d %v on edge: %v",
						it, tri,
						mesh.Triangles[it][j], mesh.InternalModel.Triangles[mesh.Triangles[it][j]],
						materials)
				}
			}
		}

		if Log {
			box := func(ps ...Point) (xmin, xmax, ymax, ymin float64) {
				xmin = +math.MaxFloat64
				xmax = -math.MaxFloat64
				ymin = +math.MaxFloat64
				ymax = -math.MaxFloat64
				for i := range ps {
					xmin = math.Min(xmin, ps[i].X)
					xmax = math.Max(xmax, ps[i].X)
					ymin = math.Min(ymin, ps[i].Y)
					ymax = math.Max(ymax, ps[i].Y)
				}
				return
			}
			for it, tri := range mesh.InternalModel.Triangles {
				if mesh.InternalModel.Triangles[it][0] == Removed {
					continue
				}
				xmin, xmax, ymax, ymin := box(
					mesh.InternalModel.Points[tri[0]],
					mesh.InternalModel.Points[tri[1]],
					mesh.InternalModel.Points[tri[2]],
				)
				if xmin <= p.X && p.X <= xmax &&
					ymin <= p.Y && p.Y <= ymax {
					log.Printf("Triangle %d %#v in box", it, tri)
					for s := 0; s < 3; s++ {
						log.Printf("Point %d: %#v", s, mesh.InternalModel.Points[tri[s]])
					}
				}
			}
		}
	}

	if 0 < len(materials) {
		return
	}

	// possible point is outside of triangulation
	err = fmt.Errorf("point %.3f is outside of triangulation", ps)
	return
}

// Materials indentify all triangles splitted by lines, only if points
// sliceis empty.
// If points slice is not empty, then return material mark number for
// each point
func (mesh *Mesh) Materials() (err error) {
	if Log {
		log.Printf("Materials")
	}
	defer func() {
		if err != nil {
			et := eTree.New("Materials")
			_ = et.Add(err)
			err = et
		}
	}()

	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		mesh.InternalModel.Triangles[i][3] = Undefined
	}

	marks := make([]bool, len(mesh.InternalModel.Triangles))

	var mark func(from, to, counter int) error
	mark = func(from, to, counter int) (err error) {
		if Debug {
			if to == Removed {
				err = fmt.Errorf("triangle `to` is removed")
				return
			}
		}
		// boundary
		if to == Boundary {
			return
		}
		// triangle is mark
		if marks[to] {
			return
		}
		// find line between 2 triangles
		var points []int
		points = append(points, mesh.InternalModel.Triangles[from][:3]...)
		points = append(points, mesh.InternalModel.Triangles[to][:3]...)
		sort.Ints(points)
		var uniq []int
		for i := 1; i < len(points); i++ {
			if points[i-1] == points[i] {
				uniq = append(uniq, points[i])
			}
		}
		if Debug {
			if len(uniq) != 2 {
				err = fmt.Errorf("not 2 points: %v. %v", uniq, points)
				return
			}
		}
		// have line and fixed
		for _, line := range mesh.InternalModel.Lines {
			if line[0] != uniq[0] && line[0] != uniq[1] {
				continue
			}
			if line[1] != uniq[0] && line[1] != uniq[1] {
				continue
			}
			if line[2] == Fixed {
				return
			}
		}
		if mesh.InternalModel.Triangles[to][3] != Undefined {
			err = fmt.Errorf("double mark: %v %v",
				mesh.InternalModel.Triangles[from][3],
				mesh.InternalModel.Triangles[to][3],
			)
			return
		}
		// mark
		marks[from] = true
		marks[to] = true
		mesh.InternalModel.Triangles[to][3] = counter
		for side := 0; side < 3; side++ {
			err = mark(to, mesh.Triangles[to][side], counter)
			if err != nil {
				return
			}
		}
		return nil
	}

	counter := 50
	for i := range mesh.InternalModel.Triangles {
		if mesh.InternalModel.Triangles[i][0] == Removed {
			continue
		}
		if marks[i] {
			continue
		}
		if mesh.InternalModel.Triangles[i][3] != Undefined {
			err = fmt.Errorf("unmarked undefined triangle for triangle: %#v",
				mesh.InternalModel.Triangles[i])
			return
		}
		mesh.InternalModel.Triangles[i][3] = counter
		for side := 0; side < 3; side++ {
			from := i
			to := mesh.Triangles[i][side]
			err = mark(from, to, counter)
			if err != nil {
				return
			}
		}
		counter++
	}
	return
}

// Smooth move all movable point by average distance
func (mesh *Mesh) Smooth(pts ...int) (err error) {
	if Log {
		log.Printf("Smooth")
	}
	defer func() {
		if err != nil {
			et := eTree.New("Smooth")
			_ = et.Add(fmt.Errorf("input point list: %v", pts))
			_ = et.Add(err)
			err = et
		}
	}()
	// for acceptable movable points calculate all side distances from that
	// point to points near triangles and move to average distance.
	//
	// split sides with maximal side distance

	type Store struct {
		index         int   // point index
		nearPoints    []int // index of near points
		nearTriangles []int // index of near triangles
	}
	var store []Store

	if len(pts) == 0 {
		pts = make([]int, len(mesh.InternalModel.Points))
		for i := range pts {
			pts[i] = i
		}
	}

	if len(pts) == 0 {
		err = fmt.Errorf("points list is empty")
		return
	}

	// create list of all movable points
	nearPoints := make([]int, 0, 20)
	nearTriangles := make([]int, 0, 20)
	for _, p := range pts {
		if mesh.Points[p] != Movable {
			continue
		}
		if mesh.Points[p] == Fixed {
			continue
		}
		{ // point is not on fixed line
			fix := false
			for _, line := range mesh.InternalModel.Lines {
				if line[0] != p && line[1] != p {
					continue
				}
				if line[2] == Fixed {
					fix = true
				}
			}
			if fix {
				continue
			}
		}
		// find near triangles
		nearPoints = nearPoints[:0]
		nearTriangles = nearTriangles[:0]
		for index, tri := range mesh.InternalModel.Triangles {
			if p != tri[0] && p != tri[1] && p != tri[2] {
				continue
			}
			nearPoints = append(nearPoints, tri[0:3]...)
			nearTriangles = append(nearTriangles, index)
		}
		{ // point is not on boundary triangle side
			onBoundary := false
			for _, tr := range nearTriangles {
				switch p {
				case mesh.InternalModel.Triangles[tr][0]:
					if mesh.Triangles[tr][0] == Boundary {
						onBoundary = true
					}
					if mesh.Triangles[tr][2] == Boundary {
						onBoundary = true
					}

				case mesh.InternalModel.Triangles[tr][1]:
					if mesh.Triangles[tr][0] == Boundary {
						onBoundary = true
					}
					if mesh.Triangles[tr][1] == Boundary {
						onBoundary = true
					}

				case mesh.InternalModel.Triangles[tr][2]:
					if mesh.Triangles[tr][1] == Boundary {
						onBoundary = true
					}
					if mesh.Triangles[tr][2] == Boundary {
						onBoundary = true
					}
				}
			}
			if onBoundary {
				continue
			}
		}
		if len(nearPoints) == 0 {
			continue
		}
		// uniq points
		sort.Ints(nearPoints)
		uniq := []int{nearPoints[0]}
		for i := 1; i < len(nearPoints); i++ {
			if nearPoints[i] == p {
				continue
			}
			if nearPoints[i-1] != nearPoints[i] {
				uniq = append(uniq, nearPoints[i])
			}
		}
		store = append(store, Store{
			index:         p,
			nearPoints:    uniq, // nearPoints,
			nearTriangles: nearTriangles,
		})
	}

	if len(store) == 0 {
		return
	}

	max := 1.0
	iter := 0

	for ; iter < 10 && Eps < max; iter++ {
		max = 0.0
		for _, st := range store {
			var x, y float64
			for _, n := range st.nearPoints {
				x += mesh.InternalModel.Points[n].X
				y += mesh.InternalModel.Points[n].Y
			}
			x /= float64(len(st.nearPoints))
			y /= float64(len(st.nearPoints))
			// move only if all triangles will be clockwise
			last := Point{
				X: mesh.InternalModel.Points[st.index].X,
				Y: mesh.InternalModel.Points[st.index].Y,
			}
			mesh.InternalModel.Points[st.index].X = x
			mesh.InternalModel.Points[st.index].Y = y
			isValid := true
			for _, index := range st.nearTriangles {
				if ClockwisePoints != Orientation(
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[index][0]],
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[index][1]],
					mesh.InternalModel.Points[mesh.InternalModel.Triangles[index][2]],
				) {
					isValid = false
					break
				}
			}
			if !isValid {
				mesh.InternalModel.Points[st.index] = last
				//store = append(store[:indexStore], store[indexStore+1:]...)
				continue
			}
			max = math.Max(max, Distance(last, Point{x, y}))
		}
	}
	// typically amount iter is 1
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("end of func: %v", err)
			return
		}
	}
	return
}

// Split all triangles edge on distance `factor`
func (mesh *Mesh) Split(factor float64) (err error) {
	factor = math.Abs(factor)
	if factor == 0 {
		err = fmt.Errorf("zero distance is not valid")
		return
	}
	return mesh.SplitFunc(func(p1, p2 Point) bool {
		d := Distance(p1, p2)
		return factor < d
	})
}

// SplitFunc split all triangles edge on distance only if function argument return true.
// Example of factorFunc:
//
//	factorFunc = func(p1, p2 Point) bool {
//		d := gog.Distance(p1, p2)
//		return factor < d
//	}
//
// If factorFunc is not valid, then splitting will be infinite.
func (mesh *Mesh) SplitFunc(factorFunc func(p1, p2 Point) bool) (err error) {
	if Log {
		log.Printf("Split")
	}
	defer func() {
		if err != nil {
			et := eTree.New("Split")
			_ = et.Add(err)
			err = et
		}
	}()
	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("input: %v", err)
			return
		}
	}

	// only for debug
	var chains []Point

	counter := 0
	addpoint := func(p1, p2 Point, tag int, triIndexes ...int) (err error) {
		if !factorFunc(p1, p2) {
			return
		}
		counter++
		// add middle point
		mid := MiddlePoint(p1, p2)
		// add all points of model
		defer func() {
			if err != nil {
				et := eTree.New("addpoint")
				_ = et.Add(fmt.Errorf("left  point: %.9f", p1))
				_ = et.Add(fmt.Errorf("mid   point: %.9f", mid))
				_ = et.Add(fmt.Errorf("right point: %.9f", p2))
				_ = et.Add(err)
				err = et
			}
		}()
		if Debug {
			if Orientation(p1, mid, p2) != CollinearPoints {
				et := eTree.New("MiddlePoint")
				_ = et.Add(fmt.Errorf("p1   = %.10f", p1))
				_ = et.Add(fmt.Errorf("mid  = %.10f", mid))
				_ = et.Add(fmt.Errorf("p2   = %.10f", p2))
				err = et
				return
			}
		}
		idp, err := mesh.AddPoint(mid, tag, triIndexes...)
		if err != nil {
			return
		}
		if tag == Movable {
			err = mesh.Smooth(idp)
		}
		if Debug {
			counter := 0
			var ch Point
			for i := range chains {
				if dist := Distance(chains[i], mid); dist < Eps {
					ch = chains[i]
					counter++
				}
			}
			if 200 < counter {
				err = fmt.Errorf("found same points(propably infine loop) %.9f %.9f %.9e",
					mid, ch, Distance128(mid, ch))
				return
			}
			chains = append(chains, mid)
		}
		return
	}

	iteration := func(do func() error) (err error) {
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("begin in loop: %v", err)
				return
			}
		}
		for iter := 0; ; iter++ {
			counter = 0
			// action
			if err = do(); err != nil {
				return
			}
			if iter == 10000 {
				err = fmt.Errorf("too many iterations")
				return
			}
			if counter == 0 {
				break
			}
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				err = fmt.Errorf("begin in loop: %v", err)
				return
			}
		}
		return nil
	}

	if err = iteration(func() (err error) {
		// split fixed lines
		for _, line := range mesh.InternalModel.Lines {
			if line[2] == Removed {
				continue
			}
			err = addpoint(
				mesh.InternalModel.Points[line[0]],
				mesh.InternalModel.Points[line[1]],
				Fixed,
			)
			if err != nil {
				et := eTree.New("split fixed lines")
				_ = et.Add(err)
				err = et
				return
			}
		}
		return nil
	}); err != nil {
		return
	}

	if err = iteration(func() (err error) {
		// split big triangle edges with boundary
		for i := range mesh.InternalModel.Triangles {
			if mesh.InternalModel.Triangles[i][0] == Removed {
				continue
			}
			// add point on triangle edge
			t := mesh.InternalModel.Triangles[i]
			if t[0] == Removed {
				continue
			}
			p0 := mesh.InternalModel.Points[t[0]]
			p1 := mesh.InternalModel.Points[t[1]]
			p2 := mesh.InternalModel.Points[t[2]]
			switch {
			case mesh.Triangles[i][0] == Boundary && factorFunc(p0, p1):
				err = addpoint(p0, p1, Fixed, i)
			case mesh.Triangles[i][1] == Boundary && factorFunc(p1, p2):
				err = addpoint(p1, p2, Fixed, i)
			case mesh.Triangles[i][2] == Boundary && factorFunc(p2, p0):
				err = addpoint(p2, p0, Fixed, i)
			}
			if err != nil {
				et := eTree.New("split boundary triangle edge")
				_ = et.Add(fmt.Errorf("counter: %d", counter))
				_ = et.Add(fmt.Errorf("triangle: %d", i))
				_ = et.Add(fmt.Errorf("points :%.8f", p0))
				_ = et.Add(fmt.Errorf("points :%.8f", p1))
				_ = et.Add(fmt.Errorf("points :%.8f", p2))
				_ = et.Add(err)
				err = et
				return
			}
		}
		return nil
	}); err != nil {
		return
	}

	if err = iteration(func() (err error) {
		// split big triangle edges
		for i := range mesh.InternalModel.Triangles {
			if mesh.InternalModel.Triangles[i][0] == Removed {
				continue
			}
			// add point on triangle edge
			t := mesh.InternalModel.Triangles[i]
			if t[0] == Removed {
				continue
			}
			p0 := mesh.InternalModel.Points[t[0]]
			p1 := mesh.InternalModel.Points[t[1]]
			p2 := mesh.InternalModel.Points[t[2]]
			d01 := Distance(p0, p1)
			d12 := Distance(p1, p2)
			d20 := Distance(p2, p0)
			maxd := math.Max(math.Max(d01, d12), d20)
			switch {
			case maxd == d12 && factorFunc(p1, p2):
				err = addpoint(p1, p2, Movable, i)
			case maxd == d01 && factorFunc(p0, p1):
				err = addpoint(p0, p1, Movable, i)
			case maxd == d20 && factorFunc(p2, p0):
				err = addpoint(p2, p0, Movable, i)
			}
			if err != nil {
				et := eTree.New("split big triangle edge")
				_ = et.Add(fmt.Errorf("counter: %d", counter))
				_ = et.Add(fmt.Errorf("triangle: %d", i))
				_ = et.Add(fmt.Errorf("points :%.8f", p0))
				_ = et.Add(fmt.Errorf("points :%.8f", p1))
				_ = et.Add(fmt.Errorf("points :%.8f", p2))
				_ = et.Add(fmt.Errorf("distance 01 :%.8f", d01))
				_ = et.Add(fmt.Errorf("distance 12 :%.8f", d12))
				_ = et.Add(fmt.Errorf("distance 20 :%.8f", d20))
				_ = et.Add(err)
				err = et
				return
			}
		}
		return nil
	}); err != nil {
		return
	}

	//
	// Delanay is no need, because inside function AddPoint
	//
	//	err = mesh.Delanay()
	//	if err != nil {
	//		err = fmt.Errorf("at Delanay: %v", err)
	//		return
	//	}

	err = mesh.Smooth()
	if err != nil {
		return
	}

	if Debug {
		err = mesh.Check()
		if err != nil {
			err = fmt.Errorf("at the end: %v", err)
			return
		}
	}

	return
}

// AddLine is add line in triangulation with tag
func (mesh *Mesh) AddLine(inp1, inp2 Point) (err error) {
	if Log {
		log.Printf("AddLine")
	}
	defer func() {
		if err != nil {
			et := eTree.New("AddLine")
			_ = et.Add(err)
			err = et
		}
	}()

	if SamePoints(inp1, inp2) {
		err = fmt.Errorf("AddLine: points are same")
		return
	}
	var list []int
	{
		// add points of points
		var idp1 int
		idp1, err = mesh.AddPoint(inp1, Fixed)
		if err != nil {
			et := eTree.New("add p1")
			_ = et.Add(err)
			err = et
			return
		}
		var idp2 int
		idp2, err = mesh.AddPoint(inp2, Fixed)
		if err != nil {
			et := eTree.New("add p2")
			_ = et.Add(err)
			err = et
			return
		}
		// put fixed lines
		mesh.InternalModel.AddLine(inp1, inp2, Fixed)
		// triangle edges on line
		list = []int{idp1, idp2}
	}
	// triangle edges on line
again:
	if Debug {
		for _, idp := range list {
			found := false
			for _, tps := range mesh.InternalModel.Triangles {
				if tps[0] == Removed {
					continue
				}
				for i := range 3 {
					if tps[i] == idp {
						found = true
					}
				}
				if found {
					break
				}
			}
			if !found {
				err = fmt.Errorf("no found point %d in triangle list", idp)
				panic(err)
			}
		}
		for _, idp := range list {
			for _, tps := range mesh.InternalModel.Triangles {
				if tps[0] == Removed {
					continue
				}
				found := false
				for i := range 3 {
					if idp != tps[i] {
						continue
					}
					if SamePoints(mesh.InternalModel.Points[idp], mesh.InternalModel.Points[tps[i]]) {
						found = true
					}
				}
				if found {
					continue
				}
				res, lineIntersect, err := TriangleSplitByPoint(
					mesh.InternalModel.Points[idp],
					mesh.InternalModel.Points[tps[0]],
					mesh.InternalModel.Points[tps[1]],
					mesh.InternalModel.Points[tps[2]],
				)
				if err != nil {
					continue
				}
				if len(res) == 0 {
					continue
				}
				panic(fmt.Errorf(":RRR: %v %v", lineIntersect, res))
			}
		}
	}
	for i := 1; i < len(list); i++ {
		// find triangle with that points
		idp1 := list[i-1]
		idp2 := list[i]
		if idp1 == idp2 {
			et := eTree.New("equal point index")
			_ = et.Add(fmt.Errorf("idp = %d", idp1))
			_ = et.Add(fmt.Errorf("list = %v", list))
			err = et
			return
		}
		if SamePoints(mesh.InternalModel.Points[idp1], mesh.InternalModel.Points[idp2]) {
			et := eTree.New("same point")
			_ = et.Add(fmt.Errorf("idp = %d", idp1))
			_ = et.Add(fmt.Errorf("list = %v", list))
			err = et
			return
		}
		{
			found := false
			for _, tri := range mesh.InternalModel.Triangles {
				if tri[0] == Removed {
					continue
				}
				if (idp1 == tri[0] || idp1 == tri[1] || idp1 == tri[2]) &&
					(idp2 == tri[0] || idp2 == tri[1] || idp2 == tri[2]) {
					found = true
					break
				}
			}
			if found {
				continue
			}
		}
		{
			found := false
			for _, tri := range mesh.InternalModel.Triangles {
				if tri[0] == Removed {
					continue
				}
				start := idp1
				last := idp2
				var pi []Point
				var stA, stB State
				if idp1 == tri[0] {
					pi, stA, stB = LineLine(
						mesh.InternalModel.Points[start],
						mesh.InternalModel.Points[last],
						mesh.InternalModel.Points[tri[1]],
						mesh.InternalModel.Points[tri[2]],
					)
				}
				if idp1 == tri[1] {
					pi, stA, stB = LineLine(
						mesh.InternalModel.Points[start],
						mesh.InternalModel.Points[last],
						mesh.InternalModel.Points[tri[0]],
						mesh.InternalModel.Points[tri[2]],
					)
				}
				if idp1 == tri[2] {
					pi, stA, stB = LineLine(
						mesh.InternalModel.Points[start],
						mesh.InternalModel.Points[last],
						mesh.InternalModel.Points[tri[0]],
						mesh.InternalModel.Points[tri[1]],
					)
				}
				if len(pi) != 1 {
					continue
				}
				if stA.Has(OnSegment) || stB.Has(OnSegment) {
					mid := pi[0]
					var idp int
					idp, err = mesh.AddPoint(mid, Fixed)
					if err != nil {
						et := eTree.New("add intersection")
						_ = et.Add(fmt.Errorf("mid = %e", mid))
						_ = et.Add(fmt.Errorf("len(list) = %d", len(list)))
						_ = et.Add(fmt.Errorf("list = %v", list))
						_ = et.Add(err)
						err = et
						return
					}
					list = append(list[:i], append([]int{idp}, list[i:]...)...)
					found = true
					goto again
				}
			}
			if found {
				continue
			}
		}
		var (
			p1  = mesh.InternalModel.Points[idp1]
			p2  = mesh.InternalModel.Points[idp2]
			mid = MiddlePoint(p1, p2)
			idp int
		)
		idp, err = mesh.AddPoint(mid, Fixed)
		if err != nil {
			et := eTree.New("add mid")
			_ = et.Add(fmt.Errorf("mid = %e", mid))
			_ = et.Add(fmt.Errorf("len(list) = %d", len(list)))
			_ = et.Add(fmt.Errorf("list = %v", list))
			_ = et.Add(err)
			err = et
			return
		}
		list = append(list[:i], append([]int{idp}, list[i:]...)...)
		if 1000 < len(list) {
			err = fmt.Errorf("too big list")
			return
		}
		goto again
	}
	return
}
