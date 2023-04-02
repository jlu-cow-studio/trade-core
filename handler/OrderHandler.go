package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/trade-core/biz"
	"github.com/sanity-io/litter"
)

func (h *Handler) Order(ctx context.Context, req *trade_core.OrderRequest) (res *trade_core.OrderResponse, err error) {

	res = &trade_core.OrderResponse{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
	}

	cmd := redis.DB.Get(redis.GetUserTokenKey(req.Base.Token))
	if cmd.Err() != nil {
		res.Base.Message = cmd.Err().Error()
		res.Base.Code = "401"
		return res, nil
	}

	info := &redis_model.UserInfo{}
	fmt.Println("get user info :", cmd.Val())

	if err := json.Unmarshal([]byte(cmd.Val()), info); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "402"
		return res, nil
	}

	transaction, order, err := biz.Order(info.Uid, strconv.FormatInt(req.ItemId, 10), int(req.Count))

	if err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "403"
		return res, nil
	}

	litter.Dump(transaction)
	litter.Dump(order)

	res.Base.Message = ""
	res.Base.Code = "200"

	return
}
