package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// seqWriter реализует интерфейс io.Writer для отправки логов в Seq
type seqWriter struct {
	url    string
	client *http.Client
}

func (s *seqWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/events/raw?clef", s.url), bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/vnd.serilog.clef")

	resp, err := s.client.Do(req)
	if err != nil {
		// Если Seq недоступен, не паникуем, просто игнорируем (или пишем в stderr)
		return len(p), nil
	}
	defer resp.Body.Close()

	return len(p), nil
}

func New(level string, seqURL string) *zerolog.Logger {
	var l zerolog.Level

	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	if seqURL != "" {
		seqW := &seqWriter{
			url:    seqURL,
			client: &http.Client{Timeout: 2 * time.Second},
		}

		multi := zerolog.MultiLevelWriter(consoleWriter, seqW)
		logger := zerolog.New(multi).With().Timestamp().Logger()
		return &logger
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	return &logger
}
