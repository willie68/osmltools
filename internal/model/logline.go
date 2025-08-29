package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/willie68/osmltools/internal/osmlnmea"
)

const (
	fmtLoggerTime = "15:04:05.000"
)

type LogLine struct {
	Timestamp   time.Time
	Channel     string
	Unknown     string
	NMEAMessage nmea.Sentence
}

func ParseLogLine(line string) (ll *LogLine, ok bool, err error) {
	sl := strings.SplitN(line, ";", 3)
	if len(sl) < 3 {
		return nil, false, errors.New("to less message parts")
	}
	ll = &LogLine{
		Timestamp: convertTime(sl[0]),
		Channel:   sl[1],
		Unknown:   sl[2],
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

func convertTime(st string) time.Time {
	t, err := time.Parse(fmtLoggerTime, st)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (l *LogLine) String() string {
	var msg string
	if l.NMEAMessage != nil {
		msg = l.NMEAMessage.String()
	} else {
		msg = l.Unknown
	}
	return fmt.Sprintf("%s;%s;%s", l.Timestamp.Format(fmtLoggerTime), l.Channel, msg)
}
