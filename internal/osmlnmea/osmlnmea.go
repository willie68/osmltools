package osmlnmea

import (
	"regexp"

	"github.com/adrianmo/go-nmea"
)

const nmeaRegex = `^\$[A-Za-z]{5}([0-9A-Za-z]*,)*\*[0-9A-Fa-f]{2}\r\n$`

// $POSMST,Start NMEA Logger,V 0.1.15*06
type OSMST struct {
	nmea.BaseSentence
	Message string
	Version string
}

// $POSMCFG,255,255,255,255,ffff,65535*73
type OSMCFG struct {
	nmea.BaseSentence
	BaudA      int64
	BaudB      int64
	Seatalk    int64
	Outputs    int64
	VesselID   int64
	BootLoader int64
}

//$POSMGYR,-328,-131,-37*42
//$POSMACC,112,10372,14156*5E

var (
	sp      nmea.SentenceParser
	actual  *nmea.BaseSentence
	nmeaReg *regexp.Regexp
)

func init() {
	sp = nmea.SentenceParser{
		CustomParsers: make(map[string]nmea.ParserFunc),
	}

	sp.CustomParsers["OSMST"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		// This example uses the package builtin parsing helpers
		// you can implement your own parsing logic also
		p := nmea.NewParser(s)
		return OSMST{
			BaseSentence: s,
			Message:      p.String(0, "message"),
			Version:      p.String(1, "version"),
		}, p.Err()
	}

	sp.CustomParsers["OSMCFG"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		// This example uses the package builtin parsing helpers
		// you can implement your own parsing logic also
		p := nmea.NewParser(s)
		return OSMCFG{
			BaseSentence: s,
			BaudA:        p.Int64(0, "bauda"),
			BaudB:        p.Int64(0, "baudb"),
			Seatalk:      p.Int64(0, "seatalk"),
			Outputs:      p.Int64(0, "outputs"),
			VesselID:     p.HexInt64(0, "vesselid"),
			BootLoader:   p.Int64(0, "bootloader"),
		}, p.Err()
	}
	sp.OnBaseSentence = func(sentence *nmea.BaseSentence) error {
		actual = sentence
		return nil
	}

	// Compile the regex
	nmeaReg = regexp.MustCompile(nmeaRegex)

}

func ParseNMEA(line string) (nmea.Sentence, error) {
	actual = nil
	nm, err := sp.Parse(line)
	if err != nil {
		return actual, err
	}
	return nm, nil
}

func IsNMEASentence(sentence string) bool {
	return nmeaReg.MatchString(sentence)
}
