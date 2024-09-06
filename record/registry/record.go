package registry

type Record interface {
	FromMap(map[string]any) error
	RecordType() string
}
