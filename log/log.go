package log

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func NewNopLogger() logrus.FieldLogger {
	l := logrus.New()
	l.Out = ioutil.Discard

	return l
}
