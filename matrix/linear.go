package matrix

import "fmt"

type C float64

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

// Matrix is a matrix of points interpreted as (j, k) points, supporting linear operations.
type Matrix[T Val[T]] struct {
	fields []T
	jRange int
	kRange int
}

func New[T Val[T]](jRange, kRange int) *Matrix[T] {
	return &Matrix[T]{
		fields: make([]T, jRange*kRange),
		jRange: jRange,
		kRange: kRange,
	}
}

func SameSize[T Val[T], Q Val[Q]](m *Matrix[T], t *Matrix[Q]) bool {
	return m.jRange == t.jRange && m.kRange == t.kRange
}

func (m *Matrix[T]) Fill(value T) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, value)
	}
}

// SetSym sets the value symmetrically w.r.t the diagonal.
func (m *Matrix[T]) SetSym(j, k int, p T) {
	m.Set(j, k, p)
	m.Set(k, j, p)
}

func (m *Matrix[T]) Set(j, k int, value T) {
	i := k*m.jRange + j
	if i < 0 || i >= len(m.fields) {
		panic(fmt.Sprintf("index outside range i=%d, len(m.fields)=%d", i, len(m.fields)))
	}
	m.fields[i] = value
}

func (m *Matrix[T]) SetIt(it *Iter, v T) {
	m.fields[it.index()] = v
}

func (m *Matrix[T]) GetIt(it *Iter) T {
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
	return fmt.Sprintf("%dx%d", m.jRange, m.kRange)
}

type Iter struct {
	J, K, jRange, kRange int
}

func (m *Matrix[T]) Iter() *Iter {
	return &Iter{jRange: m.jRange, kRange: m.kRange}
}

func (it *Iter) Next() {
	if !it.HasNext() {
		return
	}
	it.J++
	if it.J >= it.jRange {
		it.J = 0
		it.K++
	}
}

func (it *Iter) HasNext() bool {
	return it.K < it.kRange
}

func (it *Iter) index() int {
	return it.J + it.K*it.jRange
}

func (it *Iter) String() string {
	return fmt.Sprintf("%d,%d", it.J, it.K)
}
