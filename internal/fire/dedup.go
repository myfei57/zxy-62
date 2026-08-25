package fire

import "strconv"

var alarmSeq int

func (s *Service) RaiseAlarm(eventKey string) bool {
	key := eventKey + "-" + strconv.Itoa(alarmSeq)
	alarmSeq++
	if !s.deduper.Record(key) {
		return false
	}
	s.alarms = append(s.alarms, eventKey)
	return true
}
