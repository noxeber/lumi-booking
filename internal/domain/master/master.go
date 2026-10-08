package master

import (
	"errors"
	"slices"
	"strings"
	"uuid"
)

var (
	ErrEmptyName               = errors.New("master name is empty")
	ErrServiceIDIsAlreadyExist = errors.New("master already provides this service")
)

type Master struct {
	id       uuid.UUID
	name     string
	services []uuid.UUID
}

func CreateMaster(name string) (*Master, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}
	return &Master{id: uuid.New(), name: name}, nil
}

func (m *Master) AddService(id uuid.UUID) error {
	for _, service := range m.services {
		if service == id {
			return ErrServiceIDIsAlreadyExist
		}
	}
	m.services = append(m.services, id)
	return nil
}

func (m *Master) RemoveService(id uuid.UUID) {
	for i, service := range m.services {
		if service == id {
			m.services = slices.Delete(m.services, i, i+1)
		}
	}
}
