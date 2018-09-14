package arukas

type ServicePlan struct {
	Category string `jsonapi:"attr,category"`
	Code     string `jsonapi:"attr,code"`
	Cpus     string `jsonapi:"attr,cpus"`
	ID       string `jsonapi:"primary,service-plans"`
	Memory   int    `jsonapi:"attr,memory"`
	Name     string `jsonapi:"attr,name"`
	Price    int    `jsonapi:"attr,price"`
	RegionID int    `jsonapi:"attr,region-id"`
	Version  int    `jsonapi:"attr,version"`

	CreatedAt string `jsonapi:"attr,created-at"`
	UpdatedAt string `jsonapi:"attr,updated-at"`
}
