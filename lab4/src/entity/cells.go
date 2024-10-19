package entity

type Cell int32

const (
	Empty     Cell = 0x000000
	Apple     Cell = 0xFF0000
	SnakeBody Cell = 0x006400
	SnakeHead Cell = 0x00FF00
)
