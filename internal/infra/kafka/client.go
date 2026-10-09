package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"slices"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Options struct {
	Brokers  []string
	TLS      bool
	SASL     string // "", "plain", "scram-sha-256", or "scram-sha-512"
	Username string
	Password string
	ClientID string
}

func (opts Options) security() (*tls.Config, sasl.Mechanism, error) {
	if len(opts.Brokers) == 0 {
		return nil, nil, errors.New("kafka brokers are required")
	}
	if slices.Contains(opts.Brokers, "") {
		return nil, nil, errors.New("kafka broker address cannot be empty")
	}
	var tlsConfig *tls.Config
	if opts.TLS {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	if opts.SASL == "" {
		if opts.Username != "" || opts.Password != "" {
			return nil, nil, errors.New("kafka SASL mechanism is required with credentials")
		}
		return tlsConfig, nil, nil
	}
	if opts.Username == "" || opts.Password == "" {
		return nil, nil, errors.New("kafka SASL username and password are required")
	}
	if !opts.TLS {
		return nil, nil, errors.New("kafka SASL requires TLS")
	}
	switch opts.SASL {
	case "plain":
		return tlsConfig, plain.Mechanism{Username: opts.Username, Password: opts.Password}, nil
	case "scram-sha-256":
		mechanism, err := scram.Mechanism(scram.SHA256, opts.Username, opts.Password)
		return tlsConfig, mechanism, err
	case "scram-sha-512":
		mechanism, err := scram.Mechanism(scram.SHA512, opts.Username, opts.Password)
		return tlsConfig, mechanism, err
	default:
		return nil, nil, fmt.Errorf("unsupported Kafka SASL mechanism %q", opts.SASL)
	}
}

// Ping checks one broker at startup or during an explicit readiness probe.
func Ping(ctx context.Context, opts Options) error {
	tlsConfig, mechanism, err := opts.security()
	if err != nil {
		return err
	}
	dialer := &kafkago.Dialer{Timeout: 3 * time.Second, TLS: tlsConfig, SASLMechanism: mechanism, ClientID: opts.ClientID}
	for _, broker := range opts.Brokers {
		conn, dialErr := dialer.DialContext(ctx, "tcp", broker)
		if dialErr == nil {
			return conn.Close()
		}
		err = errors.Join(err, fmt.Errorf("broker %s: %w", broker, dialErr))
	}
	return fmt.Errorf("connect Kafka: %w", err)
}

// NewWriter returns a synchronous producer. The caller owns Writer.Close().
// WriteMessages errors must be handled; a canceled write may have reached Kafka.
func NewWriter(opts Options, topic string) (*kafkago.Writer, error) {
	if topic == "" {
		return nil, errors.New("Kafka topic is required")
	}
	tlsConfig, mechanism, err := opts.security()
	if err != nil {
		return nil, err
	}
	return &kafkago.Writer{
		Addr: kafkago.TCP(opts.Brokers...), Topic: topic,
		Balancer: &kafkago.Hash{}, RequiredAcks: kafkago.RequireAll,
		Transport: &kafkago.Transport{TLS: tlsConfig, SASL: mechanism, ClientID: opts.ClientID},
	}, nil
}

// NewReader returns a consumer group reader. The caller owns Reader.Close().
// Use FetchMessage and CommitMessages to acknowledge only processed messages.
func NewReader(opts Options, topic, groupID string) (*kafkago.Reader, error) {
	if topic == "" || groupID == "" {
		return nil, errors.New("Kafka topic and consumer group are required")
	}
	tlsConfig, mechanism, err := opts.security()
	if err != nil {
		return nil, err
	}
	return kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: opts.Brokers, Topic: topic, GroupID: groupID,
		Dialer: &kafkago.Dialer{Timeout: 3 * time.Second, TLS: tlsConfig, SASLMechanism: mechanism, ClientID: opts.ClientID},
	}), nil
}
