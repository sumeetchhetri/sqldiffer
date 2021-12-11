package table

import (
	"bytes"
	sql "database/sql"
	"fmt"
	"regexp"
	c "github.com/sumeetchhetri/sqldiffer/common"
	"strings"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//MysqlTable -
type MysqlTable struct {
}

//GenerateNew -
func (db *MysqlTable) GenerateNew(tb *pb2.Table, context interface{}) string {
	var b strings.Builder
	b.WriteString("\nCREATE TABLE ")
	b.WriteString(*tb.Name)
	b.WriteString("(")
	rgxp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	f := false
	for _, c := range tb.Columns {
		b.WriteString("\n\t")
		if rgxp.MatchString(*c.Name) {
			b.WriteString(*c.Name)
		} else {
			b.WriteString("\"")
			b.WriteString(*c.Name)
			b.WriteString("\"")
		}
		b.WriteString(" ")
		b.WriteString(*c.Type)
		b.WriteString(" ")
		if c.Notnull != nil && *c.Notnull == true {
			b.WriteString("NOT NULL ")
		}
		if c.DefVal != nil && strings.TrimSpace(*c.DefVal) != "" {
			b.WriteString("DEFAULT ")
			b.WriteString(*c.DefVal)
		}
		b.WriteString(",")
		f = true
	}
	if f {
		t := strings.TrimSuffix(b.String(), ",")
		t += "\n);\n"
		return t
	}
	b.WriteString("\n);\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlTable) GenerateUpd(tb *pb2.Table, context interface{}) string {
	return ""
}

//GenerateDel -
func (db *MysqlTable) GenerateDel(tb *pb2.Table, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP TABLE ")
	b.WriteString(*tb.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *MysqlTable) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT table_name,'N' FROM information_schema.tables 
		where table_schema = '%s' and table_type = 'BASE TABLE'`, args[1].(string))
}

//FromResult -
func (db *MysqlTable) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetTableFromRow(rows, context)
}
