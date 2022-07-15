package convert

type StrTo string

func (s StrTo) String() string {
	return string(s)
}
