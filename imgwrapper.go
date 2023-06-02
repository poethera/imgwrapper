package imgwrapper

import (
	"context"
	"imgwrapper/pkg/api/types"
	"imgwrapper/pkg/containerd_operations"
)

func ImgWInitialize(imgwOpts *types.ImgWOptions) *types.ImgWContext {

	//
	//imgctx := types.ImgWContext{}
	ctx := context.Background()
	imgwCtx, err := containerd_operations.Create(ctx, imgwOpts)
	if err != nil {
		println(err.Error())
		return nil
	}

	return imgwCtx
}

func ImgWCommitContainer(imgwCtx *types.ImgWContext) {

}

func ImgWBuild(imgwCtx *types.ImgWContext) {

}

