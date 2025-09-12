package backup

import (
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
