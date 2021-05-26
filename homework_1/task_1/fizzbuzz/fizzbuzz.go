package fizzbuzz

import (
	"errors"
	"strconv"
	"strings"
)

type Rule struct {
	Divisor int
	Value   string
}

func GetSequence(start, end int, rules []Rule) ([]string, error) {
	if end < 0 || start < 0 {
		return nil, errors.New("start and end must be positive integers")
	}

	if end <= start {
		return nil, errors.New("the end must be higher than the start")
	}

	mapper := makeMapper(rules)
	return mapRange(start, end, mapper), nil
}

func mapRange(start, end int, mapper func(int) string) []string {
	result := make([]string, 0, end-start)
	for i := start; i <= end; i++ {
		result = append(result, mapper(i))
	}
	return result
}

func makeMapper(rules []Rule) func(int) string {
	return func(number int) string {
		var sb strings.Builder
		for _, rule := range rules {
			if number%rule.Divisor == 0 {
				sb.WriteString(rule.Value)
			}
		}

		if sb.Len() != 0 {
			return sb.String()
		} else {
			return strconv.Itoa(number)
		}
	}
}
