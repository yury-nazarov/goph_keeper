package tools

type Err401 struct {
	message string
}

func (e Err401) Error() string {
	return e.message
}
func NewErr401(message string) *Err401 {
	return &Err401{
		message: message,
	}
}
