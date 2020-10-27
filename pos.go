package main

type Pos struct {
	m_x int
	m_y int
}

func newPos() *Pos {
	return &Pos{}
}

func (this *Pos) init(x, y int) {
	this.m_x = x
	this.m_y = y
}

type ValuePos struct {
	m_value  int
	m_ptrPos *Pos
}

func newValuePos() *ValuePos {
	return &ValuePos{
		m_ptrPos: newPos(),
	}
}

func (this *ValuePos) init(value int, ptrPos *Pos) {
	this.m_value = value
	*this.m_ptrPos = *ptrPos
}
