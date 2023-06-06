package imgwrapper

import (
	"context"
	"fmt"
	"imgwrapper/pkg/api/types"
	"imgwrapper/pkg/containerd_operations"
)

func ImgWCommitContainer(imgwOpts *types.ImgWOptions, commitOpts *types.CommitOperationOptions) error {
	imgwCtx, err := initialize(imgwOpts)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	defer imgwCtx.Cancel()
	imgwCtx.ImgWOps.Image_Commit_For_Container(imgwCtx, imgwOpts, commitOpts)

	//TODO: push

	return nil
}

func ImgWBuild(imgwCtx *types.ImgWContext) {
	//TODO: build

	//TODO: push
}

func initialize(imgwOpts *types.ImgWOptions) (*types.ImgWContext, error) {

	//
	//imgctx := types.ImgWContext{}
	ctx := context.Background()
	imgwCtx, err := containerd_operations.Create(ctx, imgwOpts)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	return imgwCtx, nil
}
