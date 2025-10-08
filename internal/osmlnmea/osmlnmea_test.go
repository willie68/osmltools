package osmlnmea

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/stretchr/testify/suite"
)

type OsmlnmeaSuite struct {
	suite.Suite
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(OsmlnmeaSuite))
}

func (s *OsmlnmeaSuite) SetupTest() {
}

func TestDateTime(t *testing.T) {
	tt := time.Time{}
	js, _ := json.Marshal(tt)
	fmt.Print(string(js))
}

func (s *OsmlnmeaSuite) TestRegistrationCheck() {
	s.NotNil(sp)

	s.Equal(8, len(sp.CustomParsers))

	_, ok := sp.CustomParsers["OSMST"]
	s.True(ok)
	_, ok = sp.CustomParsers["OSMSO"]
	s.True(ok)
	_, ok = sp.CustomParsers["OSMCFG"]
	s.True(ok)
	_, ok = sp.CustomParsers["OSMGYR"]
	s.True(ok)
	_, ok = sp.CustomParsers["OSMACC"]
	s.True(ok)
	_, ok = sp.CustomParsers["OSMVCC"]
	s.True(ok)
	_, ok = sp.CustomParsers["GRMM"]
	s.True(ok)
	_, ok = sp.CustomParsers["GRMZ"]
	s.True(ok)
}

func (s *OsmlnmeaSuite) TestNMEASentenceBasic() {
	line := "$GPGGA,101313,4721.182,N,00832.161,E,1,03,2.3,269.3,M,48.0,M,,*47"
	sen, err := ParseNMEA(line)
	s.NoError(err)
	_, ok := sen.(nmea.GGA)
	s.True(ok)
}

func (s *OsmlnmeaSuite) TestOSMLSentence() {
	var myTests = []struct {
		line     string
		nmeatype any
	}{
		{line: "$POSMST,Start NMEA Logger,V 0.1.15*06", nmeatype: OSMST{}},
		{line: "$POSMCFG,255,255,255,255,ffff,65535*73", nmeatype: OSMCFG{}},
		{line: "$POSMGYR,-340,-107,-78*42", nmeatype: OSMGYR{}},
		{line: "$POSMACC,168,10428,13928*5D", nmeatype: OSMACC{}},
		{line: "$POSMVCC,4940*72", nmeatype: OSMVCC{}},
		{line: "$POSMSO,Reason: times up*4C", nmeatype: OSMSO{}},
	}

	for _, tt := range myTests {
		sen, err := ParseNMEA(tt.line)

		s.NoError(err)
		s.IsType(tt.nmeatype, sen)
		s.Truef(IsNMEASentence(tt.line), "IsNMEASentence: %s", tt.line)
	}
}

func (s *OsmlnmeaSuite) TestVCCSentence1Param() {
	posvcc := "$POSMVCC,4940*72"
	sen, err := ParseNMEA(posvcc)

	s.NoError(err)
	s.IsType(OSMVCC{}, sen)
	vcc := sen.(OSMVCC)
	s.Equal(int64(4940), vcc.Voltage)
	s.Equal(int64(0), vcc.NormVoltage)
}

func (s *OsmlnmeaSuite) TestVCCSentence2Param() {
	posvcc := "$POSMVCC,5073,4873*5E"
	sen, err := ParseNMEA(posvcc)

	s.NoError(err)
	s.IsType(OSMVCC{}, sen)
	vcc := sen.(OSMVCC)
	s.Equal(int64(5073), vcc.Voltage)
	s.Equal(int64(4873), vcc.NormVoltage)
}

func (s *OsmlnmeaSuite) TestUnknownNMEA() {
	var myTests = []struct {
		line     string
		nmeatype any
	}{
		{line: "$PGRMN,WGS 84*05", nmeatype: &nmea.BaseSentence{}},
		{line: "$GPRTE,1,1,c,0*07", nmeatype: &nmea.BaseSentence{}},
	}

	for _, tt := range myTests {
		sen, err := ParseNMEA(tt.line)

		s.NotNil(sen)
		s.Error(err)
		s.IsType(tt.nmeatype, sen)

		s.True(IsNMEASentence(tt.line))
	}
}

func (s *OsmlnmeaSuite) TestErrorline() {
	var myTests = []struct {
		line string
	}{
		{line: "I��b��b��b��b�ºb��b��b��b�ʪb��R��j"},
		{line: "$GPGLL,,,,,101221,*51$GPRMC,101224,V,,,,,,,110916,,*3B"},
	}

	for _, tt := range myTests {
		sen, err := ParseNMEA(tt.line)

		s.Nil(sen)
		s.Error(err)

		s.False(IsNMEASentence(tt.line))
	}
}
