package container_runtime

import (
	"context"
	"io"
)

type BuildDockerfileOptions struct {
	ContextTar io.Reader
	Target     string
	BuildArgs  []string // {"key1=value1", "key2=value2", ... }
	AddHost    []string
	Network    string
}

//type StapelBuildOptions struct {
//	ServiceRunCommands []string
//	RunCommands []string
//	Volumes []string
//	VolumesFrom []string
//	Exposes []string
//	Envs map[string]string
//	Labels map[string]string
//}

type ContainerRuntime interface {
	//GetImageInspect(ctx context.Context, ref string) (image.Info, error)
	//Pull(ctx context.Context, ref string) error
	//Tag(ctx, ref, newRef string)
	//Rmi(ctx, ref string)
	//Push(ctx, ref string)

	BuildDockerfile(dockerfile []byte, opts BuildDockerfileOptions) string
	//StapelBuild(opts StapelBuildOptions) string

	String() string

	// Legacy
	RefreshImageObject(ctx context.Context, img Image) error
	PullImageFromRegistry(ctx context.Context, img Image) error
	RenameImage(ctx context.Context, img Image, newImageName string, removeOldName bool) error
	RemoveImage(ctx context.Context, img Image) error
}

// Legacy image handle
type Image interface{}

/*
 * Stapel + docker server
   * container_runtime.Image — конструктор аргументов к docker run + docker tag + docker push + docker commit
     * метод Image.Build и пр.
 * Dockerfile + docker server
   * container_runtime.DockerfileImageBuilder — конструктор аргументов к docker build
      * метод DockerfileImageBuilder.Build
 * DockerServerRuntime
 * Stapel|Dockerfile + docker-server|buildah

type DockerfileImageBuilder struct {
	ContainerRuntime ContainerRuntime
	Dockerfile []byte
	Opts BuildDockerfileOptions

	builtID string
}

func (builder *DockerfileImageBuilder) Build() error {
	builder.builtID = ContainerRuntime.BuildDockerfile(...)
}

func (builder *DockerfileImageBuilder) GetBuiltID() string {
	return builder.builtID
}

func (builder *DockerfileBuidler) Cleanup() error {
}

type StapelImageBuilder struct {
	Opts StapelBuildOptions
	...
}

func (builder *StapelImageBuilder) Build() error {
	builder.builtID = ContainerRuntime.StapelBuild(...)
}

func (builder *StapelImageBuilder) GetBuiltID() string {
	return builder.builtID
}

func (builder *StapelImageBuilder) Cleanup() error {
}

*/
