package repositories

import "errors"

var (
	ErrBadMetric       = errors.New("wrong metric format")
	ErrWrongMetricType = errors.New("wrong metric type")
	ErrWrongTarget     = errors.New("wrong target")
)
