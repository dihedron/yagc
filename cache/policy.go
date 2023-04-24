package cache

type Policy interface {
	Trigger() bool
}

type Always struct{}

func (*Always) Trigger() bool {
	return true
}

type Never struct{}

func (*Never) Trigger() bool {
	return false
}

type Batched struct {
	Size  int
	count int
}

func (b *Batched) Trigger() bool {
	b.count++
	if b.count == b.Size {
		b.count = 0
		return true
	}
	return false
}
