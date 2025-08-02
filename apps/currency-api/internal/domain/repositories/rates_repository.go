package repositories

import "context"

type RatesRepository interface {
	GetRates(ctx context.Context, currencies []string) (map[string]float64, string, error)
}
