package observer

type Event string

type Observer[T any] interface {
	Notify(event Event, data T)
}
