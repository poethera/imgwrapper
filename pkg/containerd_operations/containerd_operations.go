package containerd_operations

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"imgwrapper/pkg/api/types"
	"os"
	"path/filepath"

	"github.com/containerd/containerd"
	ndtypes "github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/buildkitutil"
	"github.com/containerd/nerdctl/pkg/clientutil"
	"github.com/containerd/nerdctl/pkg/cmd/builder"
	"github.com/containerd/nerdctl/pkg/cmd/container"
	"github.com/containerd/nerdctl/pkg/cmd/login"
)

type ContainerdOperations struct {
	Engine string
}

func Create(ctx context.Context, imgwOpts *types.ImgWOptions) (*types.ImgWContext, error) {
	client, ctx, cancel, err := clientutil.NewClient(ctx, imgwOpts.Namespace, imgwOpts.Address)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	imgwCtx := &(types.ImgWContext{
		Client:  client,
		Ctx:     ctx,
		Cancel:  cancel,
		ImgWOps: &ContainerdOperations{imgwOpts.Engine},
	})

	return imgwCtx, err
}

func (c *ContainerdOperations) Registry_Login(imgwCtx *types.ImgWContext, imgwOpts *types.ImgWOptions, serverAddr string, id string, passwd string) error {
	Stdout := new(bytes.Buffer)
	defer func() {
		Stdout = nil
	}()

	options := getLoginOptions(getGlobalOptions(imgwOpts), serverAddr, id, passwd)

	err := login.Login(imgwCtx.Ctx, options, Stdout)
	//TODO: error
	if err != nil {
		fmt.Println("login error", err.Error())
	}

	//TODO: Stdout
	fmt.Println("login output: ", Stdout.String())

	return err
}

func (c *ContainerdOperations) Registry_Logout(imgCtx *types.ImgWContext) error {
	//TODO:
	return nil
}

// image_pull(imgctx *ImgWContext) error
func (c *ContainerdOperations) Image_Push(imgwCtx *types.ImgWContext) error {
	//TODO:
	return nil
}

func (c *ContainerdOperations) Image_Build(imgwCtx *types.ImgWContext, imgwOpts *types.ImgWOptions, buildOpts *types.BuildOperationOptions) error {
	client, ok := imgwCtx.Client.(*containerd.Client)
	if !ok {
		//TODO: error
		fmt.Println("invalid client object")
		return errors.New("invalid client object")
	}

	//preprocess - buildctx
	if dockerfile_folder, err := os.MkdirTemp("", "new_folder_"); err == nil {
		if err := os.WriteFile(filepath.Join(dockerfile_folder, "Dockerfile"), []byte(buildOpts.BuildCtx), 0644); err != nil {
			defer os.RemoveAll(dockerfile_folder)
			return err
		}
		defer os.RemoveAll(dockerfile_folder)
		buildOpts.BuildCtx = dockerfile_folder
	} else {
		//TODO: preprocess error
		fmt.Println("preprocess error", err.Error())
		return err
	}

	options, err := getBuilderBuildOption(
		getGlobalOptions(imgwOpts),
		buildOpts,
	)
	if err != nil {
		//TODO: buildkit error
		fmt.Println("buildkit error", err.Error())
		return err
	}

	if err := builder.Build(imgwCtx.Ctx, client, options); err != nil {
		//TODO:
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (c *ContainerdOperations) Image_Commit_For_Container(imgwCtx *types.ImgWContext, imgwOpts *types.ImgWOptions, commitOpts *types.CommitOperationOptions) error {

	client, ok := imgwCtx.Client.(*containerd.Client)
	if !ok {
		//TODO: error
		fmt.Println("invalid client object")
		return errors.New("invalid client object")
	}

	options := getContainerCommitOptions(
		getGlobalOptions(imgwOpts),
		commitOpts,
	)
	options.Stdout = new(bytes.Buffer)
	defer func() {
		options.Stdout = nil
	}()

	if err := container.Commit(imgwCtx.Ctx, client, commitOpts.Rawref, commitOpts.Req, options); err != nil {
		//TODO:
		fmt.Println(err.Error())
	}
	//bufImgId := options.Stdout.(*bytes.Buffer)
	//return err, bufImgId.String()	//wrong imgid .. bug ? or for me no idea ?	//TODO: delete
	return nil
}

func getGlobalOptions(imgwOpts *types.ImgWOptions) ndtypes.GlobalCommandOptions {
	debug := false
	debugFull := false
	address := imgwOpts.Address
	namespace := imgwOpts.Namespace
	snapshotter := ""
	cniPath := ""
	cniConfigPath := ""
	dataRoot := ""
	cgroupManager := ""
	insecureRegistry := false
	hostsDir := []string{}
	experimental := false

	return ndtypes.GlobalCommandOptions{
		Debug:            debug,
		DebugFull:        debugFull,
		Address:          address,
		Namespace:        namespace,
		Snapshotter:      snapshotter,
		CNIPath:          cniPath,
		CNINetConfPath:   cniConfigPath,
		DataRoot:         dataRoot,
		CgroupManager:    cgroupManager,
		InsecureRegistry: insecureRegistry,
		HostsDir:         hostsDir,
		Experimental:     experimental,
	}
}

func getContainerCommitOptions(ndGOpts ndtypes.GlobalCommandOptions, commitOpts *types.CommitOperationOptions) ndtypes.ContainerCommitOptions {
	author := commitOpts.Author
	message := commitOpts.Message
	pause := false
	change := commitOpts.Change

	return ndtypes.ContainerCommitOptions{
		Stdout:   os.Stdout, //new(bytes.Buffer), //os.Stdout, //TODO: change
		GOptions: ndGOpts,
		Author:   author,
		Message:  message,
		Pause:    pause,
		Change:   change,
	}
}

func getBuilderBuildOption(ndGOpts ndtypes.GlobalCommandOptions, buildOpts *types.BuildOperationOptions) (ndtypes.BuilderBuildOptions, error) {
	buildKitHost, err := buildkitutil.GetBuildkitHost(ndGOpts.Namespace)
	if err != nil {
		return ndtypes.BuilderBuildOptions{}, err
	}

	platform := []string{"amd64"}
	buildContext := buildOpts.BuildCtx
	output := ""
	tagValue := []string{buildOpts.Tag}
	progress := ""
	filename := ""
	target := ""
	buildArgs := []string{}
	label := []string{}
	noCache := false
	secret := []string{}
	ssh := []string{}
	cacheFrom := []string{}
	cacheTo := []string{}
	rm := false
	iidfile := ""
	quiet := false

	return ndtypes.BuilderBuildOptions{
		GOptions:     ndGOpts,
		BuildKitHost: buildKitHost,
		BuildContext: buildContext,
		Output:       output,
		Tag:          tagValue,
		Progress:     progress,
		File:         filename,
		Target:       target,
		BuildArgs:    buildArgs,
		Label:        label,
		NoCache:      noCache,
		Secret:       secret,
		SSH:          ssh,
		CacheFrom:    cacheFrom,
		CacheTo:      cacheTo,
		Rm:           rm,
		IidFile:      iidfile,
		Quiet:        quiet,
		Platform:     platform,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
		Stdin:        os.Stdin,
	}, nil
}

func getImageGlobalOption(gopts ndtypes.GlobalCommandOptions) ndtypes.ImageListOptions {

	var filters []string

	quiet := false
	noTrunc := false
	format := ""
	var inputFilters []string
	digests := false
	names := false

	return ndtypes.ImageListOptions{
		GOptions:         gopts,
		Quiet:            quiet,
		NoTrunc:          noTrunc,
		Format:           format,
		Filters:          inputFilters,
		NameAndRefFilter: filters,
		Digests:          digests,
		Names:            names,
		All:              true,
		Stdout:           os.Stdout,
	}
}

func getLoginOptions(gopts ndtypes.GlobalCommandOptions, serverAddr string, id string, passwd string) ndtypes.LoginCommandOptions {
	return ndtypes.LoginCommandOptions{
		GOptions:      gopts,
		ServerAddress: serverAddr,
		Username:      id,
		Password:      passwd,
	}
}
