package types

import "context"

type ImgWContext struct {
	client  interface{}
	ctx     context.Context
	cancel  context.CancelFunc
	imgwOps IImgWOperations
}
