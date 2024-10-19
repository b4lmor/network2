package utils

// Node представляет узел двусвязного списка
type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

// Deque представляет двустороннюю очередь
type Deque[T any] struct {
	front *Node[T]
	back  *Node[T]
	size  int
}

// NewDeque создает новую пустую двустороннюю очередь
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{}
}

// Top возвращает элемент в начале очереди
func (d *Deque[T]) Top() (T, bool) {
	if d.front == nil {
		var zero T
		return zero, false
	}
	return d.front.value, true
}

// Back возвращает элемент в конце очереди
func (d *Deque[T]) Back() (T, bool) {
	if d.back == nil {
		var zero T
		return zero, false
	}
	return d.back.value, true
}

// PopTop удаляет элемент из начала очереди
func (d *Deque[T]) PopTop() (T, bool) {
	if d.front == nil {
		var zero T
		return zero, false
	}
	value := d.front.value
	d.front = d.front.next
	if d.front != nil {
		d.front.prev = nil
	} else {
		d.back = nil // Если очередь пуста
	}
	d.size--
	return value, true
}

// PopBack удаляет элемент из конца очереди
func (d *Deque[T]) PopBack() (T, bool) {
	if d.back == nil {
		var zero T
		return zero, false
	}
	value := d.back.value
	d.back = d.back.prev
	if d.back != nil {
		d.back.next = nil
	} else {
		d.front = nil // Если очередь пуста
	}
	d.size--
	return value, true
}

// Size возвращает количество элементов в очереди
func (d *Deque[T]) Size() int {
	return d.size
}

// IsEmpty проверяет, пуста ли очередь
func (d *Deque[T]) IsEmpty() bool {
	return d.size == 0
}

// PushFront добавляет элемент в начало очереди
func (d *Deque[T]) PushFront(value T) {
	newNode := &Node[T]{value: value}
	if d.size == 0 {
		d.front = newNode
		d.back = newNode
	} else {
		newNode.next = d.front
		d.front.prev = newNode
		d.front = newNode
	}
	d.size++
}

// PushBack добавляет элемент в конец очереди
func (d *Deque[T]) PushBack(value T) {
	newNode := &Node[T]{value: value}
	if d.size == 0 {
		d.front = newNode
		d.back = newNode
	} else {
		newNode.prev = d.back
		d.back.next = newNode
		d.back = newNode
	}
	d.size++
}

func (d *Deque[T]) Get(i int) (T, bool) {
	if i < 0 || i >= d.size {
		var zero T
		return zero, false // Индекс вне границ
	}
	currentNode := d.front
	for j := 0; j < i; j++ {
		currentNode = currentNode.next
	}
	return currentNode.value, true
}
