package main

type MemStorage struct {
	storage map[string]string
}

type Storage interface {
	Save() error
	Remove() error
	IdExist() (bool, error)
}
