package sequence

import (
	sql "database/sql"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//MysqlSequence -
type MysqlSequence struct {
}

//GenerateNew -
func (db *MysqlSequence) GenerateNew(seq *pb2.Sequence, context interface{}) string {
	return ""
}

//GenerateUpd -
func (db *MysqlSequence) GenerateUpd(seq *pb2.Sequence, context interface{}) string {
	return ""
}

//GenerateDel -
func (db *MysqlSequence) GenerateDel(seq *pb2.Sequence, context interface{}) string {
	return ""
}

//Query -
func (db *MysqlSequence) Query(context interface{}) string {
	return ""
}

//FromResult -
func (db *MysqlSequence) FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence {
	return nil
}

//ExQuery -
func (db *MysqlSequence) ExQuery(context interface{}) string {
	return ""
}

//UpdateSequence -
func (db *MysqlSequence) UpdateSequence(rows *sql.Rows, s *pb2.Sequence) {
}
