package llz_log

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// check struct data and supplement.
func (l *LogStruct) checkStruct() {
	if l == nil {
		l.Cache = true
		l.CacheSize = 1024 * 1024
		l.FileName = "log"
		l.FileTime = TimeDay
		l.Level = LevelError
		l.TimeFormat = "2006-01-02 15:04:05"
	} else {
		if l.FileName == "" {
			l.FileName = "log"
		}
		if l.FileTime == "" {
			l.FileTime = TimeDay
		}
		if l.Level == 0 {
			l.Level = LevelError
		}
		if l.TimeFormat == "" {
			l.TimeFormat = "2006-01-02 15:04:05"
		}
		if l.Cache && l.CacheSize == 0 {
			l.CacheSize = 1024 * 1024
		}
		if l.CacheSize != 0 && !l.Cache {
			l.Cache = true
		}
		if l.FilePath != "" {
			path := l.FilePath[len(l.FilePath)-1:]
			if path != "/" || path != `\` {
				if runtime.GOOS == "windows" {
					l.FilePath += `\`
				} else {
					l.FilePath += "/"
				}
			}
		}
	}
}

// open file and put in struct with *os.file
func (l *LogStruct) open() {
	var err error
	name := l.FilePath + l.FileName + "." + time.Now().Format(fmt.Sprint(l.FileTime))

	if l.Dir {
		d := l.dir()
		name = l.FilePath + d + l.FileName + "." + time.Now().Format(fmt.Sprint(l.FileTime))
	}

	l.file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("Open log file error: " + err.Error())
	}
}

// sleep time to make new file open.
func (l *LogStruct) upFile() {
	var last string
	var format = fmt.Sprint(l.FileTime)

	switch l.FileTime {
	case TimeMonth:
		last = time.Now().UTC().AddDate(0, 1, 0).Format(format)
	case TimeDay:
		last = time.Now().UTC().Add(time.Hour * 24).Format(format)
	case TimeHour:
		last = time.Now().UTC().Add(time.Hour * 1).Format(format)
	case TimeMinute:
		last = time.Now().UTC().Add(time.Minute).Format(format)
	}

	if stamp, err := time.Parse(format, last); err != nil {
		panic("Time parse error: " + err.Error())
	} else {
		l.stamp = stamp.UTC().Unix()
		if sleep := stamp.Sub(time.Now().UTC()).Seconds(); sleep > 5 {
			time.Sleep(time.Second * time.Duration(sleep-5))
		}
		l.tc = true
	}
}

// put log data and level in buffer.
func (l *LogStruct) put(level string, msg ...interface{}) error {
	d := make([]string, len(msg)+1)
	d[0] = time.Now().Format(l.TimeFormat) + level
	for i, r := range msg {
		d[i+1] = fmt.Sprint(r)
	}

	f := []byte(strings.Join(d, " ") + "\n")

	return l.putByte(f)
}

// put byte in cache or file.
func (l *LogStruct) putByte(bts []byte) error {
	var err error
	if l.tc {
		if err = l.check(); err != nil {
			return err
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Cache {
		l.buf = append(l.buf, bts...)
		if len(l.buf) >= l.CacheSize {
			_, err = l.file.Write(l.buf)
			l.buf = l.buf[:0]
		}
	} else {
		_, err = l.file.Write(bts)
	}
	return err
}

// check new file open.
func (l *LogStruct) check() error {
	if l.stamp <= time.Now().UTC().Unix() {
		l.mu.Lock()
		if l.Cache {
			if _, err := l.file.Write(l.buf); err != nil {
				return err
			}
			l.buf = l.buf[:0]
		}
		l.file.Close()

		l.open()
		l.tc = false
		l.mu.Unlock()

		go l.upFile()
	}
	return nil
}

// make dir about FileTime.
func (l *LogStruct) dir() string {
	// Create log file dir with year.
	l.create(time.Now().Format("2006"))

	// Create log file dir with month.
	if l.FileTime != TimeMonth {
		l.create(time.Now().Format("2006/01"))
	} else {
		return time.Now().Format("2006/")
	}
	// Create log file dir with day.
	if l.FileTime != TimeDay {
		l.create(time.Now().Format("2006/01/02"))
	} else {
		return time.Now().Format("2006/01/")
	}
	// Create log file dir with hour.
	if l.FileTime != TimeHour {
		l.create(time.Now().Format("2006/01/02/15"))
	} else {
		return time.Now().Format("2006/01/02/")
	}
	return time.Now().Format("2006/01/02/15/")
}

// check dir exist and create.
func (l *LogStruct) create(path string) {
	if _, err := os.Stat(l.FilePath + path); err != nil {
		if os.IsNotExist(err) {
			if exec.Command("sh", "-c", "mkdir "+l.FilePath+path).Run() != nil {
				panic("Create log file dir error!")
			}
		}
	}
}
