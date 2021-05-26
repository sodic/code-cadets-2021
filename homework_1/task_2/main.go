package main

import (
	"fmt"
	"log"

	"code-cadets-2021/homework_1/task_2/tax"
)

var taxBuckets = []tax.Bucket{
	{
		AppliedAbove: 1_000,
		TaxRate:      0.1,
	},
	{
		AppliedAbove: 5_000,
		TaxRate:      0.2,
	},
	{
		AppliedAbove: 10_000,
		TaxRate:      0.3,
	},
}

func main() {
	taxCalculator, err := tax.NewCalculator(taxBuckets)
	if err != nil {
		log.Fatal(err)
	}

	taxAmount, err := fmt.Println(taxCalculator.Calculate(7000))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(taxAmount)
}
