package types

type IImgWOperations interface {
	Registry_Login(imgwCtx *ImgWContext, id string, password string) error
	Registry_Logout(imgCtx *ImgWContext) error
	//image_pull(imgctx *ImgWContext) error
	Image_Push(imgwCtx *ImgWContext) error
	Image_Build(imgwCtx *ImgWContext, imgwOpts *ImgWOptions, buildOpts *BuildOperationOptions) error
	Image_Commit_For_Container(imgwCtx *ImgWContext, imgwOpts *ImgWOptions, commitOpts *CommitOperationOptions) error
}
