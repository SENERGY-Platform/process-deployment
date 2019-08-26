package kafka

import (
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

func InitTopic(zkUrl string, topics ...string) (err error) {
	return InitTopicWithConfig(zkUrl, 1, 1, topics...)
}

func InitTopicWithConfig(zkUrl string, numPartitions int, replicationFactor int, topics ...string) (err error) {
	controller, err := GetKafkaController(zkUrl)
	if err != nil {
		log.Println("ERROR: unable to find controller", err)
		return err
	}
	if controller == "" {
		log.Println("ERROR: unable to find controller")
		return errors.New("unable to find controller")
	}
	initConn, err := kafka.Dial("tcp", controller)
	if err != nil {
		log.Println("ERROR: while init topic connection ", err)
		return err
	}
	defer initConn.Close()
	for _, topic := range topics {
		err = initConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
			ConfigEntries: []kafka.ConfigEntry{
				{ConfigName: "cleanup.policy", ConfigValue: "compact"},
				{ConfigName: "delete.retention.ms", ConfigValue: "100"},
				{ConfigName: "segment.ms", ConfigValue: "100"},
				{ConfigName: "min.cleanable.dirty.ratio", ConfigValue: "0.01"},
			},
		})
		if err != nil {
			return
		}
	}
	return nil
}
