package model

import (
	"errors"
	"strings"

	"github.com/adrianmo/go-nmea"
	"github.com/willie68/osmltools/internal/osmlnmea"
)

type LogLine struct {
	Timestamp   string
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
		Timestamp: sl[0],
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
