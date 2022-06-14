package configs

import (
	"China_Telecom_Monitor/models"
	"github.com/LambdaExpression/surf/browser"
	"go.uber.org/zap"
)

var Prot string
var DockerProt string
var Username string
var Password string
var LoginIntervalTime int
var TimeOut int64
var IntervalsTime int
var DockerWaitTime int

var DataPath string

var LogLevel string
var LogEncoding string

var Dev bool

var PrintVersion bool

var Browser *browser.Browser
var Summary models.Summary

var Logger *zap.SugaredLogger
