package containerd_operations

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"imgwrapper/pkg/api/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd"
	ndtypes "github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/buildkitutil"
	"github.com/containerd/nerdctl/pkg/clientutil"
	"github.com/containerd/nerdctl/pkg/cmd/builder"
	"github.com/containerd/nerdctl/pkg/cmd/container"
	"github.com/containerd/nerdctl/pkg/cmd/image"
	"github.com/containerd/nerdctl/pkg/cmd/login"
	"github.com/containerd/nerdctl/pkg/imgutil/dockerconfigresolver"
	dockercliconfig "github.com/docker/cli/cli/config"
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

func (c *ContainerdOperations) Registry_Login(imgwCtx *types.ImgWContext, imgwOpts *types.ImgWOptions) error {
	Stdout := new(bytes.Buffer)
	defer func() {
		Stdout = nil
	}()

	options := getLoginOptions(imgwOpts)

	err := login.Login(imgwCtx.Ctx, options, Stdout)
	//TODO: error
	if err != nil {
		fmt.Println("login error", err.Error())
	}

	//TODO: Stdout
	fmt.Println("login output: ", Stdout.String())

	return err
}

func (c *ContainerdOperations) Registry_Logout(imgCtx *types.ImgWContext, imgwOpts *types.ImgWOptions) error {
	Stdout := new(bytes.Buffer)
	defer func() {
		Stdout = nil
	}()

	serverAddress := dockerconfigresolver.IndexServer
	isDefaultRegistry := true

	if len(imgwOpts.ServerAddress) >= 1 {
		if !strings.Contains(imgwOpts.ServerAddress, "index.docker.io") {
			serverAddress = imgwOpts.ServerAddress
			isDefaultRegistry = false
		}
	}

	var (
		regsToLogout    = []string{}
		hostnameAddress = serverAddress
	)

	if !isDefaultRegistry {
		hostnameAddress = dockerconfigresolver.ConvertToHostname(serverAddress)
		// the tries below are kept for backward compatibility where a user could have
		// saved the registry in one of the following format.

		//TODO: changed the order
		regsToLogout = append(regsToLogout, "https://"+hostnameAddress, "http://"+hostnameAddress, hostnameAddress)
	} else {
		regsToLogout = append(regsToLogout, serverAddress)
	}

	fmt.Fprintf(Stdout, "Removing login credentials for %s\n", hostnameAddress)

	dockerConfigFile, err := dockercliconfig.Load("")
	if err != nil {
		return err
	}
	errs := make(map[string]error)
	for _, r := range regsToLogout {
		if err := dockerConfigFile.GetCredentialsStore(r).Erase(r); err != nil {
			errs[r] = err
		}
	}

	// if at least one removal succeeded, report success. Otherwise report errors
	if len(errs) == len(regsToLogout) {
		fmt.Fprintln(Stdout, "WARNING: could not erase credentials:")
		for k, v := range errs {
			fmt.Fprintf(Stdout, "%s: %s\n", k, v)
		}
	}

	//TODO: Stdout
	fmt.Println("logout output: \n", Stdout.String())

	return nil
}

// image_pull(imgctx *ImgWContext) error
func (c *ContainerdOperations) Image_Push(imgwCtx *types.ImgWContext, imgwOpts *types.ImgWOptions, pushOpts *types.PushOperationOptions) error {
	client, ok := imgwCtx.Client.(*containerd.Client)
	if !ok {
		//TODO: error
		fmt.Println("invalid client object")
		return errors.New("invalid client object")
	}

	options := getImagePushOption(imgwOpts, pushOpts)
	_stdout := new(bytes.Buffer)
	options.Stdout = _stdout
	defer func() {
		options.Stdout = nil
	}()
	if err := image.Push(imgwCtx.Ctx, client, pushOpts.Rawref, options); err != nil {
		fmt.Println("Image Push error", err.Error())
		return err
	}

	//TODO:
	fmt.Println("push output: \n", _stdout.String())

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
		imgwOpts,
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
		imgwOpts,
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

func getContainerCommitOptions(imgwOpts *types.ImgWOptions, commitOpts *types.CommitOperationOptions) ndtypes.ContainerCommitOptions {

	gopts := getGlobalOptions(imgwOpts)

	author := commitOpts.Author
	message := commitOpts.Message
	pause := false
	change := commitOpts.Change

	return ndtypes.ContainerCommitOptions{
		Stdout:   os.Stdout, //new(bytes.Buffer), //os.Stdout, //TODO: change
		GOptions: gopts,
		Author:   author,
		Message:  message,
		Pause:    pause,
		Change:   change,
	}
}

func getBuilderBuildOption(imgwOpts *types.ImgWOptions, buildOpts *types.BuildOperationOptions) (ndtypes.BuilderBuildOptions, error) {

	gopts := getGlobalOptions(imgwOpts)

	buildKitHost, err := buildkitutil.GetBuildkitHost(gopts.Namespace)
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
		GOptions:     gopts,
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

func getLoginOptions(imgwOpts *types.ImgWOptions) ndtypes.LoginCommandOptions {

	gopts := getGlobalOptions(imgwOpts)

	return ndtypes.LoginCommandOptions{
		GOptions:      gopts,
		ServerAddress: imgwOpts.ServerAddress,
		Username:      imgwOpts.UserId,
		Password:      imgwOpts.UserPasswd,
	}
}

func getImagePushOption(imgwOpts *types.ImgWOptions, pushOpts *types.PushOperationOptions) ndtypes.ImagePushOptions {

	gopts := getGlobalOptions(imgwOpts)

	platform := []string{"amd64"}
	allPlatforms := false
	estargz := false
	ipfsEnsureImage := false
	ipfsAddress := ""
	quiet := false
	allowNonDist := false
	signOptions := ndtypes.ImageSignOptions{
		Provider:        "none",
		CosignKey:       "",
		NotationKeyName: "",
	}

	return ndtypes.ImagePushOptions{
		GOptions:                       gopts,
		SignOptions:                    signOptions,
		Platforms:                      platform,
		AllPlatforms:                   allPlatforms,
		Estargz:                        estargz,
		IpfsEnsureImage:                ipfsEnsureImage,
		IpfsAddress:                    ipfsAddress,
		Quiet:                          quiet,
		AllowNondistributableArtifacts: allowNonDist,
		Stdout:                         os.Stdout,
	}
}

// reserved
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
