package constant

const DB_CONNECT_TEMPLATE string = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"

type Code uint32

const TotalTraffic uint32 = 1000

const (
	OK                      Code = 0
	Error                   Code = 1
	FeatureValueNotExist    Code = 2
	TrafficOverflow         Code = 3
	LackDefaultFeatureValue Code = 4
	FeatureNameEmpty        Code = 5
	FeatureKeyEmpty         Code = 6
)

const (
	ReasonTypeBlacklist uint8 = 1
	ReasonTypeWhitelist uint8 = 2
	ReasonTypeHash      uint8 = 3
)
