package service

import "errors"

var ErrHallNotFound = errors.New("hall not found")

type Hall struct {
	Id       int
	Name     string
	Capacity int
}

func NewHallEntity(id int, name string, capacity int) Hall {
	return Hall{
		Id:       id,
		Name:     name,
		Capacity: capacity,
	}
}
