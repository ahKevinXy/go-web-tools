package gorm_tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"strings"
	"time"
)

type BaseModel struct {
	Id        int64                 `json:"id,string"`                              //`json:"id,string"`                                   //主键ID
	CreatedAt int64                 `json:"created_at"`                             //创建时间
	UpdatedAt int64                 `json:"updated_at"`                             //更新时间
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_name"` //软删除
}

// InsertToTable 通用插入数据库
func InsertToTable(tableName string, data interface{}, db *gorm.DB) error {

	err := db.Table(tableName).Create(data).Error

	if err != nil {
		//logs.Errorf("创建表:%+v,数据:%+v;失败:%+v", tableName, data, err)
		return err
	}
	return nil
}

// DelByTableNameAndId 删除 通过主键删除
func DelByTableNameAndId(tableName string, ids []int64, db *gorm.DB) error {

	err := db.Table(tableName).Where("id in (?)", ids).Delete(nil).Error

	if err != nil {
		//	logs.Errorf("删除表:%+v;id:%v;失败:%+v", tableName, ids, err)
		return err
	}

	return nil

}

// UpdateByTableNameAndId 通过 主键更新
func UpdateByTableNameAndId(tableName string, id int64, updates map[string]interface{}, db *gorm.DB) error {

	err := db.Table(tableName).Where("id=?", id).Updates(updates).Error

	if err != nil {
		//logs.Errorf("删除表:%+v;id:%v;失败:%+v", tableName, id, err)
		return err
	}
	return nil
}

func UpdateByTableNameAndIds(tableName string, ids []int64, updates map[string]interface{}, db *gorm.DB) error {

	err := db.Table(tableName).Where("id in (?)", ids).Updates(updates).Error

	if err != nil {
		//logs.Errorf("删除表:%+v;id:%v;失败:%+v", tableName, ids, err)
		return err
	}
	return nil
}

func UpdateByTableNameAndWhere(tableName string, where, updates map[string]interface{}, db *gorm.DB) (err error) {

	err = db.Table(tableName).Where(where).Updates(updates).Error

	if err != nil {
		//logs.Errorf("删除表:%+v;%v失败:%+v", tableName, where, err)
		return err
	}

	return nil

}

// SoftDelByTableNameAndId
// @Description:   软删除
// @Author ahKevinXy
// @Date 2023-03-29 11:58:04
func SoftDelByTableNameAndId(tableName string, id int64, db *gorm.DB) error {

	updates := map[string]interface{}{}

	updates["updated_at"] = time.Now().Unix()
	updates["deleted_at"] = time.Now().Unix()
	err := db.Table(tableName).Where("id=?", id).Updates(updates).Error
	if err != nil {
		//logs.Errorf("软删除失败 table:%s,%d,err:%+v", tableName, id, err)
		return err
	}
	return nil
}

// SplitInsertTableRecords 分片插入 数据库 需要 事务
func SplitInsertTableRecords(db *gorm.DB, table string, records ...interface{}) error {
	if len(records) == 0 {
		return nil
	}

	// 获取record长度
	recordMap, err := dumpStructToMapForInsert(records[0], table)
	if err != nil {
		return err
	}
	if len(recordMap) == 0 {
		return nil
	}
	// 由于数据库占位符长度限制,需要动态限制
	maxInsertLen := 65000 / len(recordMap)

	for len(records) != 0 {
		insertLen := len(records)
		if len(records) > maxInsertLen {
			insertLen = maxInsertLen
		}

		err := BatchInsertTableRecords(db, table, records[:insertLen]...)
		if err != nil {
			return err
		}
		records = records[insertLen:]
	}

	return nil
}

//BatchInsertTableRecords 批量插入表单记录
func BatchInsertTableRecords(db *gorm.DB, table string, records ...interface{}) error {
	if len(records) == 0 {
		return nil
	}

	//全部转换成map
	recordMap := make([]map[string]interface{}, len(records))
	for index, item := range records {
		recordMapItem, err := dumpStructToMapForInsert(item, table)
		if err != nil {
			return fmt.Errorf("BatchInsertTableRecords failed to DumpSructToMap for %s, err:%s", table, err.Error())
		}
		recordMap[index] = recordMapItem
	}

	//以一定的顺序收集字段名字
	fieldName := make([]string, 0, len(recordMap[0]))
	for key := range recordMap[0] {
		fieldName = append(fieldName, key)
	}

	//按上面收集到的字段名字顺序，查找对应的value
	values := make([]interface{}, 0, len(records))
	for _, item := range recordMap {
		valueItem := make([]interface{}, 0, len(item))
		for _, name := range fieldName {
			valueItem = append(valueItem, item[name])
		}

		values = append(values, valueItem)
	}

	for index, field := range fieldName {
		fieldName[index] = "`" + field + "`"
	}

	sql := fmt.Sprintf("insert into %s(%s) values ", table, strings.Join(fieldName, ","))
	for range values {
		sql += "(?),"
	}
	sql = strings.TrimSuffix(sql, ",")

	dbRet := db.Exec(sql, values...)
	errString := fmt.Sprintf("%v", dbRet.Error)
	if strings.Contains(errString, "Duplicate entry") {
		return fmt.Errorf("duplicate record exists")
	}
	if dbRet.Error != nil {
		return fmt.Errorf("BatchInsertTableRecords failed to insert to %s records, err:%v", table, dbRet.Error)
	}
	if dbRet.RowsAffected == 0 {
		return fmt.Errorf("BatchInsertTableRecords insert to %s record row affected is zero", table)
	}

	return nil
}

// 分片切入
func dumpStructToMapForInsert(dbStruct interface{}, table string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	data, err := json.Marshal(dbStruct)
	if err != nil {
		return nil, fmt.Errorf("DumpStructToMap failed to json.Marshal, err:%v", err)
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("DumpStructToMap failed to json.Unmarshal, err:%v", err)
	}

	return result, nil
}
