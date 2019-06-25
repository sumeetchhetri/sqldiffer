package db

import (
	"bytes"
	"fmt"
	//proto "github.com/golang/protobuf/proto"
	c "sqldiffer/common"
	pb2 "sqldiffer/protos"
)

//PgDb -
type PgDb struct {
}

//GenerateURL -
func (db *PgDb) GenerateURL(det *c.SchemaDiffAction) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		*det.Host, det.Port, *det.User, *det.Password, *det.DatabaseName, *det.SchemaName)
}

//Version -
func (db *PgDb) Version() string {
	return "SELECT current_setting('server_version_num')"
}

//Preface -
func (db *PgDb) Preface(dbe *pb2.Db) string {
	var b bytes.Buffer
	b.WriteString("--")
	b.WriteString("\n-- PostgreSQL database dump")
	b.WriteString("\n--")
	b.WriteString("\nSET statement_timeout = 0;")
	b.WriteString("\nSET lock_timeout = 0;")
	b.WriteString("\nSET idle_in_transaction_session_timeout = 0;")
	b.WriteString("\nSET client_encoding = 'SQL_ASCII';")
	b.WriteString("\nSET standard_conforming_strings = on;")
	b.WriteString("\nSET check_function_bodies = false;")
	b.WriteString("\nSET client_min_messages = warning;")
	b.WriteString("\nSET row_security = off;")
	b.WriteString("\n--")
	b.WriteString("\nSET search_path = ")
	b.WriteString(*dbe.SchemaName)
	b.WriteString(", pg_catalog;\n\n")
	return b.String()
}

//Create -
func (db *PgDb) Create(dbe *pb2.Db) string {
	return fmt.Sprintf("\nCREATE DATABASE %s;\n", *dbe.Name)
}

//Connect -
func (db *PgDb) Connect(dbe *pb2.Db) string {
	return fmt.Sprintf("\n\\CONNECT %s;\n", *dbe.Name)
}

//CreateSchema -
func (db *PgDb) CreateSchema(dbe *pb2.Db) string {
	return fmt.Sprintf("\nCREATE SCHEMA %s;\n", *dbe.Name)
}
