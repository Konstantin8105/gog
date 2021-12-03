package gog

type Mesh struct {
}

func (m Mesh) ConvexHull() {
	// TODO
}

func (m *Mesh) AddPoint() {
	// TODO
}

func (m *Mesh) AddSide() {
	// TODO
}

func (m *Mesh) Delanay() {
	// TODO
}

func (m *Mesh) Smooth() {
	// TODO
}

func (m *Mesh) Sweep() {
	// TODO
}

func (m *Mesh) MaxArea() {
	// TODO
}

func (m *Mesh) MinAngle() {
	// TODO
}

//
// public class FlipStructure {
//     public final TriangleStructure triangle;
//     public final int side;
//
//     public FlipStructure(TriangleStructure triangle, int side) {
//         this.triangle = triangle;
//         this.side = side;
//     }
// }
//
// public interface Fliper {
//     void add(TriangleStructure triangle, int index);
//
//     void run();
// }
//
//     public enum PointTriangleState {
//         POINT_ON_LINE_0,
//         POINT_ON_LINE_1,
//         POINT_ON_LINE_2,
//         POINT_ON_CORNER,
//         POINT_INSIDE,
//         POINT_OUTSIDE_LINE_0,
//         POINT_OUTSIDE_LINE_1,
//         POINT_OUTSIDE_LINE_2,
//     }
//
//
//
// public interface Searcher {
//     TriangleStructure getSearcher();
//
//     void setSearcher(TriangleStructure searcher);
//
//     void chooseSearcher(Point point);
//
//     GeometryPointTriangle.PointTriangleState movingByConvexHull(Point point);
// }
//
//
// public class FastSearcher implements Searcher {
//
//     public static double AMOUNT_SEARCHER_FACTOR = 0.5D;
//
//     private final TriangleStructure[] searcher;
//     private final double[] elevations;
//     private int positionSearcher = 0;
//     private TriangulationDelaunay triangulation = new TriangulationDelaunay();
//
//     public FastSearcher(TriangulationDelaunay triangulation, TriangleStructure init, BorderBox box, int amountOfPoints) {
//         this.triangulation = triangulation;
//         searcher = new TriangleStructure[(int) Math.max(1.0D, AMOUNT_SEARCHER_FACTOR * Math.sqrt((double) amountOfPoints))];
//         for (int i = 0; i < searcher.length; i++) {
//             searcher[i] = init;
//         }
//         double heightStep = (box.getY_max() - box.getY_min()) / (double) searcher.length;
//         elevations = new double[searcher.length];
//         for (int i = 0; i < elevations.length; i++) {
//             elevations[i] = box.getY_min() + i * heightStep;
//         }
//     }
//
//     Override
//     public TriangleStructure getSearcher() {
//         return searcher[positionSearcher];
//     }
//
//     Override
//     public void setSearcher(TriangleStructure searcher) {
//         this.searcher[positionSearcher] = searcher;
//     }
//
//     Override
//     public void chooseSearcher(Point point) {
//         for (int i = searcher.length - 1; i >= 0; i--) {
//             if (point.Y > elevations[i] - Precision.epsilon()) {
//                 positionSearcher = i;
//                 break;
//             }
//         }
//
//         if (searcher[positionSearcher].triangles != null)
//             return;
//
//         for (int i = 0; i < searcher.length; i++) {
//             if (searcher[ArrayIndexCorrection.normalize(positionSearcher + i, searcher.length)].triangles != null) {
//                 searcher[positionSearcher] = searcher[ArrayIndexCorrection.normalize(positionSearcher + i, searcher.length)];
//                 return;
//             }
//         }
//     }
//
//     /**
//      * Found next triangle
//      * Performance - O(n) in worst case and O(sqrt(n)) is average case.
//      *
//      * param point - next point
//      * return GeometryPointTriangle.PointTriangleState
//      * see Point
//      * see GeometryPointTriangle.PointTriangleState
//      */
//     private Geometry.POINT_ON_LINE[] value = new Geometry.POINT_ON_LINE[3];
//     private Point[] trianglePoint = new Point[3];
//
//     Override
//     public GeometryPointTriangle.PointTriangleState movingByConvexHull(Point point) {
//         TriangleStructure beginTriangle = getSearcher();
//         while (true) {
//             //add reserve searching
//             value[0] = Geometry.calculateValuePointOnLine(triangulation.getNode(beginTriangle.iNodes[0]), triangulation.getNode(beginTriangle.iNodes[1]), point);
//             if (Geometry.isAtRightOf(value[0])) {
//                 beginTriangle = beginTriangle.triangles[0];
//             } else {
//                 int whichOp = 0;
//                 value[1] = Geometry.calculateValuePointOnLine(triangulation.getNode(beginTriangle.iNodes[1]), triangulation.getNode(beginTriangle.iNodes[2]), point);
//                 if (Geometry.isAtRightOf(value[1])) {
//                     whichOp += 1;
//                 }
//                 value[2] = Geometry.calculateValuePointOnLine(triangulation.getNode(beginTriangle.iNodes[2]), triangulation.getNode(beginTriangle.iNodes[0]), point);
//                 if (Geometry.isAtRightOf(value[2])) {
//                     whichOp += 2;
//                 }
//                 if (whichOp == 0) {
//                     break;
//                 } else if (whichOp == 1) {
//                     beginTriangle = beginTriangle.triangles[1];
//                 } else if (whichOp == 2) {
//                     beginTriangle = beginTriangle.triangles[2];
//                 } else {
//                     if (Geometry.distanceLineAndPoint(triangulation.getNode(beginTriangle.iNodes[1]), triangulation.getNode(beginTriangle.iNodes[2]), point) >
//                             Geometry.distanceLineAndPoint(triangulation.getNode(beginTriangle.iNodes[2]), triangulation.getNode(beginTriangle.iNodes[0]), point)) {
//                         beginTriangle = beginTriangle.triangles[1];
//                     } else {
//                         beginTriangle = beginTriangle.triangles[2];
//                     }
//                 }
//             }
//         }
//         trianglePoint = new Point[]{
//                 triangulation.getNode(beginTriangle.iNodes[0]),
//                 triangulation.getNode(beginTriangle.iNodes[1]),
//                 triangulation.getNode(beginTriangle.iNodes[2])
//         };
//         setSearcher(beginTriangle);
//         return GeometryPointTriangle.statePointInTriangle(point, trianglePoint, value);
//     }
// }
//
//
//
// public class FliperDelaunay implements Fliper {
//     private final Stack<FlipStructure> buffer = new Stack<>();
//     private TriangulationDelaunay triangulation = new TriangulationDelaunay();
//
//     public FliperDelaunay(TriangulationDelaunay triangulation) {
//         this.triangulation = triangulation;
//     }
//
//     Override
//     public void add(TriangleStructure triangle, int index) {
//         buffer.add(new FlipStructure(triangle, index));
//     }
//
//     Override
//     public void run() {
//         while (!buffer.empty()) {
//             FlipStructure next = buffer.pop();
//
//             if (next.triangle.triangles == null)
//                 continue;
//
//             if (next.triangle.triangles[next.side] == null)
//                 continue;
//
//             if (!isPointInCircle(
//                     new Point[]{
//                             triangulation.getNode(next.triangle.triangles[next.side].iNodes[0]),
//                             triangulation.getNode(next.triangle.triangles[next.side].iNodes[1]),
//                             triangulation.getNode(next.triangle.triangles[next.side].iNodes[2])},
//                     triangulation.getNode(next.triangle.iNodes[normalizeSizeBy3(next.side - 1)])))
//                 continue;
//             triangulation.flipTriangles(next.triangle, next.side);
//
//             // TODO: 20.09.2016 add to another class
// //            Point p0 = triangulation.getNode(next.triangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(next.side - 1)]);
// //            Point p1 = triangulation.getNode(next.triangle.triangles[next.side].iNodes[0]);
// //            Point p2 = triangulation.getNode(next.triangle.triangles[next.side].iNodes[1]);
// //            Point p3 = triangulation.getNode(next.triangle.triangles[next.side].iNodes[2]);
// //
// //            double s1 = (p0.X - p1.X) * (p0.Y - p3.Y) - (p0.X - p3.X) * (p0.Y - p1.Y);
// //            double s2 = (p2.X - p3.X) * (p2.X - p1.X) + (p2.Y - p3.Y) * (p2.Y - p1.Y);
// //            double s3 = (p0.X - p1.X) * (p0.X - p3.X) + (p0.Y - p1.Y) * (p0.Y - p3.Y);
// //            double s4 = (p2.X - p3.X) * (p2.Y - p1.Y) - (p2.X - p1.X) * (p2.Y - p3.Y);
// //            if (s1 * s2 + s3 * s4 >= 0)
// //                continue;
// //            else
// //                triangulation.flipTriangles(next.triangle, next.side);
//         }
//     }
//
// // Triangulation data structure  "Nodes, ribs и triangles"
// // book "Algoritm building and analyse triangulation", A.B.Skvorcov.
// public class TriangleStructure {
//     // indexes of triangle points
//     public int[] iNodes;
//     // indexes of near triangles
//     public TriangleStructure[] triangles;
//     // indexes of triangle ribs
//     public int[] iRibs;
//
//     public void changeClockwise() {
//         int temp;
//         temp = iNodes[0];
//         iNodes[0] = iNodes[1];
//         iNodes[1] = temp;
//         temp = iRibs[1];
//         iRibs[1] = iRibs[2];
//         iRibs[2] = temp;
//         TriangleStructure tri = triangles[1];
//         triangles[1] = triangles[2];
//         triangles[2] = tri;
//     }
// }
//
//
// public class TriangulationDelaunay {
//     // Array of nodes - type: Point
//     private final List<Point> nodes = new ArrayList<>();
//
//     protected Point getNode(int index) {
//         return nodes.get(index);
//     }
//
//
//     private Fliper flipper;
//     private Searcher searcher;
//     private final TriangleList triangleList = new TriangleList();
//
//     public static double AMOUNT_CLEANING_FACTOR_TRIANGLE_STRUCTURE = 2.4D;
//     public static final double RATIO_DELETING_CONVEX_POINT_FROM_POINT_LIST = 0.2D;
//     public static int MINIMAL_POINTS_FOR_CLEANING = 10000;
//
//     // constructor for create convexHull region at the base on points
//     public TriangulationDelaunay(Point[] points) {
//         run(points);
//     }
//
//     public TriangulationDelaunay() {
//     }
//
//     public void run(Point[] input) {
//         flipper = new FliperDelaunay(this);
//         List<Point>[] pointArray = convexHullDouble(input);
//         if (pointArray == null)
//             return;
//         List<Point> convexPoints = pointArray[0];
//         BorderBox box = createConvexHullTriangles(convexPoints);
//         searcher = new FastSearcher(this, triangleList.getfirstNotNullableElement(), box, pointArray[1].size());
//
//         if (pointArray[1].size() >= MINIMAL_POINTS_FOR_CLEANING) {
//             int amount = (int) (AMOUNT_CLEANING_FACTOR_TRIANGLE_STRUCTURE * pointArray[1].size());
//             amount = amount < 1 ? 1 : amount;
//             triangleList.setMaxAmountNullableElements(amount);
//         }
//         for (int i = 0; i < pointArray[1].size(); i++) {
//             addNextPoint(pointArray[1].get(i));
//             flipper.run();
// //            if (i % 1000 == 0)
// //                System.err.println(i);
//         }
//         flipper.run();
//     }
//
//     private void addNextPoint(Point nextPoint) {
//
//         searcher.chooseSearcher(nextPoint);
//
//         GeometryPointTriangle.PointTriangleState state = searcher.movingByConvexHull(nextPoint);
//
//         switch (state) {
//             case POINT_INSIDE:
//                 addNextPointInTriangle(nextPoint);
//                 break;
//             case POINT_ON_LINE_0:
//                 addNextPointOnLine(nextPoint, 0);
//                 break;
//             case POINT_ON_LINE_1:
//                 addNextPointOnLine(nextPoint, 1);
//                 break;
//             case POINT_ON_LINE_2:
//                 addNextPointOnLine(nextPoint, 2);
//                 break;
//             case POINT_ON_CORNER:
//                 break;
//             default:
//                 System.out.println("STRANGE POINT : " + nextPoint);
//         }
//     }
//
//
//     // Linked list of triangles
//     private class TriangleList {
//
//         private int amountNullableElements = 0;
//         private final List<TriangleStructure> triangleStructureList = new LinkedList<>();
//         private int maxAmountNullableElements = Integer.MAX_VALUE / 2;
//
//         public void add(TriangleStructure triangle) {
//             triangleStructureList.add(triangle);
//         }
//
//         public void addAll(TriangleStructure[] triangles) {
//             Collections.addAll(triangleStructureList, triangles);
//         }
//
//         private void removeNullTriangles() {
//             amountNullableElements = 0;
//             Iterator<TriangleStructure> iterator = triangleStructureList.iterator();
//             while (iterator.hasNext()) {
//                 if (iterator.next().triangles == null)
//                     iterator.remove();
//             }
//         }
//
//         public int size() {
//             return triangleStructureList.size();
//         }
//
//         public List<TriangleStructure> get() {
//             removeNullTriangles();
//             return triangleStructureList;
//         }
//
//         public void setMaxAmountNullableElements(int maxAmountNullableElements) {
//             this.maxAmountNullableElements = maxAmountNullableElements;
//         }
//
//         public TriangleStructure getfirstNotNullableElement() {
//             for (TriangleStructure triangle : triangleStructureList)
//                 if (triangle.triangles != null)
//                     return triangle;
//             return null;
//         }
//
//         private void NullableTriangle(TriangleStructure triangle) {
//             triangle.triangles = null;
//             amountNullableElements++;
//             if (amountNullableElements > maxAmountNullableElements) {
//                 removeNullTriangles();
//             }
//         }
//     }
//
//     private int idMaximalRibs = 0;
//
//     private int getIdRib() {
//         return idMaximalRibs++;
//     }
//
//     private BorderBox createConvexHullTriangles(List<Point> points) {
//         int i = 0;
//         i++;
//         nodes.add(points.get(0));
//         int indexPoint0 = nodes.size() - 1;
//         i++;
//         nodes.add(points.get(1));
//         int indexPoint1 = nodes.size() - 1;
//         int commonRib = getIdRib();
//         TriangleStructure commonTriangle = null;
//
//         int k = 0;
//         while (i + k < points.size()) {
//             i++;
//             nodes.add(points.get(i - 1));
//             int indexPoint2 = nodes.size() - 1;
//             int rib12 = getIdRib();
//             int rib20 = getIdRib();
//
//             TriangleStructure triangle = new TriangleStructure();
//             triangle.iNodes = new int[]{indexPoint0, indexPoint1, indexPoint2};
//             triangle.iRibs = new int[]{commonRib, rib12, rib20};
//             triangle.triangles = new TriangleStructure[]{commonTriangle, null, null};
//             if (commonTriangle != null) {
//                 commonTriangle.triangles[1] = triangle;
//             }
//
//             triangleList.add(triangle);
//
//             if (i + k >= points.size())
//                 break;
//
//             int indexPoint0_next = indexPoint0;
//             int indexPoint1_next = indexPoint2;
//             k++;
//             nodes.add(points.get(points.size() - k));
//             int indexPoint2_next = nodes.size() - 1;
//
//             int rib12_next = getIdRib();
//             int rib20_next = getIdRib();
//
//             TriangleStructure triangle2 = new TriangleStructure();
//             triangle2.iNodes = new int[]{indexPoint0_next, indexPoint1_next, indexPoint2_next};
//             triangle2.iRibs = new int[]{rib20, rib12_next, rib20_next};
//             triangle2.triangles = new TriangleStructure[]{
//                     triangle, null, null
//             };
//             triangle.triangles[2] = triangle2;
//             triangleList.add(triangle2);
//
//
//             indexPoint0 = indexPoint2_next;
//             indexPoint1 = indexPoint1_next;
//             commonRib = rib12_next;
//             commonTriangle = triangle2;
//         }
//
//         BorderBox borderBox = new BorderBox();
//         for (Point point : points) {
//             borderBox.addPoint(point);
//         }
//         return borderBox;
//     }
//
//
//     protected void flipTriangles(TriangleStructure triangle, int indexTriangle) {
//         TriangleStructure[] region = new TriangleStructure[4];
//
//         int pointNewTriangle = triangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexTriangle - 1)];
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
//                     internalTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib + 1)]
//             };
//             t1.iNodes = new int[]{
//                     internalTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib + 1)],
//                     internalTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib - 1)]
//             };
//             t1.triangles = new TriangleStructure[]{
//                     internalTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib + 1)]
//             };
//             region[position++] = t1;
//             TriangleStructure t2 = new TriangleStructure();
//             t2.iRibs = new int[]{
//                     internalTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib - 1)]
//             };
//             t2.iNodes = new int[]{
//                     internalTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib - 1)],
//                     internalTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib)]
//             };
//             t2.triangles = new TriangleStructure[]{
//                     internalTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexCommonRib - 1)]
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
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle)],
//                 pointIndex,
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[1].iNodes = new int[]{
//                 pointIndex,
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[0].iRibs = new int[]{
//                 rib0,
//                 rib2,
//                 beginTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[1].iRibs = new int[]{
//                 rib1,
//                 beginTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
//                 rib2
//         };
//
//
//         triangles[0].triangles = new TriangleStructure[]{
//                 null,
//                 triangles[1],
//                 beginTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[1].triangles = new TriangleStructure[]{
//                 null,
//                 beginTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
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
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//         triangles[3].iNodes = new int[]{
//                 beginTriangle.iNodes[indexLineInTriangle],
//                 pointIndex,
//                 beginTriangle.iNodes[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//
//
//         triangles[2].iRibs = new int[]{
//                 rib0,
//                 beginTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
//                 rib3
//         };
//         triangles[3].iRibs = new int[]{
//                 rib1,
//                 rib3,
//                 beginTriangle.iRibs[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
//         };
//
//         triangles[2].triangles = new TriangleStructure[]{
//                 triangles[0],
//                 beginTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle + 1)],
//                 triangles[3]
//         };
//
//         triangles[3].triangles = new TriangleStructure[]{
//                 triangles[1],
//                 triangles[2],
//                 beginTriangle.triangles[ArrayIndexCorrection.normalizeSizeBy3(indexLineInTriangle - 1)]
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
//     private void addNextPointInTriangle(Point nextPoint) {
//
//         TriangleStructure beginTriangle = searcher.getSearcher();
//
//         nodes.add(nextPoint);
//         int pointIndex = nodes.size() - 1;
//         int rib0 = getIdRib();
//         int rib1 = getIdRib();
//         int rib2 = getIdRib();
//
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
//         triangleList.NullableTriangle(beginTriangle);
//
//         searcher.setSearcher(triangles[0]);
//         addInverseLinkOnTriangle(triangles);
//
//         for (int i = 0; i < 3; i++) {
//             flipper.add(triangles[i], 0);
//             triangleList.add(triangles[i]);
//         }
//     }
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
//
//
//
//     public List<Point[]> getTriangles() {
//         List<TriangleStructure> triangles = triangleList.get();
//         List<Point[]> trianglesPoints = new ArrayList<>();
//         for (TriangleStructure tri : triangles) {
//             Point[] points = new Point[3];
//             for (int i = 0; i < 3; i++) {
//                 points[i] = new Point(nodes.get(tri.iNodes[i]));
//             }
//             trianglesPoints.add(points);
//         }
//         return trianglesPoints;
//     }
//
// }

//Performance O(n*log(n)) in worst case
// Point[0][] - convex points
// Point[1][] - sorted list of all points
// func ConvexHull(ps []Point) (res []Point) {
// 	// 	        List<Integer> removedIndex = new ArrayList<>();
// 	// 	        for (int i = 1; i < array.size(); i++) {
// 	// 	            if (array.get(i - 1).equals(array.get(i))) {
// 	// 	                removedIndex.add(i);
// 	// 	            }
// 	// 	        }
// 	// 	        for (int i = removedIndex.size() - 1; i >= 0; i--) {
// 	// 	            int position = removedIndex.get(i);
// 	// 	            array.remove(position);
// 	// 	        }
// 	// 	        List<Point> convexPoints = new ArrayList<>();
// 	// 	        if (k > 1) {
// 	// 	            H = Arrays.copyOfRange(H, 0, k - 1); // remove non-hull vertices after k; remove k - 1 which is a duplicate
// 	// 	            boolean[] removed = new boolean[k - 1];
// 	// 	            for (int position0 = 0; position0 < removed.length; position0++) {
// 	// 	                int position1 = position0 + 1 >= removed.length ? position0 + 1 - removed.length : position0 + 1;
// 	// 	                int position2 = position0 + 2 >= removed.length ? position0 + 2 - removed.length : position0 + 2;
// 	// 	                if (Geometry.is3pointsCollinear(
// 	// 	                        H[position0],
// 	// 	                        H[position1],
// 	// 	                        H[position2])) {
// 	// 	                    removed[position1] = true;
// 	// 	                }
// 	// 	            }
// 	// 	            for (int i = 0; i < removed.length; i++) {
// 	// 	                if (!removed[i])
// 	// 	                    convexPoints.add(H[i]);
// 	// 	            }
// 	// 	            if (array.size() > 5) {
// 	// 	                if ((double) convexPoints.size() / (double) array.size() > RATIO_DELETING_CONVEX_POINT_FROM_POINT_LIST) {
// 	// 	                    boolean[] delete = new boolean[array.size()];
// 	// 	                    int position = 0;
// 	// 	                    for (int i = 0; i < array.size(); i++) {
// 	// 	                        if (array.get(i).equals(convexPoints.get(position))) {
// 	// 	                            delete[i] = true;
// 	// 	                            position++;
// 	// 	                        }
// 	// 	                    }
// 	// 	                    for (int i = array.size() - 1; i >= 0; i--) {
// 	// 	                        if (position >= convexPoints.size())
// 	// 	                            break;
// 	// 	                        if (array.get(i).equals(convexPoints.get(position))) {
// 	// 	                            delete[i] = true;
// 	// 	                            position++;
// 	// 	                        }
// 	// 	                    }
// 	// 	                    ArrayList<Point> newList = new ArrayList<>();
// 	// 	                    for (int i = 0; i < array.size(); i++) {
// 	// 	                        if (!delete[i])
// 	// 	                            newList.add(array.get(i));
// 	// 	                    }
// 	// 	                    array = newList;
// 	// 	                }
// 	// 	            }
// 	// 	        }
// 	// 	//
// 	// 	        return new List[]{convexPoints, array};
// }
