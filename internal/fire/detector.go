package fire

import "sort"

type Detector struct {
	ID   string `json:"id"`
	Zone string `json:"zone"`
}

func (s *Service) AddDetector(detector Detector) {
	s.detectors[detector.ID] = detector
}

func (s *Service) Detectors() []Detector {
	ids := make([]string, 0, len(s.detectors))
	for id := range s.detectors {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Detector, 0, len(ids))
	for _, id := range ids {
		out = append(out, s.detectors[id])
	}
	return out
}

func (s *Service) ZoneOfDetector(id string) string {
	return s.detectors[id].Zone
}

func (s *Service) Detect(detectorID string) error {
	zone := s.ZoneOfDetector(detectorID)
	return s.HandleSmoke(zone)
}

func (s *Service) DetectorCount() int {
	return len(s.detectors)
}
