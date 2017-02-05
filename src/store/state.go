package store

type State string

const (
	StateAdded     State = ""
	StateStarted   State = "started"
	StateFailed    State = "failed"
	StateCompleted State = "completed"
)

func (s State) Human() string {
	switch s {
	case StateCompleted:
		return "Выполнен"
	case StateFailed:
		return "Провален"
	case StateStarted:
		return "Запущен"
	}
	return "Запланирован"
}

func (s State) Html() string {
	switch s {
	case StateStarted:
		return "warning"
	case StateFailed:
		return "danger"
	case StateCompleted:
		return "success"
	default:
		return "info"
	}
}
