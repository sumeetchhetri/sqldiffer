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

//OrclTable -
type OrclTable struct {
}

//GenerateNew -
func (db OrclTable) GenerateNew(tb *pb2.Table, context interface{}) string {
	var b strings.Builder
	if !*tb.IsTemp {
		b.WriteString("\nCREATE TABLE \"")
	} else {
		b.WriteString("\nCREATE GLOBAL TEMPORARY TABLE \"")
	}
	b.WriteString(*tb.Name)
	b.WriteString("\" (")
	f := false
	for _, c := range tb.Columns {
		b.WriteString("\n\t")
		b.WriteString("\"")
		b.WriteString(*c.Name)
		b.WriteString("\"")
		b.WriteString(" ")
		b.WriteString(*c.Type)
		b.WriteString(" ")
		if c.Notnull != nil && *c.Notnull == true {
			//b.WriteString("NOT NULL ")
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
		if !*tb.IsTemp {
			t += "\n);\n/"
		} else {
			t += "\n) ON COMMIT DELETE ROWS;\n/"
		}
		return t
	}
	if !*tb.IsTemp {
		b.WriteString("\n);\n/")
	} else {
		b.WriteString("\n) ON COMMIT DELETE ROWS;\n/")
	}
	return b.String()
}

//GenerateUpd -
func (db OrclTable) GenerateUpd(tb *pb2.Table, context interface{}) string {
	return ""
}

//GenerateDel -
func (db OrclTable) GenerateDel(tb *pb2.Table, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP TABLE ")
	b.WriteString(*tb.Name)
	b.WriteString(";\n/")
	return b.String()
}

//Query -
func (db OrclTable) Query(context interface{}) string {
	return "SELECT table_name,temporary from user_tables"
}

//FromResult -
func (db OrclTable) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetTableFromRow(rows, context)
}
