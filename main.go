package main

import (
	"fmt"
	"time"

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

type User struct {
	ID     uint   `gorm:"primarykey;autoIncrement;"`
	Name   string `gorm:"size:24"`
	Active bool
	Age    uint `gorm:"size:11"`
}

type result struct {
	Date  time.Time
	Total int
}

type Email struct {
	UserId uint
	Name   string
	Email  string
}

type Result struct {
	Name string
	Age  int
}

var db *gorm.DB

func test_user() {
	var user User
	var users []User
	// Get first matched record
	db.Where("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	// Get all matched records
	db.Where("name <> ?", "jinzhu").Find(&users)
	// SELECT * FROM users WHERE name <> 'jinzhu';

	// IN
	db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

	// LIKE
	db.Where("name LIKE ?", "%jin%").Find(&users)
	// SELECT * FROM users WHERE name LIKE '%jin%';

	// AND
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

	// Time
	//db.Where("updated_at > ?", lastWeek).Find(&users)
	// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

	// BETWEEN
	//db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
	//Struct & Map 条件

	// Struct
	db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

	// Map
	db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

	// Slice of primary keys
	db.Where([]int64{20, 21, 22}).Find(&users)

	//To include zero values in the query conditions, you can use a map, which will include all key-values as query conditions, for example:

	db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

	//#### 指定结构体查询字段

	db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

	db.Where(&User{Name: "jinzhu"}, "Age").Find(&users)
	// SELECT * FROM users WHERE age = 0;

	//#### Not 条件

	db.Not("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

	// Not In
	db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

	// Struct
	db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
	// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

	// Not In slice of primary keys
	db.Not([]int64{1, 2, 3}).First(&user)
	// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;

	//#### Or 条件
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
	// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

	// Struct
	db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

	// Map
	db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

	//#### 选择特定字段

	db.Select("name", "age").Find(&users)
	// SELECT name, age FROM users;

	db.Select([]string{"name", "age"}).Find(&users)
	// SELECT name, age FROM users;

	db.Table("users").Select("COALESCE(age,?)", 42).Rows()
	// SELECT COALESCE(age,'42') FROM users;

	//#### 排序

	db.Order("age desc, name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	// Multiple orders
	db.Order("age desc").Order("name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	//#### Limit & Offset

	db.Limit(3).Find(&users)
	// SELECT * FROM users LIMIT 3;

	var users1 []User
	var users2 []User
	// Cancel limit condition with -1
	db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
	// SELECT * FROM users LIMIT 10; (users1)
	// SELECT * FROM users; (users2)

	db.Offset(3).Find(&users)
	// SELECT * FROM users OFFSET 3;

	db.Limit(10).Offset(5).Find(&users)
	// SELECT * FROM users OFFSET 5 LIMIT 10;

	// Cancel offset condition with -1
	db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
	// SELECT * FROM users OFFSET 10; (users1)
	// SELECT * FROM users; (users2)

	//#### Joins
	db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id
	rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
	if err == nil {
		for rows.Next() {

		}
	}
	var results []result
	db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

	//#### Joins 预加载
	db.Joins("Company").Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;

	db.InnerJoins("Company").Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` INNER JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;

	//#### Joins 一个衍生表
	query := db.Table("order").Select("MAX(order.finished_at) as latest").Joins("left join user user on order.user_id = user.id").Where("user.age > ?", 18).Group("order.user_id")
	fmt.Println(query)
	db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)
	// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest

	//#### Scan
	var result Result
	db.Table("users").Select("name", "age").Where("name = ?", "Antonio").Scan(&result)

	// Raw SQL
	db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
}

type Order struct {
	UserId     int
	FinishedAt *time.Time
}

func test_product() {
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

	//db.First(&product)
	//db.First(&product)
	// 获取最后一条记录（主键降序）
	res2 := db.Last(&product)
	fmt.Println(res2.RowsAffected)

	// works with Take
	res3 := map[string]interface{}{}
	db.Table("products").Take(&res3)

	fmt.Println(res3)

	var products []Product

	res4 := db.Find(&products, []int{1, 2, 3})
	fmt.Println(res4)
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:123456@tcp(localhost:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:123456@tcp(localhost:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 256,
	}))
	db = _db
	if err != nil {
		fmt.Println(err)
	}
	// 迁移 schema
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Email{})

	test_product()
	test_user()

	// Get all records
	//result := db.Find(&products)
	// SELECT * FROM users WHERE id IN (1,2,3);

	//batch Insert
	//var users = []User{user1,user2,user3}
	//db.Create(&users)

	//db.Delete(&product, 1)

}

type User3 struct {
	ID     uint
	Name   string
	Age    int
	Gender string
	// 假设后面还有几百个字段...
}

type APIUser struct {
	ID   uint
	Name string
}

func test3() {
	var user User
	// 查询时会自动选择 `id`, `name` 字段
	db.Model(&User3{}).Limit(10).Find(&APIUser{})
	// SELECT `id`, `name` FROM `users` LIMIT 10

	//db2, err := gorm.Open(mysql.Open("gorm.db"), &gorm.Config{
	//	QueryFields: true,
	//})

	//db2.Find(&user)
	// SELECT `users`.`name`, `users`.`age`, ... FROM `users` // 带上这个选项

	// Session Mode
	//db2.Session(&gorm.Session{QueryFields: true}).Find(&user)
	// SELECT `users`.`name`, `users`.`age`, ... FROM `users`

	//#### FirstOrInit
	//获取第一条匹配的记录，或者根据给定的条件初始化一个实例（仅支持 sturct 和 map 条件）

	// 未找到 user，则根据给定的条件初始化一条记录
	db.FirstOrInit(&user, User{Name: "non_existing"})
	// user -> User{Name: "non_existing"}

	// 找到了 `name` = `jinzhu` 的 user
	db.Where(User{Name: "jinzhu"}).FirstOrInit(&user)
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// 找到了 `name` = `jinzhu` 的 user
	db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// ###### 如果没有找到记录，可以使用包含更多的属性的结构体初始化 user，Attrs 不会被用于生成查询 SQL

	// 未找到 user，则根据给定的条件以及 Attrs 初始化 user
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// 未找到 user，则根据给定的条件以及 Attrs 初始化 user
	db.Where(User{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// 找到了 `name` = `jinzhu` 的 user，则忽略 Attrs
	db.Where(User{Name: "Jinzhu"}).Attrs(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	//###### 不管是否找到记录，Assign 都会将属性赋值给 struct，但这些属性不会被用于生成查询 SQL，也不会被保存到数据库

	// 未找到 user，根据条件和 Assign 属性初始化 struct
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)
	// user -> User{Name: "non_existing", Age: 20}

	// 找到 `name` = `jinzhu` 的记录，依然会更新 Assign 相关的属性
	db.Where(User{Name: "Jinzhu"}).Assign(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 20}

}

func test4() {
	//#### FirstOrCreate
	// 未找到 User，根据给定条件创建一条新纪录
	var user User
	result := db.FirstOrCreate(&user, User{Name: "non_existing"})
	fmt.Println(result)
	// INSERT INTO "users" (name) VALUES ("non_existing");
	// user -> User{ID: 112, Name: "non_existing"}
	// result.RowsAffected // => 1

	// 找到 `name` = `jinzhu` 的 User
	result2 := db.Where(User{Name: "jinzhu"}).FirstOrCreate(&user)
	fmt.Println(result2)
	// user -> User{ID: 111, Name: "jinzhu", "Age": 18}
	// result.RowsAffected // => 0

	//###### 如果没有找到记录，可以使用包含更多的属性的结构体创建记录，Attrs 不会被用于生成查询 SQL 。

	// 未找到 user，根据条件和 Assign 属性创建记录
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// 找到了 `name` = `jinzhu` 的 user，则忽略 Attrs
	db.Where(User{Name: "jinzhu"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "jinzhu", Age: 18}

}

func test5() {
	//#### 迭代
	//GORM 支持通过行进行迭代

	rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Rows()
	defer rows.Close()
	fmt.Println(err)
	for rows.Next() {
		var user User
		// ScanRows 方法用于将一行记录扫描至结构体
		db.ScanRows(rows, &user)

		// 业务逻辑...
	}
}

func test6() {
	//#### FindInBatches
	// 每次批量处理 100 条
	var results []User
	result := db.Where("age = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			// 批量处理找到的记录
			fmt.Println(result)
		}

		tx.Save(&results)

		//tx.RowsAffected // 本次批量操作影响的记录数

		//batch // Batch 1, 2, 3

		// 如果返回错误会终止后续批量操作
		return nil
	})
	fmt.Println(result)

	//result.Error        // returned error
	//result.RowsAffected // 整个批量操作影响的记录数
}

func test7() {
	//#### Count
	//Count 用于获取匹配的记录数
	var count int64
	db.Model(&User{}).Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'

	db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)

	db.Table("deleted_users").Count(&count)
	// SELECT count(1) FROM deleted_users;

	// Count with Distinct
	db.Model(&User{}).Distinct("name").Count(&count)
	// SELECT COUNT(DISTINCT(`name`)) FROM `users`

	db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
	// SELECT count(distinct(name)) FROM deleted_users

	// Count with Group
	users := []User{
		{Name: "name1"},
		{Name: "name2"},
		{Name: "name3"},
		{Name: "name3"},
	}
	fmt.Println(users)

	db.Model(&User{}).Group("name").Count(&count)
	//count // => 3
}

//#### 保存所有字段

func test8() {
	var user User
	db.First(&user)

	user.Name = "jinzhu 2"
	user.Age = 100
	db.Save(&user)
	// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;

	//Save is a combination function. If save value does not contain primary key, it will execute Create, otherwise it will execute Update (with all fields).
	db.Save(&User{Name: "jinzhu", Age: 100})
	// INSERT INTO `users` (`name`,`age`,`birthday`,`update_at`) VALUES ("jinzhu",100,"0000-00-00 00:00:00","0000-00-00 00:00:00")

	db.Save(&User{ID: 1, Name: "jinzhu", Age: 100})
	// UPDATE `users` SET `name`="jinzhu",`age`=100,`birthday`="0000-00-00 00:00:00",`update_at`="0000-00-00 00:00:00" WHERE `id` = 1

}

//#### 更新单个列

func test9() {
	var user User
	// Update with conditions
	db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

	// User's ID is `111`:
	db.Model(&user).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

	// Update with conditions and model value
	db.Model(&user).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

}

//#### 更新多列

func test10() {
	var user User
	// Update attributes with `struct`, will only update non-zero fields
	db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
	// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

	// Update attributes with `map`
	db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

}
