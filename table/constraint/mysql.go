package constraint

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	"strings"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//MysqlConstraint -
type MysqlConstraint struct {
}

//GenerateNew -
func (db *MysqlConstraint) GenerateNew(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE ")
	b.WriteString(*cn.TableName)
	if *cn.Type == "NOT NULL" {
		b.WriteString("MODIFY COLUMN ")
		b.WriteString(cn.Columns[0])
		b.WriteString(" NOT NULL")
	} else if *cn.Type == "UNIQUE" {
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" UNIQUE(")
		b.WriteString(strings.Join(cn.Columns, ","))
		b.WriteString(")")
	} else if *cn.Type == "PRIMARY KEY" {
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" PRIMARY KEY(")
		b.WriteString(strings.Join(cn.Columns, ","))
		b.WriteString(")")
	} else if *cn.Type == "FOREIGN KEY" {
		b.WriteString(" ADD CONSTRAINT ")
		b.WriteString(*cn.Name)
		b.WriteString(" FOREIGN KEY(")
		b.WriteString(cn.Columns[0])
		b.WriteString(") REFERENCES ")
		b.WriteString(*cn.TargetTableName)
		b.WriteString("(")
		b.WriteString(*cn.TargetColumnName)
		b.WriteString(")")
	}
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlConstraint) GenerateUpd(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(cn, context))
	b.WriteString(db.GenerateNew(cn, context))
	return b.String()
}

//GenerateDel -
func (db *MysqlConstraint) GenerateDel(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE \"")
	b.WriteString(*cn.TableName)
	b.WriteString("\" DROP CONSTRAINT \"")
	b.WriteString(*cn.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *MysqlConstraint) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`select c1.TABLE_NAME,c1.CONSTRAINT_NAME,'','',c1.CONSTRAINT_TYPE,
		COLUMN_NAME,0,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME,'' 
		from information_schema.table_constraints c1 inner join information_schema.KEY_COLUMN_USAGE c2
		on c1.CONSTRAINT_NAME = c2.CONSTRAINT_NAME 
		and c1.CONSTRAINT_schema = c2.CONSTRAINT_schema
		where c1.CONSTRAINT_schema = '%s'
		limit %d offset %d`, args[2].(string), args[1].(int), args[0].(int))
}

//FromResult -
func (db *MysqlConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetConstraintFromRow(rows, context)
}
