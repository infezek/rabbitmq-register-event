package server

import (
	"fmt"
	"rabbitmain/internal"
	"rabbitmain/pkg/entity"
	"rabbitmain/pkg/rabbit"

	"github.com/spf13/cobra"
)

func runStart() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start simulation service",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("1231231")
			mongo, err := internal.NewMongo()
			if err != nil {
				panic(err)
			}
			event, err := rabbit.NewRabbit("amqp://admin:admin@localhost:5672", mongo)
			if err != nil {
				panic(err)
			}
			value := entity.Message{
				Value: "123123",
			}
			err = event.Publish(value)
			if err != nil {
				panic(err)
			}
			fmt.Println("000")
		},
	}
}
