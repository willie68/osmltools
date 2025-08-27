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
	sl := strings.SplitAfterN(line, ";", 3)
	if len(sl) < 3 {
		return nil, false, errors.New("to less message parts")
	}
	ll = &LogLine{
		Timestamp: sl[0],
		Channel:   sl[1],
		Unknown:   sl[2],
	}

	ll.NMEAMessage, err = osmlnmea.ParseNMEA(sl[2])
	ok = ll.NMEAMessage != nil
	return
}
