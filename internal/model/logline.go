package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/willie68/osmltools/internal/osmlnmea"
)

type LogLine struct {
	Duration         time.Duration
	CorrectTimeStamp time.Time
	Channel          string
	Unknown          string
	NMEAMessage      nmea.Sentence
}

func ParseLogLine(line string) (ll *LogLine, ok bool, err error) {
	sl := strings.SplitN(line, ";", 3)
	if len(sl) < 3 {
		return nil, false, errors.New("to less message parts")
	}
	td, err := parseLoggerTime(sl[0])
	if err != nil {
		return nil, false, fmt.Errorf("invalid time format: %w", err)
	}
	ll = &LogLine{
		Duration: td,
		Channel:  sl[1],
		Unknown:  sl[2],
	}

	ok = true
	msg, err := osmlnmea.ParseNMEA(sl[2])
	var asError *nmea.NotSupportedError
	if errors.As(err, &asError) {
		if osmlnmea.IsNMEASentence(ll.Unknown) {
			return
		}
	}
	ll.NMEAMessage = msg
	if err != nil {
		ok = false
	}
	return
}

func (l *LogLine) String() string {
	var msg string
	if l.NMEAMessage != nil {
		msg = l.NMEAMessage.String()
	} else {
		msg = l.Unknown
	}
	return fmt.Sprintf("%s;%s;%s", formatLoggerDuration(l.Duration), l.Channel, msg)
}

func (l *LogLine) NMEAString() string {
	var msg string
	if l.NMEAMessage != nil {
		msg = l.NMEAMessage.String()
	} else {
		msg = l.Unknown
	}
	return fmt.Sprintf("%s: %s", formatNMEATime(l.CorrectTimeStamp), msg)
}

func parseLoggerTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	secParts := strings.SplitN(parts[2], ".", 2)
	seconds, err := strconv.Atoi(secParts[0])
	if err != nil {
		return 0, err
	}
	milliseconds := 0
	if len(secParts) == 2 {
		milliseconds, err = strconv.Atoi(secParts[1])
		if err != nil {
			return 0, err
		}
	}
	d := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond
	return d, nil
}

func formatLoggerDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func formatNMEATime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
