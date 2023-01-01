package tools

type Err404 struct {
	message string
}

func (e Err404) Error() string {
	return e.message
}
func NewErr404(message string) *Err404 {
	return &Err404{
		message: message,
	}
}
