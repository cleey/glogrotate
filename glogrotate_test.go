package glogrotate

import (
	"testing"
	"time"

	"github.com/golang/glog"
)

func TestGlogrotate(t *testing.T) {
	// start rotate log file
	go func() {
		Start(RotateOption{
			RotateInterval: time.Duration(time.Second * 1),
			CleanInterval:  time.Duration(time.Second * 1),
			Remain:         time.Duration(time.Second * 2),
		})
	}()

	// generate log file
	var times = 3
	for {
		times = times - 1
		if times < 0 {
			break
		}
		glog.Info("glogrotate test log")
		time.Sleep(time.Second)
	}

	// wait for rotate file end
	time.Sleep(time.Second * 3)
}
