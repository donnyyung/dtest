package dtest

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/datawire/dlib/dlog"
)

func TestMain(m *testing.M) {
	// Skip the test if running in CI and not on linux, because we won't be able to run the container other than on linux
	isCi := os.Getenv("CI")
	if isCi == "true" && runtime.GOOS != "linux" {
		return
	}

	// we get the lock to make sure we are the only thing running
	// because the nat tests interfere with docker functionality
	WithMachineLock(context.TODO(), func(ctx context.Context) {
		os.Exit(m.Run())
	})
}

func TestContainer(t *testing.T) {
	ctx := dlog.NewTestContext(t, false)
	id := dockerUp(ctx, "dtest-test-tag", "nginx")

	running := dockerPs(ctx)
	assert.Contains(t, running, id)

	dockerKill(ctx, id)

	running = dockerPs(ctx)
	assert.NotContains(t, running, id)
}
