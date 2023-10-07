package main

import "fmt"

func main() {
	m := NewMatrix(10, 5)
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, float32(it.Y*100+it.X))
	}
	m.MulConst(0.5)
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %f\n", it, m.GetIt(it))
	}
}
