package column

import (
	"bytes"
	sql "database/sql"
	"fmt"
	"regexp"
	"strings"

	c "github.com/sumeetchhetri/sqldiffer/common"

	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//PgColumn -
type PgColumn struct {
}

//GenerateNew -
func (db *PgColumn) GenerateNew(co *pb2.Column, context interface{}) string {
	var b strings.Builder
	b.WriteString("\nALTER TABLE ")
	b.WriteString(*co.TableName)
	b.WriteString(" ADD COLUMN ")
	rgxp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	if rgxp.MatchString(*co.Name) {
		b.WriteString(*co.Name)
	} else {
		b.WriteString("\"")
		b.WriteString(*co.Name)
		b.WriteString("\"")
	}
	b.WriteString(" ")
	b.WriteString(*co.Type)
	if !*co.Notnull {
		b.WriteString(" NOT NULL ")
	}
	if co.DefVal != nil && *co.DefVal != "" {
		b.WriteString(" DEFAULT ")
		b.WriteString(*co.DefVal)
	}
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgColumn) GenerateUpd(co *pb2.Column, context interface{}) string {
	var b strings.Builder
	dstcol := context.(*pb2.Column)
	rgxp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	if dstcol != nil && *dstcol.Type != *co.Type {
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*co.TableName)
		b.WriteString(" ALTER COLUMN ")
		if rgxp.MatchString(*co.Name) {
			b.WriteString(*co.Name)
		} else {
			b.WriteString("\"")
			b.WriteString(*co.Name)
			b.WriteString("\"")
		}
		b.WriteString(" TYPE ")
		b.WriteString(*co.Type)
		b.WriteString(" USING ")
		b.WriteString(*co.Name)
		t := *co.Type
		if strings.Contains(t, "(") {
			t = t[0:strings.Index(t, "(")]
		}
		b.WriteString("::")
		b.WriteString(t)
		b.WriteString(";\n")
	} else if dstcol != nil && ((co.DefVal != nil && dstcol.DefVal != nil && *co.DefVal != *dstcol.DefVal) ||
		(dstcol.DefVal != nil && co.DefVal == nil)) {
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*co.TableName)
		b.WriteString(" ALTER COLUMN ")
		if rgxp.MatchString(*co.Name) {
			b.WriteString(*co.Name)
		} else {
			b.WriteString("\"")
			b.WriteString(*co.Name)
			b.WriteString("\"")
		}
		b.WriteString(" SET DEFAULT ")
		b.WriteString(*co.DefVal)
		b.WriteString(";\n")
	} else if dstcol != nil && dstcol.DefVal == nil {
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*co.TableName)
		b.WriteString(" ALTER COLUMN ")
		if rgxp.MatchString(*co.Name) {
			b.WriteString(*co.Name)
		} else {
			b.WriteString("\"")
			b.WriteString(*co.Name)
			b.WriteString("\"")
		}
		b.WriteString(" DROP DEFAULT ")
		//b.WriteString(*co.DefVal)
		b.WriteString(";\n")
	}
	return b.String()
}

//GenerateDel -
func (db *PgColumn) GenerateDel(co *pb2.Column, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE ")
	b.WriteString(*co.TableName)
	b.WriteString(" DROP COLUMN ")
	rgxp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	if rgxp.MatchString(*co.Name) {
		b.WriteString(*co.Name)
	} else {
		b.WriteString("\"")
		b.WriteString(*co.Name)
		b.WriteString("\"")
	}
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *PgColumn) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`
		SELECT ordinal_position AS position, 
			column_name, 
			CASE 
				WHEN is_nullable = 'YES' THEN 'Y' 
				ELSE 'N' 
			END              AS nullable, 
			data_type, 
			CASE 
				WHEN character_maximum_length IS NOT NULL THEN character_maximum_length 
				ELSE numeric_precision 
			END              AS data_length, 
			CASE 
				WHEN character_maximum_length IS NULL THEN NULL 
				ELSE numeric_precision_radix 
			END              AS data_precision, 
			CASE 
				WHEN character_maximum_length IS NULL THEN NULL 
				ELSE numeric_scale 
			END              AS data_scale, 
			column_default   AS default_value, 
			table_name, 0
		FROM   information_schema.columns 
		WHERE  table_schema = ANY ( current_schemas(false) ) 
		ORDER  BY table_name, ordinal_position limit %d offset %d`, args[1].(int), args[0].(int))
}

//FromResult -
func (db *PgColumn) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetColumnFromRow(rows, context)
}
