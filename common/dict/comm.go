package dict

const (
	StatusOff    = 0 // 未启用/待审核
	StatusOn     = 1 // 已启用/审核通过
	StatusNoPass = 2 // 审核未通过
	DeleteNot    = 0

	CategoryArticle = 1
	CategoryLabel   = 2
	CategoryProduct = 3

	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)