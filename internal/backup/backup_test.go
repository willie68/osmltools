package backup

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

const (
	testZIP = "../../testdata/bck/bck_20250913160522.zip"
)

func TestBackup(t *testing.T) {
	ast := assert.New(t)

	inj := do.New()
	Init(inj)

	bck := do.MustInvoke[Backup](inj)
	ast.NotNil(bck)

	filename, err := bck.Backup("../../testdata/sdCard", "../../testdata/bck/")
	ast.NoError(err)
	ast.NotEmpty(filename)
}

func TestRestore(t *testing.T) {
	ast := assert.New(t)

	inj := do.New()
	Init(inj)

	bck := do.MustInvoke[Backup](inj)
	ast.NotNil(bck)

	err := os.MkdirAll("../../testdata/rst", os.ModePerm)
	ast.NoError(err)

	filename, err := bck.Restore(testZIP, "../../testdata/rst")
	ast.NoError(err)
	ast.Equal(testZIP, filename)
}

func TestRestore1(t *testing.T) {
	ast := assert.New(t)

	inj := do.New()
	Init(inj)

	bck := do.MustInvoke[Backup](inj)
	ast.NotNil(bck)

	err := os.MkdirAll("../../testdata/rst", os.ModePerm)
	ast.NoError(err)
	zip := "../../testdata/bck/bck_20250913160205.zip"

	filename, err := bck.Restore(zip, "../../testdata/rst1")
	ast.NoError(err)
	ast.Equal(zip, filename)
}
