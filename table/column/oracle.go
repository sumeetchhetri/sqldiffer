package column

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
	//"strconv"
)

//OrclColumn -
type OrclColumn struct {
	Column *pb2.Column
}

//SetColumn -
func (db *OrclColumn) SetColumn(col *pb2.Column) {
}

//GenerateNew -
func (db *OrclColumn) GenerateNew(co *pb2.Column, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE \"")
	b.WriteString(*co.TableName)
	b.WriteString("\" ADD ")
	b.WriteString("\"")
	b.WriteString(*co.Name)
	b.WriteString("\"")
	b.WriteString(" ")
	b.WriteString(*co.Type)
	if !*co.Notnull {
		//b.WriteString(" NOT NULL ");
	}
	if db.Column.DefVal != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(*co.DefVal)
	}
	b.WriteString(";\n/")
	return b.String()
}

//GenerateUpd -
func (db *OrclColumn) GenerateUpd(co *pb2.Column, context interface{}) string {
	var b bytes.Buffer
	dstcol := context.(pb2.Column)
	if *dstcol.Type != *co.Type {
		b.WriteString("\nALTER TABLE \"")
		b.WriteString(*co.TableName)
		b.WriteString("\" MODIFY ")
		b.WriteString("\"")
		b.WriteString(*co.Name)
		b.WriteString("\"")
		b.WriteString(" ")
		b.WriteString(*co.Type)
		b.WriteString(";\n/")
	} else if co.DefVal != nil && *dstcol.DefVal != *co.DefVal {
		b.WriteString("\nALTER TABLE \"")
		b.WriteString(*co.TableName)
		b.WriteString("\" MODIFY ")
		b.WriteString("\"")
		b.WriteString(*co.Name)
		b.WriteString("\"")
		b.WriteString(" DEFAULT ")
		b.WriteString(*co.DefVal)
		b.WriteString(";\n/")
	} else if dstcol.DefVal != nil && *dstcol.DefVal != *co.DefVal {
		b.WriteString("\nALTER TABLE \"")
		b.WriteString(*co.TableName)
		b.WriteString("\" MODIFY ")
		b.WriteString("\"")
		b.WriteString(*co.Name)
		b.WriteString("\"")
		if "null" == *co.DefVal {
			b.WriteString(" NOT NULL ENABLE ")
		} else {
			b.WriteString(" DROP DEFAULT ")
			b.WriteString(*co.DefVal)
		}
		b.WriteString(";\n/")
	}
	return b.String()
}

//GenerateDel -
func (db *OrclColumn) GenerateDel(co *pb2.Column, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE \"")
	b.WriteString(*co.TableName)
	b.WriteString("\" DROP COLUMN ")
	b.WriteString("\"")
	b.WriteString(*co.Name)
	b.WriteString("\"")
	b.WriteString(";\n/")
	return b.String()
}

//Query -
func (db *OrclColumn) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT column_id,column_name,nullable,data_type,data_length,data_precision,data_scale,data_default,table_name,0 rn 
		from user_tab_columns offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *OrclColumn) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetColumnFromRow(rows, context)
}
