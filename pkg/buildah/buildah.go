package buildah

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/werf/werf/pkg/buildah/types"
)

const (
	DefaultShmSize              = "65536k"
	BuildahImage                = "ghcr.io/werf/buildah:v1.22.3-1"
	BuildahStorageContainerName = "werf-buildah-storage"
)

type CommonOpts struct {
	LogWriter io.Writer
}

type BuildFromDockerfileOpts struct {
	CommonOpts
	ContextTar io.Reader
}

type RunCommandOpts struct {
	CommonOpts
	BuildArgs []string
}

type FromCommandOpts struct {
	CommonOpts
}

type PullOpts struct {
	CommonOpts
}

type Buildah interface {
	BuildFromDockerfile(ctx context.Context, dockerfile []byte, opts BuildFromDockerfileOpts) (string, error)
	RunCommand(ctx context.Context, container string, command []string, opts RunCommandOpts) error
	FromCommand(ctx context.Context, container string, image string, opts FromCommandOpts) error
	Pull(ctx context.Context, ref string, opts PullOpts) error
	Inspect(ctx context.Context, ref string) (types.BuilderInfo, error)
}

type Mode string

const (
	ModeAuto           Mode = "auto"
	ModeNativeRootless Mode = "native-rootless"
	ModeDockerWithFuse Mode = "docker-with-fuse"
)

func InitProcess(initModeFunc func() (Mode, error)) (bool, Mode, error) {
	if v := os.Getenv("_BUILDAH_PROCESS_INIT_MODE"); v != "" {
		mode := Mode(v)
		shouldTerminate, err := doInitProcess(mode)
		return shouldTerminate, mode, err
	}

	mode, err := initModeFunc()
	if err != nil {
		return false, "", fmt.Errorf("unable to init buildah mode: %s", err)
	}
	os.Setenv("_BUILDAH_PROCESS_INIT_MODE", string(mode))

	shouldTerminate, err := doInitProcess(mode)
	return shouldTerminate, mode, err
}

func doInitProcess(mode Mode) (bool, error) {
	switch resolveMode(mode) {
	case ModeNativeRootless:
		return InitNativeRootlessProcess()
	case ModeDockerWithFuse:
		return false, nil
	default:
		return false, fmt.Errorf("unsupported mode %q", mode)
	}
}

func NewBuildah(mode Mode) (b Buildah, err error) {
	switch resolveMode(mode) {
	case ModeNativeRootless:
		switch runtime.GOOS {
		case "linux":
			b, err = NewNativeRootlessBuildah()
			if err != nil {
				return nil, fmt.Errorf("unable to create new Buildah instance with mode %d: %s", mode, err)
			}
		default:
			panic("ModeNativeRootless can't be used on this OS")
		}
	case ModeDockerWithFuse:
		b, err = NewDockerWithFuseBuildah()
		if err != nil {
			return nil, fmt.Errorf("unable to create new Buildah instance with mode %d: %s", mode, err)
		}
	default:
		return nil, fmt.Errorf("unsupported mode %q", mode)
	}

	return b, nil
}

func resolveMode(mode Mode) Mode {
	switch mode {
	case ModeAuto:
		switch runtime.GOOS {
		case "linux":
			return ModeNativeRootless
		default:
			return ModeDockerWithFuse
		}
	default:
		return mode
	}
}

func debug() bool {
	return os.Getenv("WERF_BUILDAH_DEBUG") == "1"
}
