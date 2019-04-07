package sequence

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclSequence -
type OrclSequence struct {
}

//GenerateNew -
func (db *OrclSequence) GenerateNew(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE SEQUENCE \"")
	b.WriteString(*seq.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//GenerateUpd -
func (db *OrclSequence) GenerateUpd(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(seq, context))
	b.WriteString(db.GenerateNew(seq, context))
	return b.String()
}

//GenerateDel -
func (db *OrclSequence) GenerateDel(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP SEQUENCE \"")
	b.WriteString(*seq.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//Query -
func (db *OrclSequence) Query(context interface{}) string {
	return `SELECT sequence_name,LAST_NUMBER,cycle_flag,min_value,max_value,increment_by,
			CACHE_SIZE,ORDER_FLAG,null FROM user_sequences`
}

//FromResult -
func (db *OrclSequence) FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence {
	return c.GetSequenceFromRow(rows)
}

//ExQuery -
func (db *OrclSequence) ExQuery(name string) string {
	return ""
}

//UpdateSequence -
func (db *OrclSequence) UpdateSequence(rows *sql.Rows, s *pb2.Sequence) {
}
