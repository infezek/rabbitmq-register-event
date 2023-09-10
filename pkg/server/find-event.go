package server

import (
	"fmt"
	"rabbitmain/internal"

	"rabbitmain/pkg/rabbit"

	"github.com/spf13/cobra"
)

func runFind() *cobra.Command {
	return &cobra.Command{
		Use:   "find",
		Short: "find",
		Run: func(cmd *cobra.Command, args []string) {
			arg := args[0]
			mongo, err := internal.NewMongo()
			if err != nil {
				panic(err)
			}
			messagesEvent, err := mongo.FindEvent(arg)
			if err != nil {
				panic(err)
			}
			rabbitmq, err := rabbit.NewRabbit("amqp://admin:admin@localhost:5672", mongo)
			for _, message := range messagesEvent {

				rabbitmq.Publish(message.Body)
			}
			fmt.Println("end", messagesEvent)
		},
	}
}
