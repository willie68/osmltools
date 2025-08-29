package model

import (
	"encoding/json"
	"time"
)

type CheckResult struct {
	Created      time.Time              `json:"created"`
	ErrorCount   int                    `json:"errorCount"`
	WarningCount int                    `json:"warningCount"`
	Files        map[string]*FileResult `json:"files"`
}

type FileResult struct {
	Filename     string    `json:"filename"`
	Origin       string    `json:"origin"`
	Created      time.Time `json:"created"`
	VesselID     int64     `json:"vesselID"`
	ErrorCount   int       `json:"errorCount"`
	Errors       []string  `json:"errors"`
	WarningCount int       `json:"warningCount"`
	Warnings     []string  `json:"warnings"`
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

func (c *CheckResult) String() string {
	c.ErrorCount = 0
	c.WarningCount = 0
	for _, ll := range c.Files {
		ll.Calc()
		c.ErrorCount += ll.ErrorCount
		c.WarningCount += ll.WarningCount
	}
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

func AddError(fr *FileResult, msg string) {
	if fr != nil {
		fr.Errors = append(fr.Errors, msg)
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

func (f *FileResult) WithErros(errs []string) *FileResult {
	f.Errors = errs
	f.ErrorCount = len(errs)
	return f
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
