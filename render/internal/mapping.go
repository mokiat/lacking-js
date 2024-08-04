package internal

func newMapper[T any]() *Mapper[T] {
	return &Mapper[T]{
		mapping: make(map[uint32]T),
		freeID:  1,
	}
}

type Mapper[T any] struct {
	mapping map[uint32]T
	freeID  uint32
}

func (m *Mapper[T]) Allocate(v T) uint32 {
	id := m.freeID
	m.freeID++
	m.mapping[id] = v
	return id
}

func (m *Mapper[T]) Release(id uint32) {
	delete(m.mapping, id)
}

func (m *Mapper[T]) Get(id uint32) T {
	return m.mapping[id]
}

var (
	framebuffers = newMapper[*Framebuffer]()
	programs     = newMapper[*Program]()
	textures     = newMapper[*Texture]()
	samplers     = newMapper[*Sampler]()
	buffers      = newMapper[*Buffer]()
	vertexArrays = newMapper[*VertexArray]()
)
