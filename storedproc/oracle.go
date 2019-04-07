package storedproc

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclStoredProcedure -
type OrclStoredProcedure struct {
	SchemaName string
}

//GenerateNew -
func (db *OrclStoredProcedure) GenerateNew(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\n")
	b.WriteString(*sp.Definition)
	b.WriteString("\n/")
	return b.String()
}

//GenerateUpd -
func (db *OrclStoredProcedure) GenerateUpd(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(sp, context))
	b.WriteString(*sp.Definition)
	return b.String()
}

//GenerateDel -
func (db *OrclStoredProcedure) GenerateDel(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP PROECURE \"")
	b.WriteString(*sp.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//Query -
func (db *OrclStoredProcedure) Query(context interface{}) string {
	return `SELECT object_name,dbms_metadata.get_ddl(object_type, object_name),'','',0
		FROM user_objects WHERE object_type IN ('FUNCTION','PROCEDURE','PACKAGE')`
}

//FromResult -
func (db *OrclStoredProcedure) FromResult(rows *sql.Rows, context interface{}) *pb2.StoredProcedure {
	return c.GetProcedureFromRow(rows, context)
}

//DefineQuery -
func (db *OrclStoredProcedure) DefineQuery(context interface{}) string {
	return ""
}

//Definition -
func (db *OrclStoredProcedure) Definition(rows *sql.Rows) string {
	return ""
}
