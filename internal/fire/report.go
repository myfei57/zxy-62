package fire

type Report struct {
	DetectorCount int      `json:"detector_count"`
	AlarmCount    int      `json:"alarm_count"`
	ZoneCount     int      `json:"zone_count"`
	Alarms        []string `json:"alarms"`
}

func (s *Service) Report() Report {
	return Report{
		DetectorCount: len(s.detectors),
		AlarmCount:    len(s.alarms),
		ZoneCount:     len(s.zones),
		Alarms:        append([]string{}, s.alarms...),
	}
}
