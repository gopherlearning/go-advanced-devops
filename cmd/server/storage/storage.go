package storage

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/gopherlearning/go-advanced-devops/internal/metrics"
	"github.com/gopherlearning/go-advanced-devops/internal/repositories"
)

type Storage struct {
	// map[type]map[metric_name]map[target]value
	mu sync.RWMutex
	v  map[metrics.MetricType]map[string]map[string]interface{}
}

// NewStorage
func NewStorage() *Storage {
	return &Storage{
		v: map[metrics.MetricType]map[string]map[string]interface{}{
			metrics.CounterType: make(map[string]map[string]interface{}),
			metrics.GaugeType:   make(map[string]map[string]interface{}),
		},
	}
}

var rMetricURL = regexp.MustCompile(`^.*\/(\w+)\/(\w+)\/(-?[0-9\.]+)$`)
var _ repositories.Repository = new(Storage)

func (s *Storage) Update(target string, metric string) error {
	if len(target) == 0 {
		return repositories.ErrWrongTarget
	}
	match := rMetricURL.FindStringSubmatch(metric)
	if len(match) == 0 {
		return repositories.ErrBadMetric
	}
	metricType := metrics.MetricType(match[1])
	metricName := match[2]
	metricValue := match[3]

	if _, ok := s.v[metricType]; !ok {
		return repositories.ErrWrongMetricType
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.v[metricType][metricName]; !ok {
		s.v[metricType][metricName] = make(map[string]interface{})
	}
	if _, ok := s.v[metricType][metricName][target]; !ok {
		s.v[metricType][metricName][target] = nil
	}
	switch metricType {
	case metrics.CounterType:
		m, err := strconv.Atoi(metricValue)
		if err != nil {
			return repositories.ErrWrongMetricType
		}
		if s.v[metricType][metricName][target] == nil {
			s.v[metricType][metricName][target] = make([]int, 0)
		}
		mm, ok := s.v[metricType][metricName][target].([]int)
		if !ok {
			return repositories.ErrWrongMetricType
		}
		s.v[metricType][metricName][target] = append(mm, m)
	case metrics.GaugeType:
		m, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return repositories.ErrWrongMetricType
		}
		s.v[metricType][metricName][target] = m
	}
	return nil
}

func (s *Storage) List(targets ...string) map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string][]string)
	for mType := range s.v {

		for mName := range s.v[mType] {
			for target, value := range s.v[mType][mName] {
				if _, ok := res[target]; !ok {
					res[target] = make([]string, 0)
				}
				res[target] = append(res[target], fmt.Sprintf(`%s - %s - %v`, mType, mName, value))
			}
		}
	}
	return res
}

func (s *Storage) ListProm(targets ...string) []byte {
	panic("not implemented") // TODO: Implement
}
