package errorcode

const (
	GetStockError                 = 500100
	OutOfStockError               = 500101
	OutOfStockErrorMsg            = "库存不足"
	CategoryIdError               = 402101
	CategoryIdErrorMessage        = "类目 ID 错误"
	SPUNameEmpty                  = 402102
	SPUNameEmptyMessage           = "商品名空白"
	SPUPriceInvalid               = 402103
	SPUPriceInvalidMessage        = "价格无效"
	SPUUnitTooLarge               = 402104
	SPUUnitTooLargeMessage        = "单位文字过长"
	SPUImagesCountTooLarge        = 402105
	SPUImagesCountTooLargeMessage = "商品图片数量"
	SPUImagesRequired             = 402106
	SPUImagesRequiredMessage      = "需要至少一个商品图片"
	SPUStockInvalid               = 402107
	SPUStockInvalidMessage        = "库存数无效"
	ShopIdInvalid                 = 402108
	ShopIdInvalidMessage          = "店铺 ID 无效"
	SortInvalid                   = 402109
	SortInvalidMessage            = "排序必须是 0 - 9999 之间的数字"
)
