package timer

import (
	"errors"
	"strconv"
	"time"
)

// NewTimer ：make new ticker function
// stamp -> Time timing: 15:04:05; 04:05; 05;
// stamp -> time interval: s-m-h-d:  10s; 30m; 60h; 7d;
// times: 	run times [-1:forever, 0:return not run]
// run:  	defalut: running next time, if true run one times now.
func NewTimer(stamp string, times int, run bool, msg interface{}, function func(interface{})) error {
	if next, interval, err := checkTime(stamp); err != nil {
		return err
	} else {
		if times < 0 {
			times = -1
		} else if times == 0 {
			return errors.New("ticker run times can not be zero")
		}

		if run {
			switch times {
			case 1:
				function(msg)
				return nil
			case -1:
				function(msg)
			default:
				times--
				function(msg)
			}
		}

		putInto(&TimerFunc{
			function: function,
			times:    times,
			next:     next,
			interval: interval,
			msg:      msg,
		})
	}
	return nil
}

// NewRunDuration : Make a new function run
// times: [-1 meas forever], [0 meas not run]
func NewRunDuration(duration time.Duration, times int, msg interface{}, function func(interface{})) {
	if times == 0 {
		return
	} else if times < 0 {
		times = -1
	}

	var data = &TimerFunc{
		next:     time.Now().Add(duration).UnixNano(),
		times:    times,
		interval: int64(duration),
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// NewRunTime : Make a new function run time just one times
func NewRunTime(timestamp time.Time, msg interface{}, function func(interface{})) {
	var data = &TimerFunc{
		next:     timestamp.UnixNano(),
		times:    1,
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// check timerstamp value
func checkTime(stamp string) (int64, int64, error) {
	var err error
	var temp int
	var interval int64
	var now = time.Now()
	var next = now.Unix()

	switch stamp[len(stamp)-1:] {
	case "s", "S":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * SecondTimeUnit)

	case "m", "M":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * MinuteTimeUnit)

	case "h", "H":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * HourTimeUnit)

	case "d", "D":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * DayTimeUnit)

	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Solve the time difference of time zone
		timeString := now.Format(TimeFormatString)
		ts, _ := time.Parse(TimeFormatString, timeString)
		tc := now.Unix() - ts.Unix()

		var t time.Time
		switch len(stamp) {
		case 2: // second
			t, err = time.Parse(TimeFormatString, timeString[:17]+stamp)
			next = t.Unix() + tc
			interval = MinuteTimeUnit

		case 5: // min
			t, err = time.Parse(TimeFormatString, timeString[:14]+stamp)
			next = t.Unix() + tc
			interval = HourTimeUnit

		case 8: // hour
			t, err = time.Parse(TimeFormatString, timeString[:11]+stamp)
			next = t.Unix() + tc
			interval = DayTimeUnit

		default:
			err = errors.New("Can't parst time, please check it")
		}

	default:
		err = errors.New("Can't parst stamp value, please check it")
	}

	if err == nil && next <= now.Unix() {
		next += interval / SecondTimeUnit
	}
	return next * SecondTimeUnit, interval, err
}