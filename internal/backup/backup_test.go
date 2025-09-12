package backup

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	ast := assert.New(t)

	inj := do.New()
	Init(inj)

	bck := do.MustInvoke[Backup](inj)
	ast.NotNil(bck)

	err := bck.Backup("../../testdata/sdCard", "../../testdata/bck/")
	ast.NoError(err)
}

func TestRestore(t *testing.T) {
	ast := assert.New(t)

	inj := do.New()
	Init(inj)

	bck := do.MustInvoke[Backup](inj)
	ast.NotNil(bck)

	err := os.MkdirAll("../../testdata/rst", os.ModePerm)
	ast.NoError(err)

	err = bck.Restore("../../testdata/bck/bck_20250912223031.zip", "../../testdata/rst")
	ast.NoError(err)
}
