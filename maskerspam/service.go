package maskerspam

import (
	"fmt"
	"sync"
)

type producer interface {
	produce() ([]string, error)
}

type presenter interface {
	present(data []string) error
}

// Service is structure for masking url service
// Including inside 2 fields:
// producer - for data provider unit
// presenter - for data presenter unit.

type Service struct {
	prod producer
	pres presenter
}

// NewService is constructor of Service

func NewService(prod producer, pres presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}

// Run is method for start Service working

func (s *Service) Run() error {
	data, err := s.prod.produce()
	if err != nil {
		return fmt.Errorf("service.producer.produce: %w", err)
	}

	data = s.process(data)

	if err = s.pres.present(data); err != nil {
		return fmt.Errorf("service.presentor.present: %w", err)
	}

	return nil
}

func (s *Service) process(data []string) []string {
	var wg sync.WaitGroup
	maxRoutineCount := 10

	resultData := make([]string, 0, cap(data))
	results := make(chan string, maxRoutineCount)

	for i := range data {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- s.maskingURL(data[i])
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		resultData = append(resultData, result)
	}

	return resultData
}

func (s *Service) maskingURL(str string) string {
	const symbolsDetectedCount int = 7

	startURLIndex := 0
	isMasking := false
	buffer := []byte(str)

	for index := range buffer {
		if buffer[index] == 'h' && string(buffer[index:index+7]) == "http://" {
			startURLIndex = index + symbolsDetectedCount
			isMasking = true
		}

		if startURLIndex != 0 && index >= startURLIndex && isMasking {
			if buffer[index] == ' ' {
				isMasking = false

				continue
			}

			buffer[index] = '*'
		}
	}

	return string(buffer)
}
