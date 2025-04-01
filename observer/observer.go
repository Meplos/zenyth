package observer

type Event string

const (
	Create     Event = "create"
	Update     Event = "update"
	Errored    Event = "errored"
	Terminated       = "terminated"
)

type Observer[T any] interface {
	Notify(event Event, data T)
}
