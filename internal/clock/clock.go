package clock

type Clock interface{ Now() string }
type Fixed struct{ Value string }

func (f Fixed) Now() string  { return f.Value }
func New(value string) Fixed { return Fixed{Value: value} }
