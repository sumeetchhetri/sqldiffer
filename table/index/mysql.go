package index

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "sqldiffer/protos"
	"strings"
)

//MysqlIndex -
type MysqlIndex struct {
}

//GenerateNew -
func (db *MysqlIndex) GenerateNew(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE ")
	var typ, uniq string
	v, ok := in.Props["Arg1"]
	if ok {
		typ = v
	}
	v, ok = in.Props["Arg2"]
	if ok {
		uniq = v
	}
	if typ == "FULLTEXT" || typ == "SPATIAL" {
		b.WriteString(typ)
		b.WriteString(" ")
		typ = ""
	}
	if uniq == "Unique" {
		b.WriteString(uniq)
		b.WriteString(" ")
		typ = ""
	}
	b.WriteString("INDEX ")
	b.WriteString(*in.Name)
	if typ != "" {
		b.WriteString("USING ")
		b.WriteString(v)
		b.WriteString(" ")
	}
	b.WriteString("ON ")
	b.WriteString(*in.TableName)
	b.WriteString("(")
	b.WriteString(strings.Join(in.Columns, ","))
	b.WriteString(");\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlIndex) GenerateUpd(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(in, context))
	b.WriteString(db.GenerateNew(in, context))
	return b.String()
}

//GenerateDel -
func (db *MysqlIndex) GenerateDel(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP INDEX ")
	b.WriteString(*in.Name)
	b.WriteString(";\n")
	return b.String()
}

//CountQuery -
func (db *MysqlIndex) CountQuery(context interface{}) string {
	return ""
}

//Query -
func (db *MysqlIndex) Query(context interface{}) string {
	return fmt.Sprintf(`select table_name,
				index_name,
				'',
				group_concat(column_name order by seq_in_index) as index_columns,
				index_type,
				case non_unique
					when 1 then ''
					else 'Unique'
					end as is_unique
			from information_schema.statistics
			where table_schema = '%s'
			group by index_schema,
				index_name,
				index_type,
				non_unique,
				table_name
			order by index_schema,
				index_name;`, context.(string))
}

//FromResult -
func (db *MysqlIndex) FromResult(rows *sql.Rows, context interface{}) *pb2.Index {
	return c.GetIndexFromRow(rows, context)
}
