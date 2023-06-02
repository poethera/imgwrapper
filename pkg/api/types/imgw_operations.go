package types

type IImgWOperations interface {
	registry_login(imgwCtx *ImgWContext, id string, password string) bool
	registry_logout(imgCtx *ImgWContext) bool
	//image_pull(imgctx *ImgWContext) bool
	image_push(imgwCtx *ImgWContext) bool
	image_build(imgwCtx *ImgWContext) bool
	image_commit_for_container(imgctx *ImgWContext) bool
}
