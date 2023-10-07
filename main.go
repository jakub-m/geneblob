package main

import (
	"fmt"
	"geneblob/matrix"
	"math"
)

type XY struct {
	X, Y float32
}

func (p XY) Add(q XY) XY {
	return XY{
		X: p.X + q.X,
		Y: p.Y + q.Y}
}

func (p XY) Mul(c matrix.C) XY {
	return XY{
		X: p.X * float32(c),
		Y: p.Y * float32(c),
	}
}

func (p XY) Abs() float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
}

var _ matrix.Val[XY] = (*XY)(nil)

func main() {
	//m := matrix.New[XY](10, 5)
	//for it := m.Iter(); it.HasNext(); it.Next() {
	//	m.SetIt(it, XY{float32(it.X), float32(it.Y)})
	//}
	//for it := m.Iter(); it.HasNext(); it.Next() {
	//	fmt.Printf("%s %f\n", it, m.GetIt(it))
	//}

	n := 4
	vertices := []XY{
		{0, 0},
		{0.1, 0},
		{0, 1},
		{2, 2},
	}
	assert(len(vertices) == n)
	edges := matrix.New[Bool](n, n)
	setBoth(edges, 0, 1, Bool(true))
	setBoth(edges, 1, 2, Bool(true))
	setBoth(edges, 2, 3, Bool(true))
	setBoth(edges, 3, 0, Bool(true))

	cForce := float32(0.1)
	d0 := float32(100.0)
	forces := matrix.New[XY](n, n)
	assert(matrix.SameSize(edges, forces))

	for it := edges.Iter(); it.HasNext(); it.Next() {
		f := XY{}
		if edges.GetIt(it) {
			pointA := vertices[it.X]
			pointB := vertices[it.Y]
			f = calculateForce(pointA, pointB, cForce, d0)
		}
		forces.SetIt(it, f)
	}

	printMatrix(forces)

	// The proportion of the force-to-movement of the point is already determined by cForce.
	for it := forces.Iter(); it.HasNext(); it.Next() {
	}
}

func printMatrix[T matrix.Val[T]](m *matrix.Matrix[T]) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %v\n", it, m.GetIt(it))
	}
}

func setBoth[T matrix.Val[T]](m *matrix.Matrix[T], x, y int, p T) {
	m.Set(x, y, p)
	m.Set(y, x, p)
}

func calculateForce(a, b XY, c, d0 float32) XY {
	// F_vec_ab =
	//   (b_vec-a_vec)/dist(a, b) * c_ab * (dist(a, b) - d0) =
	//   (b_vec-a_vec) * c_ab * (1 - d0/dist(a, b))
	xx := b.X - a.X
	yy := b.Y - a.Y
	dist_ab := float32(math.Sqrt(float64(xx*xx + yy*yy)))
	return b.Add(a.Mul(-1)).Mul(matrix.C(c)).Mul(matrix.C(1 - d0/dist_ab))
}

type Bool bool

func (f Bool) Add(g Bool) Bool {
	panic("Add not implemented for Bool")
}

func (f Bool) Mul(g matrix.C) Bool {
	panic("Mul not implemented for Bool")
}

func assert(cond bool) {
	if !cond {
		panic("condition not met")
	}
}
