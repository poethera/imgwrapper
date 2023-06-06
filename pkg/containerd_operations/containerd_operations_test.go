package containerd_operations

import (
	"context"
	"fmt"
	"imgwrapper/pkg/api/types"
	"reflect"
	"testing"

	"github.com/containerd/containerd"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "default", //"k8s.io",
		Address:   "/run/containerd/containerd.sock",
	}

	imgwCtx, err := Create(ctx, &imgwOpts)
	if err != nil {
		t.Errorf("Create Error: %s", err.Error())
	}

	assert.NotNil(imgwCtx.Client, "imgWContext.Client is nil")
	assert.NotNil(imgwCtx.Ctx, "imgWContext.Ctx is nil")
	assert.NotNil(imgwCtx.Cancel, "imgWContext.Cancel is nil")
	assert.NotNil(imgwCtx.ImgWOps, "imgWContext.ImgWOps is nil")
	defer imgwCtx.Cancel()

	if _, ok := imgwCtx.Client.(*containerd.Client); ok != true {
		t.Errorf("Client type mismatch *containerd.Client = %s", reflect.TypeOf(imgwCtx.Client))
	}
	if _, ok := imgwCtx.ImgWOps.(*ContainerdOperations); ok != true {
		t.Errorf("ImgWOpts real type mismatch ContainerdOperation = %s", reflect.TypeOf(imgwCtx.ImgWOps))
	}

	assert.Equal(imgwOpts.Engine, imgwCtx.ImgWOps.(*ContainerdOperations).Engine, "Engine must be same.")
}

func Test_Image_Build(t *testing.T) {
	//test. buildkitd check -> 확인
	//test. error dockerfile check -> 확인

	assert := assert.New(t)

	ctx := context.Background()
	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "default", //"k8s.io",
		Address:   "/run/containerd/containerd.sock",
	}

	imgwCtx, err := Create(ctx, &imgwOpts)
	if err != nil {
		t.Errorf("Create Error: %s", err.Error())
	}

	/*
		Tag      string
		BuildCtx string
	*/
	dockerfile := fmt.Sprintln(`
		FROM busybox
		CMD ["echo", "build-test-echo"]
	`)
	buildOpts := types.BuildOperationOptions{
		Tag:      "build-test-tag:v0.1",
		BuildCtx: dockerfile,
	}

	if err := imgwCtx.ImgWOps.Image_Build(imgwCtx, &imgwOpts, &buildOpts); err != nil {
		assert.Nil(err, "Build Operation has failed")
		fmt.Println("build error: ", err.Error())
	}

	assert.True(true, "empty test")
}

func Test_Image_Commit_For_Container(t *testing.T) {
	//TODO: run 등 추가 구현 or 외부 프로세스 실행으로 테스트케이스 보완
	//test. default namepsace -> 매뉴얼 테스트 완료
	//test. k8s.io namepsace -> 매뉴얼 테스트 완료

	assert := assert.New(t)

	ctx := context.Background()
	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "default", //"k8s.io",
		Address:   "/run/containerd/containerd.sock",
	}

	imgwCtx, err := Create(ctx, &imgwOpts)
	if err != nil {
		t.Errorf("Create Error: %s", err.Error())
	}

	/*
		Req    string
		Rawref string

		Author  string
		Message string
		Change  []string
	*/
	commitOpts := types.CommitOperationOptions{
		Req:    "nginx-72d9a", //"nginx-4d023", //"72d9a2b27b62",
		Rawref: "nginx:alpine-cp1",

		Author:  "link",
		Message: "commit container test1",
		Change:  []string{},
	}

	err = imgwCtx.ImgWOps.Image_Commit_For_Container(imgwCtx, &imgwOpts, &commitOpts)
	assert.Nil(err, "Commit Operation has failed")
	//TODO: need to verify ?
}
