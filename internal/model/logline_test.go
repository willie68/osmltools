package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoglineSuite struct {
	suite.Suite
	ast *assert.Assertions
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(LoglineSuite))
}

func (s *LoglineSuite) SetupTest() {
	s.ast = assert.New(s.T())
}

func (s *LoglineSuite) TestLogline() {
	var myTests = []struct {
		line string
		ok   bool
		err  bool
	}{
		{line: "01:45:59.695;B;$PGRMM,WGS 84*06", ok: true, err: true},
		{line: "02:00:02.540;B;$GPGGA,121133,4721.463,N16,000.2,E*7D", ok: false, err: true},
		{line: "02:00:02.540;I;$POSMST,Start NMEA Logger,V 0.1.15*06", ok: true, err: false},
		{line: "00:59:59.405;B;$PGRME,6.3,M,,M,6.3,M*00", ok: true, err: false},
		{line: "01:45:57.794;B;$GPRTE,1,1,c,0*07", ok: false, err: true},
	}

	for _, tt := range myTests {
		_, ok, err := ParseLogLine(tt.line)
		if tt.err {
			s.ast.Error(err)
		} else {
			s.ast.NoError(err)
		}
		s.ast.Equal(ok, tt.ok)
	}
}
