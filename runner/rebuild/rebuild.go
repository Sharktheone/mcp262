package rebuild

import (
	"context"
	"os/exec"
	"sync/atomic"
)

const RELEASE_BUILD_THRESHOLD uint32 = 5000

type EngineLocation struct {
	ReleasePath string
	DebugPath   string

	UseDebug atomic.Bool
}

func (engine *EngineLocation) GetPath() string {
	if engine.UseDebug.Load() {
		return engine.DebugPath
	}
	return engine.ReleasePath
}

func RebuildEngine(repoRoot string, numTests uint32, rebuild bool) (*EngineLocation, context.CancelFunc, error) {
	if rebuild {
		debugErr := rebuildDebugEngine(repoRoot)
		if debugErr != nil {
			return nil, nil, debugErr
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	engine := &EngineLocation{
		ReleasePath: repoRoot + "/target/release/yavashark_test262",
		DebugPath:   repoRoot + "/target/debug/yavashark_test262",
	}

	if numTests > RELEASE_BUILD_THRESHOLD && rebuild {
		go func() {
			releaseErr := rebuildReleaseEngine(repoRoot, ctx)
			if releaseErr != nil {
				cancel()
				return
			}

			engine.UseDebug.Store(true)
		}()

	}

	engine.UseDebug.Store(false)

	return engine, cancel, nil

}

func rebuildDebugEngine(repoRoot string) error {
	cmd := exec.Command("cargo", "build")

	cmd.Dir = repoRoot

	return cmd.Run()
}

func rebuildReleaseEngine(repoRoot string, ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "cargo", "build", "--release")

	cmd.Dir = repoRoot

	return cmd.Run()
}
