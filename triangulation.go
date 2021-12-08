package gog

import (
	"bytes"
	"fmt"
)

type Mesh struct {
	model     Model
	Triangles []Triangle
	// TODO
}

const (
	Boundary  = -1
	Removed   = -2
	Undefined = -3
)

func New(model Model) (mesh *Mesh, err error) {
	// create a new Mesh
	mesh = new(Mesh)
	// convex
	cps := ConvexHull(model.Points) // points on convex hull
	// prepare mesh triangles
	for i := 3; i < len(cps); i++ {
		mesh.model.AddTriangle(cps[0], cps[i-2], cps[i-1], Boundary)
		if i == 3 {
			mesh.Triangles = append(mesh.Triangles, Triangle{
				tr: [3]int{Boundary, Boundary, 1},
			})
		} else {
			mesh.Triangles = append(mesh.Triangles, Triangle{
				tr: [3]int{i - 4, Boundary, i - 2},
			})
		}
	}
	// last not exist triangle and mark as boundary
	mesh.Triangles[len(mesh.Triangles)-1].tr[2] = Boundary
	// clockwise all triangles
	mesh.Clockwise()
	// add all points of model
	for i := range model.Points {
		// TODO remove
		err = mesh.Check()
		if err != nil {
			return
		}
		err = mesh.AddPoint(model.Points[i])
		if err != nil {
			return
		}
		// TODO remove
		err = mesh.Check()
		if err != nil {
			return
		}
		// TODO remove
		// if i == 4 {
		// 	break
		// }
	}
	// TODO remove
	err = mesh.Check()
	if err != nil {
		return
	}
	return
}

func (mesh Mesh) Check() error {
	// amount of triangles
	if len(mesh.model.Triangles) != len(mesh.Triangles) {
		return fmt.Errorf("sizes is not same")
	}
	// undefined triangles
	for i := range mesh.model.Triangles {
		for j := 0; j < 3; j++ {
			if mesh.model.Triangles[i][j] == Undefined {
				return fmt.Errorf("undefined point of triangle")
			}
			if mesh.Triangles[i].tr[j] == Undefined {
				return fmt.Errorf("undefined triangle of triangle")
			}
		}
	}
	// clockwise triangles
	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		or := Orientation(
			mesh.model.Points[mesh.model.Triangles[i][0]],
			mesh.model.Points[mesh.model.Triangles[i][1]],
			mesh.model.Points[mesh.model.Triangles[i][2]],
		)
		if or != ClockwisePoints {
			return fmt.Errorf("not clockwise: %d. IsCounterClock %v",
				i,
				or == CounterClockwisePoints)
		}
	}
	// same triangles - self linked
	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i].tr[j] == i {
				return fmt.Errorf("self linked triangle: %v %d", mesh.Triangles[i].tr, i)
			}
		}
	}
	// double triangles
	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		tri := mesh.Triangles[i].tr
		if tri[0] == tri[1] && tri[0] != Boundary {
			return fmt.Errorf("double triangles 0: %v", tri)
		}
		if tri[1] == tri[2] && tri[1] != Boundary {
			return fmt.Errorf("double triangles 1: %v", tri)
		}
		if tri[2] == tri[0] && tri[2] != Boundary {
			return fmt.Errorf("double triangles 2: %v", tri)
		}
	}
	// near triangle
	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i].tr[j] == Boundary {
				continue
			}
			if i != mesh.Triangles[mesh.Triangles[i].tr[j]].tr[0] &&
				i != mesh.Triangles[mesh.Triangles[i].tr[j]].tr[1] &&
				i != mesh.Triangles[mesh.Triangles[i].tr[j]].tr[2] {
				return fmt.Errorf("near-near triangle: %d, %v, %v",
					i,
					mesh.Triangles[i].tr,
					mesh.Triangles[mesh.Triangles[i].tr[j]].tr,
				)
			}
		}
	}
	// near triangles
	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		for j := 0; j < 3; j++ {
			if mesh.Triangles[i].tr[j] == Boundary {
				continue
			}
			arr := mesh.model.Triangles[mesh.Triangles[i].tr[j]]
			found := false
			if arr[0] == mesh.model.Triangles[i][j] ||
				arr[1] == mesh.model.Triangles[i][j] ||
				arr[2] == mesh.model.Triangles[i][j] {
				found = true
			}
			if !found {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "triangle have not same side.\n")
				fmt.Fprintf(&buf, "tr   = %d\n", i)
				fmt.Fprintf(&buf, "side = %d\n", j)
				fmt.Fprintf(&buf, "triangle points      = %+2d\n", mesh.model.Triangles[i])
				fmt.Fprintf(&buf, "link triangles       = %+2d\n", mesh.Triangles[i].tr)
				fmt.Fprintf(&buf, "near triangle points = %+2d\n", mesh.model.Triangles[mesh.Triangles[i].tr[j]])
				fmt.Fprintf(&buf, "1: %v\n", mesh.model.Triangles)
				fmt.Fprintf(&buf, "2: %v\n", mesh.Triangles)
				return fmt.Errorf(buf.String())
			}
		}
	}

	// TODO

	// no error
	return nil
}

func (model *Model) Get(mesh *Mesh) {
	for _, tr := range mesh.model.Triangles {
		if tr[0] == Removed {
			continue
		}
		model.AddTriangle(
			mesh.model.Points[tr[0]],
			mesh.model.Points[tr[1]],
			mesh.model.Points[tr[2]],
			0,
		)
	}
}

func (mesh *Mesh) Clockwise() {
	for i := range mesh.model.Triangles {
		switch Orientation(
			mesh.model.Points[mesh.model.Triangles[i][0]],
			mesh.model.Points[mesh.model.Triangles[i][1]],
			mesh.model.Points[mesh.model.Triangles[i][2]],
		) {
		case CounterClockwisePoints:
			mesh.Triangles[i].tr[0], mesh.Triangles[i].tr[2] =
				mesh.Triangles[i].tr[2], mesh.Triangles[i].tr[0]
			mesh.model.Triangles[i][1], mesh.model.Triangles[i][2] =
				mesh.model.Triangles[i][2], mesh.model.Triangles[i][1]
		case CollinearPoints:
			panic(fmt.Errorf("collinear triangle: %#v", mesh.model))
		}
	}
}

func (mesh *Mesh) AddPoint(p Point) (err error) {
	fmt.Println("add point ", p)
	defer func() {
		fmt.Println("end add point ", p)
		if err != nil {
			err = fmt.Errorf("Add point: %v\n%v", p, err)
		}
	}()
	// ignore points if on corner
	for _, pt := range mesh.model.Points {
		if Distance(p, pt) < Eps {
			return
		}
	}

	// TODO : add to delanay flip linked list
	size := len(mesh.Triangles)
	for i := 0; i < size; i++ {
		// ignore removed triangle
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		if mesh.model.Triangles[i][1] == Removed {
			continue
		}
		if mesh.model.Triangles[i][2] == Removed {
			continue
		}
		// split triangle
		res, state, err := TriangleSplitByPoint(
			p,
			mesh.model.Points[mesh.model.Triangles[i][0]],
			mesh.model.Points[mesh.model.Triangles[i][1]],
			mesh.model.Points[mesh.model.Triangles[i][2]],
		)
		if err != nil {
			panic(err)
		}
		if len(res) == 0 {
			continue
		}
		// index of new point
		idp := mesh.model.AddPoint(p)
		// removed triangles
		removedTriangles := []int{i}

		// repair near triangles

		// find intersect side and near triangle if exist
		// only for point on side
		if len(res) == 2 {
			// point on some line
			if mesh.Triangles[i].tr[state] != Boundary {
				removedTriangles = append(removedTriangles, mesh.Triangles[i].tr[state])
				err = mesh.repairTriangles(idp, removedTriangles, 100+state)
				if err != nil {
					err = fmt.Errorf("%v\npreliminary split triangles: %v", err, res)
					return err
				}
				return nil
			}
			return mesh.repairTriangles(idp, removedTriangles, 200+state)
		}
		return mesh.repairTriangles(idp, removedTriangles, 300)
	}
	// outside of triangles or on corners
	return nil
}

// func (mesh *Mesh) addPointInTriangle(
// 	p Point,
// 	triangle int,
// ) {
// 	// index of new point
// 	ip := m.model.AddPoint(p)
// 	// create a new triangles
// 	m.model.AddTriangle(
// 		m.model.Points[m.Triangles[triangle].tr[0]],
// 		p,
// 		m.model.Points[m.Triangles[triangle].tr[1]],
// 		Undefined,
// 	)
// 	mesh.Triangles = append(mesh.Triangles, Triangle{
// 		nodes: [3]int{p0, p1, p2},
// 		tr:    [3]int{Boundary, Boundary, 1},
// 	})
// 	mesh.model.AddTriangle(
// 		m.model.Points[m.Triangles[triangle].tr[1]],
// 		p,
// 		m.model.Points[m.Triangles[triangle].tr[2]],
// 		Undefined,
// 	)
// 	mesh.model.AddTriangle(
// 		m.model.Points[m.Triangles[triangle].tr[2]],
// 		p,
// 		m.model.Points[m.Triangles[triangle].tr[0]],
// 		Undefined,
// 	)
// }

// shiftTriangle is shift numbers on right on one
func (mesh *Mesh) shiftTriangle(i int) {
	mesh.Triangles[i].tr[0], mesh.Triangles[i].tr[1], mesh.Triangles[i].tr[2] =
		mesh.Triangles[i].tr[2], mesh.Triangles[i].tr[0], mesh.Triangles[i].tr[1]
	mesh.model.Triangles[i][0], mesh.model.Triangles[i][1], mesh.model.Triangles[i][2] =
		mesh.model.Triangles[i][2], mesh.model.Triangles[i][0], mesh.model.Triangles[i][1]
}

// repairTriangles
//
// ap is index of added point
//
// rt is removed triangles
//
// state:
//	100 - point on line with 2 triangles
//	200 - point on line with 1 boundary triangle
//	300 - point in triangle
func (mesh *Mesh) repairTriangles(ap int, rt []int, state int) error {
	fmt.Println("repairTriangles ", ap, rt, state)
	// fmt.Println("mesh ", mesh)

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
		from, to Point // point of triangle base
		// left, right int   // triangles index at left anf right
		in, out int // inside/outside index triangles
	}
	var chains []chain
	tc := [2]int{Undefined, Undefined} // index of corner triangle

	// amount triangles before added
	size := len(mesh.Triangles)

	// create chain
	switch state {
	case 100:
		// debug checking
		if len(rt) != 2 {
			return fmt.Errorf("removed triangles: %v", rt)
		}
		// debug: point in not on line
		{
			_, _, stB := PointLine(
				mesh.model.Points[ap],
				mesh.model.Points[mesh.model.Triangles[rt[0]][0]],
				mesh.model.Points[mesh.model.Triangles[rt[0]][1]],
			)
			if stB.Not(OnSegment) {
				panic("point is not on line")
			}
		}
		// rotate second triangle to point intersect line 0
		if mesh.model.Triangles[rt[0]][0] != mesh.model.Triangles[rt[1]][1] ||
			mesh.model.Triangles[rt[0]][1] != mesh.model.Triangles[rt[1]][0] {
			mesh.shiftTriangle(rt[1])

			// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>")
			// var err bytes.Buffer
			// fmt.Fprintf(&err, "cannot find first point: %v %v\n",
			// 	mesh.model.Triangles[rt[0]],
			// 	mesh.model.Triangles[rt[1]],
			// )
			// fmt.Fprintf(&err, "coordinates of first point\n")
			// for p := 0; p < 3; p++ {
			// 	fmt.Fprintf(&err, "%.3f\n",
			// 		mesh.model.Points[mesh.model.Triangles[rt[0]][p]])
			// }
			// fmt.Fprintf(&err, "coordinates of second point\n")
			// for p := 0; p < 3; p++ {
			// 	fmt.Fprintf(&err, "%.3f\n",
			// 		mesh.model.Points[mesh.model.Triangles[rt[1]][p]])
			// }
			// fmt.Println(err.String())

			return mesh.repairTriangles(ap, rt, state)
		}
		// debug checking
		if a := mesh.model.Triangles[rt[0]][0]; a != mesh.model.Triangles[rt[1]][0] &&
			a != mesh.model.Triangles[rt[1]][1] &&
			a != mesh.model.Triangles[rt[1]][2] {
			var err bytes.Buffer
			fmt.Fprintf(&err, "cannot find first point: %v %v\n",
				mesh.model.Triangles[rt[0]],
				mesh.model.Triangles[rt[1]],
			)
			fmt.Fprintf(&err, "coordinates of first point\n")
			for p := 0; p < 3; p++ {
				fmt.Fprintf(&err, "%.3f\n",
					mesh.model.Points[mesh.model.Triangles[rt[0]][p]])
			}
			fmt.Fprintf(&err, "coordinates of second point\n")
			for p := 0; p < 3; p++ {
				fmt.Fprintf(&err, "%.3f\n",
					mesh.model.Points[mesh.model.Triangles[rt[1]][p]])
			}
			return fmt.Errorf(err.String())
		}
		if a := mesh.model.Triangles[rt[0]][1]; a != mesh.model.Triangles[rt[1]][0] &&
			a != mesh.model.Triangles[rt[1]][1] &&
			a != mesh.model.Triangles[rt[1]][2] {
			return fmt.Errorf("cannot find second point: %v %v",
				mesh.model.Triangles[rt[0]],
				mesh.model.Triangles[rt[1]],
			)
		}
		// point on triangle 0 line 0 with 2 triangles
		chains = []chain{{
			from: mesh.model.Points[mesh.model.Triangles[rt[1]][1]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[1]][2]],
			in:   size,
			out:  mesh.Triangles[rt[1]].tr[1],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[1]][2]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[1]][0]],
			in:   size + 1,
			out:  mesh.Triangles[rt[1]].tr[2],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[1]][0]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			in:   size + 2,
			out:  mesh.Triangles[rt[0]].tr[1],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][0]],
			in:   size + 3,
			out:  mesh.Triangles[rt[0]].tr[2],
		}}
		tc = [2]int{size, size + 3}

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
			return fmt.Errorf("removed triangles: %v", rt)
		}
		if mesh.Triangles[rt[0]].tr[0] != Boundary {
			return fmt.Errorf("not valid boundary")
		}
		// point on triangle boundary line 0
		chains = []chain{{
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][1]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			in:   size,
			out:  mesh.Triangles[rt[0]].tr[1],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][0]],
			in:   size + 1,
			out:  mesh.Triangles[rt[0]].tr[2],
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
			return fmt.Errorf("removed triangles: %v", rt)
		}
		// point in triangle
		chains = []chain{{
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][0]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][1]],
			in:   size,
			out:  mesh.Triangles[rt[0]].tr[0],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][1]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			in:   size + 1,
			out:  mesh.Triangles[rt[0]].tr[1],
		}, {
			from: mesh.model.Points[mesh.model.Triangles[rt[0]][2]],
			to:   mesh.model.Points[mesh.model.Triangles[rt[0]][0]],
			in:   size + 2,
			out:  mesh.Triangles[rt[0]].tr[2],
		}}
		tc = [2]int{size, size + 2}

	default:
		return fmt.Errorf("not clear state %v", state)
	}
	// debug checking
	if tc[0] == Undefined || tc[1] == Undefined {
		panic(fmt.Errorf("undefined corner triangle"))
	}
	fmt.Println("chains ", chains)

	// create triangles
	for i := range chains {
		mesh.model.AddTriangle(
			chains[i].from,
			chains[i].to,
			mesh.model.Points[ap],
			Undefined, // TODO for case with 2 triangles - not clear tag
		)
		tr := [3]int{Undefined, Undefined, Undefined}

		tr[0] = chains[i].out
		if chains[i].out != Boundary {
			for k := range mesh.Triangles[chains[i].out].tr {
				for _, rem := range rt {
					if mesh.Triangles[chains[i].out].tr[k] == rem {
						mesh.Triangles[chains[i].out].tr[k] = chains[i].in
					}
				}
			}
		}
		if i == len(chains)-1 {
			tr[1] = tc[0]
		} else {
			tr[1] = chains[i+1].in
		}
		if i == 0 {
			tr[2] = tc[1]
		} else {
			tr[2] = chains[i-1].in
		}
		mesh.Triangles = append(mesh.Triangles, Triangle{tr: tr})
	}

	// remove triangles
	for _, rem := range rt {
		mesh.model.Triangles[rem][0] = Removed
		mesh.model.Triangles[rem][1] = Removed
		mesh.model.Triangles[rem][2] = Removed
		mesh.Triangles[rem].tr[0] = Removed
		mesh.Triangles[rem].tr[1] = Removed
		mesh.Triangles[rem].tr[2] = Removed
	}

	// ntr - near triangles
	// 	ntr := []int{}
	// 	for _, rt := range removedTriangles {
	// 		ntr = append(ntr,
	// 			mesh.Triangles[rt].tr[0],
	// 			mesh.Triangles[rt].tr[1],
	// 			mesh.Triangles[rt].tr[2],
	// 		)
	// 		for i := 0; i < 3; i++ {
	// 			for j := 0; j < 3; j++ {
	// 				if mesh.Triangles[rt].tr[i] == Boundary {
	// 					continue
	// 				}
	// 				tr := &mesh.Triangles[mesh.Triangles[rt].tr[i]].tr[j]
	// 				if *tr == rt {
	// 					*tr = Undefined
	// 				}
	// 			}
	// 		}
	// 	}
	// 	// remove boundary near triangles
	// 	for i := len(ntr) - 1; 0 <= i; i-- {
	// 		if ntr[i] == Boundary {
	// 			ntr = append(ntr[:i], ntr[i+1:]...)
	// 		}
	// 	}
	// 	// create a new triangle only if ot not on one line
	// 	for _, rem := range removedTriangles {
	// 		for _, n := range [][3]int{
	// 			{mesh.model.Triangles[rem][0], pointIndex, mesh.model.Triangles[rem][1]},
	// 			{mesh.model.Triangles[rem][1], pointIndex, mesh.model.Triangles[rem][2]},
	// 			{mesh.model.Triangles[rem][2], pointIndex, mesh.model.Triangles[rem][0]},
	// 		} {
	// 			if Orientation(
	// 				mesh.model.Points[n[0]], mesh.model.Points[n[1]], mesh.model.Points[n[2]],
	// 			) == CollinearPoints {
	// 				// do not add triangles with points on one line
	// 				continue
	// 			}
	// 			id := mesh.model.AddTriangle(
	// 				mesh.model.Points[n[0]], mesh.model.Points[n[1]], mesh.model.Points[n[2]],
	// 				Undefined,
	// 			)
	// 			mesh.Triangles = append(mesh.Triangles, Triangle{
	// 				// nodes: [3]int{n[0], n[1], n[2]},
	// 				tr: [3]int{Undefined, Undefined, Undefined},
	// 			})
	// 			ntr = append(ntr, id)
	// 		}
	// 	}
	// 	// repair near triangles
	// 	// Example:
	// 	//	triangle      point
	// 	//	index         indexes
	// 	//	1             0  1  2
	// 	//	2
	// 	repair := func(
	// 		ti int, sidei int, si [2]int, // triangle index and side
	// 		tj int, sidej int, sj [2]int, // triangle index and side
	// 	) {
	// 		if !((Distance(
	// 			mesh.model.Points[mesh.model.Triangles[ti][si[0]]],
	// 			mesh.model.Points[mesh.model.Triangles[tj][sj[0]]],
	// 		) < Eps && Distance(
	// 			mesh.model.Points[mesh.model.Triangles[ti][si[1]]],
	// 			mesh.model.Points[mesh.model.Triangles[tj][sj[1]]],
	// 		) < Eps) || (Distance(
	// 			mesh.model.Points[mesh.model.Triangles[ti][si[0]]],
	// 			mesh.model.Points[mesh.model.Triangles[tj][sj[1]]],
	// 		) < Eps && Distance(
	// 			mesh.model.Points[mesh.model.Triangles[ti][si[1]]],
	// 			mesh.model.Points[mesh.model.Triangles[tj][sj[0]]],
	// 		) < Eps)) {
	// 			return
	// 		}
	// 		// find valid triangles
	// 		tr := []int{ti, tj}
	// 	again:
	// 		for _, rem := range removedTriangles {
	// 			for j := range tr {
	// 				if rem == tr[j] {
	// 					tr = append(tr[:j], tr[j+1:]...)
	// 					goto again
	// 				}
	// 			}
	// 		}
	// 		if len(tr) == 0 {
	// 			return
	// 		}
	// 		if 1 < len(tr) {
	// 			// 			var debug bytes.Buffer
	// 			// 			fmt.Fprintf(&debug, "2 triangles\n")
	// 			// 			fmt.Fprintf(&debug, "point %d triangle: %v\n",ti, mesh.model.Triangles[ti])
	// 			// 			fmt.Fprintf(&debug, "point %d triangle: %v\n",tj, mesh.model.Triangles[tj])
	// 			// 			fmt.Fprintf(&debug, "removed triangles:\n")
	// 			// 			for _,rem := range removedTriangles {
	// 			// 				fmt.Fprintf(&debug, "point %d triangle: %v\n",
	// 			// 					rem , mesh.model.Triangles[rem])
	// 			// 			}
	// 			// 			panic(debug.String())
	// 			return
	// 		}
	// 		// 2 triangles have same side
	// 		mesh.Triangles[ti].tr[sidei] = tr[0]
	// 		mesh.Triangles[tj].tr[sidej] = tr[0]
	// 	}
	// 	sides := [][2]int{{0, 1}, {1, 2}, {2, 0}}
	// 	for _, ti := range ntr {
	// 		for _, tj := range ntr {
	// 			if ti == tj {
	// 				continue
	// 			}
	// 			for sidei, si := range sides {
	// 				for sidej, sj := range sides {
	// 					repair(ti, sidei, si, tj, sidej, sj)
	// 				}
	// 			}
	// 		}
	// 	}
	// 	// remove triangles
	// 	for _, rem := range removedTriangles {
	// 		for i := range mesh.model.Triangles[rem] {
	// 			mesh.model.Triangles[rem][i] = Removed
	// 		}
	// 	}
	return nil
}

func (mesh *Mesh) AddSide() {
	// TODO
}

func (mesh *Mesh) Delanay() (err error) {
	// triangle is success by delanay, if all points is outside of circle
	// from 3 triangle points
	delanay := func(tr, side int) (flip bool, err error) {
		if mesh.model.Triangles[tr][0] == Removed {
			return
		}
		neartr := mesh.Triangles[tr].tr[side]
		if neartr == Boundary {
			return
		}
		for iter := 0; ; iter++ {
			if iter == 50 {
				err = fmt.Errorf("delanay infinite loop 1")
				return
			}
			if mesh.model.Triangles[tr][side] == mesh.model.Triangles[neartr][0] {
				break
			}
			mesh.shiftTriangle(neartr)
		}
		if PointInCircle(
			mesh.model.Points[mesh.model.Triangles[neartr][1]],
			[3]Point{
				mesh.model.Points[mesh.model.Triangles[tr][0]],
				mesh.model.Points[mesh.model.Triangles[tr][1]],
				mesh.model.Points[mesh.model.Triangles[tr][2]],
			},
		) {
			// flip
			flip = true
			for iter := 0; ; iter++ {
				if iter == 50 {
					err = fmt.Errorf("delanay infinite loop 2")
					return
				}
				if mesh.model.Triangles[tr][0] == mesh.model.Triangles[neartr][0] {
					break
				}
				mesh.shiftTriangle(tr)
			}
			// 			fmt.Println(">>>>>>", tr, neartr)
			// 			fmt.Println(">>>>>>", mesh.model.Triangles[tr], mesh.model.Triangles[neartr])
			// 			fmt.Println("++++++", mesh.Triangles[tr].tr, mesh.Triangles[neartr].tr)
			// 			fmt.Println(":::", Orientation(
			// 				mesh.model.Points[mesh.model.Triangles[tr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][2]],
			// 			) == ClockwisePoints)
			// 			fmt.Println(":::", Orientation(
			// 				mesh.model.Points[mesh.model.Triangles[neartr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][2]],
			// 			) == ClockwisePoints)
			//
			// 			fmt.Println(
			// 				mesh.model.Points[mesh.model.Triangles[tr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][2]],
			// 			)
			// 			fmt.Println(
			// 				mesh.model.Points[mesh.model.Triangles[neartr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][2]],
			// 			)
			// 			fmt.Println("flip")

			swap := func(elem, from, to int) {
				if elem == Boundary {
					return
				}
				// fmt.Println(">", elem, from, to)
				for h := 0; h < 3; h++ {
					if from == mesh.Triangles[elem].tr[h] {
						mesh.Triangles[elem].tr[h] = to
					}
				}
			}

			switch {
			case mesh.model.Triangles[tr][1] == mesh.model.Triangles[neartr][2]:
				// 				fmt.Println(">>> 1")
				{
					// flip points
					red := &mesh.model.Triangles[tr]
					blu := &mesh.model.Triangles[neartr]
					red[1], blu[0] =
						blu[1], red[2]
				}
				{
					// flip near triangles
					red := &mesh.Triangles[tr].tr
					blu := &mesh.Triangles[neartr].tr
					swap(red[1], tr, neartr)
					swap(blu[0], neartr, tr)
					red[0], red[1], blu[0], blu[2] =
						blu[0], red[0], blu[2], red[1]
				}

			case mesh.model.Triangles[tr][2] == mesh.model.Triangles[neartr][1]:
				// 				fmt.Println(">>> 2")
				{
					// flip points
					red := &mesh.model.Triangles[tr]
					blu := &mesh.model.Triangles[neartr]
					red[0], blu[1] =
						blu[2], red[1]
				}
				{
					// flip near triangles
					red := &mesh.Triangles[tr].tr
					blu := &mesh.Triangles[neartr].tr
					swap(red[0], tr, neartr)
					swap(blu[1], neartr, tr)
					red[0], red[2], blu[0], blu[1] =
						red[2], blu[1], red[0], blu[0]
				}

			default:
				panic("not clear")
			}
			// 			fmt.Println(">>>>>>", mesh.model.Triangles[tr], mesh.model.Triangles[neartr])
			// 			fmt.Println("++++++", mesh.Triangles[tr].tr, mesh.Triangles[neartr].tr)
			// 			fmt.Println(":::", Orientation(
			// 				mesh.model.Points[mesh.model.Triangles[tr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[tr][2]],
			// 			) == ClockwisePoints)
			// 			fmt.Println(":::", Orientation(
			// 				mesh.model.Points[mesh.model.Triangles[neartr][0]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][1]],
			// 				mesh.model.Points[mesh.model.Triangles[neartr][2]],
			// 			) == ClockwisePoints)
			// 			fmt.Println("}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}")
		}
		return
	}
	// TODO remove
	err = mesh.Check()
	if err != nil {
		return
	}

	// loop of triangles
	for iter := 0; ; iter++ {
		counter := 0
		flipTr := []int{}
		for tr := range mesh.model.Triangles {
			if mesh.model.Triangles[tr][0] == Removed {
				continue
			}
			// fmt.Println("> flip counter ", counter)
			var flip bool
			for side := 0; side < 3; side++ {
				flip, err = delanay(tr, side)
				if err != nil {
					return
				}
				// TODO remove
				err = mesh.Check()
				if err != nil {
					return
				}
				if flip {
					flipTr = append(flipTr, tr)
					counter++
					break
				}
			}
		}
		fmt.Println(">>", flipTr)
		if counter == 0 {
			break
		}
		if iter == 5000 {
			return fmt.Errorf("global delanay infinite loop")
		}
		// TODO remove
		err = mesh.Check()
		if err != nil {
			return
		}
	}
	return nil
}

func (mesh *Mesh) Smooth() {
	// for acceptable movable points calculate all side distances from that
	// point to points near triangles and move to average distance.
	//
	// split sides with maximal side distance
	// TODO
}

func (mesh *Mesh) MaxArea() {
	// TODO
}

func (mesh *Mesh) MinAngle() {
	//
	// TODO
}

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
//
type Triangle struct {
	//nodes [3]int // indexes of triangle points
	// ribs  [3]int // indexes of triangle ribs
	tr [3]int // indexes of near triangles
}

// func (t *Triangle) swap() {
// 	// 	t.nodes[0], t.nodes[1] = t.nodes[1], t.nodes[0]
// 	// 	t.ribs[1], t.ribs[2] = t.ribs[2], t.ribs[1]
// 	t.tr[1], t.tr[2] = t.tr[2], t.tr[1]
// }
