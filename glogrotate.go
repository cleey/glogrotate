package glogrotate

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
)

/**
Clean and Force rotating glog

1. Clean glog:
We will clean up the glog every hour, and keep logs according to your configuration.


2. Force rotating glog：
If you use glog package in Go as your logger, one thing you'll notice is that the only way it rotates is by size. There's MaxSize defined and exported, and glog will rotate the log file when the current file exceeding it (1,800 MiB). There's no way to rotate by time, and there's no way to do manually log rotating.

Actually that's not true. There's a trick. You can see it at : glogRotate.rotate()
*/

var (
	// initailize log fill info
	FILL_CHARACTOR     = "."
	FILL_SIZE_BYTE int = 1024

	DEFAULT_CLEAN_INTERNAL = time.Minute
	EMPTY_DURATION         = time.Duration(0)
	EMPTY_STRING           = ""
)

// define the glogRotate options:
type glogRotate struct {
	dir            string        // Log files will be clean to this directory
	prefix         string        // Log files name prefix
	remain         time.Duration // Log files remain time
	rotateInterval time.Duration // Log files rotate interval default is 1 hour
	cleanInterval  time.Duration // Log files clean interval default is 1 minute
	fillInfo       string        // Log fill info while log rotate
}

// define the glogRotate init options:
type RotateOption struct {
	Dir            string        // Log files will be clean to this directory
	Prefix         string        // Log files name prefix
	Remain         time.Duration // Log files remain time
	RotateInterval time.Duration // Log files rotate interval default is every hour
	CleanInterval  time.Duration // Log files clean interval default is every hour
}

// start a glogRotate with RotateOption
func Start(option RotateOption) *glogRotate {
	r := &glogRotate{
		dir:            option.Dir,
		prefix:         option.Prefix,
		remain:         option.Remain,
		rotateInterval: option.RotateInterval,
		cleanInterval:  option.CleanInterval,
	}

	if EMPTY_STRING == r.dir {
		var logDir, err = filepath.Abs(flag.Lookup("log_dir").Value.String())
		if err != nil {
			glog.Errorln(err)
			panic(err)
		}
		r.dir = logDir
	}

	if EMPTY_STRING == r.prefix {
		program := filepath.Base(os.Args[0])
		r.prefix = program
	}

	if EMPTY_DURATION == r.cleanInterval {
		r.cleanInterval = DEFAULT_CLEAN_INTERNAL
	}

	r.fillInfo = strings.Repeat(FILL_CHARACTOR, FILL_SIZE_BYTE)

	go r.rotater()
	go r.cleaner()
	return r
}

// rotater provides regular rotate function by given log files rotate
func (r *glogRotate) rotater() {
	for {
		if EMPTY_DURATION == r.rotateInterval {
			// default rotate glog per hour
			t := time.Now()
			if t.Minute() == 0 && t.Second() == 0 {
				r.rotate()
			}

			time.Sleep(time.Duration(time.Second))
		} else {
			// rotate glog with user rotateInterval
			time.Sleep(r.rotateInterval)
			r.rotate()
		}
	}
}

// cleaner provides regular clean function by given log files clean
func (r *glogRotate) cleaner() {
	for {
		r.clean()
		time.Sleep(r.cleanInterval)
	}
}

// rotate glog, The idea is simple: we change glog.MaxSize to a very small value, so that the next write will definitely makes it to rotate. After that, we just restore the default size value.
func (r *glogRotate) rotate() {
	prevMaxSize := glog.MaxSize
	glog.MaxSize = uint64(FILL_SIZE_BYTE)
	glog.Info(r.fillInfo)
	glog.Flush()
	glog.MaxSize = prevMaxSize
}

// clean log files in dir
func (r *glogRotate) clean() {
	// 1. ensure dir exist
	if _, err := os.Stat(r.dir); err != nil {
		if !os.IsNotExist(err) {
			glog.Errorln(err)
		}
		return
	}

	// 2. get files from available dir
	files, err := ioutil.ReadDir(r.dir)
	if err != nil {
		glog.Errorln(err)
		return
	}

	// 3. scan files and drop all of the overtime files
	for _, f := range files {
		prefix := strings.HasPrefix(f.Name(), r.prefix)
		str := strings.Split(f.Name(), `.`)
		// drop glog format log files
		if prefix && len(str) >= 7 {
			r.dropIfOvertime(f)
		}
	}
}

// drop file if over the remain time
func (r *glogRotate) dropIfOvertime(f os.FileInfo) {
	if time.Since(f.ModTime()) > r.remain {
		if err := os.Remove(r.dir + "/" + f.Name()); err != nil {
			glog.Errorln(err)
		}
	}
}
