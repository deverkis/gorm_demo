package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
*
gorm:"primary_key"	字段设置为主键
gorm:"AUTO_INCREMENT"	字段设置为自增
gorm:"size:20	字段长度设置为20
gorm:"index:idx_user_id	字段设置普通索引，名称为idx_user_id
gorm:"not null	设置字段为非空
gorm:"type:varchar(64)"	设置字段为varchar类型，长度为64
gorm:"column:remark"	设置数据库字段名为remark
gorm:"-"	忽略此字段，不在表中创建该字段
gorm:"default:'默认'"	设置字段的默认值
*/
type Product struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"default:'';type:varchar(30);not null;comment:编码"`
	Price       uint   `gorm:"default:0;type:int(11);not null;comment:价格"`
	CreatedTime uint   `gorm:"default:0;type:int(11);not null;comment:创建日期"`
	UpdatedTime uint   `gorm:"default:0;type:int(11);not null;comment:更新日期"`
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:123456@tcp(localhost:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:123456@tcp(localhost:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 256,
	}))
	if err != nil {
		fmt.Println(err)
	}
	// 迁移 schema
	db.AutoMigrate(&Product{})

	// Create
	result := db.Create(&Product{Code: "200", Price: 100})
	fmt.Println(result.Error)        // nil
	fmt.Println(result.RowsAffected) // 1
	// Read
	var product Product
	db.First(&product, 1) // 根据整型主键查找

	fmt.Println(product)

	db.Model(&product).Update("price", 300)

	db.Model(&product).Updates(Product{Price: 400, Code: "F34"})

	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F33"})

	//batch Insert
	//var users = []User{user1,user2,user3}
	//db.Create(&users)

	//db.Delete(&product, 1)

}
