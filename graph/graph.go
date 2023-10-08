// graph module models a graph of the forces acting along the edges on the vertices.

package graph

import (
	"geneblob/matrix"
	"image"
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

type Bool bool

func (f Bool) Add(g Bool) Bool {
	panic("Add not implemented for Bool")
}

func (f Bool) Mul(g matrix.C) Bool {
	panic("Mul not implemented for Bool")
}

var _ matrix.Val[Bool] = (*Bool)(nil)

type Graph struct {
	Vertices   []XY
	Edges      *matrix.Matrix[Bool]
	forces     *matrix.Matrix[XY]
	ForceConst float64
	BaseDist   float64
}

func New(n int) *Graph {
	return &Graph{
		Vertices:   make([]XY, n),
		Edges:      matrix.New[Bool](n, n),
		forces:     matrix.New[XY](n, n),
		ForceConst: 0.1,
		BaseDist:   100,
	}

}

func (g *Graph) UpdateForces() {
	for it := g.Edges.Iter(); it.HasNext(); it.Next() {
		f := XY{}
		if g.Edges.GetIt(it) {
			pointA := g.Vertices[it.J]
			pointB := g.Vertices[it.K]
			f = calculateForce(pointA, pointB, g.ForceConst, g.BaseDist)
		}
		g.forces.SetIt(it, f)
	}

}

// UpdatePoints updates the positions of the points w.r.t the forces. The proportion of the force-to-movement
// of the point is already determined by the forceConst.
func (g *Graph) UpdatePoints() {
	for it := g.forces.Iter(); it.HasNext(); it.Next() {
		i := it.J
		g.Vertices[i] = g.Vertices[i].Add(g.forces.GetIt(it))
	}
}

func (g *Graph) DrawImage() image.Image {
	dc := g.drawNewImage()
	return dc.Image()
}

func (g *Graph) SavePNG(path string) {
	dc := g.drawNewImage()
	dc.SavePNG(path)
}

func (g *Graph) drawNewImage() *gg.Context {
	lineWidth := 1.0
	dotSize := 2.0
	colors := []struct{ r, g, b float64 }{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 0},
		{0, 1, 1},
		{1, 0, 1},
	}

	dc := gg.NewContext(300, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetLineWidth(lineWidth)
	iColor := 0
	for it := g.Edges.Iter(); it.HasNext(); it.Next() {
		if !g.Edges.GetIt(it) {
			continue
		}
		a := g.Vertices[it.J]
		b := g.Vertices[it.K]
		color := colors[iColor]
		dc.SetRGB(color.r, color.g, color.b)
		iColor = (iColor + 1) % len(colors)
		dc.DrawLine(a.X, a.Y, b.X, b.Y)
		dc.Stroke()
		// fmt.Printf("draw %v->%v\n", a, b)
	}
	for _, v := range g.Vertices {
		dc.DrawPoint(v.X, v.Y, dotSize)
		dc.Stroke()
	}
	return dc
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
