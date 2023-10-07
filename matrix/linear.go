package matrix

import "fmt"

type C float32

type Val[T any] interface {
	Add(T) T
	Mul(C) T
}

type Float32 float32

func (f Float32) Add(g Float32) Float32 {
	return f + g
}

func (f Float32) Mul(g C) Float32 {
	return Float32(float32(f) * float32(g))
}

// Matrix is a matrix of points interpreted as (x, y) points, supporting linear operations.
type Matrix[T Val[T]] struct {
	fields []T
	yrange int
	xrange int
}

func New[T Val[T]](xrange, yrange int) *Matrix[T] {
	return &Matrix[T]{
		fields: make([]T, xrange*yrange),
		xrange: xrange,
		yrange: yrange,
	}
}

func SameSize[T Val[T], Q Val[Q]](m *Matrix[T], t *Matrix[Q]) bool {
	return m.xrange == t.xrange && m.yrange == t.yrange
}

func (m *Matrix[T]) Fill(value T) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, value)
	}
}

func (m *Matrix[T]) Set(x, y int, value T) {
	i := y*m.xrange + x
	if i < 0 || i >= len(m.fields) {
		panic(fmt.Sprintf("index outside range i=%d, len(m.fields)=%d", i, len(m.fields)))
	}
	m.fields[i] = value
}

func (m *Matrix[T]) SetIt(it *Iter[T], v T) {
	m.fields[it.index()] = v
}

func (m *Matrix[T]) GetIt(it *Iter[T]) T {
	return m.fields[it.index()]
}

func (m *Matrix[T]) Add(t *Matrix[T]) error {
	if err := m.checkSameSize(t); err != nil {
		return err
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it).Add(t.GetIt(it)))
	}
	return nil
}

func (m *Matrix[T]) AddVal(c T) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it).Add(c))
	}
	return nil
}

func (m *Matrix[T]) MulConst(c C) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it).Mul(c))
	}
	return nil
}

func (m *Matrix[T]) checkSameSize(t *Matrix[T]) error {
	if !SameSize(m, t) {
		return fmt.Errorf("matrices differ in size: %s vs %s", m.getSizeString(), t.getSizeString())
	}
	return nil
}

func (m *Matrix[T]) getSizeString() string {
	return fmt.Sprintf("%dx%d", m.xrange, m.yrange)
}

type Iter[T Val[T]] struct {
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
