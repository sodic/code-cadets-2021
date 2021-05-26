package fizzbuzz_test

import (
	"testing"

	"code-cadets-2021/homework_1/task_1/fizzbuzz"
)

type TestCase struct {
	start     int
	end       int
	overrides []fizzbuzz.Rule

	expectedOutput []string
	expectingError bool
}

var testCases = []TestCase{
	{
		start:     1,
		end:       1,
		overrides: []fizzbuzz.Rule{},

		expectedOutput: nil,
		expectingError: true,
	},
	{
		start:     -2,
		end:       -1,
		overrides: []fizzbuzz.Rule{},

		expectedOutput: nil,
		expectingError: true,
	},
	{
		start:     1,
		end:       2,
		overrides: []fizzbuzz.Rule{},

		expectedOutput: []string{
			"1",
			"2",
		},
		expectingError: false,
	},
	{
		start:     1,
		end:       3,
		overrides: []fizzbuzz.Rule{},

		expectedOutput: []string{
			"1",
			"2",
			"3",
		},
		expectingError: false,
	},
	{
		start: 1,
		end:   3,
		overrides: []fizzbuzz.Rule{
			{
				Divisor: 2,
				Value:   "Test",
			},
		},

		expectedOutput: []string{
			"1",
			"Test",
			"3",
		},
		expectingError: false,
	},
	{
		start: 1,
		end:   8,
		overrides: []fizzbuzz.Rule{
			{
				Divisor: 2,
				Value:   "Fizz",
			},
			{
				Divisor: 3,
				Value:   "Buzz",
			},
		},

		expectedOutput: []string{
			"1",
			"Fizz",
			"Buzz",
			"Fizz",
			"5",
			"FizzBuzz",
			"7",
			"Fizz",
		},
		expectingError: false,
	},
}

func TestGetSequence(t *testing.T) {
	for _, tc := range testCases {
		actualOutput, actualError := fizzbuzz.GetSequence(tc.start, tc.end, tc.overrides)
		if tc.expectingError {
			if actualError == nil {
				t.Errorf(
					"Expected an error for input (%v, %v, %v) but got 'nil'",
					tc.start,
					tc.end,
					tc.overrides,
				)
			}
		} else {
			if actualError != nil {
				t.Errorf("Expected no error but got an error != nil")
			} else {
				if !areSlicesEqual(tc.expectedOutput, actualOutput) {
					t.Errorf(
						"Actual output different than expected - actual: %v, expected: %v",
						actualOutput,
						tc.expectedOutput)
				}
			}
		}
	}
}

func areSlicesEqual(first []string, second []string) bool {
	if len(first) != len(second) {
		return false
	}

	for idx, x := range first {
		if x != second[idx] {
			return false
		}
	}

	return true
}
