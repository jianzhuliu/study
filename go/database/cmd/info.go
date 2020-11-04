package main

import (
	"fmt"
	"flag"
	"reflect"
	"database/sql"
	"database/sql/driver"
	
	"gitee.com/jianzhuliu/tools/common"
	"gitee.com/jianzhuliu/tools/conf"
	"gitee.com/jianzhuliu/tools/logger"
)

func init(){
	//初始化定义相关目录，比如日志文件夹
	err := common.DefaultInit()
	if err != nil {
		common.ExitErr("common.DefaultInit()", err)
	}

	flag.Parse()
	
	//初始化数据库
	err = common.InitDb()
	if err != nil {
		common.ExitErr("common.InitDb()", err)
	}
}

func main(){
	showDatabases()
}

func showDatabases(){
	db := conf.V_db
	
	_,err := db.Exec("use yiyan")
	if err != nil {
		common.ExitErr("db.Exec()",err)
	}
	
	//rows, err := db.Query("show databases")
	//rows, err := db.Query("desc t_employee")
	rows, err := db.Query("select * from yiyan limit 5")
	
	if err != nil {
		common.ExitErr("db.Query()",err)
	}
	
	defer rows.Close()
	
	//字段列表
	columns, err := rows.Columns()
	
	if err != nil {
		common.ExitErr("rows.Columns()", err)
	}
	
	logger.PrintfWithTime("columns === %v", columns)
	
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		common.ExitErr("rows.ColumnTypes()", err)
	}
	
	//字段类型
	for _, columnType := range columnTypes {
		logger.Printf("columntype === name:%s, database_type:%s",columnType.Name(), columnType.DatabaseTypeName())
	}
	
	values := make([]interface{}, len(columns))
	
	for idx, columnType := range columnTypes {
		if columnType.ScanType() != nil {
			values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
		} else {
			values[idx] = new(interface{})
		}
	}
	
	for rows.Next(){
		rows.Scan(values...)
		//fmt.Println(values)
		mapValue := map[string]interface{}{}
		scanIntoMap(mapValue, values, columns)
		fmt.Printf("%#v \n",mapValue)
		//fmt.Printf("%+v \n",mapValue)
		//break
	}
}

func scanIntoMap(mapValue map[string]interface{}, values []interface{}, columns []string) {
	for idx, column := range columns {
		if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(values[idx]))); reflectValue.IsValid() {
			mapValue[column] = reflectValue.Interface()
			if valuer, ok := mapValue[column].(driver.Valuer); ok {
				mapValue[column], _ = valuer.Value()
			} else if b, ok := mapValue[column].(sql.RawBytes); ok {
				mapValue[column] = string(b)
			}
		} else {
			mapValue[column] = nil
		}
	}
}