package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

type Consumer struct {
	consumer       *cluster.Consumer
	sig            chan int
	OnMsgReceiver  func(msg Message)
	OnMsgError     func(error)
	OnMsgRebalance func(ntf Notification)
	OnClosed       func()
}

func (c *Consumer) Prepare(cfg *ConsumerConfig) error {
	var err error

	clusterCfg := cluster.NewConfig()

	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = cfg.UserName
	clusterCfg.Net.SASL.Password = cfg.Password
	clusterCfg.Net.SASL.Handshake = true

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(globalConfig.Cert)
	if !ok {
		return errors.New("kafka consumer failed to parse root certificate")
	}

	clusterCfg.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	clusterCfg.Net.TLS.Enable = true
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	clusterCfg.Group.Return.Notifications = true

	clusterCfg.Version = sarama.V0_10_0_0
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		return errors.New(msg)
	}

	c.consumer, err = cluster.NewConsumer(cfg.Servers, cfg.ConsumerGroup, cfg.Topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		return errors.New(msg)
	}

	c.sig = make(chan int, 1)

	return nil
}

func (c *Consumer) consume() {
	for {
		select {
		case msg, more := <-c.consumer.Messages():
			if more {
				m := Message{}
				m.Key = string(msg.Key)
				m.Msg = string(msg.Value)
				m.Time = msg.Timestamp
				m.Topic = msg.Topic
				c.OnMsgReceiver(m)
				c.consumer.MarkOffset(msg, "") // mark message as processed
			}
		case err, more := <-c.consumer.Errors():
			if more {
				c.OnMsgError(err)
			}
		case ntf, more := <-c.consumer.Notifications():
			if more {
				n := Notification{}
				n.Claimed = ntf.Claimed
				n.Current = ntf.Current
				n.Released = ntf.Released
				c.OnMsgRebalance(n)
			}
		case <-c.sig:
			c.consumer.Close()
			c.OnClosed()
			return
		}
	}
}

func (c *Consumer) Stop() {
	c.sig <- 0
}

func (c *Consumer) Start() {
	go c.consume()
}
