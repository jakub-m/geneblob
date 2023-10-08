package main

import (
	"fmt"
	"geneblob/graph"
	"geneblob/matrix"
	"log"
	"math/rand"
)

func main() {
	//g := newGraph(4)
	//g.vertices = []XY{
	//	{100 - 50, 200 + 50},
	//	{200, 200},
	//	{200, 100},
	//	{100, 100},
	//}
	//g.forceConst = 0.1
	//g.baseDist = 100

	//setBoth(g.edges, 0, 1, Bool(true))
	//setBoth(g.edges, 1, 2, Bool(true))
	//setBoth(g.edges, 2, 3, Bool(true))
	//setBoth(g.edges, 3, 0, Bool(true))
	//setBoth(g.edges, 0, 2, Bool(true))

	rand.Seed(0)

	n := 100
	iterCount := 200
	edgeProb := 0.07

	g := graph.New(n)
	for i := range g.Vertices {
		g.Vertices[i] = graph.XY{
			X: 50 + 200*rand.Float64(),
			Y: 50 + 200*rand.Float64(),
		}
	}
	for i := range g.Vertices {
		for k := range g.Vertices {
			if i < k {
				if rand.Float64() < edgeProb {
					g.Edges.SetSym(i, k, true)
				}
			}
		}
	}

	for i := 0; i < iterCount; i++ {
		log.Printf("iter %d", i)
		g.UpdateForces()
		// printMatrix(g.forces)
		g.UpdatePoints()
		// fmt.Println(g.vertices)
		g.PrintPNG(fmt.Sprintf("tmp_%0.3d.png", i))
		// for it := g.forces.Iter(); it.HasNext(); it.Next() {
		// 	fmt.Printf("%s %v\n", it, g.forces.GetIt(it).Abs())
		// }
	}
}

func printMatrix[T matrix.Val[T]](m *matrix.Matrix[T]) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %v\n", it, m.GetIt(it))
	}
}
