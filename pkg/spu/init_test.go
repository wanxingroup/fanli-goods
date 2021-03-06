package spu

import (
	"os"
	"testing"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/test"
)

func TestMain(t *testing.M) {
	test.Init()

	code := t.Run()

	if code == 0 {
		test.Release()
	}

	os.Exit(code)
}
