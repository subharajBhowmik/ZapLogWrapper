package model

import (
	"github.com/subharajBhowmik/ZapLogWrapper/utils"
	"go.uber.org/zap"
)

type ZapLogError struct {
	ErrSrc     string
	ErrCode    string
	ErrSubCode string
	ErrCause   error
	Message    string
	ZapFields  []zap.Field
}

func BuildNewZapLogError(errSrc string, errCode string, errSubCode string, errCause error, additionalFields ...interface{}) *ZapLogError {
	zapLogError := ZapLogError{
		ErrSrc:     errSrc,
		ErrCode:    errCode,
		ErrSubCode: errSubCode,
		ErrCause:   errCause,
	}
	zapLogError.BuildZapLogMessageAndFields(additionalFields)
	return &zapLogError
}

func (l *ZapLogError) Error() string {
	return l.Message
}

func (l *ZapLogError) BuildZapLogMessageAndFields(addedFields ...interface{}) {
	l.Message, l.ZapFields = utils.GenerateZapLog(l.ErrSrc, l.ErrCode, l.ErrSubCode, l.ErrCause, addedFields...)
}

func (l *ZapLogError) WithError(alternateError error) *ZapLogError {
	l.ErrCause = alternateError
	return l
}
