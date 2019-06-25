package constraint

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "sqldiffer/protos"
	"strings"
)

//SqlsConstraint -
type SqlsConstraint struct {
}

//GenerateNew -
func (db *SqlsConstraint) GenerateNew(cn *pb2.Constraint, context interface{}) string {
	var b strings.Builder
	b.WriteString("\nGO\nALTER TABLE [")
	b.WriteString(*cn.TableName)
	b.WriteString("] ADD CONSTRAINT [")
	b.WriteString(*cn.Name)
	b.WriteString("] ")
	if *cn.Type == "Primary key" {
		b.WriteString("PRIMARY KEY([")
		b.WriteString(strings.Join(cn.Columns, "],["))
		b.WriteString("])")
	} else if *cn.Type == "Unique constraint" {
		b.WriteString("UNIQUE([")
		b.WriteString(strings.Join(cn.Columns, "],["))
		b.WriteString("])")
	} else if *cn.Type == "Foreign key" {
		b.WriteString("FOREIGN KEY(")
		b.WriteString(cn.Columns[0])
		b.WriteString(") REFERENCES ")
		b.WriteString(*cn.TargetTableName)
		b.WriteString("(")
		b.WriteString(*cn.TargetColumnName)
		b.WriteString(")")
	} else if *cn.Type == "Check constraint" {
		b.WriteString("CHECK")
		b.WriteString(*cn.Condition)
	} else if *cn.Type == "Default constraint" {
		b.WriteString("DEFAULT ")
		b.WriteString(*cn.Condition)
		b.WriteString(" FOR ")
		b.WriteString(cn.Columns[0])
	} /*else if *cn.Type == "Unique clustered index" {

	} else if *cn.Type == "Unique index" {

	}*/
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *SqlsConstraint) GenerateUpd(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(cn, context))
	b.WriteString(db.GenerateNew(cn, context))
	return b.String()
}

//GenerateDel -
func (db *SqlsConstraint) GenerateDel(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nALTER TABLE [")
	b.WriteString(*cn.TableName)
	b.WriteString("] DROP CONSTRAINT [")
	b.WriteString(*cn.Name)
	b.WriteString("];\n")
	return b.String()
}

//CountQuery -
func (db *SqlsConstraint) CountQuery(context interface{}) string {
	return ""
}

//Query -
func (db *SqlsConstraint) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`select table_view,
			constraint_name, '' a,details,
			constraint_type,coln,0,tgt,tgc,'' d
		from (
			select t.[name] as table_view, 
				case when t.[type] = 'U' then 'Table'
					when t.[type] = 'V' then 'View'
					end as [object_type],
				case when c.[type] = 'PK' then 'Primary key'
					when c.[type] = 'UQ' then 'Unique constraint'
					when i.[type] = 1 then 'Unique clustered index'
					when i.type = 2 then 'Unique index'
					end as constraint_type, 
				isnull(c.[name], i.[name]) as constraint_name,
				'' [details],
				substring(column_names, 1, len(column_names)-1) as coln,
				'' tgt,'' tgc
			from sys.objects t
				left outer join sys.indexes i
					on t.object_id = i.object_id
				left outer join sys.key_constraints c
					on i.object_id = c.parent_object_id 
					and i.index_id = c.unique_index_id
			cross apply (select col.[name] + ', '
								from sys.index_columns ic
									inner join sys.columns col
										on ic.object_id = col.object_id
										and ic.column_id = col.column_id
								where ic.object_id = t.object_id
									and ic.index_id = i.index_id
										order by col.column_id
										for xml path ('') ) D (column_names)
			where is_unique = 1 and t.[type] = 'U'
			and t.is_ms_shipped <> 1
			union all 
			SELECT  tab1.name AS [table],
				'Table',
				'Foreign key',
				obj.name AS FK_NAME,
				'',
				col1.name AS [coln],
				tab2.name AS [referenced_table],
				col2.name AS [referenced_column]
			FROM sys.foreign_key_columns fkc
			INNER JOIN sys.objects obj
				ON obj.object_id = fkc.constraint_object_id
			INNER JOIN sys.tables tab1
				ON tab1.object_id = fkc.parent_object_id
			INNER JOIN sys.schemas sch
				ON tab1.schema_id = sch.schema_id
			INNER JOIN sys.columns col1
				ON col1.column_id = parent_column_id AND col1.object_id = tab1.object_id
			INNER JOIN sys.tables tab2
				ON tab2.object_id = fkc.referenced_object_id
			INNER JOIN sys.columns col2
				ON col2.column_id = referenced_column_id AND col2.object_id = tab2.object_id
			union all
			select t.[name],
				'Table',
				'Check constraint',
				con.[name] as constraint_name, 
				con.[definition],'' coln,'' tgt,'' tgc
			from sys.check_constraints con
				left outer join sys.objects t
					on con.parent_object_id = t.object_id
				left outer join sys.all_columns col
					on con.parent_column_id = col.column_id
					and con.parent_object_id = col.object_id
			union all
			select t.[name],
				'Table',
				'Default constraint',
				con.[name],
				con.[definition],col.[name] coln, '' tgt,'' tgc
			from sys.default_constraints con
				left outer join sys.objects t
					on con.parent_object_id = t.object_id
				left outer join sys.all_columns col
					on con.parent_column_id = col.column_id
					and con.parent_object_id = col.object_id) a
		where table_view <> 'sysdiagrams'
		order by table_view, constraint_type, constraint_name
		offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *SqlsConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetConstraintFromRow(rows, context)
}
