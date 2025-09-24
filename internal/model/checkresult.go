package model

import (
	"encoding/json"
	"time"
)

type GeneralResult struct {
	Result   bool     `json:"result"`
	Messages []string `json:"message"`
}

type CheckResult struct {
	Created      time.Time              `json:"created"`
	ErrorCount   int                    `json:"errorCount"`
	WarningCount int                    `json:"warningCount"`
	Files        map[string]*FileResult `json:"files"`
}

type FileResult struct {
	Filename       string    `json:"filename"`
	Origin         string    `json:"origin"`
	Created        time.Time `json:"created"`
	VesselID       int64     `json:"vesselID"`
	DatagramCount  int       `json:"datagramCount"`
	Version        string    `json:"version"`
	FirstTimestamp time.Time `json:"firstTimestamp"`
	LastTimestamt  time.Time `json:"lastTimestamp"`
	ErrorCount     int       `json:"errorCount"`
	Errors         []string  `json:"errors"`
	WarningCount   int       `json:"warningCount"`
	Warnings       []string  `json:"warnings"`
	ErrorA         int       `json:"errorA"`
	ErrorB         int       `json:"errorB"`
	ErrorI         int       `json:"errorI"`
}

func NewGeneralResult() *GeneralResult {
	return &GeneralResult{
		Result:   true,
		Messages: make([]string, 0),
	}
}

func NewCheckResult() *CheckResult {
	return &CheckResult{
		Created: time.Now(),
		Files:   make(map[string]*FileResult),
	}
}

func NewFileResult() *FileResult {
	return &FileResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}
}

func (c *CheckResult) Calc() {
	c.ErrorCount = 0
	c.WarningCount = 0
	for _, ll := range c.Files {
		ll.Calc()
		c.ErrorCount += ll.ErrorCount
		c.WarningCount += ll.WarningCount
	}
}

func (c *CheckResult) String() string {
	return c.JSON()
}

func (c *CheckResult) JSON() string {
	c.Calc()
	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(js)
}

func (c *CheckResult) WithFileResult(fn string, fr *FileResult) *CheckResult {
	c.Files[fn] = fr
	return c
}

func AddWarning(fr *FileResult, msg string) {
	if fr != nil {
		fr.Warnings = append(fr.Warnings, msg)
	}
}

func AddError(channel string, fr *FileResult, msg string) {
	if fr != nil {
		fr.AddErrors(channel, msg)
	}
}

func (f *FileResult) WithFilename(fn string) *FileResult {
	f.Filename = fn
	return f
}

func (f *FileResult) WithVesselID(vid int64) *FileResult {
	f.VesselID = vid
	return f
}

func (f *FileResult) WithCreated(cr time.Time) *FileResult {
	f.Created = cr
	return f
}

func (f *FileResult) WithOrigin(fn string) *FileResult {
	f.Origin = fn
	return f
}

func (f *FileResult) WithErrors(errs []string) *FileResult {
	f.Errors = errs
	f.ErrorCount = len(errs)
	return f
}

func (f *FileResult) AddErrors(channel string, errs ...string) {
	f.Errors = append(f.Errors, errs...)
	f.ErrorCount = len(f.Errors)
	switch channel {
	case "A":
		f.ErrorA += len(errs)
	case "B":
		f.ErrorB += len(errs)
	case "I":
		f.ErrorI += len(errs)
	}
}

func (f *FileResult) WithWarnings(wrns []string) *FileResult {
	f.Warnings = wrns
	f.WarningCount = len(wrns)
	return f
}

func (f *FileResult) Calc() {
	f.ErrorCount = len(f.Errors)
	f.WarningCount = len(f.Warnings)
}

func (f *FileResult) WithVersion(v string) *FileResult {
	f.Version = v
	return f
}

func (f *FileResult) WithDatagramCount(c int) *FileResult {
	f.DatagramCount = c
	return f
}

func (f *FileResult) WithFirstTimestamp(ts time.Time) *FileResult {
	f.FirstTimestamp = ts
	return f
}

func (f *FileResult) WithLastTimestamp(ts time.Time) *FileResult {
	f.LastTimestamt = ts
	return f
}

func (g GeneralResult) JSON() string {
	js, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(js)
}
