package column

import (
	sql "database/sql"
	"fmt"

	c "github.com/sumeetchhetri/sqldiffer/common"

	//c "github.com/sumeetchhetri/sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	"strings"

	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

// SqlsColumn -
type SqlsColumn struct {
	Column pb2.Column
}

// GenerateNew -
func (db *SqlsColumn) GenerateNew(co *pb2.Column, context interface{}) string {
	var b strings.Builder
	b.WriteString("\nGO\nALTER TABLE [")
	b.WriteString(*co.TableName)
	b.WriteString("] ADD COLUMN [")
	b.WriteString(*co.Name)
	b.WriteString("] ")
	b.WriteString(*co.Type)
	if *co.Notnull {
		b.WriteString("NOT NULL")
	} else {
		b.WriteString("NULL")
	}
	b.WriteString(";\n")
	return b.String()
}

// GenerateUpd -
func (db *SqlsColumn) GenerateUpd(co *pb2.Column, context interface{}) string {
	var b strings.Builder
	dstcol := context.(*pb2.Column)
	if dstcol != nil && *dstcol.Type != *co.Type {
		b.WriteString("\nGO\nALTER TABLE [")
		b.WriteString(*co.TableName)
		b.WriteString("] ALTER COLUMN [")
		b.WriteString(*co.Name)
		b.WriteString("] ")
		b.WriteString(*co.Type)
		b.WriteString(";\n")
	}
	return b.String()
}

// GenerateDel -
func (db *SqlsColumn) GenerateDel(co *pb2.Column, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nALTER TABLE [")
	b.WriteString(*co.TableName)
	b.WriteString("] DROP COLUMN [")
	b.WriteString(*co.Name)
	b.WriteString("];\n")
	return b.String()
}

// Query -
func (db *SqlsColumn) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`
		SELECT ordinal_position AS position, 
			column_name, 
			CASE 
				WHEN is_nullable = 'YES' THEN 'Y'
				ELSE 'N'
			END              AS notnull, 
			data_type, 
			CASE 
				WHEN character_maximum_length IS NOT NULL THEN character_maximum_length 
				ELSE numeric_precision 
			END              AS data_length, 
			CASE 
				WHEN character_maximum_length IS NOT NULL THEN NULL 
				ELSE numeric_precision_radix 
			END              AS data_precision, 
			CASE 
				WHEN character_maximum_length IS NOT NULL THEN NULL 
				ELSE numeric_scale 
			END              AS data_scale, 
			column_default   AS default_value, 
			table_name, 0
		FROM   information_schema.columns 
		where table_name <> 'sysdiagrams'
		ORDER  BY table_name, ordinal_position
		offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

// FromResult -
func (db *SqlsColumn) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetColumnFromRow(rows, context)
}
