package main

import "flag"

type NginxConfig struct {
	LogFile    string
	Config     string
	ConfigFile string
	TimeFormat string
}

type SentryConfig struct {
	Debug      bool
	Dsn        string
	Env        string
	Message    string
	ServerName string
}

var nginxConfig NginxConfig
var sentryConfig SentryConfig

func init() {
	flag.StringVar(&nginxConfig.LogFile, "file", "/var/log/nginx/access.log", "Nginx access log to follow.")
	flag.StringVar(&nginxConfig.ConfigFile, "config-file", "/var/log/nginx/access.log", "Nginx access log to follow.")
	flag.StringVar(&nginxConfig.Config, "config", "", "Nginx config contents instead of `config-file`.")
	flag.StringVar(&nginxConfig.TimeFormat, "time-format", "02/Jan/2006:15:04:05 -0700", "Nginx log time format.")

	flag.BoolVar(&sentryConfig.Debug, "debug", false, "Debug Sentry.")
	flag.StringVar(&sentryConfig.Dsn, "dsn", "", "Sentry DSN. If not specified, get from SENTRY_DNS env (recommended).")
	flag.StringVar(&sentryConfig.Env, "env", "", "Environment to use in event.")
	flag.StringVar(&sentryConfig.Message, "message", "500", "Issue message.")
	flag.StringVar(&sentryConfig.ServerName, "server-name", "", "Server name to use in event, default to current host name.")
}
