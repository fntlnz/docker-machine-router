package events

type Listener interface {
	Listen() error
}
