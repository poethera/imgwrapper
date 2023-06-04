package containerd_operations

import (
	"context"
	"imgwrapper/pkg/api/types"

	//ndtypes "github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/clientutil"
)

type ContainerdOperations struct {
	Engine string
}

func Create(ctx context.Context, imgwOpts *types.ImgWOptions) (*types.ImgWContext, error) {
	client, ctx, cancel, err := clientutil.NewClient(ctx, imgwOpts.Namespace, imgwOpts.Address)
	if err != nil {
		println(err.Error())
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

func (c *ContainerdOperations) Registry_Login(imgwCtx *types.ImgWContext, id string, password string) bool {
	return false
}
func (c *ContainerdOperations) Registry_Logout(imgCtx *types.ImgWContext) bool {
	return false
}

// image_pull(imgctx *ImgWContext) bool
func (c *ContainerdOperations) Image_Push(imgwCtx *types.ImgWContext) bool {
	return false
}

func (c *ContainerdOperations) Image_Build(imgwCtx *types.ImgWContext) bool {
	return false
}
func (c *ContainerdOperations) Image_Commit_For_Container(imgctx *types.ImgWContext) bool {
	

	return false
}
