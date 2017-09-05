package main

import "github.com/sirupsen/logrus"

func logger() logrus.FieldLogger {
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}}

	l.SetLevel(logrus.DebugLevel)

	return l.WithFields(nil)
}
