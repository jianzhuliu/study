package db

import (
	"fmt"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var (
	Db *gorm.DB
)

//根据参数，建立数据库连接
func Setup(db_user, db_passwd, db_host string, db_port int) (err error){
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
	db_user, db_passwd, db_host, db_port)
	
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		//DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		//DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	return
}

//获取数据库列表
func GetDatabaseList() ([]string, error){
	sql := "show databases"
	rows, err := Db.Raw(sql).Rows()
	
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	type Database struct {
		Database string
	}
	
	var result Database
	data := []string{}
	
	for rows.Next() {
		Db.ScanRows(rows, &result)
		data = append(data, result.Database)
	}
	
	return data,nil
}

//根据数据库名，获取所有表
func GetTableList(dbname string) ([]string, error){
	Db.Exec(fmt.Sprintf("use %s", dbname))

	sql := "show tables"
	rows, err := Db.Raw(sql).Rows()
	
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	
	var tblname string
	data := []string{}
	
	for rows.Next() {
		rows.Scan(&tblname)
		data = append(data, tblname)
	}
	
	return data,nil
}


type TblDesc struct {
	Field string 
	Type string 
	Null string 
	Key string
	Default string 
	Extra string 
}

func GetTableDesc(dbname, tblname string) ([]TblDesc, error){
	Db.Exec(fmt.Sprintf("use %s", dbname))

	sql := fmt.Sprintf("desc %s",tblname)
	rows, err := Db.Raw(sql).Rows()
	
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	
	data := []TblDesc{}
	var tmp TblDesc
	
	for rows.Next() {
		Db.ScanRows(rows, &tmp)
		data = append(data, tmp)
	}
	
	return data,nil
}



