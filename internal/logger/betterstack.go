package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// BetterStackSink implements zapcore.WriteSyncer interface for sending logs to BetterStack
type BetterStackSink struct {
	token        string
	url          string
	client       *http.Client
	batchSize    int
	flushTimeout time.Duration
	buffer       [][]byte
	mutex        sync.Mutex
	timer        *time.Timer
	lastFlush    time.Time
}

// NewBetterStackSink creates a new BetterStack log sink with batch capabilities
func NewBetterStackSink(token, url string, batchSize int, flushInterval time.Duration) *BetterStackSink {
	sink := &BetterStackSink{
		token: token,
		url:   url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		batchSize:    batchSize,
		flushTimeout: flushInterval,
		buffer:       make([][]byte, 0, batchSize),
		lastFlush:    time.Now(),
	}

	// Start the timer for periodic flushing
	sink.resetTimer()

	return sink
}

// resetTimer creates a new timer for periodic flushing
func (s *BetterStackSink) resetTimer() {
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timer = time.AfterFunc(s.flushTimeout, func() {
		s.Flush()
		s.resetTimer()
	})
}

// Write implements io.Writer interface
func (s *BetterStackSink) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Make a copy of the log entry to avoid potential memory issues
	logCopy := make([]byte, len(p))
	copy(logCopy, p)

	// Add to buffer
	s.buffer = append(s.buffer, logCopy)

	// If we've reached the batch size, flush
	if len(s.buffer) >= s.batchSize {
		if err := s.flush(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

// flush sends the current batch of logs to BetterStack
func (s *BetterStackSink) flush() error {
	if len(s.buffer) == 0 {
		return nil
	}

	// Create a batch JSON array with the buffered logs
	var batch bytes.Buffer
	batch.WriteString("[")
	for i, log := range s.buffer {
		if i > 0 {
			batch.WriteString(",")
		}
		batch.Write(log)
	}
	batch.WriteString("]")

	// Clear the buffer before sending to prevent duplicate logs if this takes time
	s.buffer = make([][]byte, 0, s.batchSize)
	s.lastFlush = time.Now()

	// Send the batch to BetterStack
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://%s", s.url),
		&batch,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("User-Agent", "RemoteJobsWebScraper/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Flush sends any buffered logs immediately
func (s *BetterStackSink) Flush() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.flush()
}

// Sync implements zapcore.WriteSyncer interface
func (s *BetterStackSink) Sync() error {
	// Stop the timer to avoid concurrent flushes
	if s.timer != nil {
		s.timer.Stop()
	}

	return s.Flush()
}
