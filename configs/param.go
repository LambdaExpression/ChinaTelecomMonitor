package configs

import (
	"China_Telecom_Monitor/models"
	"go.uber.org/zap"
)

var Prot string
var Username string
var Password string
var LoginIntervalTime int
var TimeOut int64
var IntervalsTime int

var DataPath string

var LogLevel string
var LogEncoding string

var Dev bool

var PrintVersion bool

var ClientVersion string

var Summary models.Summary

var Logger *zap.SugaredLogger
