package types

type IImgWOperations interface {
	Registry_Login(imgwCtx *ImgWContext, imgwOpts *ImgWOptions, serverAddr string, id string, passwd string) error
	Registry_Logout(imgwCtx *ImgWContext) error
	//image_pull(imgwCtx *ImgWContext) error
	Image_Push(imgwCtx *ImgWContext) error
	Image_Build(imgwCtx *ImgWContext, imgwOpts *ImgWOptions, buildOpts *BuildOperationOptions) error
	Image_Commit_For_Container(imgwCtx *ImgWContext, imgwOpts *ImgWOptions, commitOpts *CommitOperationOptions) error
}
