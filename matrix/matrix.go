package matrix

import "fmt"

type Val float32

// Matrix is a matrix of points interpreted as (x, y) points, supporting linear operations.
type Matrix[T Val] struct {
	fields []float32
	yrange int
	xrange int
}

func New[T Val](xrange, yrange int) *Matrix[T] {
	return &Matrix[T]{
		fields: make([]float32, xrange*yrange),
		xrange: xrange,
		yrange: yrange,
	}
}

func (m *Matrix[T]) Fill(value float32) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, value)
	}
}

func (m *Matrix[T]) SetIt(it *Iter[T], v float32) {
	m.fields[it.index()] = v
}

func (m *Matrix[T]) GetIt(it *Iter[T]) float32 {
	return m.fields[it.index()]
}

func (m *Matrix[T]) Add(t *Matrix[T]) error {
	if err := m.checkSameSize(t); err != nil {
		return err
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)+t.GetIt(it))
	}
	return nil
}

func (m *Matrix[T]) AddVal(c float32) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)+c)
	}
	return nil
}

func (m *Matrix[T]) MulConst(c float32) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)*c)
	}
	return nil
}

func (m *Matrix[T]) checkSameSize(t *Matrix[T]) error {
	if !(m.xrange == t.xrange && m.yrange == t.yrange) {
		return fmt.Errorf("matrices differ in size: %s vs %s", m.getSizeString(), t.getSizeString())
	}
	return nil
}

func (m *Matrix[T]) getSizeString() string {
	return fmt.Sprintf("%dx%d", m.xrange, m.yrange)
}

type Iter[T Val] struct {
	X, Y int
	m    *Matrix[T]
}

func (m *Matrix[T]) Iter() *Iter[T] {
	return &Iter[T]{m: m}
}

func (it *Iter[T]) Next() {
	if !it.HasNext() {
		return
	}
	it.X++
	if it.X >= it.m.xrange {
		it.X = 0
		it.Y++
	}
}

func (it *Iter[T]) HasNext() bool {
	return it.Y < it.m.yrange
}

func (it *Iter[T]) index() int {
	return it.X + it.Y*it.m.xrange
}

func (it *Iter[T]) String() string {
	return fmt.Sprintf("%d,%d", it.X, it.Y)
}
