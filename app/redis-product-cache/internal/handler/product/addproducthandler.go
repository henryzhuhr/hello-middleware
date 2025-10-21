package product

import (
	"net/http"

	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/logic/product"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/svc"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddProductReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewAddProductLogic(r.Context(), svcCtx)
		resp, err := l.AddProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
