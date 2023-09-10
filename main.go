package main

type Event interface {
	Consume() error
	Public() error
	Reprocess() error
}
