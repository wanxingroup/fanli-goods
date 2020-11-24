package constants

const RedisConfigKey = "common"

const (
	Prefix_Proj      = "goods:"
	RedisKeySKUStock = Prefix_Proj + "sku_stock:{%v}:%v" // goods:sku_stock:{${spuid}}:${skuid}
)
