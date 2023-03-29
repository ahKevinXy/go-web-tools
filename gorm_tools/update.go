package gorm_tools

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// SplitUpdateTableFieldWithoutVersionByID
// @Description:  批量更新
// @Author ahKevinXy
// @Date 2023-03-29 12:01:28
func SplitUpdateTableFieldWithoutVersionByID(db *gorm.DB, tableName string, idUpdateMap map[int64]map[string]interface{}) error {
	maxUpdateLen := 500

	index, totalNum := 0, 0
	sliceUpdate := make([]map[int64]map[string]interface{}, 0, len(idUpdateMap)/maxUpdateLen+1)

	var tmpMap map[int64]map[string]interface{}
	for id, updateMap := range idUpdateMap {
		if index == 0 {
			tmpMap = make(map[int64]map[string]interface{})
		}
		tmpMap[id] = updateMap
		index++
		totalNum++
		if index == maxUpdateLen || totalNum == len(idUpdateMap) {
			index, sliceUpdate = 0, append(sliceUpdate, tmpMap)
		}
	}

	for _, item := range sliceUpdate {
		err := BatchUpdateTableFieldWithoutVersionByID(db, tableName, item)
		if err != nil {
			return err
		}
	}

	return nil
}

// BatchUpdateTableFieldWithoutVersionByID ...
func BatchUpdateTableFieldWithoutVersionByID(db *gorm.DB, tableName string, idUpdateMap map[int64]map[string]interface{}) error {
	if len(idUpdateMap) == 0 {
		return nil
	}

	updateField := make([]string, 0, 10)
	for _, item := range idUpdateMap {
		for field := range item {
			updateField = append(updateField, field)
		}
		break
	}

	caseSQL, beginSQL := "", "update "+tableName+" set "

	args, allID := make([]interface{}, 0, len(idUpdateMap)*2), make([]int64, 0, len(idUpdateMap))
	for _, field := range updateField {
		caseSQL += field + "= case "
		for id, updateMap := range idUpdateMap {
			value, ok := updateMap[field]
			if !ok {
				return fmt.Errorf("update field must equal")
			}
			caseSQL += "when id = ? then ? "
			args = append(args, id, value)
		}
		caseSQL += "end,"
	}
	wholeSQL := beginSQL + strings.TrimSuffix(caseSQL, ",") + " where id in (?)"

	for id := range idUpdateMap {
		allID = append(allID, id)
	}
	args = append(args, allID)

	dbRet := db.Exec(wholeSQL, args...)
	if dbRet.Error != nil {
		return fmt.Errorf("BatchUpdateTableFieldWithoutVersionByID failed to batch update table:%s, err:%v", tableName, dbRet.Error)
	}

	return nil
}
