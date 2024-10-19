package mapper

import (
	"lab4/entity"
	"utils"
)

func MapList[T any, U any](input []T, mapper func(T) U) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = mapper(v)
	}
	return output
}

func MapListToSet[T comparable](input []*T) utils.Set[T] {
	res := utils.NewSet[T]()
	for _, v := range input {
		res.Add(*v)
	}
	return res
}

func MapMapToList[K comparable, T any, U any](input map[K]T, mapper func(T) U) []U {
	output := make([]U, len(input))
	i := 0
	for _, v := range input {
		output[i] = mapper(v)
		i++
	}
	return output
}

func FilterList[T any](input []T, pred func(T) bool) []T {
	output := make([]T, 0)
	for _, v := range input {
		if pred(v) {
			output = append(output, v)
		}
	}
	return output
}

func FilterMap[K comparable, T any](input map[K]T, pred func(K, T) bool) map[K]T {
	output := make(map[K]T)
	for k, v := range input {
		if pred(k, v) {
			output[k] = v
		}
	}
	return output
}

func MapPlayersToGameFieldCache(players map[string]*entity.Player) map[entity.Point]*entity.Snake {
	res := make(map[entity.Point]*entity.Snake)
	for _, player := range players {
		for i := 0; i < player.Snake.Body.Size(); i++ {
			pos, _ := player.Snake.Body.Get(i)
			res[pos] = player.Snake
		}
	}
	return res
}
