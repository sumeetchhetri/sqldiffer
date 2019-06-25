package index

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	"strings"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclIndex -
type OrclIndex struct {
}

//GenerateNew -
func (db *OrclIndex) GenerateNew(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE ")
	var cexp, uniq string
	v, ok := in.Props["Arg1"]
	if ok {
		cexp = v
	}
	v, ok = in.Props["Arg2"]
	if ok {
		uniq = v
	}
	if uniq == "UNIQUE" {
		b.WriteString(uniq)
		b.WriteString(" ")
	}
	b.WriteString("INDEX ")
	b.WriteString(*in.Name)
	b.WriteString(" ON ")
	b.WriteString(*in.TableName)
	b.WriteString("(")
	if cexp != "" {
		b.WriteString(cexp)
	} else {
		b.WriteString(strings.Join(in.Columns, ","))
	}
	b.WriteString(");\n/")
	return b.String()
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

//CountQuery -
func (db *OrclIndex) CountQuery(context interface{}) string {
	return "select count(1) from user_indexes"
}

//Query -
func (db *OrclIndex) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT
			i.table_name,
			i.index_name,
			'',
			c.column_name,
			f.column_expression,
			i.uniqueness
		FROM user_indexes i
		INNER JOIN user_ind_columns c
			ON i.index_name = c.index_name
		LEFT JOIN user_ind_expressions f
			ON c.index_name = f.index_name
			AND c.table_name = f.table_name
			AND c.column_position = f.column_position
		ORDER BY i.table_owner, i.table_name, i.index_name, c.column_position
		offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *OrclIndex) FromResult(rows *sql.Rows, context interface{}) *pb2.Index {
	return c.GetIndexFromRow(rows, context)
}
