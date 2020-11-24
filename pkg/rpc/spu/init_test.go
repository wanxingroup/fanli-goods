package spu

import (
	"os"
	"testing"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/test"
)

func TestMain(t *testing.M) {
	test.Init()

	code := t.Run()

	test.Release()

	os.Exit(code)
}
