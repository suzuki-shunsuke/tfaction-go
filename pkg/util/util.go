package util //nolint:revive

func StrP(s string) *string {
	return &s
}

func IntP(i int) *int {
	return &i
}
