package main

import (
	"flag"
	"github.com/getsentry/sentry-go"
	"github.com/hpcloud/tail"
	"github.com/satyrius/gonx"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Parse()
	if sentryConfig.Dsn == "" {
		sentryConfig.Dsn = os.Getenv("SENTRY_DSN")
	}
	if sentryConfig.ServerName == "" {
		var err error
		if sentryConfig.ServerName, err = os.Hostname(); err != nil {
			log.Fatalf("can't get hostname from OS: %v", err)
		}
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryConfig.Dsn,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Debug:            sentryConfig.Debug,
	}); err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	log.Printf("server name: %s", sentryConfig.ServerName)
	log.Printf("env: %s", sentryConfig.Env)
	log.Printf("debug: %v", sentryConfig.Debug)
	log.Printf("file: %v", nginxConfig.LogFile)
	log.Printf("time format: %v", nginxConfig.TimeFormat)

	t, err := tail.TailFile(nginxConfig.LogFile, tail.Config{
		Location: &tail.SeekInfo{
			Whence: 2,
		},
		Follow:    true,
		MustExist: false,
		ReOpen:    true,
	})
	if err != nil {
		log.Fatalf("tail.TailFile error: %v", err)
	}

	var cr io.Reader
	if nginxConfig.Config != "" {
		cr = strings.NewReader(nginxConfig.Config)
	} else {
		if cr, err = os.Open(nginxConfig.ConfigFile); err != nil {
			log.Fatalf("error opening config file: %v", err)
		}
	}

	parser, err := gonx.NewNginxParser(cr, "main")
	if err != nil {
		log.Fatalf("gonx.NewNginxParser error: %v", err)
	}

	for line := range t.Lines {
		if line.Text == "" {
			continue
		}
		entry, err := parser.ParseString(line.Text)
		if err != nil {
			log.Fatalf("parser.ParseString error: %v", err)
		}
		status, err := entry.IntField("status")
		if err != nil {
			log.Fatalf("entry.IntField(status) error: %v", err)
		}
		if status != 500 {
			continue
		}
		timeLocalString, err := entry.Field("time_local")
		if err != nil {
			log.Fatalf("entry.IntField(time_local) error: %v", err)
		}
		timeLocal, err := time.Parse(nginxConfig.TimeFormat, timeLocalString)
		event := &sentry.Event{
			Message:     sentryConfig.Message,
			Level:       sentry.LevelFatal,
			Timestamp:   timeLocal,
			Environment: sentryConfig.Env,
			ServerName:  sentryConfig.ServerName,
			Request: &sentry.Request{
				Headers: map[string]string{},
			},
		}
		if referer := entry.Fields()["http_referer"]; referer != "" {
			event.Request.Headers["Referer"] = referer
		}
		if userAgent := entry.Fields()["http_user_agent"]; userAgent != "" {
			event.Request.Headers["User-Agent"] = userAgent
		}
		if uri := entry.Fields()["request_uri"]; uri != "" {
			event.Request.URL = uri
		}
		if method := entry.Fields()["request_method"]; method != "" {
			event.Request.Method = method
		}
		if ip := entry.Fields()["remote_addr"]; ip != "" {
			event.User.IPAddress = ip
		}
		if xff := entry.Fields()["http_x_forwarded_for"]; xff != "" {
			event.Request.Headers["X-Forwarded-For"] = xff
		}
		sentry.CaptureEvent(event)
	}

	if err := t.Wait(); err != nil {
		log.Fatalf("t.Wait error: %v", err)
	}
}
