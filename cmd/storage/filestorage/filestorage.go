package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"os"
)

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadMetric() (*memstorage.Metrics, error) {

	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}

	data := c.scanner.Bytes()

	metric := memstorage.Metrics{}
	err := json.Unmarshal(data, &metric)
	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (c *Consumer) ReadMetrics() (metrics []*memstorage.Metrics, err error) {
	metric := memstorage.Metrics{}
	for c.scanner.Scan() {
		line := c.scanner.Text()

		err := json.Unmarshal([]byte(line), &metric)
		if err != nil {
			return nil, err
		}
		metricAddress := &metric

		metrics = append(metrics, metricAddress)
	}

	return metrics, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteMetric(metric *memstorage.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}
