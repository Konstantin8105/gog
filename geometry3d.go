package gog

import (
	"math"

	"github.com/Konstantin8105/pow"
)

const Eps3D = 1e-5

// Space 3D
//
//	+-----------+-------+-------+-----------+
//	|           | Point | Line2 | Triangle3 |
//	+-----------+-------+-------+-----------+
//	| Point     |   V   |   V   |           |
//	+-----------+-------+-------+-----------+
//	| Line2     |   -   |       |           |
//	+-----------+-------+-------+-----------+
//	| Triangle3 |   -   |   -   |           |
//	+-----------+-------+-------+-----------+

type Point3d [3]float64

func Distance3d(p0, p1 Point3d) float64 {
	return math.Sqrt(pow.E2(p0[0]-p1[0]) +
		pow.E2(p0[1]-p1[1]) +
		pow.E2(p0[2]-p1[2]))
	// https://arxiv.org/pdf/1904.09481.pdf
}

func PointPoint3d(
	p0 Point3d,
	p1 Point3d,
) (
	intersect bool,
) {
	for i := range p0 {
		if Eps < math.Abs(p0[i]-p1[i]) {
			return false
		}
	}
	return Distance3d(p0, p1) < Eps3D
}

func PointLine3d(
	p Point3d,
	l0, l1 Point3d,
) (
	intersect bool,
) {
	// is point on point line
	for _, v := range [2]*Point3d{&l0, &l1} {
		if PointPoint3d(p, *v) {
			return
		}
	}
	// line zero lenght
	if ZeroLine3d(l0, l1) {
		return
	}
	// is point in line box
	for i := range p {
		if p[i] < l0[i] && p[i] < l1[i] {
			return
		}
		if l0[i] < p[i] && l1[i] < p[i] {
			return
		}
	}
	// is point on line
	return true
}

func ZeroLine3d(
	l0, l1 Point3d,
) (
	zero bool,
) {
	return Distance3d(l0, l1) < Eps3D
}

func PointLineRatio3d(
	l0, l1 Point3d,
	ratio float64,
) (
	p Point3d,
) {
	for i := 0; i < 3; i++ {
		// p[i] = l0[i] + ratio*(l1[i]-l0[i])
		p[i] = math.FMA(ratio, l1[i]-l0[i], l0[i])
	}
	return
}

const oneThree = 1.0 / 3.0

// Point on line corner ignored
func LineLine3d(
	a0, a1 Point3d,
	b0, b1 Point3d,
) (
	ratioA, ratioB float64,
	intersect bool,
) {
	// Lina a:
	//	x = xa0 + Ka * (xa1-xa0)
	//	y = ya0 + Ka * (ya1-ya0)
	//	z = za0 + Ka * (za1-za0)
	//
	// Lina b:
	//	x = xb0 + Kb * (xb1-xb0)
	//	y = yb0 + Kb * (yb1-yb0)
	//	z = zb0 + Kb * (zb1-zb0)
	//
	// Intersection point:
	//	xa0 + Ka * (xa1-xa0) = xb0 + Kb * (xb1-xb0)
	//
	// System of equations:
	//	Ka * (a1x-a0x) - Kb * (b1x-b0x) = b0x - a0x
	//	Ka * (a1y-a0y) - Kb * (b1y-b0y) = b0y - a0y
	//	Ka * (a1z-a0z) - Kb * (b1z-b0z) = b0z - a0z
	sys := [3][3]float64{
		{a1[0] - a0[0], -(b1[0] - b0[0]), b0[0] - a0[0]},
		{a1[1] - a0[1], -(b1[1] - b0[1]), b0[1] - a0[1]},
		{a1[2] - a0[2], -(b1[2] - b0[2]), b0[2] - a0[2]},
	}
	Ka := make([]float64, 0, 3)
	Kb := make([]float64, 0, 3)
	for _, v := range [3][2]int{{0, 1}, {1, 2}, {2, 0}} {
		x, y, err := Linear(
			sys[v[0]][0], sys[v[0]][1], sys[v[0]][2],
			sys[v[1]][0], sys[v[1]][1], sys[v[1]][2],
		)
		if err != nil {
			continue
		}
		Ka = append(Ka, x)
		Kb = append(Kb, y)
	}
	if len(Ka) < 1 {
		return
	}
	for i := range Ka {
		ratioA = math.FMA(Ka[i], 1.0/float64(len(Ka)), ratioA)
		ratioB = math.FMA(Kb[i], 1.0/float64(len(Kb)), ratioB)
	}
	intersect = true
	return
}

// Plane equation `A*x+B*y+C*z+D=0`
func Plane(
	p1, p2, p3 Point3d,
) (
	A, B, C, D float64,
) {
	var (
		a1 = p2[0] - p1[0]
		b1 = p2[1] - p1[1]
		c1 = p2[2] - p1[2]
		a2 = p3[0] - p1[0]
		b2 = p3[1] - p1[1]
		c2 = p3[2] - p1[2]
	)
	// algorithm float
	// A = b1*c2 - b2*c1
	// B = a2*c1 - a1*c2
	// C = a1*b2 - b1*a2
	// D = (-A*x1 - B*y1 - C*z1)

	// algorithm FMA
	A = math.FMA(b1, c2, -b2*c1)
	B = math.FMA(a2, c1, -a1*c2)
	C = math.FMA(a1, b2, -b1*a2)
	D = math.FMA(-A, p1[0], math.FMA(-B, p1[1], -C*p1[2]))
	return
}

func PointOnPlane3d(
	A, B, C, D float64,
	p Point3d,
) (
	on bool,
) {
	return math.Abs(math.FMA(A, p[0], math.FMA(B, p[1], math.FMA(C, p[2], D)))) < Eps3D
}

func ZeroTriangle3d(
	t0, t1, t2 Point3d,
) (
	zero bool,
) {
	return ZeroLine3d(t0, t1) || ZeroLine3d(t1, t2) || ZeroLine3d(t2, t0)
}

// do not check point on triangle edge
func PointTriangle3d(
	p Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
) {
	A, B, C, D := Plane(t0, t1, t2)
	if !PointOnPlane3d(A, B, C, D, p) {
		// point is not plane
		return
	}
	// point on plane
	for _, v := range [3][4]*Point3d{
		{&t0, &p, &t1, &t2},
		{&t1, &p, &t2, &t0},
		{&t2, &p, &t0, &t1},
	} {
		var rA, rB float64
		rA, rB, intersect = LineLine3d(*v[0], *v[1], *v[2], *v[3])
		if !intersect || rA < 0.0 || rB < 0.0 {
			// point is not in triangle
			return false
		}
	}
	return true
}

// line intersect triangle in one point
func LineTriangle3dI1(
	l0, l1 Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
	pi []Point3d,
) {
	A, B, C, D := Plane(t0, t1, t2)
	if PointOnPlane3d(A, B, C, D, l0) && PointOnPlane3d(A, B, C, D, l1) {
		// Lines points on Plane
		return
	}

	// Line intersect Triangle on one point
	// div := ((l1[0]-l0[0])*A + (l1[1]-l0[1])*B + (l1[2]-l0[2])*C)
	div := math.FMA(l1[0]-l0[0], A, math.FMA(l1[1]-l0[1], B, (l1[2]-l0[2])*C))
	if math.Abs(div) < Eps3D {
		return
	}
	// Ka := (A*l0[0] + B*l0[1] + C*l0[2] + D) / (-div)
	Ka := math.FMA(A, l0[0], math.FMA(B, l0[1], math.FMA(C, l0[2], D))) / (-div)

	if Ka < 0 || 1 < Ka {
		return
	}

	p := PointLineRatio3d(l0, l1, Ka)

	if !PointTriangle3d(p, t0, t1, t2) {
		return
	}
	intersect = true
	pi = append(pi, p)
	return
}

// Line on triangle plane
// Line is not zero
// ignore triangle point on line
func LineTriangle3dI2(
	l0, l1 Point3d,
	t0, t1, t2 Point3d,
) (
	intersect bool,
	pi []Point3d,
) {
	A, B, C, D := Plane(t0, t1, t2)
	if !(PointOnPlane3d(A, B, C, D, l0) && PointOnPlane3d(A, B, C, D, l1)) {
		// line not on triangle plane
		return
	}
	// intersection line inside triangle
	for _, v := range [2]*Point3d{&l0, &l1} {
		if PointTriangle3d(*v, t0, t1, t2) {
			intersect = true
			pi = append(pi, *v)
		}
	}
	// line outside triangle
	for _, v := range [3][2]*Point3d{{&t0, &t1}, {&t1, &t2}, {&t2, &t0}} {
		if rA, rB, ill := LineLine3d(l0, l1, *v[0], *v[1]); ill &&
			0 <= rA && rA <= 1 && 0 <= rB && rB <= 1 {
			intersect = true
			p := PointLineRatio3d(l0, l1, rA)
			pi = append(pi, p)
		}
	}
	return
}

// do not intersect with egdes
func TriangleTriangle3d(
	a0, a1, a2 Point3d,
	b0, b1, b2 Point3d,
) (
	intersect bool,
	pi []Point3d,
) {
	var pis []Point3d
	for i := 0; i < 2; i++ {
		if i == 1 {
			a0, a1, a2, b0, b1, b2 = b0, b1, b2, a0, a1, a2 // swap
		}
		for _, f := range [2]func(Point3d, Point3d, Point3d, Point3d, Point3d) (bool, []Point3d){
			LineTriangle3dI1,
			LineTriangle3dI2,
		} {
			for _, v := range [3][2]*Point3d{{&a0, &a1}, {&a1, &a2}, {&a2, &a0}} {
				ilt, pit := f(*v[0], *v[1], b0, b1, b2)
				if ilt {
					intersect = intersect || ilt
					pis = append(pis, pit...)
				}
			}
		}
	}
	for i := range pis {
		if PointTriangle3d(pis[i], a0, a1, a2) &&
			PointTriangle3d(pis[i], b0, b1, b2) {
			pi = append(pi, pis[i])
			continue
		}
	}
	return
}
