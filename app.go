package arukas

import "time"

type App struct {
	ID        string     `jsonapi:"primary,apps"`
	Name      string     `jsonapi:"attr,name"`
	Services  []*Service `jsonapi:"relation,services"`
	User      *User      `jsonapi:"relation,user"`
	CreatedAt time.Time  `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt time.Time  `jsonapi:"attr,updated-at,iso8601"`
}

type Apps []*App
