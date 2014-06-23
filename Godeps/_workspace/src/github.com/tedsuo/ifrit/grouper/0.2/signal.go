package grouper

const (
	Continue = Signal("continue")
)

type Signal string

func (s Signal) Signal() {}
func (s Signal) String() string {
	return string(s)
}
