package clock

type Sequence struct {
	Values []string
	index  int
}

func (s *Sequence) Now() string {
	if len(s.Values) == 0 {
		return ""
	}
	v := s.Values[s.index%len(s.Values)]
	s.index++
	return v
}
func (s *Sequence) Reset() { s.index = 0 }
