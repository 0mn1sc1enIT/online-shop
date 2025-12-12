package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type seqWriter struct {
	url    string
	client *http.Client
}

func (s *seqWriter) Write(p []byte) (n int, err error) {
	p = bytes.TrimSpace(p)
	if len(p) == 0 {
		return 0, nil
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/events/raw?clef", s.url), bytes.NewBuffer(p))
	if err != nil {
		fmt.Printf("Logger Error: failed to create seq request: %v\n", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/vnd.serilog.clef")

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Printf("Logger Error (Network): failed to send log to seq: %v\n", err)
		return len(p), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Logger Error (HTTP Status %d from Seq): %s\n", resp.StatusCode, string(bodyBytes))
	}
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

	zerolog.TimestampFieldName = "@t"
	zerolog.MessageFieldName = "@m"
	zerolog.LevelFieldName = "@l"

	var writer io.Writer

	if seqURL != "" {
		writer = &seqWriter{
			url:    seqURL,
			client: &http.Client{Timeout: 3 * time.Second},
		}
	} else {
		zerolog.TimestampFieldName = "time"
		zerolog.MessageFieldName = "message"
		zerolog.LevelFieldName = "level"
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger
}
