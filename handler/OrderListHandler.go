package handler

import (
	"context"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	"github.com/jlu-cow-studio/trade-core/biz"
)

func (h *Handler) OrderList(ctx context.Context, req *trade_core.OrderListRequest) (res *trade_core.OrderListResponse, err error) {
	res = &trade_core.OrderListResponse{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
		OrderList: []*trade_core.Order{},
	}

	offset := req.Page * req.PageSize
	limit := req.PageSize

	list, err := biz.OrderList(strconv.FormatInt(req.UserId, 10), int(offset), int(limit))

	for _, o := range list {
		res.OrderList = append(res.OrderList, &trade_core.Order{
			Id:        int64(o.OrderID),
			UserId:    int64(o.UserID),
			ItemId:    int64(o.ItemID),
			Quantity:  int32(o.PurchaseQuantity),
			CreatedAt: o.PurchaseDate.Unix(),
			UpdatedAt: o.PurchaseDate.Unix(),
		})
	}
	res.Base.Message = ""
	res.Base.Code = "200"
	return
}
