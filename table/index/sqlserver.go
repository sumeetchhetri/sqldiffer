package index

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//SqlsIndex -
type SqlsIndex struct {
}

//GenerateNew -
func (db *SqlsIndex) GenerateNew(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nCREATE INDEX ")
	b.WriteString(*in.Name)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *SqlsIndex) GenerateUpd(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(in, context))
	b.WriteString(db.GenerateNew(in, context))
	return b.String()
}

//GenerateDel -
func (db *SqlsIndex) GenerateDel(in *pb2.Index, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nDROP INDEX ")
	b.WriteString(*in.Name)
	b.WriteString(";\n")
	return b.String()
}

//CountQuery -
func (db *SqlsIndex) CountQuery(context interface{}) string {
	return ""
}

//Query -
func (db *SqlsIndex) Query(context interface{}) string {
	return `select t.[name] as table_view, i.[name] as index_name,'',
			substring(column_names, 1, len(column_names)-1) as [columns],
			case when i.[type] = 1 then 'Clustered index'
				when i.[type] = 2 then 'Nonclustered unique index'
				when i.[type] = 3 then 'XML index'
				when i.[type] = 4 then 'Spatial index'
				when i.[type] = 5 then 'Clustered columnstore index'
				when i.[type] = 6 then 'Nonclustered columnstore index'
				when i.[type] = 7 then 'Nonclustered hash index'
				end as index_type,
			case when i.is_unique = 1 then 'Unique'
				else 'Not unique' end as [unique]
		from sys.objects t
			inner join sys.indexes i
				on t.object_id = i.object_id
			cross apply (select col.[name] + ', '
							from sys.index_columns ic
								inner join sys.columns col
									on ic.object_id = col.object_id
									and ic.column_id = col.column_id
							where ic.object_id = t.object_id
								and ic.index_id = i.index_id
									order by col.column_id
									for xml path ('') ) D (column_names)
		where t.is_ms_shipped <> 1 and t.[type] = 'U' and t.[name] <> 'sysdiagrams'
		and index_id > 0
		order by i.[name]`
}

//FromResult -
func (db *SqlsIndex) FromResult(rows *sql.Rows, context interface{}) *pb2.Index {
	return c.GetIndexFromRow(rows, context)
}
