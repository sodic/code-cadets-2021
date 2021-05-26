package tax

import (
	"math"
	"sort"

	"github.com/pkg/errors"
)

type Bucket struct {
	AppliedAbove float64
	TaxRate      float64
}

func NewCalculator(buckets []Bucket) (*calculator, error) {
	err := validateBuckets(buckets)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid tax buckets")
	}

	return &calculator{
		taxRanges: mapToDetailed(buckets),
	}, nil
}

func (calculator *calculator) Calculate(amount float64) (float64, error) {
	if amount < 0 {
		return 0, errors.New("the amount must not be negative")
	}

	totalTax := 0.0
	for _, taxRange := range calculator.taxRanges {
		ceiling := math.Min(taxRange.Ceiling, amount)
		totalTax += math.Max(taxRange.Rate*(ceiling-taxRange.Floor), 0)
	}

	return totalTax, nil
}

type taxRange struct {
	Floor   float64
	Ceiling float64
	Rate    float64
}

type calculator struct {
	taxBelowRange map[float64]float64
	taxRanges     []taxRange
}

func validateBuckets(buckets []Bucket) error {
	for _, bucket := range buckets {
		if bucket.AppliedAbove < 0 {
			return errors.New("the buckets cannot have negative thresholds")
		}
	}
	return nil
}

func mapToDetailed(classes []Bucket) []taxRange {
	result := make([]taxRange, 0, len(classes))
	lastCeil := math.Inf(1)

	for _, class := range sortDescending(classes) {
		result = append(result, taxRange{
			Ceiling: lastCeil,
			Floor:   class.AppliedAbove,
			Rate:    class.TaxRate,
		})
		lastCeil = class.AppliedAbove
	}

	return result
}

func sortDescending(classes []Bucket) []Bucket {
	sortedClasses := make([]Bucket, len(classes))
	copy(sortedClasses, classes)

	sort.Slice(sortedClasses, func(i, j int) bool {
		return classes[i].AppliedAbove > classes[j].AppliedAbove
	})

	return sortedClasses
}
