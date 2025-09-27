package logger

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/willie68/gowillie68/pkg/fileutils"
)

type LoggerConfig struct {
	Seatalk  bool  `json:"seatalk"`
	BaudA    int16 `json:"baudA"`
	BaudB    int16 `json:"baudB"`
	Gyro     bool  `json:"gyro"`
	Supply   bool  `json:"supply"`
	VesselID int16 `json:"vesselID"`
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Seatalk:  false,
		BaudA:    4800,
		BaudB:    4800,
		Gyro:     true,
		Supply:   false,
		VesselID: 0,
	}
}

// WithVesselID methods to set the vesselid in a fluent style
func (c *LoggerConfig) WithVesselID(vid int16) *LoggerConfig {
	c.VesselID = vid
	return c
}

// WithSeatalk methods to set the seatalk flag for channel A
func (c *LoggerConfig) WithSeatalk(st bool) *LoggerConfig {
	c.Seatalk = st
	return c
}

// WithBaudA methods to set the baud rate for channel A
func (c *LoggerConfig) WithBaudA(baud int16) *LoggerConfig {
	c.BaudA = baud
	return c
}

// WithBaudB methods to set the baud rate for channel B
func (c *LoggerConfig) WithBaudB(baud int16) *LoggerConfig {
	c.BaudB = baud
	return c
}

// WithGyro methods to set the gyro flag, the logger then writes gyro data to the data files
func (c *LoggerConfig) WithGyro(gyro bool) *LoggerConfig {

	c.Gyro = gyro
	return c
}

// WithSupply methods to set the supply flag, the logger then writes supply data to the data files
func (c *LoggerConfig) WithSupply(supply bool) *LoggerConfig {
	c.Supply = supply
	return c
}

func (c *LoggerConfig) Write(w io.Writer) error {
	template := "%s\r\n"
	if c.Seatalk {
		template = "s%s\r\n"
	}
	fmt.Fprintf(w, template, ConvertBaudToCode(c.BaudA))

	fmt.Fprintf(w, "%s\r\n", ConvertBaudToCode(c.BaudB))

	outputs := 0
	if c.Supply {
		outputs += 1
	}
	if c.Gyro {
		outputs += 2
	}
	fmt.Fprintf(w, "%d\r\n", outputs)

	if c.VesselID > 0 {
		fmt.Fprintf(w, "%.8x\r\n", c.VesselID)
	}

	return nil
}

func Read(r io.Reader) (*LoggerConfig, error) {
	var cfg LoggerConfig
	// Scanner erstellen
	scanner := bufio.NewScanner(r)
	count := 0
	// Zeilenweise lesen
	for scanner.Scan() {
		line := scanner.Text()
		switch count {
		case 0:
			if strings.HasPrefix(line, "s") {
				cfg.Seatalk = true
				line = strings.TrimPrefix(line, "s")
			}
			cfg.BaudA = ConvertCodeToBaud(line)
		case 1:
			cfg.BaudB = ConvertCodeToBaud(line)
		case 2:
			outputs, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			cfg.Gyro = (outputs & 2) != 0
			cfg.Supply = (outputs & 1) != 0
		case 3:
			vid, err := strconv.ParseUint(line, 16, 16)
			if err != nil {
				return nil, err
			}
			cfg.VesselID = int16(vid)
		}
		count++
		if count > 3 {
			break
		}
	}

	// Fehler beim Scannen prüfen
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *LoggerConfig) JSON() (string, error) {
	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func (c *LoggerConfig) Validate() error {
	if c.BaudA != 0 && c.BaudA != 1200 && c.BaudA != 2400 && c.BaudA != 4800 && c.BaudA != 9600 && c.BaudA != 19200 {
		return errors.New("invalid baud rate for channel A")
	}
	if c.BaudB != 0 && c.BaudB != 1200 && c.BaudB != 2400 && c.BaudB != 4800 {
		return errors.New("invalid baud rate for channel B")
	}
	if c.VesselID < 0 || c.VesselID > 9999 {
		return errors.New("invalid vessel ID")
	}
	return nil
}

func (c *LoggerConfig) WriteToSDCard(sdCardFolder string) error {
	cfgFile := filepath.Join(sdCardFolder, "config.dat")
	err := c.backupCfg(sdCardFolder)
	if err != nil {
		return err
	}
	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = c.Write(f)
	if err != nil {
		return err
	}
	return nil
}

func ReadFromSDCard(sdCardFolder string) (*LoggerConfig, error) {
	cfgFile := filepath.Join(sdCardFolder, "config.dat")
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg, err := Read(f)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *LoggerConfig) String() string {
	return fmt.Sprintf("Seatalk: %t, BaudA: %d, BaudB: %d, Gyro: %t, Supply: %t, VesselID: %d",
		c.Seatalk, c.BaudA, c.BaudB, c.Gyro, c.Supply, c.VesselID)
}

func (c *LoggerConfig) backupCfg(sdCardFolder string) error {
	oldFile := filepath.Join(sdCardFolder, "config.old")
	newFile := filepath.Join(sdCardFolder, "config.dat")
	if fileutils.FileExists(oldFile) {
		err := os.Remove(oldFile)
		if err != nil {
			return err
		}
	}
	return os.Rename(newFile, oldFile)
}

func ConvertBaudToCode(baud int16) string {
	switch baud {
	case 0:
		return "0"
	case 1200:
		return "1"
	case 2400:
		return "2"
	case 4800:
		return "3"
	case 9600:
		return "4"
	case 19200:
		return "5"
	default:
		return "0"
	}
}

func ConvertCodeToBaud(code string) int16 {
	switch code {
	case "0":
		return 0
	case "1":
		return 1200
	case "2":
		return 2400
	case "3":
		return 4800
	case "4":
		return 9600
	case "5":
		return 19200
	default:
		return 0
	}
}
