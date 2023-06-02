package containerd_operations

import (
	"context"
	"fmt"
	"imgwrapper/pkg/api/types"
	"testing"
)

func TestCreate(t *testing.T) {

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

	if imgwCtx.client == nil {
		t.Error("Invalid client field")
	}
	
}
