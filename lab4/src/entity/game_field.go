package entity

import (
	"math/rand"
	"time"
	"utils"
)

type GameField struct {
	Height     int32
	Width      int32
	FoodStatic int32
	SnakeCache map[Point]*Snake
	Food       utils.Set[Point]
	IsMoving   bool
}

func NewGameField(h, w, foodStatic int32) *GameField {
	rand.NewSource(time.Now().UnixNano())
	return &GameField{
		Height:     h,
		Width:      w,
		FoodStatic: foodStatic,
		SnakeCache: make(map[Point]*Snake),
		Food:       utils.NewSet[Point](),
	}
}

func (gf *GameField) getCell(point Point) Cell {
	if gf.Food.Contains(point) {
		return Apple
	} else if gf.SnakeCache[point] != nil {
		if head, _ := gf.SnakeCache[point].Body.Top(); head == point {
			return SnakeHead
		} else {
			return SnakeBody
		}
	} else {
		return Empty
	}
}

func (gf *GameField) NextPos(point Point, direction Direction) Point {
	switch direction {
	case Up:
		point.Y = (point.Y + 1) % gf.Height
	case Down:
		point.Y = (point.Y - 1 + gf.Height) % gf.Height
	case Right:
		point.X = (point.X + 1) % gf.Width
	case Left:
		point.X = (point.X - 1 + gf.Width) % gf.Width
	}
	return point
}

func (gf *GameField) killSnake(snake *Snake) {
	for !snake.Body.IsEmpty() {
		pos, _ := snake.Body.PopTop()
		if rand.Intn(2) == 1 {
			gf.SpawnApple(pos)
		} else {
			delete(gf.SnakeCache, pos)
		}
	}
}

func (gf *GameField) SpawnApple(point Point) {
	gf.Food.Add(point)
}

func (gf *GameField) FindEmptyRandom() Point { // [!] Weak spot
	for {
		point := Point{X: rand.Int31n(gf.Width), Y: rand.Int31n(gf.Height)}
		if gf.isSquareEmpty(point) {
			return point
		}
	}
}

func (gf *GameField) MoveSnake(snake *Snake) {
	gf.IsMoving = true
	defer func() {
		gf.IsMoving = false
	}()

	if snake == nil || snake.Body.IsEmpty() {
		return
	}

	head, _ := snake.Body.Top()
	back, _ := snake.Body.Back()
	nextPos := gf.NextPos(head, snake.Direction)

	nextCell := gf.getCell(nextPos)
	switch nextCell {
	case Empty:
		gf.SnakeCache[nextPos] = snake
		gf.SnakeCache[head] = snake
		delete(gf.SnakeCache, back)
		snake.move(nextPos)
	case Apple:
		gf.SnakeCache[nextPos] = snake
		gf.SnakeCache[head] = snake
		snake.grow(nextPos)
		gf.Food.Remove(nextPos)
		gf.SpawnApple(gf.FindEmptyRandom())
	case SnakeBody:
		cached := gf.SnakeCache[nextPos]
		gf.killSnake(snake)
		if cached != snake {
			gf.SnakeCache[nextPos] = snake
		}
	case SnakeHead:
		gf.killSnake(snake)
		gf.killSnake(gf.SnakeCache[nextPos])
	}
}

func (gf *GameField) GetAllNotEmpty() map[Point]Cell {
	res := make(map[Point]Cell)
	for k := range gf.SnakeCache {
		res[k] = gf.getCell(k)
	}
	for _, k := range gf.Food.GetList() {
		res[k] = gf.getCell(k)
	}
	return res
}

func (gf *GameField) isSquareEmpty(point Point) bool {
	if gf.getCell(point) != Empty {
		return false
	}
	directions := []Direction{Left, Left, Up, Up, Right, Right, Right, Right, Down, Down, Down, Down, Left, Left, Left, Left, Up, Right, Right, Right, Up, Up, Left, Left}
	for _, dir := range directions {
		point = gf.NextPos(point, dir)
		if gf.getCell(point) != Empty {
			return false
		}
	}
	return true
}
