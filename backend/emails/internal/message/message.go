package message

type EmailMessage interface {
	Recipients() ([]string, error)
	Subject() string
	Render() (string, error)
}
