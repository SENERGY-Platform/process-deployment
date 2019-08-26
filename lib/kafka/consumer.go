package kafka

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/segmentio/kafka-go"
	"io"
	"io/ioutil"
	"log"
	"time"
)

func NewConsumer(ctx context.Context, config config.Config, topic string, listener func(delivery []byte) error) error {
	broker, err := GetBroker(config.ZookeeperUrl)
	if err != nil {
		log.Println("ERROR: unable to get broker list", err)
		return err
	}

	err = InitTopic(config.ZookeeperUrl, topic)
	if err != nil {
		log.Println("ERROR: unable to create topic", err)
		return err
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval: 0, //synchronous commits
		Brokers:        broker,
		GroupID:        config.ConsumerGroup,
		Topic:          topic,
		MaxWait:        1 * time.Second,
		Logger:         log.New(ioutil.Discard, "", 0),
		ErrorLogger:    log.New(ioutil.Discard, "", 0),
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("close kafka reader ", topic)
				return
			default:
				m, err := r.FetchMessage(ctx)
				if err == io.EOF || err == context.Canceled {
					log.Println("close consumer for topic ", topic)
					return
				}
				if err != nil {
					log.Fatal("ERROR: while consuming topic ", topic, err)
					return
				}
				err = listener(m.Value)
				if err != nil {
					log.Println("ERROR: unable to handle message (no commit)", err)
				} else {
					err = r.CommitMessages(ctx, m)
					if err != nil {
						log.Fatal("ERROR: while committing consumption ", topic, err)
						return
					}
				}
			}
		}
	}()
	return nil
}
