package app

var config = AppConfig{
	AppName: "WishlistSNS",
	Database: DatabaseConfig{
		Host:     "host.docker.internal",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "wishlist_dev",
		SSLMode:  "disable",
	},
	JWTSecret:     "your-super-secret-jwt-key-change-in-production",
	JWTExpiration: 24, // 24 hours
}
