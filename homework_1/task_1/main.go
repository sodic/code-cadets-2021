package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"code-cadets-2021/homework_1/task_1/fizzbuzz"
)

var rules = []fizzbuzz.Rule{
	{Divisor: 3, Value: "Fizz"},
	{Divisor: 5, Value: "Buzz"},
}

func main() {
	start, end := parseArgs()

	fizzBuzzSequence, err := fizzbuzz.GetSequence(start, end, rules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Join(fizzBuzzSequence, ", "))
}

func parseArgs() (int, int) {
	start := flag.Int("start", 1, "The start of the FizzBuzz range.")
	end := flag.Int("end", 20, "The end of the FizzBuzz range.")
	flag.Parse()
	return *start, *end
}
