package utils

import (
	"encoding/csv"
	"io"
	"os"
)

// ReadCsv
//  @Description:   读取CSV文件
//  @param path
//  @return *[][]string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:52:03
func ReadCsv(path string) *[][]string {
	fs, err := os.Open(path)

	if err != nil {

		return nil
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	content, err := r.ReadAll()

	if err != nil {

		return nil
	}
	return &content
}

// WriteCsv
//  @Description: 写入csv
//  @param path
//  @param content
//  @return error
//  @Author  ahKevinXy
//  @Date2023-04-04 14:52:34
func WriteCsv(path string, content []string) error {

	nfs, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer nfs.Close()

	nfs.Seek(0, io.SeekEnd)

	w := csv.NewWriter(nfs)
	w.Comma = ','
	w.UseCRLF = true
	row := []string{"1", "2", "3", "4", "5,6"}
	err = w.Write(row)
	if err != nil {
		return err
	}
	//这里必须刷新，才能将数据写入文件。
	w.Flush()

	return nil
}
