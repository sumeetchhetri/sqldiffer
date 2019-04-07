package storedproc

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	"strings"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//MysqlStoredProcedure -
type MysqlStoredProcedure struct {
	SchemaName string
}

//GenerateNew -
func (db *MysqlStoredProcedure) GenerateNew(sp *pb2.StoredProcedure, context interface{}) string {
	var b strings.Builder
	b.WriteString("\nDELIMITER $$\nCREATE PROCEDURE ")
	b.WriteString(*sp.Name)
	b.WriteString(" (")
	ps := ""
	for _, p := range sp.Params {
		ps += *p.Mode
		ps += " "
		ps += *p.Name
		ps += " "
		ps += *p.Type
		ps += ","
	}
	ps = strings.TrimSuffix(ps, ",")
	b.WriteString(ps)
	b.WriteString(")\n")
	b.WriteString(*sp.Definition)
	b.WriteString("\nEND $$\nDELIMETER ;\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlStoredProcedure) GenerateUpd(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(sp, context))
	b.WriteString(*sp.Definition)
	return b.String()
}

//GenerateDel -
func (db *MysqlStoredProcedure) GenerateDel(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP PROECURE ")
	b.WriteString(*sp.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *MysqlStoredProcedure) Query(context interface{}) string {
	dbName := context.(string)
	return fmt.Sprintf(`select ROUTINE_NAME,ROUTINE_DEFINITION,'','',0 from information_schema.routines
		where routine_schema = '%s'`, dbName)
}

//FromResult -
func (db *MysqlStoredProcedure) FromResult(rows *sql.Rows, context interface{}) *pb2.StoredProcedure {
	return c.GetProcedureFromRow(rows, context)
}

//DefineQuery -
func (db *MysqlStoredProcedure) DefineQuery(context interface{}) string {
	return ""
}

//Definition -
func (db *MysqlStoredProcedure) Definition(rows *sql.Rows) string {
	return ""
}
