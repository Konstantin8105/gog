package gog

type Mesh struct {
	model     Model
	Triangles []Triangle
	// TODO
}

func New(model Model) (mesh *Mesh) {
	// create a new Mesh
	mesh = new(Mesh)
	// convex
	cps := ConvexHull(model.Points) // points on convex hull
	// prepare mesh triangles
	for i := 2; i < len(cps); i++ {
		var (
			p0 = mesh.model.AddPoint(cps[0])
			p1 = mesh.model.AddPoint(cps[i-2])
			p2 = mesh.model.AddPoint(cps[i-1])
		)
		mesh.model.AddTriangle(cps[0], cps[i-2], cps[i-1], -1)
		if i == 2 {
			mesh.Triangles = append(mesh.Triangles, Triangle{
				nodes: [3]int{p0, p1, p2},
				tr:    [3]int{-1, -1, 1},
			})
		} else {
			mesh.Triangles = append(mesh.Triangles, Triangle{
				nodes: [3]int{p0, p1, p2},
				tr:    [3]int{i - 2, -1, i - 1},
			})
		}
	}
	mesh.Triangles[len(mesh.Triangles)-1].tr[2] = -1 // last not exist triangle
	// add all points of model
	for i := range model.Points {
		mesh.AddPoint(model.Points[i])
	}

	// TODO
	return
}

func (m *Mesh) AddPoint(p Point) {
	// TODO : add to delanay flip linked list
	for i := range m.Triangles {
		res, err := TriangleSplitByPoint(
			p,
			m.model.Points[m.Triangles[i].nodes[0]],
			m.model.Points[m.Triangles[i].nodes[1]],
			m.model.Points[m.Triangles[i].nodes[2]],
		)
		if err != nil {
			panic(err)
		}
		switch len(res) {
		case 2:
			// point on some line
			// find intersect side and near triangle if exist
			// TODO
			return

		case 3:
			// point inside triangle


//         TriangleStructure[] triangles = new TriangleStructure[3];
//         for (int i = 0; i < 3; i++) {
//             triangles[i] = new TriangleStructure();
//         }
//
//         triangles[0].iNodes = new int[]{beginTriangle.iNodes[0], beginTriangle.iNodes[1], pointIndex};
//         triangles[0].iRibs = new int[]{beginTriangle.iRibs[0], rib1, rib0};
//
//         triangles[1].iNodes = new int[]{beginTriangle.iNodes[1], beginTriangle.iNodes[2], pointIndex};
//         triangles[1].iRibs = new int[]{beginTriangle.iRibs[1], rib2, rib1};
//
//         triangles[2].iNodes = new int[]{beginTriangle.iNodes[2], beginTriangle.iNodes[0], pointIndex};
//         triangles[2].iRibs = new int[]{beginTriangle.iRibs[2], rib0, rib2};
//
//         triangles[0].triangles = new TriangleStructure[]{beginTriangle.triangles[0], triangles[1], triangles[2]};
//         triangles[1].triangles = new TriangleStructure[]{beginTriangle.triangles[1], triangles[2], triangles[0]};
//         triangles[2].triangles = new TriangleStructure[]{beginTriangle.triangles[2], triangles[0], triangles[1]};
//
//         addInverseLinkOnTriangle(triangles);
//
//         for (int i = 0; i < 3; i++) {
//             flipper.add(triangles[i], 0);
//             triangleList.add(triangles[i]);
//         }
//     }



			// TODO
			return
		}
	}

	panic("point outside triangles")
}

func (m *Mesh) AddSide() {
	// TODO
}

func (m *Mesh) Delanay() {
	// triangle is success by delanay, if all points is outside of circle
	// from 3 triangle points

	// if (!isPointInCircle(
	//         new Point[]{
	//                 triangulation.getNode(next.triangle.triangles[next.side].iNodes[0]),
	//                 triangulation.getNode(next.triangle.triangles[next.side].iNodes[1]),
	//                 triangulation.getNode(next.triangle.triangles[next.side].iNodes[2])},
	//         triangulation.getNode(next.triangle.iNodes[normalizeSizeBy3(next.side - 1)])))
	//     continue;
	// triangulation.flipTriangles(next.triangle, next.side);

	// TODO
}

func (m *Mesh) Smooth() {
	// for acceptable movable points calculate all side distances from that
	// point to points near triangles and move to average distance.
	//
	// split sides with maximal side distance
	// TODO
}

func (m *Mesh) MaxArea() {
	// TODO
}

func (m *Mesh) MinAngle() {
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
	nodes [3]int // indexes of triangle points
	// ribs  [3]int // indexes of triangle ribs
	tr [3]int // indexes of near triangles
}

func (t *Triangle) swap() {
	// 	t.nodes[0], t.nodes[1] = t.nodes[1], t.nodes[0]
	// 	t.ribs[1], t.ribs[2] = t.ribs[2], t.ribs[1]
	t.tr[1], t.tr[2] = t.tr[2], t.tr[1]
}

//     protected void flipTriangles(TriangleStructure triangle, int indexTriangle) {
//         TriangleStructure[] region = new TriangleStructure[4];
//
//         int pointNewTriangle = triangle.iNodes[normSizeBy3(indexTriangle - 1)];
//         int commonRib = triangle.iRibs[indexTriangle];
//         TriangleStructure[] baseTriangles = new TriangleStructure[]{
//                 triangle,
//                 triangle.triangles[indexTriangle]
//         };
//
//         // global region
//         int position = 0;
//         for (int i = 0; i < 2; i++) {
//             TriangleStructure internalTriangle = baseTriangles[i];
//             int indexCommonRib = -1;
//             for (int j = 0; j < 3; j++) {
//                 if (commonRib == internalTriangle.iRibs[j]) {
//                     indexCommonRib = j;
//                 }
//             }
//             TriangleStructure t1 = new TriangleStructure();
//             t1.iRibs = new int[]{
//                     internalTriangle.iRibs[normSizeBy3(indexCommonRib + 1)]
//             };
//             t1.iNodes = new int[]{
//                     internalTriangle.iNodes[normSizeBy3(indexCommonRib + 1)],
//                     internalTriangle.iNodes[normSizeBy3(indexCommonRib - 1)]
//             };
//             t1.triangles = new TriangleStructure[]{
//                     internalTriangle.triangles[normSizeBy3(indexCommonRib + 1)]
//             };
//             region[position++] = t1;
//             TriangleStructure t2 = new TriangleStructure();
//             t2.iRibs = new int[]{
//                     internalTriangle.iRibs[normSizeBy3(indexCommonRib - 1)]
//             };
//             t2.iNodes = new int[]{
//                     internalTriangle.iNodes[normSizeBy3(indexCommonRib - 1)],
//                     internalTriangle.iNodes[normSizeBy3(indexCommonRib)]
//             };
//             t2.triangles = new TriangleStructure[]{
//                     internalTriangle.triangles[normSizeBy3(indexCommonRib - 1)]
//             };
//             region[position++] = t2;
//         }
//
//
//         if (Geometry.isCounterClockwise(
//                 nodes.get(region[1].iNodes[0]),
//                 nodes.get(region[1].iNodes[1]),
//                 nodes.get(region[2].iNodes[1])))
//             return;
//
//
//         if (Geometry.isCounterClockwise(
//                 nodes.get(region[3].iNodes[0]),
//                 nodes.get(region[3].iNodes[1]),
//                 nodes.get(region[0].iNodes[1])
//         ))
//             return;
//
//         //checking base point
//         if (region[0].iNodes[1] != pointNewTriangle) {
//             return;
//         }
//
//         //separate on 2 triangles
//         int newCommonRib = getIdRib();
//
//         TriangleStructure[] triangles = new TriangleStructure[2];
//         for (int i = 0; i < 2; i++) {
//             triangles[i] = new TriangleStructure();
//         }
//
//         triangles[0].iNodes = new int[]{
//                 region[1].iNodes[0],
//                 region[1].iNodes[1],
//                 region[2].iNodes[1]
//         };
//         triangles[0].iRibs = new int[]{
//                 region[1].iRibs[0],
//                 region[2].iRibs[0],
//                 newCommonRib
//         };
//         triangles[0].triangles = new TriangleStructure[]{
//                 region[1].triangles[0],
//                 region[2].triangles[0],
//                 triangles[1]
//         };
//
//         triangles[1].iNodes = new int[]{
//                 region[3].iNodes[0],
//                 region[3].iNodes[1],
//                 region[0].iNodes[1]
//         };
//         triangles[1].iRibs = new int[]{
//                 region[3].iRibs[0],
//                 region[0].iRibs[0],
//                 newCommonRib
//         };
//         triangles[1].triangles = new TriangleStructure[]{
//                 region[3].triangles[0],
//                 region[0].triangles[0],
//                 triangles[0]
//         };
//
//         //inverse link on triangle
//         addInverseLinkOnTriangle(triangles);
//
//         triangleList.add(triangles[0]);
//         triangleList.add(triangles[1]);
//
//         //move beginTriangle
//         searcher.setSearcher(triangles[0]);
//
//         //add null in old triangles
//         for (TriangleStructure base : baseTriangles) {
//             triangleList.NullableTriangle(base);
//         }
//
//         flipper.add(triangles[0], 0);
//         flipper.add(triangles[0], 1);
//
//         flipper.add(triangles[1], 0);
//         flipper.add(triangles[1], 1);
//     }
//
//
//     private void addNextPointOnLine(Point nextPoint, int indexLineInTriangle) {
//
//         TriangleStructure beginTriangle = searcher.getSearcher();
//
//         nodes.add(nextPoint);
//         int pointIndex = nodes.size() - 1;
//
//         int rib0 = getIdRib();
//         int rib1 = getIdRib();
//         int rib2 = getIdRib();
//         int rib3 = getIdRib();
//
//         TriangleStructure triangles[] = new TriangleStructure[4];
//         for (int i = 0; i < 4; i++) {
//             triangles[i] = new TriangleStructure();
//         }
//
//
//         triangles[0].iNodes = new int[]{
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle)],
//                 pointIndex,
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[1].iNodes = new int[]{
//                 pointIndex,
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle + 1)],
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[0].iRibs = new int[]{
//                 rib0,
//                 rib2,
//                 beginTriangle.iRibs[normSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[1].iRibs = new int[]{
//                 rib1,
//                 beginTriangle.iRibs[normSizeBy3(indexLineInTriangle + 1)],
//                 rib2
//         };
//
//
//         triangles[0].triangles = new TriangleStructure[]{
//                 null,
//                 triangles[1],
//                 beginTriangle.triangles[normSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[1].triangles = new TriangleStructure[]{
//                 null,
//                 beginTriangle.triangles[normSizeBy3(indexLineInTriangle + 1)],
//                 triangles[0]
//         };
//
//         if (beginTriangle.triangles[indexLineInTriangle] == null) {
//             addInverseLinkOnTriangle(new TriangleStructure[]{triangles[0], triangles[1]});
//
//             triangleList.NullableTriangle(beginTriangle);
//
//             searcher.setSearcher(triangles[0]);
//
//             flipper.add(triangles[0], 2);
//             flipper.add(triangles[1], 1);
//
//             triangleList.add(triangles[0]);
//             triangleList.add(triangles[1]);
//             return;
//         }
//
//         int ribConnectId = beginTriangle.iRibs[indexLineInTriangle];
//         TriangleStructure nextTriangle = beginTriangle.triangles[indexLineInTriangle];
//         triangleList.NullableTriangle(beginTriangle);
//         beginTriangle = nextTriangle;
//         for (int i = 0; i < 3; i++) {
//             if (beginTriangle.iRibs[i] == ribConnectId) {
//                 indexLineInTriangle = i;
//             }
//         }
//
//         triangles[0].triangles[0] = triangles[2];
//         triangles[1].triangles[0] = triangles[3];
//
//         triangles[2].iNodes = new int[]{
//                 pointIndex,
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle + 1)],
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[3].iNodes = new int[]{
//                 beginTriangle.iNodes[indexLineInTriangle],
//                 pointIndex,
//                 beginTriangle.iNodes[normSizeBy3(indexLineInTriangle - 1)]
//         };
//
//
//         triangles[2].iRibs = new int[]{
//                 rib0,
//                 beginTriangle.iRibs[normSizeBy3(indexLineInTriangle + 1)],
//                 rib3
//         };
//         triangles[3].iRibs = new int[]{
//                 rib1,
//                 rib3,
//                 beginTriangle.iRibs[normSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[2].triangles = new TriangleStructure[]{
//                 triangles[0],
//                 beginTriangle.triangles[normSizeBy3(indexLineInTriangle + 1)],
//                 triangles[3]
//         };
//
//         triangles[3].triangles = new TriangleStructure[]{
//                 triangles[1],
//                 triangles[2],
//                 beginTriangle.triangles[normSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         addInverseLinkOnTriangle(triangles);
//
//         triangleList.NullableTriangle(beginTriangle);
//         searcher.setSearcher(triangles[0]);
//
//         flipper.add(triangles[0], 2);
//         flipper.add(triangles[1], 1);
//         flipper.add(triangles[2], 1);
//         flipper.add(triangles[3], 2);
//
//         triangleList.addAll(triangles);
//     }
//
//
//
//     private void addInverseLinkOnTriangle(TriangleStructure[] triangles) {
//         for (TriangleStructure triangle : triangles) {
//             if (triangle == null)
//                 continue;
//             for (int j = 0; j < 3; j++) {
//                 TriangleStructure externalTriangle = triangle.triangles[j];
//                 int commonRib = triangle.iRibs[j];
//                 if (externalTriangle != null) {
//                     for (int k = 0; k < 3; k++) {
//                         if (externalTriangle.iRibs[k] == commonRib) {
//                             externalTriangle.triangles[k] = triangle;
//                         }
//                     }
//                 }
//             }
//         }
//         boolean inverseAgain = false;
//         for (TriangleStructure triangle : triangles) {
//             if (Geometry.isCounterClockwise(
//                     nodes.get(triangle.iNodes[0]),
//                     nodes.get(triangle.iNodes[1]),
//                     nodes.get(triangle.iNodes[2]))) {
//                 triangle.changeClockwise();
//                 inverseAgain = true;
//             }
//         }
//         if (inverseAgain) {
//             addInverseLinkOnTriangle(triangles);
//         }
//     }
