package gog

import (
	"bytes"
	"fmt"
	"math"
	"sort"
)

type Mesh struct {
	model     Model
	Points    []int // tags for points
	Triangles []Triangle
	// TODO
}

var Debug = false

const (
	Boundary  = -1
	Removed   = -2
	Undefined = -3
	Fixed     = 100
	Movable   = 200
)

func New(model Model) (mesh *Mesh, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("New: %v", err)
		}
	}()
	// create a new Mesh
	mesh = new(Mesh)
	// convex
	cps := ConvexHull(model.Points) // points on convex hull
	if len(cps) < 3 {
		err = fmt.Errorf("not enought points for convex")
		return
	}
	// add last point for last triangle
	cps = append(cps, cps[0])
	// prepare mesh triangles
	for i := 3; i < len(cps); i++ {
		// TODO : triangles is not boundary, side is boundary
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
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
		err = mesh.AddPoint(model.Points[i], Fixed)
		if err != nil {
			return
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
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
			return
		}
	}
	// add fixed tags
	if Debug {
		if len(mesh.Points) != len(mesh.model.Points) {
			err = fmt.Errorf("not equal points size")
			return
		}
	}
	for i := range model.Points {
		mesh.Points[i] = Fixed
	}

	// add fixed lines
	for i := range model.Lines {
		err = mesh.AddLine(
			model.Points[model.Lines[i][0]],
			model.Points[model.Lines[i][1]],
			Fixed,
		)
		if err != nil {
			return
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
	}

	return
}

func (mesh Mesh) Check() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Check: %v", err)
		}
	}()
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
	// correct remove
	for i := range mesh.Triangles {
		if mesh.Triangles[i].tr[0] == Removed && mesh.model.Triangles[i][0] != Removed {
			return fmt.Errorf("uncorrect removing")
		}
	}
	// double triangles
	for i := range mesh.Triangles {
		if mesh.Triangles[i].tr[0] == Removed {
			continue
		}
		tri := mesh.Triangles[i].tr
		if tri[0] == tri[1] && tri[0] != Boundary {
			return fmt.Errorf("double triangles 0: %d %v %v", i, tri, mesh.Triangles[tri[0]])
		}
		if tri[1] == tri[2] && tri[1] != Boundary {
			return fmt.Errorf("double triangles 1: %d %v %v", i, tri, mesh.Triangles[tri[1]])
		}
		if tri[2] == tri[0] && tri[2] != Boundary {
			return fmt.Errorf("double triangles 2: %d %v %v", i, tri, mesh.Triangles[tri[2]])
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

	// TODO add error for undefined mesh.Points

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
			tr[3],
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

func (mesh *Mesh) AddPoint(p Point, tag int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("AddPoint: %v", err)
		}
	}()
	// ignore points if on corner
	for _, pt := range mesh.model.Points {
		if Distance(p, pt) < Eps {
			return
		}
	}

	// add points on line
	for i, size := 0, len(mesh.model.Lines); i < size; i++ {
		if mesh.model.Lines[i][2] == Removed {
			continue
		}
		_, _, stB := PointLine(
			p,
			mesh.model.Points[mesh.model.Lines[i][0]],
			mesh.model.Points[mesh.model.Lines[i][1]],
		)
		if !stB.Has(OnSegment) {
			continue
		}
		// replace point tag
		tag = Fixed
		// index of new point
		idp := mesh.model.AddPoint(p)
		for i := len(mesh.Points) - 1; i < idp; i++ {
			mesh.Points = append(mesh.Points, Undefined)
		}
		mesh.Points[idp] = tag
		// add new lines
		mesh.model.AddLine(mesh.model.Points[mesh.model.Lines[i][0]], p,
			mesh.model.Lines[i][2])
		mesh.model.AddLine(mesh.model.Points[mesh.model.Lines[i][1]], p,
			mesh.model.Lines[i][2])
	}

	// TODO : add to delanay flip linked list
	for i, size := 0, len(mesh.Triangles); i < size; i++ {
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
		for i := len(mesh.Points) - 1; i < idp; i++ {
			mesh.Points = append(mesh.Points, Undefined)
		}
		mesh.Points[idp] = tag
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
func (mesh *Mesh) repairTriangles(ap int, rt []int, state int) (err error) {
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

	// amount triangles before added
	size := len(mesh.Triangles)

	// create chain
	switch state {
	case 100:
		// point on triangle 0 line 0 with 2 triangles

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
			// repair triangles sides
			return mesh.repairTriangles(ap, rt, state)
		}
		// 2 ------------- 0 1 ------------ 2 //
		//  \         2    | |    1        /  //
		//   \             | |            /   //
		//     --   rt[1]  | | rt[0]    --    //
		//       \         | |         /      //
		//        --      0| |0      --       //
		//          \ 1    * *    2 /         //
		//           --    | |    --          //
		//             \   | |   /            //
		//              -- | | --             //
		//                \| |/               //
		//                 1 0                //

		// create chains
		chains = []chain{{
			from:   mesh.model.Triangles[rt[0]][1],
			to:     mesh.model.Triangles[rt[0]][2],
			in:     size,
			out:    mesh.Triangles[rt[0]].tr[1],
			before: rt[0],
		}, {
			from:   mesh.model.Triangles[rt[0]][2],
			to:     mesh.model.Triangles[rt[0]][0],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]].tr[2],
			before: rt[0],
		}, {
			from:   mesh.model.Triangles[rt[1]][1],
			to:     mesh.model.Triangles[rt[1]][2],
			in:     size + 2,
			out:    mesh.Triangles[rt[1]].tr[1],
			before: rt[1],
		}, {
			from:   mesh.model.Triangles[rt[1]][2],
			to:     mesh.model.Triangles[rt[1]][0],
			in:     size + 3,
			out:    mesh.Triangles[rt[1]].tr[2],
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
			return fmt.Errorf("removed triangles: %v", rt)
		}
		if mesh.Triangles[rt[0]].tr[0] != Boundary {
			return fmt.Errorf("not valid boundary")
		}
		// point on triangle boundary line 0
		chains = []chain{{
			from:   mesh.model.Triangles[rt[0]][1],
			to:     mesh.model.Triangles[rt[0]][2],
			in:     size,
			out:    mesh.Triangles[rt[0]].tr[1],
			before: rt[0],
		}, {
			from:   mesh.model.Triangles[rt[0]][2],
			to:     mesh.model.Triangles[rt[0]][0],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]].tr[2],
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
			return fmt.Errorf("removed triangles: %v", rt)
		}
		// point in triangle
		chains = []chain{{
			from:   mesh.model.Triangles[rt[0]][0],
			to:     mesh.model.Triangles[rt[0]][1],
			in:     size,
			out:    mesh.Triangles[rt[0]].tr[0],
			before: rt[0],
		}, {
			from:   mesh.model.Triangles[rt[0]][1],
			to:     mesh.model.Triangles[rt[0]][2],
			in:     size + 1,
			out:    mesh.Triangles[rt[0]].tr[1],
			before: rt[0],
		}, {
			from:   mesh.model.Triangles[rt[0]][2],
			to:     mesh.model.Triangles[rt[0]][0],
			in:     size + 2,
			out:    mesh.Triangles[rt[0]].tr[2],
			before: rt[0],
		}}
		tc = [2]int{size + 2, size}

	default:
		return fmt.Errorf("not clear state %v", state)
	}
	// debug checking
	if tc[0] == Undefined || tc[1] == Undefined {
		panic(fmt.Errorf("undefined corner triangle"))
	}

	// create triangles
	for i := range chains {
		mesh.model.AddTriangle(
			mesh.model.Points[chains[i].from],
			mesh.model.Points[chains[i].to],
			mesh.model.Points[ap],
			Undefined, // TODO for case with 2 triangles - not clear tag
		)
		tr := [3]int{Undefined, Undefined, Undefined}
		if chains[i].before == Undefined {
			panic("undefined")
		}

		tr[0] = chains[i].out
		mesh.Swap(chains[i].out, chains[i].before, chains[i].in)
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
	return nil
}

func (mesh *Mesh) Swap(elem, from, to int) {
	if elem == Boundary {
		return
	}
	counter := 0
	for h := 0; h < 3; h++ {
		if from == mesh.Triangles[elem].tr[h] {
			counter++
			mesh.Triangles[elem].tr[h] = to
		}
	}
	if 1 < counter {
		panic("swap")
	}
}

// TODO delanay only for some triangles, if list empty then for  all triangles
func (mesh *Mesh) Delanay() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Delanay: %v", err)
		}
	}()
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
		// rotate near triangle
		for iter := 0; ; iter++ {
			if iter == 50 {
				err = fmt.Errorf("delanay infinite loop 1")
				return
			}
			if mesh.model.Triangles[tr][side] == mesh.model.Triangles[neartr][1] {
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
		if !PointInCircle(
			mesh.model.Points[mesh.model.Triangles[neartr][2]],
			[3]Point{
				mesh.model.Points[mesh.model.Triangles[tr][0]],
				mesh.model.Points[mesh.model.Triangles[tr][1]],
				mesh.model.Points[mesh.model.Triangles[tr][2]],
			},
		) {
			return
		}

		// flip only if middle side is not fixed
		{
			idp0 := mesh.model.AddPoint(mesh.model.Points[mesh.model.Triangles[neartr][0]])
			idp1 := mesh.model.AddPoint(mesh.model.Points[mesh.model.Triangles[neartr][1]])
			for _, line := range mesh.model.Lines {
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
			if mesh.model.Triangles[tr][0] == mesh.model.Triangles[neartr][1] {
				break
			}
			mesh.shiftTriangle(tr)
		}

		// flip
		flip = true

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
			red := &mesh.model.Triangles[neartr]
			blu := &mesh.model.Triangles[tr]
			red[1], blu[1] =
				blu[2], red[2]
		}
		{
			// flip near triangles
			red := &mesh.Triangles[neartr].tr
			blu := &mesh.Triangles[tr].tr
			red[0], red[1], blu[0], blu[1] =
				blu[1], red[0], red[1], blu[0]
			mesh.Swap(red[0], tr, neartr)
			mesh.Swap(blu[0], neartr, tr)
		}
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			return
		}
	}

	// loop of triangles
	for iter := 0; ; iter++ {
		counter := 0
		for tr := range mesh.model.Triangles {
			if mesh.model.Triangles[tr][0] == Removed {
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
							return
						}
					}
					break
				}
			}
		}
		if counter == 0 {
			break
		}
		if iter == 5000 {
			return fmt.Errorf("global delanay infinite loop")
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (mesh *Mesh) Materials() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Materials: %v", err)
		}
	}()
	marks := make([]bool, len(mesh.model.Triangles))

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
		points = append(points, mesh.model.Triangles[from][:3]...)
		points = append(points, mesh.model.Triangles[to][:3]...)
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
		for _, line := range mesh.model.Lines {
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
		// mark
		marks[to] = true
		mesh.model.Triangles[to][3] = counter
		for side := 0; side < 3; side++ {
			err = mark(to, mesh.Triangles[to].tr[side], counter)
			if err != nil {
				return
			}
		}
		return nil
	}

	counter := 50
	for i := range mesh.model.Triangles {
		if marks[i] {
			continue
		}
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		mesh.model.Triangles[i][3] = counter
		for side := 0; side < 3; side++ {
			from := i
			to := mesh.Triangles[i].tr[side]
			err = mark(from, to, counter)
			if err != nil {
				return
			}
		}
		counter++
	}
	return
}

func (mesh *Mesh) Smooth() {
	// for acceptable movable points calculate all side distances from that
	// point to points near triangles and move to average distance.
	//
	// split sides with maximal side distance

	type Store struct {
		index int   // point index
		near  []int // index of near points
	}
	var store []Store

	// create list of all movable points
	for i := range mesh.model.Points {
		if mesh.Points[i] != Movable {
			continue
		}
		if mesh.Points[i] == Fixed {
			continue
		}
		{
			fix := false
			for _, line := range mesh.model.Lines {
				if line[0] != i && line[1] != i {
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
		var near []int
		for _, tri := range mesh.model.Triangles {
			if i != tri[0] && i != tri[1] && i != tri[2] {
				continue
			}
			near = append(near, tri[0:3]...)
		}
		// uniq points
		sort.Ints(near)
		uniq := []int{near[0]}
		for i := 1; i < len(near); i++ {
			if near[i-1] != near[i] {
				uniq = append(uniq, near[i])
			}
		}
		store = append(store, Store{
			index: i,
			near:  near,
		})
	}

	max := 1.0
	for iter := 0; iter < 100 && Eps < max; iter++ {
		max = 0.0
		for _, st := range store {
			var x, y float64
			for _, n := range st.near {
				x += mesh.model.Points[n].X
				y += mesh.model.Points[n].Y
			}
			x /= float64(len(st.near))
			y /= float64(len(st.near))
			max = math.Max(max, Distance(mesh.model.Points[st.index], Point{x, y}))
			mesh.model.Points[st.index].X = x
			mesh.model.Points[st.index].Y = y
		}
	}
}

func (mesh *Mesh) middlePoint(p1, p2 Point) Point {
	return Point{
		X: p1.X*0.5 + p2.X*0.5,
		Y: p1.Y*0.5 + p2.Y*0.5,
	}
}

func (mesh *Mesh) Split(d float64) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Split: %v", err)
		}
	}()
	var pnts []Point

	addpoint := func(p1, p2 Point) {
		dist := Distance(p1, p2)
		if dist < d {
			return
		}
		// add middle point
		pnts = append(pnts, mesh.middlePoint(p1, p2))
		// TODO points free or fixed
	}

	for i := range mesh.model.Triangles {
		if mesh.model.Triangles[i][0] == Removed {
			continue
		}
		addpoint(
			mesh.model.Points[mesh.model.Triangles[i][0]],
			mesh.model.Points[mesh.model.Triangles[i][1]],
		)
		addpoint(
			mesh.model.Points[mesh.model.Triangles[i][1]],
			mesh.model.Points[mesh.model.Triangles[i][2]],
		)
		addpoint(
			mesh.model.Points[mesh.model.Triangles[i][2]],
			mesh.model.Points[mesh.model.Triangles[i][0]],
		)
	}

	// add all points of model
	for i := range pnts {
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
		err = mesh.AddPoint(pnts[i], Movable)
		if err != nil {
			return
		}
		if Debug {
			err = mesh.Check()
			if err != nil {
				return
			}
		}
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			return
		}
	}

	err = mesh.Delanay()
	if err != nil {
		return
	}

	if 0 < len(pnts) {
		return mesh.Split(d)
	}

	return
}

func (mesh *Mesh) AddLine(p1, p2 Point, tag int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("AddLine: %v", err)
		}
	}()
	// get point index
	idp1 := mesh.model.AddPoint(p1)
	idp2 := mesh.model.AddPoint(p2)
	// find triangle with that points
	for _, tri := range mesh.model.Triangles {
		if idp1 != tri[0] && idp1 != tri[1] && idp1 != tri[2] {
			continue
		}
		if idp2 != tri[0] && idp2 != tri[1] && idp2 != tri[2] {
			continue
		}
		mesh.model.AddLine(p1, p2, tag)
		return
	}
	// possible a few triangles on line

	// add middle point
	mid := mesh.middlePoint(p1, p2)
	err = mesh.AddPoint(mid, tag)
	if err != nil {
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
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
			return
		}
	}

	// add both lines
	err = mesh.AddLine(p1, mid, tag)
	if err != nil {
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			return
		}
	}
	err = mesh.AddLine(mid, p2, tag)
	if err != nil {
		return
	}
	if Debug {
		err = mesh.Check()
		if err != nil {
			return
		}
	}

	return
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
