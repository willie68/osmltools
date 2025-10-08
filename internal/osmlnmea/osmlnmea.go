package osmlnmea

import (
	"regexp"
	"strings"

	"github.com/adrianmo/go-nmea"
)

const nmeaRegex = `^(\$[A-Za-z]{5,7}(?:,[0-9A-Za-z\-\:\.\s]*)+\*[0-9A-Fa-f]{2}){1}$`

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

// $POSMGYR,-328,-131,-37*42
type OSMGYR struct {
	nmea.BaseSentence
	XAxis int64
	YAxis int64
	ZAxis int64
}

// $POSMACC,112,10372,14156*5E
type OSMACC struct {
	nmea.BaseSentence
	XAcc int64
	YAcc int64
	ZAcc int64
}

// $POSMVCC,5073,4873*5E
type OSMVCC struct {
	nmea.BaseSentence
	Voltage     int64
	NormVoltage int64
}

// $POSMSO,Reason: times up*4C
type OSMSO struct {
	nmea.BaseSentence
	Message string
}

// $PGRMM,WGS 84*06
type GRMM struct {
	nmea.BaseSentence
	Mapdate string
}

// $PGRMZ,20,f,3*29
type GRMZ struct {
	nmea.BaseSentence
	Altitude int64
	Unit     string
	Fixtype  int64
}

var (
	sp      nmea.SentenceParser
	actual  nmea.Sentence
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

	sp.CustomParsers["OSMSO"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		// This example uses the package builtin parsing helpers
		// you can implement your own parsing logic also
		p := nmea.NewParser(s)
		return OSMSO{
			BaseSentence: s,
			Message:      p.String(0, "message"),
		}, p.Err()
	}

	sp.CustomParsers["OSMCFG"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		// This example uses the package builtin parsing helpers
		// you can implement your own parsing logic also
		p := nmea.NewParser(s)
		return OSMCFG{
			BaseSentence: s,
			BaudA:        p.Int64(0, "bauda"),
			BaudB:        p.Int64(1, "baudb"),
			Seatalk:      p.Int64(2, "seatalk"),
			Outputs:      p.Int64(3, "outputs"),
			VesselID:     p.HexInt64(4, "vesselid"),
			BootLoader:   p.Int64(5, "bootloader"),
		}, p.Err()
	}

	sp.CustomParsers["OSMGYR"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		return OSMGYR{
			BaseSentence: s,
			XAxis:        p.Int64(0, "xaxis"),
			YAxis:        p.Int64(1, "yaxis"),
			ZAxis:        p.Int64(2, "zaxis"),
		}, p.Err()
	}

	sp.CustomParsers["OSMACC"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		return OSMACC{
			BaseSentence: s,
			XAcc:         p.Int64(0, "xacc"),
			YAcc:         p.Int64(1, "yacc"),
			ZAcc:         p.Int64(2, "zacc"),
		}, p.Err()
	}

	sp.CustomParsers["OSMVCC"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		v := p.Int64(0, "voltage")
		nv := int64(0)
		if len(p.Fields) > 1 {
			nv = p.Int64(1, "normvoltage")
		}
		return OSMVCC{
			BaseSentence: s,
			Voltage:      v,
			NormVoltage:  nv,
		}, p.Err()
	}

	sp.CustomParsers["GRMM"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		return GRMM{
			BaseSentence: s,
			Mapdate:      p.String(0, "mapdate"),
		}, p.Err()
	}

	sp.CustomParsers["GRMZ"] = func(s nmea.BaseSentence) (nmea.Sentence, error) {
		p := nmea.NewParser(s)
		return GRMZ{
			BaseSentence: s,
			Altitude:     p.Int64(0, "altitude"),
			Unit:         p.String(1, "unit"),
			Fixtype:      p.Int64(2, "fixtype"),
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
	nm, err := sp.Parse(strings.TrimSpace(line))
	if err != nil {
		return actual, err
	}
	return nm, nil
}

func IsNMEASentence(sentence string) bool {
	return nmeaReg.MatchString(sentence)
}
