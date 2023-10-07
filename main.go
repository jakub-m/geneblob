package main

import (
	"fmt"
	"geneblob/matrix"
	"math"

	"github.com/fogleman/gg"
)

type XY struct {
	X, Y float64
}

func (p XY) Add(q XY) XY {
	return XY{
		X: p.X + q.X,
		Y: p.Y + q.Y}
}

func (p XY) Mul(c matrix.C) XY {
	return XY{
		X: p.X * float64(c),
		Y: p.Y * float64(c),
	}
}

func (p XY) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

var _ matrix.Val[XY] = (*XY)(nil)

func main() {
	g := newGraph(4)
	g.vertices = []XY{
		{0, 0},
		{10, 0},
		{0, 100},
		{200, 200},
	}
	g.forceConst = 0.1
	g.baseDist = 100

	// setBoth(g.edges, 0, 1, Bool(true))
	// setBoth(g.edges, 1, 2, Bool(true))
	// setBoth(g.edges, 2, 3, Bool(true))
	setBoth(g.edges, 3, 0, Bool(true))

	fmt.Println(g.vertices)
	for i := 0; i < 10; i++ {
		g.updateForces()
		// printMatrix(g.forces)
		g.updatePoints()
		fmt.Println(g.vertices)

		g.printPng(fmt.Sprintf("tmp_%d.png", i))
		for it := g.forces.Iter(); it.HasNext(); it.Next() {
			fmt.Printf("%s %v\n", it, g.forces.GetIt(it).Abs())
		}
		// fmt.Println(
		// 	lo.Map(g.forces, func(p XY, i int) float32 {p.Abs()})
		// )
	}
}

type graph struct {
	vertices   []XY
	edges      *matrix.Matrix[Bool]
	forces     *matrix.Matrix[XY]
	forceConst float64
	baseDist   float64
}

func newGraph(n int) *graph {
	return &graph{
		vertices:   make([]XY, n),
		edges:      matrix.New[Bool](n, n),
		forces:     matrix.New[XY](n, n),
		forceConst: 1.0,
		baseDist:   100,
	}

}

func (g *graph) updateForces() {
	for it := g.edges.Iter(); it.HasNext(); it.Next() {
		f := XY{}
		if g.edges.GetIt(it) {
			pointA := g.vertices[it.X]
			pointB := g.vertices[it.Y]
			f = calculateForce(pointA, pointB, g.forceConst, g.baseDist)
		}
		g.forces.SetIt(it, f)
	}

}

// updatePoints updates the positions of the points w.r.t the forces. The proportion of the force-to-movement
// of the point is already determined by the forceConst.
func (g *graph) updatePoints() {
	for it := g.forces.Iter(); it.HasNext(); it.Next() {
		i := it.X
		g.vertices[i] = g.vertices[i].Add(g.forces.GetIt(it))
	}
}

func (g *graph) printPng(path string) {
	dc := gg.NewContext(300, 300)
	for it := g.edges.Iter(); it.HasNext(); it.Next() {
		if !g.edges.GetIt(it) {
			continue
		}
		a := g.vertices[it.X]
		b := g.vertices[it.Y]
		dc.DrawLine(a.X, a.Y, b.X, b.Y)
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

// calculateForce calculates force vector working on A towards (because of) B.
func calculateForce(a, b XY, c, d0 float64) XY {
	// F_vec_ab =
	//   (b_vec-a_vec)/dist(a, b) * c_ab * (dist(a, b) - d0) =
	//   (b_vec-a_vec) * c_ab * (1 - d0/dist(a, b))
	xx := b.X - a.X
	yy := b.Y - a.Y
	dist_ab := math.Sqrt(xx*xx + yy*yy)
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
