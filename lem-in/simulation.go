package main

import (
	"fmt"
	"strings"
)

type Ant struct {
	ID    int
	Path  []*Room
	Index int // -1 means at Start, len(Path) means End
}

func Simulate(graph *Graph, paths [][]*Room, inputLines string) {
	// 1. جهّز المسارات (أضف end)
	fullPaths := make([][]*Room, len(paths))
	for i, p := range paths {
		fullPaths[i] = append(p, graph.End)
	}

	// 2. وزّع النمل (سريع جدًا)
	antsPerPath := DistributeAnts(fullPaths, int64(graph.NumAnts))

	// 3. محاكاة الحركة
	activeAnts := make([][]*Ant, len(fullPaths))
	var nextID int64 = 1
	var finished int64 = 0

	for finished < int64(graph.NumAnts) {
		var moves []string

		for pIdx, path := range fullPaths {

			// حرّك النمل الموجود
			for _, ant := range activeAnts[pIdx] {
				ant.Index++
				if ant.Index >= len(path) {
					continue
				}

				room := path[ant.Index]
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, room.Name))

				if room == graph.End {
					finished++
				}
			}

			// أضف نملة جديدة إذا متاح
			if antsPerPath[pIdx] > 0 {
				ant := &Ant{
					ID:    int(nextID),
					Path:  path,
					Index: -1,
				}
				nextID++
				antsPerPath[pIdx]--

				ant.Index++
				room := path[ant.Index]
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, room.Name))

				if room != graph.End {
					activeAnts[pIdx] = append(activeAnts[pIdx], ant)
				} else {
					finished++
				}
			}

			// احذف النمل الذي وصل للنهاية
			tmp := activeAnts[pIdx][:0]
			for _, ant := range activeAnts[pIdx] {
				if ant.Index < len(path)-1 {
					tmp = append(tmp, ant)
				}
			}
			activeAnts[pIdx] = tmp
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}

func DistributeAnts(paths [][]*Room, numAnts int64) []int64 {
	n := len(paths)
	ants := make([]int64, n)

	// أطوال المسارات
	lengths := make([]int64, n)
	for i := 0; i < n; i++ {
		lengths[i] = int64(len(paths[i]))
	}

	// ابحث عن أقل T
	var T int64
	for {
		var sum int64
		for i := 0; i < n; i++ {
			if T >= lengths[i] {
				sum += T - lengths[i] + 1
			}
		}
		if sum >= numAnts {
			break
		}
		T++
	}

	// وزّع النمل
	remaining := numAnts
	for i := 0; i < n && remaining > 0; i++ {
		if T >= lengths[i] {
			ants[i] = T - lengths[i] + 1
			if ants[i] > remaining {
				ants[i] = remaining
			}
			remaining -= ants[i]
		}
	}

	return ants
}