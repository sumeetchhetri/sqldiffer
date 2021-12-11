package constraint

import (
	"bytes"
	sql "database/sql"
	"fmt"
	//"regexp"
	c "github.com/sumeetchhetri/sqldiffer/common"
	"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//OrclConstraint -
type OrclConstraint struct {
}

//GenerateNew -
func (db *OrclConstraint) GenerateNew(cn *pb2.Constraint, context interface{}) string {
	if *cn.Type == "U" {
		var b bytes.Buffer
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*cn.TableName)
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" UNIQUE(")
		b.WriteString(strings.Join(cn.Columns, ","))
		b.WriteString(");\n/")
		return b.String()
	} else if *cn.Type == "P" {
		var b bytes.Buffer
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*cn.TableName)
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" PRIMARY KEY(")
		b.WriteString(strings.Join(cn.Columns, ","))
		b.WriteString(");\n/")
		return b.String()
	} else if *cn.Type == "R" {
		var b bytes.Buffer
		b.WriteString("\nALTER TABLE ")
		b.WriteString(*cn.TableName)
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" FOREIGN KEY(")
		b.WriteString(cn.Columns[0])
		b.WriteString(") REFERENCES ")
		b.WriteString(*cn.TargetTableName)
		b.WriteString("(")
		b.WriteString(*cn.TargetColumnName)
		b.WriteString(");\n/")
		return b.String()
	}
	if strings.Index(*cn.Condition, " IS NOT NULL") > 0 {
		var b bytes.Buffer
		b.WriteString("\nALTER TABLE \"")
		b.WriteString(*cn.TableName)
		b.WriteString("\" MODIFY COLUMN \"")
		b.WriteString(cn.Columns[0])
		b.WriteString("\" NOT NULL ENABLE;\n/")
		return b.String()
	}
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE \"")
	b.WriteString(*cn.TableName)
	b.WriteString("\" ADD CONSTRAINT \"")
	b.WriteString(*cn.Name)
	b.WriteString("\" CHECK(")
	b.WriteString(*cn.Condition)
	b.WriteString(");\n/")
	return b.String()
}

//GenerateUpd -
func (db *OrclConstraint) GenerateUpd(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(cn, context))
	b.WriteString(db.GenerateNew(cn, context))
	return b.String()
}

//GenerateDel -
func (db *OrclConstraint) GenerateDel(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE \"")
	b.WriteString(*cn.Name)
	b.WriteString("\" DROP CONSTRAINT \"")
	b.WriteString(*cn.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//CountQuery -
func (db *OrclConstraint) CountQuery(context interface{}) string {
	return "select count(1) FROM user_constraints"
}

//Query -
func (db *OrclConstraint) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`select * from (
		select cols.table_name, cons.constraint_name,'',
			COALESCE(cons.SEARCH_CONDITION_VC,' '), cons.constraint_type, cols.column_name,0 rn,'' a,'' b,'' c
			FROM user_constraints cons INNER join user_cons_columns cols ON cons.constraint_name = cols.constraint_name 
			AND cons.owner = cols.owner AND cons.status = 'ENABLED' 
			AND cols.table_name IN (SELECT table_name FROM user_tables)
			AND cons.constraint_type <> 'R'
		UNION ALL
		SELECT a.table_name, a.constraint_name, '','','R', a.column_name, 0 rn, c_pk.table_name, b.column_name, ''              
		  FROM user_cons_columns a
		  JOIN user_constraints c ON a.owner = c.owner
								AND a.constraint_name = c.constraint_name
		  JOIN user_constraints c_pk ON c.r_owner = c_pk.owner
								   AND c.r_constraint_name = c_pk.constraint_name
		  JOIN user_cons_columns b ON a.owner = b.owner
								AND c_pk.constraint_name = b.constraint_name
		 WHERE c.constraint_type = 'R') offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *OrclConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetConstraintFromRow(rows, context)
}
