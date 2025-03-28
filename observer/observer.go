package observer

type Event string

const (
	Create Event = "create"
	Update Event = "update"
)

type Observer[T any] interface {
	Notify(event Event, data T)
}
