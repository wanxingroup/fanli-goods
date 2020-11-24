package errorcode

const (
	SpuNotFoundError         = 402108 // SPU不存在
	SPUNotFoundErrorMessage  = "SPU不存在"
	SPUForbiddenError        = 402109 // SPU无权限
	SPUForbiddenErrorMessage = "无权限访问此 SPU"
	SPUNoLongerExist         = 402110 // SPU已经不存在
	SPUNoLongerExistMessage  = "SPU已经不存在"

	CreateCategoryError              = 502001 // 创建类目失败
	CreateCategoryErrorMessage       = "创建类目失败"
	UpdateCategoryError              = 502002 // 更新类目失败
	UpdateCategoryErrorMessage       = "更新类目失败"
	DeleteCategoryError              = 502003 // 删除类目失败
	DeleteCategoryErrorMessage       = "删除类目失败"
	GetCategoriesError               = 502004 // 查询类目列表失败
	GetCategoryError                 = 502005 // 查询类目失败
	CategoryNotFoundError            = 502006 // 类目不存在
	CategoryNotFoundErrorMessage     = "类目不存在"
	CategoryForbidError              = 502007 // 类目无权限
	DeleteNotNilCategoryError        = 502008 // 无法删除非空类目
	DeleteNotNilCategoryErrorMessage = "无法删除非空类目"
	CategoryForbidErrorMessage       = "类目无权限"
	CreateSpuError                   = 502101 // 创建SPU失败
	UpdateSpuError                   = 502102 // 更新SPU失败
	UpdateSpuErrorMessage            = "更新SPU失败"
	GetSpusError                     = 502104 // 查询SPU列表失败
	GetSpuError                      = 502105 // 查询SPU失败
	InvalidStockError                = 502106 // 非法的库存值
	SetStockError                    = 502107 // 设置库存失败
	UpdateSpuStatusError             = 502108 // 更新SPU状态失败
	GetSkuError                      = 502201 // 查询SKU失败
	SkuNotFoundError                 = 502202 // SKU不存在
	GetSkuStockError                 = 502203 // 获取SKU库存失败
)
