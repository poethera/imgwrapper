package containerd_operations

import (
	"context"
	"imgwrapper/pkg/api/types"
	"reflect"
	"testing"

	"github.com/containerd/containerd"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	imgwOpts := types.ImgWOptions{
		Engine:    "containerd",
		Namespace: "k8s.io",
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
