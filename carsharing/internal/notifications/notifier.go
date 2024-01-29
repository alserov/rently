package notifications

type Notifier interface {
	RentCreated()
	RentCanceled()
}

func NewNotifier() Notifier {
	return &notifier{}
}

type notifier struct {
}

func (n notifier) RentCreated() {
	//TODO implement me
	panic("implement me")
}

func (n notifier) RentCanceled() {
	//TODO implement me
	panic("implement me")
}
