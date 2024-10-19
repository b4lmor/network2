package entity

type Role byte

const (
	None Role = iota
	Master
	Default
	Deputy
	Viewer
)
