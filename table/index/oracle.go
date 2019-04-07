package index

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclIndex -
type OrclIndex struct {
}

//GenerateNew -
func (db *OrclIndex) GenerateNew(in *pb2.Index, context interface{}) string {
	return *in.Definition + ";\n/"
}

//GenerateUpd -
func (db *OrclIndex) GenerateUpd(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(in, context))
	b.WriteString(db.GenerateNew(in, context))
	return b.String()
}

//GenerateDel -
func (db *OrclIndex) GenerateDel(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP INDEX \"")
	b.WriteString(*in.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//Query -
func (db *OrclIndex) Query(context interface{}) string {
	return "SELECT table_name,index_name, dbms_metadata.get_ddl('INDEX', index_name),'','','' from user_indexes"
}

//FromResult -
func (db *OrclIndex) FromResult(rows *sql.Rows, context interface{}) *pb2.Index {
	return c.GetIndexFromRow(rows, context)
}
