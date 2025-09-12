package logger

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerConfig(t *testing.T) {
	ast := assert.New(t)
	cfg := NewLoggerConfig().
		WithBaudA(4800).
		WithBaudB(9600).
		WithGyro(true).
		WithSupply(true).
		WithSeatalk(true).
		WithVesselID(int16(1234))

	ast.Equal(int16(4800), cfg.BaudA)
	ast.Equal(int16(9600), cfg.BaudB)
	ast.True(cfg.Gyro)
	ast.True(cfg.Supply)
	ast.True(cfg.Seatalk)
	ast.Equal(int16(1234), cfg.VesselID)
}

func TestLoggerConfigDefaults(t *testing.T) {
	ast := assert.New(t)
	cfg := NewLoggerConfig()

	ast.Equal(int16(4800), cfg.BaudA)
	ast.Equal(int16(4800), cfg.BaudB)
	ast.True(cfg.Gyro)
	ast.False(cfg.Supply)
	ast.False(cfg.Seatalk)
	ast.Equal(int16(0), cfg.VesselID)
}

func TestLoggerConfigChaining(t *testing.T) {
	ast := assert.New(t)
	cfg := NewLoggerConfig().
		WithBaudA(4800).
		WithBaudB(19200).
		WithGyro(false).
		WithSupply(false).
		WithSeatalk(false).
		WithVesselID(int16(5678))

	ast.Equal(int16(4800), cfg.BaudA)
	ast.Equal(int16(19200), cfg.BaudB)
	ast.False(cfg.Gyro)
	ast.False(cfg.Supply)
	ast.False(cfg.Seatalk)
	ast.Equal(int16(5678), cfg.VesselID)
}

func TestWrite(t *testing.T) {
	ast := assert.New(t)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	cfg := NewLoggerConfig().
		WithBaudA(4800).
		WithBaudB(19200).
		WithGyro(true).
		WithSupply(true).
		WithSeatalk(true).
		WithVesselID(int16(5678))

	cfg.Write(writer)
	writer.Flush()

	ast.Equal("s3\r\n5\r\n3\r\n0000162e\r\n", buf.String())

}

func TestJson(t *testing.T) {
	ast := assert.New(t)

	cfg := NewLoggerConfig().
		WithBaudA(4800).
		WithBaudB(19200).
		WithGyro(true).
		WithSupply(true).
		WithSeatalk(true).
		WithVesselID(int16(5678))

	js, err := cfg.JSON()
	ast.NoError(err)
	ast.Equal("{\n    \"seatalk\": true,\n    \"baudA\": 4800,\n    \"baudB\": 19200,\n    \"gyro\": true,\n    \"supply\": true,\n    \"vesselID\": 5678\n}", js)
}

func TestValidate(t *testing.T) {
	ast := assert.New(t)
	cfg := NewLoggerConfig().
		WithBaudA(19200).
		WithBaudB(4800).
		WithGyro(true).
		WithSupply(true).
		WithSeatalk(true).
		WithVesselID(int16(5678))
	err := cfg.Validate()
	ast.NoError(err)
}

func TestValidateFail(t *testing.T) {
	ast := assert.New(t)
	cfg := NewLoggerConfig().
		WithBaudA(1234).
		WithBaudB(4800).
		WithGyro(true).
		WithSupply(true).
		WithSeatalk(true).
		WithVesselID(int16(5678))
	err := cfg.Validate()
	ast.Error(err)
}

func TestRead(t *testing.T) {
	ast := assert.New(t)
	input := "s5\r\n3\r\n3\r\n0000162e\r\n"
	reader := bufio.NewReader(bytes.NewBufferString(input))
	cfg, err := Read(reader)
	ast.NoError(err)
	ast.Equal(int16(19200), cfg.BaudA)
	ast.Equal(int16(4800), cfg.BaudB)
	ast.True(cfg.Seatalk)
	ast.True(cfg.Gyro)
	ast.True(cfg.Supply)
	ast.Equal(int16(5678), cfg.VesselID)
}

func TestReadFail(t *testing.T) {
	ast := assert.New(t)
	input := "s5\r\n3\r\n3\r\n0000162x\r\n" // invalid vessel id
	reader := bufio.NewReader(bytes.NewBufferString(input))
	cfg, err := Read(reader)
	ast.NotNil(err)
	ast.Nil(cfg)
}
