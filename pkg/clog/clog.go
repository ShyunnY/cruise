package clog

import (
	"go.uber.org/zap"
	"log"
)

var CL *zap.Logger

func SetLogger() error {
	var err error
	CL, err = zap.NewProduction()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
