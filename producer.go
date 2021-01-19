package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func (p *Producer) Prepare(cfg *ProducerConfig) error {
	var err error

	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = cfg.Ak
	mqConfig.Net.SASL.Password = cfg.Password
	mqConfig.Net.SASL.Handshake = true
	mqConfig.Version = sarama.V0_10_2_1

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(globalConfig.Cert)
	if !ok {
		return errors.New("kafka producer failed to parse root certificate")
	}

	mqConfig.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Successes = true

	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *cfg, err)
		return errors.New(msg)
	}

	p.producer, err = sarama.NewSyncProducer(cfg.Servers, mqConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		return errors.New(msg)
	}

	return nil
}

func (p *Producer) SendMsg(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	if len(key) > 0 {
		msg.Key = sarama.StringEncoder(key)
	}
	msg.Value = sarama.StringEncoder(content)

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		msg := fmt.Sprintf("Kafka send message error. topic: %v. key: %v. content: %v .err=%v", topic, key, content, err)
		return errors.New(msg)
	}

	return nil
}
