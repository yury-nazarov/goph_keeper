package tools

type Err409 struct {
	message string
}

func (e Err409) Error() string {
	return e.message
}
func NewErr409(message string) *Err409 {
	return &Err409{
		message: message,
	}
}
