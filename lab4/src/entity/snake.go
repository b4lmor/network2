package entity

import "utils"

type Snake struct {
	Body      *utils.Deque[Point]
	Direction Direction
}

func NewSnake(head, tail Point, direction Direction) *Snake {
	body := utils.NewDeque[Point]()
	body.PushFront(tail)
	body.PushFront(head)
	return &Snake{Body: body, Direction: direction}
}

func (s *Snake) GetScore() int32 {
	return int32(max(s.Body.Size()-2, 0))
}

func (s *Snake) ChangeDirection(direction Direction) {
	if GetOpposite(s.Direction) != direction {
		s.Direction = direction
	}
}

func (s *Snake) move(point Point) {
	s.Body.PopBack()
	s.Body.PushFront(point)
}

func (s *Snake) grow(point Point) {
	s.Body.PushFront(point)
}
