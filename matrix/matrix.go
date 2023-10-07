package matrix

import "fmt"

// Matrix is a matrix of points interpreted as (x, y) points.
type Matrix struct {
	fields []float32
	yrange int
	xrange int
}

func New(xrange, yrange int) *Matrix {
	return &Matrix{
		fields: make([]float32, xrange*yrange),
		xrange: xrange,
		yrange: yrange,
	}
}

func (m *Matrix) Fill(value float32) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, value)
	}
}

func (m *Matrix) SetIt(it *Iter, v float32) {
	m.fields[it.index()] = v
}

func (m *Matrix) GetIt(it *Iter) float32 {
	return m.fields[it.index()]
}

func (m *Matrix) Add(t *Matrix) error {
	if err := m.checkSameSize(t); err != nil {
		return err
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)+t.GetIt(it))
	}
	return nil
}

func (m *Matrix) AddConst(c float32) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)+c)
	}
	return nil
}

func (m *Matrix) Mul(t *Matrix) error {
	if err := m.checkSameSize(t); err != nil {
		return err
	}
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)*t.GetIt(it))
	}
	return nil
}

func (m *Matrix) MulConst(c float32) error {
	for it := m.Iter(); it.HasNext(); it.Next() {
		m.SetIt(it, m.GetIt(it)*c)
	}
	return nil
}

func (m *Matrix) checkSameSize(t *Matrix) error {
	if !(m.xrange == t.xrange && m.yrange == t.yrange) {
		return fmt.Errorf("matrices differ in size: %s vs %s", m.getSizeString(), t.getSizeString())
	}
	return nil
}

func (m *Matrix) getSizeString() string {
	return fmt.Sprintf("%dx%d", m.xrange, m.yrange)
}

type Iter struct {
	X, Y int
	m    *Matrix
}

func (m *Matrix) Iter() *Iter {
	return &Iter{m: m}
}

func (it *Iter) Next() {
	if !it.HasNext() {
		return
	}
	it.X++
	if it.X >= it.m.xrange {
		it.X = 0
		it.Y++
	}
}

func (it *Iter) HasNext() bool {
	return it.Y < it.m.yrange
}

func (it *Iter) index() int {
	return it.X + it.Y*it.m.xrange
}

func (it *Iter) String() string {
	return fmt.Sprintf("%d,%d", it.X, it.Y)
}
