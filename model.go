package gog

type Model struct {
	Points []Point  // Points is slice of points
	Lines  [][2]int // Lines store 2 index of Points
	Arcs   [][3]int // Arcs store 3 index of Points
}

func (m *Model) AddLine() {
}

func (m *Model) AddCircle() {
}

func (m *Model) Split() {
}
