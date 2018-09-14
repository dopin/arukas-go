package arukas

type User struct {
	ID   string `jsonapi:"primary,users"`
	Name string `jsonapi:"attr,name"`
}
