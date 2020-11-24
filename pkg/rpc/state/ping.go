package state

import (
	"context"
	"errors"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
)

func (service *Controller) Ping(ctx context.Context, req *protos.PingRequest) (*protos.PingReply, error) {

	if req == nil {
		return nil, errors.New("input parameters error")
	}

	return &protos.PingReply{
		Message: req.Message,
	}, nil
}
