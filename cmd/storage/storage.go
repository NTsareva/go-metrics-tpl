package storage

type MemStorage struct {
	storage map[string]string
}

type Storage interface {
	Save(m *MemStorage) error
	Remove(m *MemStorage) error
	IfExist(m *MemStorage) (bool, error)
}
