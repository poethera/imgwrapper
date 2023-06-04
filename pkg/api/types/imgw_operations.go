package types

type IImgWOperations interface {
	Registry_Login(imgwCtx *ImgWContext, id string, password string) bool
	Registry_Logout(imgCtx *ImgWContext) bool
	//image_pull(imgctx *ImgWContext) bool
	Image_Push(imgwCtx *ImgWContext) bool
	Image_Build(imgwCtx *ImgWContext) bool
	Image_Commit_For_Container(imgctx *ImgWContext) bool
}
