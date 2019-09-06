package saramaservice

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	MaxErrMsgLen = 256
)

var userAgent = fmt.Sprintf("Prometheus/%s", version.Version)

type SaramaService struct {
	brokers  string
	topics   string
	group    string
	version  string
	storeUrl string
}

func NewSaramaService(brokers, group, topics, version, storeUrl string) *SaramaService {
	return &SaramaService{
		brokers:  brokers,
		topics:   topics,
		group:    group,
		version:  version,
		storeUrl: storeUrl,
	}
}

func (c *SaramaService) SyncRun() {
	go c.runSaramaConsumer()
}
func (c *SaramaService) Run() error {
	return c.runSaramaConsumer()
}

func (c *SaramaService) runSaramaConsumer() error {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(c.version)
	if err != nil {
		log.Fatalf("KafkaCollector run: ParseVersion %s failed: %v", c.version, err)
		return err
	}
	config.Version = version

	client, err := sarama.NewConsumerGroup(strings.Split(c.brokers, ","), c.group, config)
	if err != nil {
		log.Fatalf("KafkaCollector run: NewConsumerGroup error: %v", err)
		return err
	}
	defer client.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		err := client.Consume(ctx, strings.Split(c.topics, ","), c)
		if err != nil {
			log.Fatalf("KafkaCollector run: client consume error: %v", err)
			return err
		}
	}
	return nil
}

func (c *SaramaService) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *SaramaService) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *SaramaService) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for rawMessage := range claim.Messages() {
		session.MarkMessage(rawMessage, "")
		log.Infof("Topic: %s Partition:%v Offset:%v\n", rawMessage.Topic, rawMessage.Partition, rawMessage.Offset)
		if err := StoreM3db(rawMessage.Value, c.storeUrl); err != nil {
			log.Infof("Store error：%v", err.Error())
		}
	}
	return nil
}

// Store M3DB
func StoreM3db(req []byte, url string) error {
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(req))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	//httpReq = httpReq.WithContext(ctx)
	httpClient, err := config.NewClientFromConfig(config.HTTPClientConfig{}, "remote_storage", false)
	if err != nil {
		return err
	}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, MaxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		return err
	}
	return err
}
