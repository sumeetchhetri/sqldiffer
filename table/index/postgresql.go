package index

import (
	"bytes"
	sql "database/sql"

	c "github.com/sumeetchhetri/sqldiffer/common"

	//"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//PgIndex -
type PgIndex struct {
}

//GenerateNew -
func (db *PgIndex) GenerateNew(in *pb2.Index, context interface{}) string {
	return *in.Definition + ";\n"
}

//GenerateUpd -
func (db *PgIndex) GenerateUpd(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(in, context))
	b.WriteString(db.GenerateNew(in, context))
	return b.String()
}

//GenerateDel -
func (db *PgIndex) GenerateDel(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP INDEX ")
	b.WriteString(*in.Name)
	b.WriteString(";\n")
	return b.String()
}

//CountQuery -
func (db *PgIndex) CountQuery(context interface{}) string {
	return ""
}

//Query -
func (db *PgIndex) Query(context interface{}) string {
	return `select tablename,indexname,indexdef,'','','' from pg_indexes where schemaname = ANY(current_schemas(false))`
}

//FromResult -
func (db *PgIndex) FromResult(rows *sql.Rows, context interface{}) *pb2.Index {
	return c.GetIndexFromRow(rows, context)
}
