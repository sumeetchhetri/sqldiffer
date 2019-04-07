package sequence

import (
	"bytes"
	sql "database/sql"
	//"fmt"
	proto "github.com/golang/protobuf/proto"
	c "sqldiffer/common"
	pb2 "sqldiffer/protos"
	//"strings"
)

//PgSequence -
type PgSequence struct {
}

//GenerateNew -
func (db *PgSequence) GenerateNew(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	//TODO
	b.WriteString("\nCREATE SEQUENCE ")
	b.WriteString(*seq.Name)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgSequence) GenerateUpd(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(seq, context))
	b.WriteString(db.GenerateNew(seq, context))
	return b.String()
}

//GenerateDel -
func (db *PgSequence) GenerateDel(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP SEQUENCE ")
	b.WriteString(*seq.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *PgSequence) Query(context interface{}) string {
	return `SELECT sequence_name, 
			start_value,
			CASE 
				WHEN cycle_option = 'YES' THEN 'Y' 
				ELSE 'N' 
			END AS cycle, 
			minimum_value, 
			maximum_value, 
			INCREMENT, 
			NULL, 
			NULL, 
			NULL 
		FROM   information_schema.sequences`
}

//FromResult -
func (db *PgSequence) FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence {
	return c.GetSequenceFromRow(rows)
}

//ExQuery -
func (db *PgSequence) ExQuery(name string) string {
	//return "select cache_value from " + name
	return ""
}

//UpdateSequence -
func (db *PgSequence) UpdateSequence(rows *sql.Rows, s *pb2.Sequence) {
	var out string
	err := rows.Scan(&out)
	if err != nil {
		c.Fatal("Error fetching sequence cache value details from database", err)
	}

	s.Cache = proto.String(out)
}
