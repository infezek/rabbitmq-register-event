package entity

var (
	ExchangeRequestCollects = "fr.direct"
	QueueRequestCollects    = "fr.request_collection"
)

type Event interface {
	Consume() error
	Public() error
	Reprocess() error
}

type Message struct {
	Value string
}
