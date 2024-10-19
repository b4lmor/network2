package utils

// Определяем тип для дженерик множества
type Set[T comparable] map[T]struct{}

// Создаем новое множество
func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

// Добавляем элемент в множество
func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Size() int32 {
	return int32(len(s))
}

// Проверяем наличие элемента в множестве
func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func (s Set[T]) GetList() []T {
	res := make([]T, 0, len(s))
	for v := range s {
		res = append(res, v)
	}
	return res
}
