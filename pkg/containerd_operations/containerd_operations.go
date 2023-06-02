package containerd_operations

import (
	"context"
	"imgwrapper/pkg/api/types"

	//ndtypes "github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/clientutil"
)

type ContainerdOperations struct {
	ctx context.Context
}

func Create(ctx context.Context, imgwOpts *types.ImgWOptions) (*types.ImgWContext, error) {
	client, ctx, cancel, err := clientutil.NewClient(ctx, imgwOpts.Namespace, imgwOpts.Address)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	//defer cancel()
	//TODO: add cancel

	imgwCtx := &(types.ImgWContext{
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		imgwOps: ContainerdOperations{},
	})

	return imgwCtx, err
	//return types.ImgWContext{client, ctx, cancel, &ContainerdOperations{}}, nil
}

func (c *ContainerdOperations) registry_login(imgwCtx *types.ImgWContext, id string, password string) bool {
	return false
}
func (c *ContainerdOperations) registry_logout(imgCtx *types.ImgWContext) bool {
	return false
}

// image_pull(imgctx *ImgWContext) bool
func (c *ContainerdOperations) image_push(imgwCtx *types.ImgWContext) bool {
	return false
}
func (c *ContainerdOperations) image_build(imgwCtx *types.ImgWContext) bool {
	return false
}
func (c *ContainerdOperations) image_commit_for_container(imgctx *types.ImgWContext) bool {
	return false
}
