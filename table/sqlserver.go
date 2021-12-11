package table

import (
	"bytes"
	sql "database/sql"
	c "github.com/sumeetchhetri/sqldiffer/common"
	"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//SqlsTable -
type SqlsTable struct {
}

//GenerateNew -
func (db *SqlsTable) GenerateNew(tb *pb2.Table, context interface{}) string {
	var b strings.Builder
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString("CREATE TABLE [")
	b.WriteString(*tb.Name)
	b.WriteString("](")
	for _, c := range tb.Columns {
		b.WriteString("\n\t[")
		b.WriteString(*c.Name)
		b.WriteString("] ")
		b.WriteString(strings.Replace(*c.Type, "-1", "MAX", 1))
		b.WriteString(" ")
		if c.Notnull != nil && *c.Notnull == true {
			b.WriteString("NOT NULL")
		} else {
			b.WriteString("NULL")
		}
		/*if c.DefVal != nil && strings.TrimSpace(*c.DefVal) != "" {
			b.WriteString("DEFAULT ")
			b.WriteString(*c.DefVal)
		}*/
		b.WriteString(",")
	}
	t := strings.TrimSuffix(b.String(), ",")
	return t + "\n);\n"
}

//GenerateUpd -
func (db *SqlsTable) GenerateUpd(tb *pb2.Table, context interface{}) string {
	return ""
}

//GenerateDel -
func (db *SqlsTable) GenerateDel(tb *pb2.Table, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString("DROP TABLE [")
	b.WriteString(*tb.Name)
	b.WriteString("];\n")
	return b.String()
}

//Query -
func (db *SqlsTable) Query(context interface{}) string {
	return "SELECT table_name,'Y' FROM information_schema.tables where table_type = 'BASE TABLE' and table_name <> 'sysdiagrams'"
}

//FromResult -
func (db *SqlsTable) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetTableFromRow(rows, context)
}
