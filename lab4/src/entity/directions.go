package entity

import "math/rand"

type Direction byte

const (
	Up Direction = iota
	Down
	Left
	Right
)

func GetOpposite(dir Direction) Direction {
	switch dir {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	default:
		return dir
	}
}

func NewDirectionRandom() Direction {
	return Direction(rand.Intn(4))
}
