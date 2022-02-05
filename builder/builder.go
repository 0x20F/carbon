package builder

type Segment struct {
	Priority int
	Key      string
	Value    string
}

type Command interface {
	Build() string
}
