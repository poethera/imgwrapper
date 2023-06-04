package types

import "context"

type ImgWContext struct {
	Client  interface{}
	Ctx     context.Context
	Cancel  context.CancelFunc
	ImgWOps IImgWOperations
}
