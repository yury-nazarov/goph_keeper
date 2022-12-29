package tools

type Err500 struct {
	message string
}

func (e Err500) Error() string {
	return e.message
}

func NewErr500(message string) *Err500 {
	return &Err500{
		message: message,
	}
}
