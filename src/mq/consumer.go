package mq

import "log"

var done chan bool

// StartConsume: Consume Message
func StartConsume(
	qName,
	cName string,
	callback func(msg []byte) bool,
) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	done = make(chan bool)

	go func() {
		for d := range msgs {
			processErr := callback(d.Body)
			if processErr {
				// TODO: Write message into Error Queue
			}
		}
	}()

	<-done

	channel.Close()
}

// StopConsume: Stop Consume
func StopConsume() {
	done <- true
}
