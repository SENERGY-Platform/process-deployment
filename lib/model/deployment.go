package model

type Deployment struct {
	Id       string        `json:"id"`
	Xml      string        `json:"xml"`
	Name     string        `json:"name"`
	Elements []Element     `json:"elements"`
	Lanes    []LaneElement `json:"lanes"`
}

type Element interface {
	GetOrder() int
}

type LaneElement interface {
	GetOrder() int
}
