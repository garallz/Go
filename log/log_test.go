package llz_log

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestLogJustLogInit(t *testing.T) {
	l := LogInit(&LogStruct{
		Cache: false,
	})

	err := l.WriteError("message")
	if err != nil {
		t.Error(err)
	}
}

func TestDirCreate(t *testing.T) {
	l := LogInit(&LogStruct{
		FileTime: TimeHour,
		FilePath: "data",
		Dir:      true,
	})

	l.WriteError("message")
}

func TestFileCreate(t *testing.T) {
	l := LogInit(&LogStruct{
		FileName:   "log.data",
		TimeFormat: "15:04:05",
		FileTime:   TimeMinute,
	})

	for i := 0; i <= 65; i++ {
		l.WriteError(i)
		time.Sleep(time.Second)
	}
}

func TestLogLevel(t *testing.T) {
	l := LogInit(&LogStruct{
		Level: LevelWarn,
	})

	err := l.WriteInfo("message")
	if err != nil {
		t.Error(err)
	}
	err = l.WriteDebug("message")
	if err != nil {
		t.Error(err)
	}
	err = l.WriteWarn("message")
	if err != nil {
		t.Error(err)
	}
	err = l.WriteError("message")
	if err != nil {
		t.Error(err)
	}
}

func TestCacheWrite(t *testing.T) {
	str := []string{
		"Stray birds of summer come to my window to sing and fly away. " +
			"And yellow leaves of autumn, which have no songs, flutter and fall there with a sigh. ",

		"Sorrow is hushed into peace in my heart like the evening among the silent trees.",

		"Listen, my heart, to the whispers of the world with which it makes love to you.",
	}

	l := LogInit(&LogStruct{
		TimeFormat: "15:04:05",
		CacheSize:  1024,
	})

	for _, s := range str {
		fmt.Println(1024/len([]byte(s)) + 1)
	}

	for i := 0; i <= 10; i++ {
		l.WriteError(str[0])
		time.Sleep(time.Millisecond * 500)
	}

	for i := 0; i <= 20; i++ {
		l.WriteError(str[1])
		time.Sleep(time.Millisecond * 500)
	}

	for i := 0; i <= 20; i++ {
		l.WriteError(str[2])
		time.Sleep(time.Millisecond * 500)
	}
}

func TestByteWrite(t *testing.T) {
	var b = struct {
		Time time.Time
		Data string
	}{Time: time.Now(), Data: "Test Data."}

	l := LogInit(&LogStruct{})

	if s, err := json.Marshal(b); err != nil {
		t.Error(err)
	} else {
		l.WriteByte(s)
		// Check \n
		l.WriteByte(s)
	}
}
