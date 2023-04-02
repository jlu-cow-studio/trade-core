package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	"github.com/jlu-cow-studio/trade-core/biz"
)

func (h *Handler) OrderList(ctx context.Context, req *trade_core.OrderListRequest) (res *trade_core.OrderListResponse, err error) {
	res = &trade_core.OrderListResponse{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
		OrderList: []*trade_core.OrderForList{},
	}

	offset := req.Page * req.PageSize
	limit := req.PageSize

	list, err := biz.OrderList(strconv.FormatInt(req.UserId, 10), int(offset), int(limit))

	for _, ofl := range list {
		oi := &trade_core.OrderForList{
			Id:        int64(ofl.OrderID),
			UserId:    int64(ofl.UserID),
			ItemId:    int64(ofl.ItemID),
			Quantity:  int32(ofl.PurchaseQuantity),
			CreatedAt: ofl.PurchaseDate.Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		if ofl.ItemInfo != nil {
			oi.ItemInfo = &product_core.ItemInfo{
				ItemId:      ofl.ItemInfo.ID,
				Name:        ofl.ItemInfo.Name,
				Description: ofl.ItemInfo.Description,
				Category:    ofl.ItemInfo.Category,
				Price:       ofl.ItemInfo.Price,
				Stock:       ofl.ItemInfo.Stock,
				ImageUrl:    ofl.ItemInfo.ImageURL,
				Province:    ofl.ItemInfo.Province,
				City:        ofl.ItemInfo.City,
				District:    ofl.ItemInfo.District,
				UserId:      ofl.ItemInfo.UserID,
				UserType:    ofl.ItemInfo.UserType,
			}
		}

		res.OrderList = append(res.OrderList, oi)
	}
	res.Base.Message = ""
	res.Base.Code = "200"
	return
}
