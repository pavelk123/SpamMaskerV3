package maskerspam

import (
	"fmt"
)

type producer interface {
	produce() ([]string, error)
}

type presenter interface {
	present(data []string) error
}

// Including inside 2 fields:
// producer - for data provider unit
// presenter - for data presenter unit.

type Service struct {
	prod producer
	pres presenter
}

func (s *Service) Run() error {
	data, err := s.prod.produce()
	if err != nil {
		return fmt.Errorf("service.producer.produce: %w", err)
	}

	for i := range data {
		data[i] = s.maskingURL(data[i])
	}

	if err = s.pres.present(data); err != nil {
		return fmt.Errorf("service.presentor.present: %w", err)
	}

	return nil
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

func NewService(prod producer, pres presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}
