package main

import (
	"fmt"
	"geneblob/matrix"
)

type Point struct {
	X, Y float32
}

//func (p Point) Add(q Point) Point {
//}
//
//func (p Point) Mul(q Point) Point {
//}
//
//
//var _ matrix.Val[Point] = (*Point)(nil)

func main() {
	m := matrix.New[matrix.Float32](10, 5)
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, matrix.Float32(it.Y*100+it.X))
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %f\n", it, m.GetIt(it))
	}
}

/*

F_vec_ab = (b_vec-a_vec) * c_ab * (1 - d0/dist(a, b))

*/
