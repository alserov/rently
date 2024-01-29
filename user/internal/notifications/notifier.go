package notifications

type Notifier interface {
	Registration(email string) error
	Login(email string) error
}

func NewNotifier() Notifier {
	return &notifier{}
}

type notifier struct {
}

func (n notifier) Login(email string) error {
	//TODO implement me
	panic("implement me")
}

func (n notifier) Registration(email string) error {
	//TODO implement me
	panic("implement me")
}
