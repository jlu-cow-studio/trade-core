package handler

import "github.com/jlu-cow-studio/common/dal/rpc/trade_core"

type Handler struct {
	trade_core.UnimplementedTreadeCoreServiceServer
}
