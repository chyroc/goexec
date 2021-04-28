package goexec_test

import (
	"testing"

	"github.com/chyroc/goexec"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	as := assert.New(t)

	{
		stdout, _, err := goexec.New("env").SetEnv("MAGIC_UUID", "x").Run()
		as.Nil(err)
		as.Contains(stdout, "MAGIC_UUID=x")
	}

	{
		stdout, _, err := goexec.New("ls").Run()
		as.Nil(err)
		as.Contains(stdout, "README.md")
	}
}
