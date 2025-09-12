package osmlnmea

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OsmlnmeaSuite struct {
	suite.Suite
	ast *assert.Assertions
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(OsmlnmeaSuite))
}

func (s *OsmlnmeaSuite) SetupTest() {
	s.ast = assert.New(s.T())
}

func TestDateTime(t *testing.T) {
	tt := time.Time{}
	js, _ := json.Marshal(tt)
	fmt.Print(string(js))
}

func (s *OsmlnmeaSuite) TestRegistrationCheck() {
	s.ast.NotNil(sp)

	s.ast.Equal(5, len(sp.CustomParsers))

	_, ok := sp.CustomParsers["OSMST"]
	s.ast.True(ok)
	_, ok = sp.CustomParsers["OSMSO"]
	s.ast.True(ok)
	_, ok = sp.CustomParsers["OSMCFG"]
	s.ast.True(ok)
	_, ok = sp.CustomParsers["OSMGYR"]
	s.ast.True(ok)
	_, ok = sp.CustomParsers["OSMACC"]
	s.ast.True(ok)
}

func (s *OsmlnmeaSuite) TestNMEASentenceBasic() {
	line := "$GPGGA,101313,4721.182,N,00832.161,E,1,03,2.3,269.3,M,48.0,M,,*47"
	sen, err := ParseNMEA(line)
	s.ast.NoError(err)
	_, ok := sen.(nmea.GGA)
	s.ast.True(ok)
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
		{line: "$POSMSO,Reason: times up*4C", nmeatype: OSMSO{}},
	}

	for _, tt := range myTests {
		sen, err := ParseNMEA(tt.line)

		s.ast.NoError(err)
		s.ast.IsType(tt.nmeatype, sen)
		s.ast.Truef(IsNMEASentence(tt.line), "IsNMEASentence: %s", tt.line)
	}
}

func (s *OsmlnmeaSuite) TestUnknownNMEA() {
	var myTests = []struct {
		line     string
		nmeatype any
	}{
		{line: "$PGRMM,WGS 84*06", nmeatype: &nmea.BaseSentence{}},
		{line: "$GPRTE,1,1,c,0*07", nmeatype: &nmea.BaseSentence{}},
	}

	for _, tt := range myTests {
		sen, err := ParseNMEA(tt.line)

		s.ast.NotNil(sen)
		s.ast.Error(err)
		s.ast.IsType(tt.nmeatype, sen)

		s.ast.True(IsNMEASentence(tt.line))
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

		s.ast.Nil(sen)
		s.ast.Error(err)

		s.ast.False(IsNMEASentence(tt.line))
	}
}
