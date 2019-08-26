package kafka

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/segmentio/kafka-go"
	"github.com/wvanbergen/kazoo-go"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ExampleKafka() {
	config, err := config.LoadConfig("../../config.json")
	if err != nil {
		log.Println(err)
		return
	}
	config.Debug = true

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Println(err)
		return
	}

	closeZk, _, zkIp, err := Zookeeper(pool)
	if err != nil {
		log.Println(err)
		return
	}
	defer closeZk()
	config.ZookeeperUrl = zkIp + ":2181"

	//kafka
	closeKafka, err := Kafka(pool, config.ZookeeperUrl)
	if err != nil {
		log.Println(err)
		return
	}
	defer closeKafka()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumed := []string{}
	mux := sync.Mutex{}

	wait, waitDone := context.WithTimeout(context.Background(), 10*time.Second)

	err = Factory.NewConsumer(ctx, config, "test", func(delivery []byte) error {
		mux.Lock()
		defer mux.Unlock()
		consumed = append(consumed, string(delivery))
		go func() {
			time.Sleep(2 * time.Second) //wait for completed commit
			waitDone()
		}()
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}

	producer, err := Factory.NewProducer(ctx, config, "test")
	if err != nil {
		log.Println(err)
		return
	}
	err = producer.Produce("foo", []byte("my message"))
	if err != nil {
		log.Println(err)
		return
	}

	<-wait.Done()
	mux.Lock()
	defer mux.Unlock()
	fmt.Println("CONSUMED: ", consumed)

	//output:
	//CONSUMED:  [my message]
}

func Kafka(pool *dockertest.Pool, zookeeperUrl string) (closer func(), err error) {
	kafkaport, err := getFreePort()
	if err != nil {
		log.Fatalf("Could not find new port: %s", err)
	}
	networks, _ := pool.Client.ListNetworks()
	hostIp := ""
	for _, network := range networks {
		if network.Name == "bridge" {
			hostIp = network.IPAM.Config[0].Gateway
		}
	}
	log.Println("host ip: ", hostIp)
	env := []string{
		"ALLOW_PLAINTEXT_LISTENER=yes",
		"KAFKA_LISTENERS=OUTSIDE://:9092",
		"KAFKA_ADVERTISED_LISTENERS=OUTSIDE://" + hostIp + ":" + strconv.Itoa(kafkaport),
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=OUTSIDE:PLAINTEXT",
		"KAFKA_INTER_BROKER_LISTENER_NAME=OUTSIDE",
		"KAFKA_ZOOKEEPER_CONNECT=" + zookeeperUrl,
	}
	log.Println("start kafka with env ", env)
	kafkaContainer, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: "bitnami/kafka", Tag: "latest", Env: env, PortBindings: map[docker.Port][]docker.PortBinding{
		"9092/tcp": {{HostIP: "", HostPort: strconv.Itoa(kafkaport)}},
	}})
	if err != nil {
		return func() {}, err
	}
	err = pool.Retry(func() error {
		log.Println("try kafka connection...")
		conn, err := kafka.Dial("tcp", hostIp+":"+strconv.Itoa(kafkaport))
		if err != nil {
			log.Println(err)
			return err
		}
		defer conn.Close()
		return nil
	})
	return func() { kafkaContainer.Close() }, err
}

func Zookeeper(pool *dockertest.Pool) (closer func(), hostPort string, ipAddress string, err error) {
	zkport, err := getFreePort()
	if err != nil {
		log.Fatalf("Could not find new port: %s", err)
	}
	env := []string{}
	log.Println("start zookeeper on ", zkport)
	zkContainer, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: "wurstmeister/zookeeper", Tag: "latest", Env: env, PortBindings: map[docker.Port][]docker.PortBinding{
		"2181/tcp": {{HostIP: "", HostPort: strconv.Itoa(zkport)}},
	}})
	if err != nil {
		return func() {}, "", "", err
	}
	hostPort = strconv.Itoa(zkport)
	err = pool.Retry(func() error {
		log.Println("try zk connection...")
		zookeeper := kazoo.NewConfig()
		zk, chroot := kazoo.ParseConnectionString(zkContainer.Container.NetworkSettings.IPAddress)
		zookeeper.Chroot = chroot
		kz, err := kazoo.NewKazoo(zk, zookeeper)
		if err != nil {
			log.Println("kazoo", err)
			return err
		}
		_, err = kz.Brokers()
		if err != nil && strings.TrimSpace(err.Error()) != strings.TrimSpace("zk: node does not exist") {
			log.Println("brokers", err)
			return err
		}
		return nil
	})
	return func() { zkContainer.Close() }, hostPort, zkContainer.Container.NetworkSettings.IPAddress, err
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
