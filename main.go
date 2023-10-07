package main

import (
	"fmt"
	"geneblob/matrix"
)

func main() {
	m := matrix.New[matrix.Int32](10, 5)
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, matrix.Int32(it.Y*100+it.X))
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %d\n", it, m.GetIt(it))
	}
}
