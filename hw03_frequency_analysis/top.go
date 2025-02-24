package hw03frequencyanalysis

import (
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type Word struct {
	Data   string
	Amount int
}

func Top10(s string) []string {
	wordsMap := map[string]Word{}
	allWords := strings.Fields(s)
	// Собираем мапу
	for _, w := range allWords {
		if wStruct, exist := wordsMap[w]; exist {
			wStruct.Amount++
			wordsMap[w] = wStruct
			continue
		}
		wordsMap[w] = Word{
			Data:   w,
			Amount: 1,
		}
	}

	// Получаем слайс
	wordsSlice := maps.Values(wordsMap)
	// Сортируем слайс
	slices.SortFunc(wordsSlice, compareWords)

	length := min(10, len(wordsSlice))
	res := make([]string, 0, length)

	// Забираем первые 10 элементов
	for _, w := range wordsSlice[:length] {
		res = append(res, w.Data)
	}

	return res
}

func compareWords(w1, w2 Word) int {
	diff := w2.Amount - w1.Amount
	if diff != 0 {
		return diff
	}
	return strings.Compare(w1.Data, w2.Data)
}
