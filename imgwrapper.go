package imgwrapper

import (
	"context"
	"fmt"
	"imgwrapper/pkg/api/types"
	"imgwrapper/pkg/containerd_operations"
)

func CommitContainerNPush(imgwOpts *types.ImgWOptions, commitOpts *types.CommitOperationOptions) error {
	var (
		pushOpts = types.PushOperationOptions{
			Rawref: commitOpts.Rawref,
		}
	)

	imgwCtx, err := initialize(imgwOpts)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	defer imgwCtx.Cancel()
	if err := imgwCtx.ImgWOps.Image_Commit_For_Container(imgwCtx, imgwOpts, commitOpts); err != nil {
		return err
	}

	if err := imgwCtx.ImgWOps.Registry_Login(imgwCtx, imgwOpts); err != nil {
		return err
	}

	if err := imgwCtx.ImgWOps.Image_Push(imgwCtx, imgwOpts, &pushOpts); err != nil {
		return err
	}

	imgwCtx.ImgWOps.Registry_Logout(imgwCtx, imgwOpts)

	return nil
}

func BuildNPush(imgwOpts *types.ImgWOptions, buildOpts *types.BuildOperationOptions) error {
	//TODO: build
	var (
		pushOpts = types.PushOperationOptions{
			Rawref: buildOpts.Tag,
		}
	)

	imgwCtx, err := initialize(imgwOpts)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	defer imgwCtx.Cancel()
	if err := imgwCtx.ImgWOps.Image_Build(imgwCtx, imgwOpts, buildOpts); err != nil {
		return err
	}

	if err := imgwCtx.ImgWOps.Registry_Login(imgwCtx, imgwOpts); err != nil {
		return err
	}

	if err := imgwCtx.ImgWOps.Image_Push(imgwCtx, imgwOpts, &pushOpts); err != nil {
		return err
	}

	imgwCtx.ImgWOps.Registry_Logout(imgwCtx, imgwOpts)

	return nil
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
