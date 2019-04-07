package constraint

import (
	"bytes"
	sql "database/sql"
	"fmt"
	"regexp"
	c "sqldiffer/common"
	"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclConstraint -
type OrclConstraint struct {
}

//GenerateNew -
func (db *OrclConstraint) GenerateNew(cn *pb2.Constraint, context interface{}) string {
	if *cn.Type == "P" || *cn.Type == "U" || *cn.Type == "R" {
		ats := *cn.Definition
		re := regexp.MustCompile("([a-zA-Z]+)([\t\n ]+)ALTER TABLE")
		if re.MatchString(ats) {
			ats = re.ReplaceAllString(ats, "$1\n\nALTER TABLE")
			ats = strings.ReplaceAll(ats, "\n\n", ";\n/\n")
			ats = ats + ";\n/"
		} else {
			ats = ats + ";\n/"
		}
		//fmt.Println(ats)
		return ats
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

//Query -
func (db *OrclConstraint) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT outer.* FROM (
		select cols.table_name, cons.constraint_name, 
			CASE cons.constraint_type
				WHEN 'P' THEN DBMS_METADATA.GET_DDL('CONSTRAINT', cons.constraint_name)
				WHEN 'U' THEN DBMS_METADATA.GET_DDL('CONSTRAINT', cons.constraint_name)
				--WHEN 'C' THEN DBMS_METADATA.GET_DDL('CONSTRAINT', cons.constraint_name)
				WHEN 'R' THEN DBMS_METADATA.get_dependent_ddl('REF_CONSTRAINT', cols.table_name)
			END AS definition,
			COALESCE(cons.SEARCH_CONDITION_VC,' '), cons.constraint_type, cols.column_name,rownum rn,'' a,'' b,'' c
					FROM user_constraints cons INNER join user_cons_columns cols ON cons.constraint_name = cols.constraint_name 
					AND cons.owner = cols.owner AND cons.status = 'ENABLED' 
					ORDER BY cols.table_name, cols.POSITION
		) outer where outer.rn >= %d and outer.rn < %d`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *OrclConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetConstraintFromRow(rows, context)
}
