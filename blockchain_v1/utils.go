package blockchain_v1

import (
	"bytes"
	"encoding/binary"
	"github.com/astaxie/beego"
)

func Int2Bytes(i int64) []byte{
	var buffer bytes.Buffer
	// func Write(w io.Writer, order ByteOrder, data interface{}) error {
	// 采用大端对齐把i写到buffer里
	err := binary.Write(&buffer, binary.BigEndian, i)
	CheckErr("Int2Bytes",err)
	return buffer.Bytes()
}
// 校验error的函数
// 参数(错误位置, 错误)
func CheckErr(position string, err error) {
	if err != nil {
		beego.Info("pos, error :",position,err)
		return
	}
}
