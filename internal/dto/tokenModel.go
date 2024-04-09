package dto

type Token struct {
	value string
	age   int
}

func (t *Token) GetValue() string {
	return t.value
}

func (t *Token) GetAge() int {
	return t.age
}
