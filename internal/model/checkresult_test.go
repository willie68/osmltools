package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	js_basic = "{\n    \"created\": \"1970-01-01T01:00:00+01:00\",\n    \"errorCount\": 0,\n    \"warningCount\": 0,\n    \"files\": {\n        \"test\": {\n            \"filename\": \"testfilename\",\n            \"origin\": \"\",\n            \"created\": \"0001-01-01T00:00:00Z\",\n            \"vesselID\": 0,\n            \"errorCount\": 0,\n            \"errors\": [],\n            \"warningCount\": 0,\n            \"warnings\": []\n        }\n    }\n}"
)

func TestCeckResultBasic(t *testing.T) {
	ast := assert.New(t)
	res := NewCheckResult()

	ast.NotNil(res)

	res.WithFileResult("test", NewFileResult().WithFilename("testfilename"))
	res.Created = time.Unix(0, 0)

	ast.Equal(1, len(res.Files))

	ast.Equal("testfilename", res.Files["test"].Filename)

	js := res.String()

	ast.Equal(js_basic, js)

	ast.Equal(0, res.ErrorCount)
	ast.Equal(0, res.WarningCount)
}

func TestCeckResultJSON(t *testing.T) {
	ast := assert.New(t)

	fs := NewFileResult().
		WithCreated(time.Time{}).
		WithFilename("filename").
		WithOrigin("origin").
		WithVesselID(1234).
		WithErros([]string{"error"}).
		WithWarnings([]string{"warning"})

	ast.Equal(time.Time{}, fs.Created)
	ast.Equal("filename", fs.Filename)
	ast.Equal("origin", fs.Origin)
	ast.Equal(int64(1234), fs.VesselID)
	ast.Equal(1, fs.ErrorCount)
	ast.Equal(1, fs.WarningCount)
	ast.Equal("error", fs.Errors[0])
	ast.Equal("warning", fs.Warnings[0])
}
