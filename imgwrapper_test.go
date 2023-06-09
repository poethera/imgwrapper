package imgwrapper

import (
	"imgwrapper/pkg/api/types"
	"os"
	"os/exec"
	"testing"

	"github.com/containerd/nerdctl/pkg/imgutil/dockerconfigresolver"
	"github.com/stretchr/testify/assert"
)

// func Test_Prepare(t *testing.T) {
// 	var stdout = new(bytes.Buffer)

// 	// path, err := exec.LookPath("nerdctl")
// 	// if errors.Is(err, exec.ErrDot) {
// 	// 	err = nil
// 	// }
// 	// fmt.Println("found path: ", path)
// 	var (
// 		target_image = "nginx:alpine"
// 		target_contr = "gotest_nginx_new"
// 	)

// 	cmd := exec.Command("nerdctl", "run", "-d", "--name", target_contr, target_image)
// 	cmd.Stdout = stdout //os.Stdout

// 	//cmd.Run()
// 	cmd.Start()
// 	cmd.Wait()

// 	fmt.Printf(stdout.String())
// }

func Test_CommitContainerNPush(t *testing.T) {

	assert := assert.New(t)

	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "default", //"k8s.io",
		Address:   "/run/containerd/containerd.sock",

		ServerAddress: dockerconfigresolver.IndexServer, //sds.redii.net:443
		UserId:        "poethera",                       //, "sangjiny.nam"
		UserPasswd:    "",
	}

	var (
		target_image = "nginx:alpine"
		target_contr = "gotest_nginx_new"
		new_image    = imgwOpts.UserId + "/" + "nginx:commit_new"
	)

	commitOpts := types.CommitOperationOptions{
		Req:    target_contr,
		Rawref: new_image,

		Author:  "test",
		Message: "test commit and push",
		Change:  []string{},
	}

	cmd := exec.Command("nerdctl", "run", "-d", "--name", target_contr, target_image)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		assert.NotNil(err, err.Error())
	}

	if err := CommitContainerNPush(&imgwOpts, &commitOpts); err != nil {
		assert.NotNil(err, err.Error())
	}

	cmd = exec.Command("nerdctl", "rm", "-f", target_contr)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		assert.NotNil(err, err.Error())
	}

	cmd = exec.Command("nerdctl", "rmi", target_image, new_image)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		assert.NotNil(err, err.Error())
	}
}

func Test_BuildNPush(t *testing.T) {
	assert := assert.New(t)

	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "default", //"k8s.io",
		Address:   "/run/containerd/containerd.sock",

		ServerAddress: dockerconfigresolver.IndexServer, //sds.redii.net:443
		UserId:        "poethera",                       //, "sangjiny.nam"
		UserPasswd:    "",
	}

	var (
		target_image = "nginx:alpine"
		new_image    = imgwOpts.UserId + "/" + "nginx:build_new"
	)

	dockerfile := `
FROM ` + target_image + `
CMD ["echo", "Test_BuildNPush"]"
`

	buildOpts := types.BuildOperationOptions{
		Tag:      new_image,
		BuildCtx: dockerfile,
	}

	if err := BuildNPush(&imgwOpts, &buildOpts); err != nil {
		assert.NotNil(err, err.Error())
	}

	cmd := exec.Command("nerdctl", "rmi", target_image, new_image)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		assert.NotNil(err, err.Error())
	}
}
