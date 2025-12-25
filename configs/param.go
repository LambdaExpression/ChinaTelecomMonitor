package configs

import (
	"China_Telecom_Monitor/models"

	"go.uber.org/zap"
)

var Port string
var Username string
var Password string
var LoginInterval int
var Timeout int64
var Interval int

var DataPath string

var LogLevel string
var LogFormat string

var Dev bool

var PrintVersion bool

var ClientVersion string

var Summary models.Summary

var Logger *zap.SugaredLogger
