package entity

import (
	"time"
)

const (
	StatusPassed    = "passed"
	StatusOnAir     = "on_air"
	StatusScheduled = "scheduled"
)

type CinemaSession struct {
	Id        int
	MovieId   int
	HallId    int
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

func New(id, movieId, hallId int, startTime, endTime time.Time, timeZone *time.Location) CinemaSession {
	session := CinemaSession{
		Id:        id,
		MovieId:   movieId,
		HallId:    hallId,
		StartTime: startTime.In(timeZone),
		EndTime:   endTime.In(timeZone),
	}
	session.setStatus()
	return session
}

func (c *CinemaSession) setStatus() {
	current := time.Now().UTC()

	if c.StartTime.Before(current) && c.EndTime.After(current) {
		c.Status = StatusOnAir
	} else if c.EndTime.Before(current) {
		c.Status = StatusPassed
	} else {
		c.Status = StatusScheduled
	}
}
