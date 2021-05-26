package tax_test

import (
	"fmt"
	"testing"

	"code-cadets-2021/homework_1/task_2/tax"

	. "github.com/smartystreets/goconvey/convey"
)

type TestCase struct {
	calculator TaxCalculator
	input      float64

	expectedOutput float64
	expectingError bool
}

var testCases = []TestCase{
	{
		calculator: newTaxCalculator(noTaxes),
		input:      -10,

		expectedOutput: 0,
		expectingError: true,
	},
	{
		calculator: newTaxCalculator(noTaxes),
		input:      0,

		expectedOutput: 0,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(noTaxes),
		input:      100_000,

		expectedOutput: 0,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(noTaxes),
		input:      100_000,

		expectedOutput: 0,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      -10,

		expectedOutput: 0,
		expectingError: true,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      1_000,

		expectedOutput: 0,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      1_010,

		expectedOutput: 0.1 * 10,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      5_000,

		expectedOutput: 0.1 * 4000,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      5000 + 1000,

		expectedOutput: 4000*0.1 + 1000*0.2,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      7123.5,

		expectedOutput: 4000*0.1 + 2123.5*0.2,
		expectingError: false,
	},
	{
		calculator: newTaxCalculator(sampleTaxes),
		input:      156_237_253.54,

		expectedOutput: 4000*0.1 + 5000*0.2 + (156_237_253.54-1000-4000-5000)*0.3,
		expectingError: false,
	},
}

var sampleTaxes = []tax.Bucket{
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

func TestNewCalculator(t *testing.T) {
	Convey("Given an array of valid tax buckets", t, func() {
		calculator, err := tax.NewCalculator(sampleTaxes)
		So(err, ShouldBeNil)
		So(calculator, ShouldNotBeNil)
	})

	Convey("Given an array of containing an invalid tax bucket", t, func() {
		invalidTaxes := append(sampleTaxes, tax.Bucket{
			AppliedAbove: -10,
			TaxRate:      0.5,
		})
		_, err := tax.NewCalculator(invalidTaxes)
		So(err, ShouldNotBeNil)
	})
}

var noTaxes []tax.Bucket

func TestCalculator_Calculate(t *testing.T) {
	for idx, tc := range testCases {
		Convey(fmt.Sprintf("Given test case #%v: %+v", idx, tc), t, func() {
			actualOutput, actualErr := tc.calculator.Calculate(tc.input)
			if tc.expectingError {
				So(actualErr, ShouldNotBeNil)
			} else {
				So(actualErr, ShouldBeNil)
				So(actualOutput, ShouldEqual, tc.expectedOutput)
			}
		})
	}
}

func newTaxCalculator(taxBuckets []tax.Bucket) TaxCalculator {
	calculator, _ := tax.NewCalculator(taxBuckets)
	return calculator
}

type TaxCalculator interface {
	Calculate(amount float64) (float64, error)
}
