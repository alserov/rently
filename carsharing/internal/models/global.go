package models

type Error struct {
	Status int
	Msg    string
}

func (e *Error) Error() string {
	return e.Msg
}

type CtxID string

const ID CtxID = "id"
