package gog

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkLine3d(b *testing.B) {
	pps := []Point3d{
		Point3d{1, 1, 0}, // 0
		Point3d{4, 4, 0}, // 1
		Point3d{0, 5, 0}, // 2
		Point3d{5, 0, 0}, // 3
	}
	ra, rb, tint := LineLine3d(
		pps[0], pps[1],
		pps[2], pps[3],
	)
	if !tint {
		panic(fmt.Errorf("%v %v %v", ra, rb, tint))
	}
	for n := 0; n < b.N; n++ {
		LineLine3d(
			pps[0], pps[1],
			pps[2], pps[3],
		)
	}
}

func ExamplePlane() {
	delta := 1.0
	for i := 0; i < 30; i++ {
		A, B, C, D := Plane(
			Point3d{-10, -1, 10}, Point3d{10, -1, 10}, Point3d{0, 1 + delta, 0},
		)
		fmt.Fprintf(os.Stdout, "%.1e\t%e %e %e %e\n", delta, A, B, C, D)
		delta /= 10.0
	}

	// Output:
	// 1.0e+00	-0.000000e+00 2.000000e+02 6.000000e+01 -4.000000e+02
	// 1.0e-01	-0.000000e+00 2.000000e+02 4.200000e+01 -2.200000e+02
	// 1.0e-02	-0.000000e+00 2.000000e+02 4.020000e+01 -2.020000e+02
	// 1.0e-03	-0.000000e+00 2.000000e+02 4.002000e+01 -2.002000e+02
	// 1.0e-04	-0.000000e+00 2.000000e+02 4.000200e+01 -2.000200e+02
	// 1.0e-05	-0.000000e+00 2.000000e+02 4.000020e+01 -2.000020e+02
	// 1.0e-06	-0.000000e+00 2.000000e+02 4.000002e+01 -2.000002e+02
	// 1.0e-07	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-08	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-09	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-10	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-11	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-12	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-13	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-14	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-15	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-16	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-17	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-18	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-19	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-20	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-21	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-22	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-23	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-24	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-25	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-26	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-27	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-28	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
	// 1.0e-29	-0.000000e+00 2.000000e+02 4.000000e+01 -2.000000e+02
}

func ExamplePointPoint3d() {
	p := Point3d{1, 1, 1}
	delta := 1.0
	for i := 0; i < 30; i++ {
		pf := Point3d{1, 1, 1 + delta}
		value := PointPoint3d(p, pf)
		fmt.Fprintf(os.Stdout, "%.1e\t%v\n", delta, value)
		delta /= 10.0
	}

	// Output:
	// 1.0e+00	false
	// 1.0e-01	false
	// 1.0e-02	false
	// 1.0e-03	false
	// 1.0e-04	false
	// 1.0e-05	false
	// 1.0e-06	false
	// 1.0e-07	false
	// 1.0e-08	false
	// 1.0e-09	false
	// 1.0e-10	false
	// 1.0e-11	true
	// 1.0e-12	true
	// 1.0e-13	true
	// 1.0e-14	true
	// 1.0e-15	true
	// 1.0e-16	true
	// 1.0e-17	true
	// 1.0e-18	true
	// 1.0e-19	true
	// 1.0e-20	true
	// 1.0e-21	true
	// 1.0e-22	true
	// 1.0e-23	true
	// 1.0e-24	true
	// 1.0e-25	true
	// 1.0e-26	true
	// 1.0e-27	true
	// 1.0e-28	true
	// 1.0e-29	true
}

func Test3D(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			switch len(tc.ps) {
			case 3:
				var (
					p0 = Point3d{tc.ps[0].X, tc.ps[0].Y, 0}
					p1 = Point3d{tc.ps[1].X, tc.ps[1].Y, 0}
					p2 = Point3d{tc.ps[2].X, tc.ps[2].Y, 0}
				)
				intersect := PointLine3d(
					p0,
					p1, p2,
				)
				tint := tc.itA.Has(OnSegment) || tc.itB.Has(OnSegment)
				if intersect != tint {
					t.Errorf("not same")
				}
			case 4:
				var (
					p0 = Point3d{tc.ps[0].X, tc.ps[0].Y, 0}
					p1 = Point3d{tc.ps[1].X, tc.ps[1].Y, 0}
					p2 = Point3d{tc.ps[2].X, tc.ps[2].Y, 0}
					p3 = Point3d{tc.ps[3].X, tc.ps[3].Y, 0}
				)
				rA, rB, intersect := LineLine3d(
					p0, p1,
					p2, p3,
				)
				if intersect &&
					(rA < 0 || 1 < rA || rB < 0 || 1 < rB) {
					intersect = false
				}
				tint := tc.itA.Has(OnSegment) || tc.itB.Has(OnSegment)
				if intersect != tint {
					t.Errorf("not same: %f %f. %v", rA, rB, tc.pi)
				}
			}
		})
	}
	for _, tc := range triaTests {
		t.Run(tc.name, func(t *testing.T) {
			p := Point3d{tc.pt.X, tc.pt.Y, 0}
			tri0 := Point3d{tc.tri[0].X, tc.tri[0].Y, 0}
			tri1 := Point3d{tc.tri[1].X, tc.tri[1].Y, 0}
			tri2 := Point3d{tc.tri[2].X, tc.tri[2].Y, 0}
			intersect := PointTriangle3d(p, tri0, tri1, tri2)

			i01 := PointLine3d(p, tri0, tri1)
			i12 := PointLine3d(p, tri1, tri2)
			i20 := PointLine3d(p, tri2, tri0)

			intersect = intersect || i01 || i12 || i20

			tint := 0 < len(tc.res)
			if intersect != tint {
				t.Errorf("not same: %v", tc)
			}
		})
	}

	t.Run("LL false", func(t *testing.T) {
		rA, rB, intersect := LineLine3d(
			Point3d{0.5, 3, 0},
			Point3d{2.5, 3, 0},
			Point3d{0.25, 2.6, 0.4330127018922196},
			Point3d{0.4330127018922194, 2.8, 0.25},
		)
		if intersect {
			t.Errorf("false: %v, %v", rA, rB)
		}
	})

	t.Run("PT0 false", func(t *testing.T) {
		intersect := PointTriangle3d(
			Point3d{0, 0, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})
	t.Run("PT0 true", func(t *testing.T) {
		intersect := PointTriangle3d(
			Point3d{0, 0, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect {
			t.Errorf("false")
		}
	})

	t.Run("PT1 false", func(t *testing.T) {
		intersect := PointTriangle3d(
			Point3d{0, 1, 0},
			Point3d{0, 0, 0}, Point3d{0, 2, 0}, Point3d{1, 2, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})

	t.Run("ZeroT", func(t *testing.T) {
		zero := ZeroTriangle3d(
			Point3d{-1, 0, 0}, Point3d{-1, 0, 0}, Point3d{0, 10, 0},
		)
		if !zero {
			t.Errorf("false")
		}
	})
	t.Run("LT1 beg true", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0, 0, 0}, Point3d{0, 0, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 mid true", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0, 0, -1}, Point3d{0, 0, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 end true", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0, 0, -1}, Point3d{0, 0, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect {
			t.Errorf("false")
		}
	})

	t.Run("LT1 not in tri false", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{3, 3, -1}, Point3d{3, 3, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 beg false", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0, 0, 0.1}, Point3d{0, 0, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 end false", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0, 0, -1}, Point3d{0, 0, -0.1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 on plane false", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0.1, 0.1, 0}, Point3d{-0.1, -0.1, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})
	t.Run("LT1 on parallel false", func(t *testing.T) {
		intersect, _ := LineTriangle3dI1(
			Point3d{0.1, 0.1, 1}, Point3d{-0.1, -0.1, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false")
		}
	})

	t.Run("LT2 inside tri true", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-0.1, -0.1, 0}, Point3d{0.1, 0.1, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect || len(ip) != 2 {
			t.Errorf("false: %v", ip)
		}
	})
	t.Run("LT2 inside_outside0 true", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-3, -2, 0}, Point3d{0.1, 0.1, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect || len(ip) != 2 {
			t.Errorf("false: %v", ip)
		}
	})
	t.Run("LT2 inside_outside1 true", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-0.1, -0.1, 0}, Point3d{3, 2, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if !intersect || len(ip) != 2 {
			t.Errorf("false: %v", ip)
		}
	})

	t.Run("LT2 outside on point true", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-3, 10, 0}, Point3d{0.1, 10, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect {
			t.Errorf("false: %v %v", intersect, ip)
		}
	})
	t.Run("LT2 outside 1 false", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-3, 11, 0}, Point3d{0.1, 11, 0},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect || len(ip) != 0 {
			t.Errorf("false: %v", ip)
		}
	})
	t.Run("LT2 outside 2 false", func(t *testing.T) {
		intersect, ip := LineTriangle3dI2(
			Point3d{-3, 3, 0}, Point3d{3, 3, 1},
			Point3d{-1, -1, 0}, Point3d{1, -1, 0}, Point3d{0, 10, 0},
		)
		if intersect || len(ip) != 0 {
			t.Errorf("false: %v", ip)
		}
	})

	t.Run("TT0 false", func(t *testing.T) {
		intersect, ip := TriangleTriangle3d(
			Point3d{0.4330127018922193, 0.4, 0.25000000000000006},
			Point3d{1.2499999999999998, 0.2, 2.165063509461097},
			Point3d{2.1650635094610964, 0.4, 1.2500000000000002},
			Point3d{0.24999999999999997, 0.2, 0.43301270189221935},
			Point3d{0.4330127018922193, 0.4, 0.25000000000000006},
			Point3d{1.2499999999999998, 0.2, 2.165063509461097},
		)
		if intersect {
			t.Errorf("false: %v %v", intersect, ip)
		}
	})
	t.Run("TT1 true", func(t *testing.T) {
		intersect, ip := TriangleTriangle3d(
			Point3d{-1.1, 0, -1}, Point3d{1.1, 0, -1}, Point3d{0, 0, 1},
			Point3d{-1.0, -1, 0}, Point3d{1.0, -1, 0}, Point3d{0, 1, 0},
		)
		if !intersect || len(ip) != 2 {
			t.Errorf("false: %v %v", intersect, ip)
		}
	})

	t.Run("ZeroTriangle3d true", func(t *testing.T) {
		zero := ZeroTriangle3d(
			Point3d{0, 0, 0},
			Point3d{0, 2, 0},
			Point3d{1, 2, 0},
		)
		if zero {
			t.Errorf("false")
		}
	})
}
