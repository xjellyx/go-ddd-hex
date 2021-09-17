package application

// Repository 存储库接口
type Repository interface {
}

// Database 数据库基础组件接口
type Database interface {
	Connect() error
	InjectEntities(en ...interface{}) error
	DB() interface{}
}
