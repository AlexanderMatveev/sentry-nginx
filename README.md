# sentry-nginx

[![Workflow](https://github.com/AlexanderMatveev/sentry-nginx/actions/workflows/go.yml/badge.svg)](https://github.com/AlexanderMatveev/sentry-nginx/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlexanderMatveev/sentry-nginx)](https://goreportcard.com/report/github.com/AlexanderMatveev/sentry-nginx)

Push Sentry events from nginx

## Usage

```sh
go run github.com/AlexanderMatveev/sentry-nginx --help
```

### Available options

```
  -dsn string
    	Sentry DSN. If not specified, get from SENTRY_DNS env (recommended).
  -config-file string
    	Nginx access log to follow. (default "/var/log/nginx/access.log")
  -file string
    	Nginx access log to follow. (default "/var/log/nginx/access.log")
  -time-format string
    	Nginx log time format. (default "02/Jan/2006:15:04:05 -0700")

  -config config-file
    	Nginx config contents instead of config-file.
  -debug
    	Debug Sentry
  -env string
    	Environment to use in event.
  -message string
    	Issue message. (default "500")
  -server-name string
    	Server name to use in event, default to current host name.
 
```
