module gorm.io/gen/tests

go 1.16

require (
	golang.org/x/sys v0.1.0 // indirect
	gorm.io/driver/mysql v1.4.3
	gorm.io/driver/sqlite v1.4.3
	gorm.io/gen v0.3.16
	gorm.io/gorm v1.24.1-0.20221019064659-5dd2bb482755
	gorm.io/plugin/dbresolver v1.3.0
)

replace gorm.io/gen => ../
