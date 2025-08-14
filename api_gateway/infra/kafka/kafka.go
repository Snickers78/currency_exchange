package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaHook struct {
	kafkaProducer sarama.AsyncProducer
	topic         string
}

func NewKafkaHook(brokers []string, topic string) *KafkaHook {
	cfg := sarama.NewConfig()

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		panic(err)
	}

	hook := &KafkaHook{
		kafkaProducer: producer,
		topic:         topic,
	}

	go func() {
		for {
			select {
			case err := <-producer.Errors():
				log.Printf("Got error while sending message to kafka broker: %v", err)
			}
		}
	}()

	return hook
}

func (hook *KafkaHook) Fire(message string) error {
	msg := &sarama.ProducerMessage{
		Topic: hook.topic,
		Value: sarama.StringEncoder(message),
	}

	hook.kafkaProducer.Input() <- msg
	return nil
}
