package main

import (
	sql "database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"

	_ "github.com/godror/godror"
	c "github.com/sumeetchhetri/sqldiffer/common"
	db "github.com/sumeetchhetri/sqldiffer/db"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
	sq "github.com/sumeetchhetri/sqldiffer/sequence"
	sp "github.com/sumeetchhetri/sqldiffer/storedproc"
	spp "github.com/sumeetchhetri/sqldiffer/storedproc/storedprocparam"
	tb "github.com/sumeetchhetri/sqldiffer/table"
	co "github.com/sumeetchhetri/sqldiffer/table/column"
	cn "github.com/sumeetchhetri/sqldiffer/table/constraint"
	in "github.com/sumeetchhetri/sqldiffer/table/index"
	tr "github.com/sumeetchhetri/sqldiffer/table/trigger"
	vw "github.com/sumeetchhetri/sqldiffer/view"
	"google.golang.org/protobuf/proto"
)

var opts struct {
	DatabaseType   string `short:"t" long:"type" description:"The database type" required:"true" choice:"postgres" choice:"oracle" choice:"mysql" choice:"sqlserver"`
	DatabaseName   string `short:"n" long:"name" description:"The database name" required:"true"`
	Host           string `short:"i" long:"host" description:"The database host" required:"true"`
	Port           int32  `short:"p" long:"port" description:"The database port" required:"true"`
	User           string `short:"u" long:"user" description:"The database user" required:"true"`
	Password       string `short:"w" long:"password" description:"The database password" required:"true"`
	SchemaName     string `short:"s" long:"sch-nam" description:"The database schema name"`
	FileName       string `short:"f" long:"fil-nam" description:"The generated schema/diff file name"`
	StoreProcFiles bool   `long:"gen-proc" description:"Generate stored procedure files"`
}
var fc int64 = 1000000000

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println("Failed to parse args: %!", err)
			os.Exit(2)
		}
	}

	action := c.SchemaDiffAction{
		Host:           &opts.Host,
		Port:           opts.Port,
		DatabaseType:   &opts.DatabaseType,
		DatabaseName:   &opts.DatabaseName,
		User:           &opts.User,
		Password:       &opts.Password,
		SchemaName:     &opts.SchemaName,
		FileName:       &opts.FileName,
		Parallel:       true,
		IsDiffNeeded:   false,
		StoreProcFiles: opts.StoreProcFiles,
	}
	//fmt.Println("%!", action)
	generateSchema(&action)
}

func generateSchema(action *c.SchemaDiffAction) {
	action.Procs = make(map[string]*pb2.StoredProcedure)
	action.Tables = make(map[string]*pb2.Table)
	action.TablesP = make(map[string]bool)
	action.Indexes = make(map[string][]*pb2.Index)
	action.Indexess = make(map[string]*pb2.Index)
	action.Triggers = make(map[string][]*pb2.Trigger)
	action.Constraints = make(map[string]*pb2.Constraint)

	if action.FileName == nil || *action.FileName == "" {
		action.FileName = proto.String("schema_" + (fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))) + ".pbf")
	}

	tdb := pb2.Db{}
	tdb.Name = proto.String(*action.DatabaseName)

	if *action.DatabaseType == "postgres" {
		if action.SchemaName == nil || *action.SchemaName == "" {
			action.SchemaName = proto.String("public")
		}
		tdb.SchemaName = action.SchemaName
		action.Db = &db.PgDb{}
		action.Table = &tb.PgTable{}
		action.Column = &co.PgColumn{}
		action.Constraint = &cn.PgConstraint{}
		action.Index = &in.PgIndex{}
		action.Trigger = &tr.PgTrigger{SchemaName: *action.SchemaName}
		action.StoredProcedure = &sp.PgStoredProcedure{SchemaName: *action.SchemaName}
		action.StoredProcedureParam = &spp.PgStoredProcedureParam{}
		action.Sequence = &sq.PgSequence{}
		action.View = &vw.PgView{}
		action.DuplicateProcNamesAllowed = true
		ps := pb2.Db_POSTGRESQL
		tdb.Type = &ps
	} else if *action.DatabaseType == "oracle" {
		tdb.SchemaName = action.User
		*action.DatabaseType = "oci8"
		action.Db = &db.OrclDb{}
		action.Table = &tb.OrclTable{}
		action.Column = &co.OrclColumn{}
		action.Constraint = &cn.OrclConstraint{}
		action.Index = &in.OrclIndex{}
		action.Trigger = &tr.OrclTrigger{SchemaName: *action.SchemaName}
		action.StoredProcedure = &sp.OrclStoredProcedure{SchemaName: *action.SchemaName}
		action.StoredProcedureParam = &spp.OrclStoredProcedureParam{}
		action.Sequence = &sq.OrclSequence{}
		action.View = &vw.OrclView{}
		action.DuplicateProcNamesAllowed = false
		ps := pb2.Db_ORACLE
		tdb.Type = &ps
	} else if *action.DatabaseType == "mysql" {
		action.Db = &db.MysqlDb{}
		action.Table = &tb.MysqlTable{}
		action.Column = &co.MysqlColumn{}
		action.Constraint = &cn.MysqlConstraint{}
		action.Index = &in.MysqlIndex{}
		action.Trigger = &tr.MysqlTrigger{SchemaName: *action.SchemaName}
		action.StoredProcedure = &sp.MysqlStoredProcedure{SchemaName: *action.SchemaName}
		action.StoredProcedureParam = &spp.MysqlStoredProcedureParam{}
		action.Sequence = &sq.MysqlSequence{}
		action.View = &vw.MysqlView{}
		action.DuplicateProcNamesAllowed = false
		ps := pb2.Db_MYSQL
		tdb.Type = &ps
	} else if *action.DatabaseType == "sqlserver" {
		tdb.SchemaName = action.DatabaseName
		action.Db = &db.SqlsDb{}
		action.Table = &tb.SqlsTable{}
		action.Column = &co.SqlsColumn{}
		action.Constraint = &cn.SqlsConstraint{}
		action.Index = &in.SqlsIndex{}
		action.Trigger = &tr.SqlsTrigger{SchemaName: *action.SchemaName}
		action.StoredProcedure = &sp.SqlsStoredProcedure{SchemaName: *action.SchemaName}
		action.StoredProcedureParam = &spp.SqlsStoredProcedureParam{}
		action.Sequence = &sq.SqlsSequence{}
		action.View = &vw.SqlsView{}
		action.DuplicateProcNamesAllowed = false
		ps := pb2.Db_SQLSERVER
		tdb.Type = &ps
	}

	versQ := action.Db.Version()
	if versQ != "" {
		db := getConn(action)
		rows, err := db.Query(versQ)
		if err != nil {
			c.Fatal("Error querying Database version", err)
		}
		for rows.Next() {
			var vers string
			rows.Scan(&vers)
			if err != nil {
				c.Fatal("Error fetching Database version", err)
			}
			tdb.Version = proto.String(vers)
		}
		db.Close()
	}

	var wg sync.WaitGroup
	wg.Add(4)

	ch := make(chan int)

	go objectifyStoredProcedures(action, &tdb, &wg)
	go objectifyViews(action, &tdb, &wg)
	go objectifyTriggers(action, &tdb, &wg)
	go objectifyTables(action, &tdb, &wg)
	go objectifySequences(action, &tdb, ch)

	wg.Wait()
	<-ch

	mergeTriggersIndexesWithTables(&tdb, action)

	f, err := os.Create(*action.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Change the type to oracle as goracle is only needed for sql
	//driver connection, the value that is persisted is always oracle
	//for oracle database
	if *action.DatabaseType == "oci8" {
		*action.DatabaseType = "oracle"
	}

	data, err := proto.Marshal(&tdb)
	if err != nil {
		c.Fatal("Marshalling error: ", err)
	}

	f.Write(data)
	f.Close()
}

func getConn(action *c.SchemaDiffAction) *sql.DB {
	url := action.Db.GenerateURL(action)
	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(20 * time.Second)
	return db
}

func objectifyStoredProcedures(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	st := time.Now()
	if action.StoreProcFiles {
		os.Mkdir("generated-procs", 0755)
	}

	db := getConn(action)

	cntxt := make([]interface{}, 0)
	cntxt = append(cntxt, *tdb.Name)
	if tdb.Version == nil {
		cntxt = append(cntxt, "")
	} else {
		cntxt = append(cntxt, *tdb.Version)
	}

	qt1 := time.Now()
	//fmt.Println(action.StoredProcedure.Query(*tdb.Name))
	rows, err := db.Query(action.StoredProcedure.Query(cntxt))
	if err != nil {
		if strings.Contains(err.Error(), " is an aggregate function") {
			c.Warn("Error querying StoredProcedures", err)
		} else {
			c.Fatal("Error querying StoredProcedures", err)
		}
	}
	qt2 := time.Now()

	qt := qt2.Sub(qt1).Nanoseconds()

	for rows.Next() {
		schn := ""
		if tdb.SchemaName != nil {
			schn = *tdb.SchemaName
		}
		tsp := action.StoredProcedure.FromResult(rows, schn)
		if strings.Index(*tsp.Name, "\"") == 0 {
			*tsp.Name = (*tsp.Name)[1 : len(*tsp.Name)-1]
		}
		if tsp.Declaration != nil && *tsp.Declaration != "" {
			*tsp.Declaration = strings.Replace(*tsp.Declaration, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tsp.Declaration = strings.Replace(*tsp.Declaration, *tdb.SchemaName+".", "", -1)
		}
		if tsp.DropDeclaration != nil && *tsp.DropDeclaration != "" {
			*tsp.DropDeclaration = strings.Replace(*tsp.DropDeclaration, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tsp.DropDeclaration = strings.Replace(*tsp.DropDeclaration, *tdb.SchemaName+".", "", -1)
		}
		if tsp.Definition != nil && *tsp.Definition != "" && tdb.SchemaName != nil {
			*tsp.Definition = strings.Replace(*tsp.Definition, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tsp.Definition = strings.Replace(*tsp.Definition, *tdb.SchemaName+".", "", -1)
		}

		tdb.StoredProcs = append(tdb.StoredProcs, tsp)
		if !action.DuplicateProcNamesAllowed {
			action.Procs[*tsp.Name] = tsp
		}
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching StoredProcedure details from database", err)
	}

	fmt.Printf("Total stored procedures %d, Query time = %ds, Total time = %ds\n", len(tdb.StoredProcs), qt/fc, time.Since(st).Nanoseconds()/fc)

	defer rows.Close()

	qt = time.Duration(0).Nanoseconds()
	st = time.Now()
	spParamCount := 0
	if action.DuplicateProcNamesAllowed {
		for _, tsp := range tdb.StoredProcs {
			//fmt.Printf("Objectifying procedure %s\n", *tsp.Name)

			cntxt := make([]interface{}, 0)
			cntxt = append(cntxt, *tsp.Name)
			cntxt = append(cntxt, *tsp.NumParams)
			cntxt = append(cntxt, *tdb.Name)

			qt1 := time.Now()
			rows, err := db.Query(action.StoredProcedureParam.Query(cntxt))
			if err != nil {
				c.Fatal("Error querying StoredProcedureParams", err)
			}
			qt2 := time.Now()

			qt += qt2.Sub(qt1).Nanoseconds()

			for rows.Next() {
				tmp := action.StoredProcedureParam.FromResult(rows, cntxt)
				switch tmp.(type) {
				case []*pb2.StoredProcedureParam:
					spParamCount += len(tmp.([]*pb2.StoredProcedureParam))
					tsp.Params = append(tsp.Params, tmp.([]*pb2.StoredProcedureParam)...)
				case *pb2.StoredProcedureParam:
					spParamCount++
					tsp.Params = append(tsp.Params, tmp.(*pb2.StoredProcedureParam))
				}
			}
			err = rows.Err()
			if err != nil {
				c.Fatal("Error fetching StoredProcedureParam details from database", err)
			}
			rows.Close()

			tsp.NumParams = proto.Int32(int32(len(tsp.Params)))

			/*spl := action.StoredProcedure.DefineQuery(cntxt)
			if spl != "" {
				rows, err := db.Query(spl)
				if err != nil {
					c.Fatal("Error querying the database", err)
				}
				for rows.Next() {
					def := action.StoredProcedure.Definition(rows)
					tsp.Definition = proto.String(def)
				}
				err = rows.Err()
				if err != nil {
					c.Fatal("Error fetching details from database", err)
				}
			}
			if tsp.Definition == nil || *tsp.Definition == "" {
				fmt.Printf("Unable to fetch definition for sp %s\n", *tsp.Name)
			}*/
		}
	} else {
		start := 0
		size := 5000
		for {
			count := 0

			cntxt := make([]interface{}, 0)
			cntxt = append(cntxt, start)
			cntxt = append(cntxt, size)
			cntxt = append(cntxt, *tdb.Name)

			qt1 := time.Now()
			rows, err := db.Query(action.StoredProcedureParam.Query(cntxt))
			if err != nil {
				c.Fatal("Error querying StoredProcedureParams", err)
			}
			qt2 := time.Now()

			qt += qt2.Sub(qt1).Nanoseconds()

			cntxt = make([]interface{}, 0)
			cntxt = append(cntxt, "")
			cntxt = append(cntxt, 0)
			for rows.Next() {
				count++
				tmp := action.StoredProcedureParam.FromResult(rows, cntxt)
				switch tmp.(type) {
				case []*pb2.StoredProcedureParam:
					spplst := tmp.([]*pb2.StoredProcedureParam)
					for _, spp := range spplst {
						tsp, ok := action.Procs[*spp.ProcName]
						if ok {
							spParamCount += len(spplst)
							tsp.Params = append(tsp.Params, spp)
						}
					}
				case *pb2.StoredProcedureParam:
					spp := tmp.(*pb2.StoredProcedureParam)
					tsp, ok := action.Procs[*spp.ProcName]
					if ok {
						spParamCount++
						tsp.Params = append(tsp.Params, spp)
					}
				}
			}
			err = rows.Err()
			if err != nil {
				c.Fatal("Error fetching StoredProcedureParam details from database", err)
			}

			for _, tsp := range tdb.StoredProcs {
				tsp.NumParams = proto.Int32(int32(len(tsp.Params)))
			}

			if count == 0 || count < size {
				break
			}
			start += count

			defer rows.Close()
		}
	}
	fmt.Printf("Total stored procedure parameters = %d, Query time = %ds, Total time = %ds\n", spParamCount, qt/fc, time.Since(st).Nanoseconds()/fc)

	if action.StoreProcFiles {
		st = time.Now()
		donefiles := make(map[string]int)
		for _, tsp := range tdb.StoredProcs {
			fileName := ""
			_, ok := donefiles[*tsp.Name]
			if !ok {
				donefiles[*tsp.Name] = 0
			} else {
				donefiles[*tsp.Name] = donefiles[*tsp.Name] + 1
				fileName = "_" + strconv.Itoa(donefiles[*tsp.Name])
			}
			fileName = "generated-procs/" + *tsp.Name + fileName + ".sql"
			spf, err := os.Create(fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				spf.WriteString(*tsp.Definition)
				defer spf.Close()
			}
		}
		fmt.Printf("Time taken for generating stored procedure files = %ds\n", time.Since(st).Nanoseconds()/fc)
	}

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifyViews(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	st := time.Now()

	db := getConn(action)

	rows, err := db.Query(action.View.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying Views", err)
	}
	qt := time.Since(st).Nanoseconds()

	for rows.Next() {
		tvw := action.View.FromResult(rows, nil)
		if tvw.Definition != nil && tdb.SchemaName != nil {
			*tvw.Definition = strings.Replace(*tvw.Definition, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tvw.Definition = strings.Replace(*tvw.Definition, *tdb.SchemaName+".", "", -1)
		}
		//fmt.Printf("Objectifying view %s\n", *tvw.Name)
		tdb.Views = append(tdb.Views, tvw)
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching View details from database", err)
	}

	fmt.Printf("Total views %d, Query time = %ds, Total time = %ds\n", len(tdb.Views), qt/fc, time.Since(st).Nanoseconds()/fc)

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifySequences(action *c.SchemaDiffAction, tdb *pb2.Db, ch chan int) {
	st := time.Now()

	cntxt := make([]interface{}, 0)
	cntxt = append(cntxt, *tdb.Name)
	if tdb.Version == nil {
		cntxt = append(cntxt, "")
	} else {
		cntxt = append(cntxt, *tdb.Version)
	}

	query := action.Sequence.Query(cntxt)
	if query == "" {
		ch <- 1
		return
	}

	db := getConn(action)

	rows, err := db.Query(query)
	if err != nil {
		c.Fatal("Error querying Sequences", err)
	}
	qtd := time.Since(st).Nanoseconds()

	uniqViewMap := make(map[string]*pb2.Sequence)
	for rows.Next() {
		tvw := action.Sequence.FromResult(rows, &uniqViewMap)
		if tvw == nil {
			continue
		}

		cntxt := make([]interface{}, 0)
		cntxt = append(cntxt, *tvw.Name)
		if tdb.Version == nil {
			cntxt = append(cntxt, "")
		} else {
			cntxt = append(cntxt, *tdb.Version)
		}

		cacheValQ := action.Sequence.ExQuery(cntxt)
		if cacheValQ != "" {
			qt1 := time.Now()
			rows, err = db.Query(cacheValQ)
			if err != nil {
				c.Fatal("Error querying Sequence cache value", err)
			}
			qt2 := time.Now()

			qtd += qt2.Sub(qt1).Nanoseconds()

			for rows.Next() {
				action.Sequence.UpdateSequence(rows, tvw)
			}
			err = rows.Err()
			if err != nil {
				c.Fatal("Error fetching Sequence cache value details from database", err)
			}
		}

		tdb.Sequences = append(tdb.Sequences, tvw)
		//if tvw.Name != nil {
		//fmt.Printf("Objectifying sequence %s\n", *tvw.Name)
		//}
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Sequence details from database", err)
	}

	fmt.Printf("Total sequences %d, Query time = %ds, Total time = %ds\n", len(tdb.Sequences), qtd/fc, time.Since(st).Nanoseconds()/fc)

	defer rows.Close()
	defer db.Close()

	ch <- 1
}

func objectifyTables(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	st := time.Now()

	db := getConn(action)

	cntxt := make([]interface{}, 0)
	cntxt = append(cntxt, action.SchemaName)
	cntxt = append(cntxt, *tdb.Name)

	rows, err := db.Query(action.Table.Query(cntxt))
	if err != nil {
		c.Fatal("Error querying Tables", err)
	}
	qt := time.Since(st).Nanoseconds()

	for rows.Next() {
		tbl := action.Table.FromResult(rows, nil)
		tdb.Tables = append(tdb.Tables, tbl)
		action.Tables[*tbl.Name] = tbl
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Table details from database", err)
	}

	fmt.Printf("Total tables %d, Query time = %ds, Total time = %ds\n", len(tdb.Tables), qt/fc, time.Since(st).Nanoseconds()/fc)

	rows.Close()
	db.Close()

	ch := make(chan int)

	go objectifyColumns(tdb, action, ch)
	go objectifyConstraints(tdb, action, ch)
	go objectifyIndexes(tdb, action, ch)

	<-ch
	<-ch
	<-ch

	for _, tdst := range tdb.Tables {
		for _, c := range tdst.Columns {
			if c.DefVal != nil && tdb.SchemaName != nil {
				*c.DefVal = strings.Replace(*c.DefVal, "\""+*tdb.SchemaName+"\".", "", -1)
				//*c.DefVal = strings.Replace(*c.DefVal, *tdb.SchemaName+".", "", -1)
			}
		}

		updtConstraints := make([]*pb2.Constraint, 0)
		for _, c := range tdst.Constraints {
			if c.Definition != nil && tdb.SchemaName != nil {
				*c.Definition = strings.Replace(*c.Definition, "\""+*tdb.SchemaName+"\".", "", -1)
				//*c.Definition = strings.Replace(*c.Definition, *tdb.SchemaName+".", "", -1)
			}

			//Check for valid column names in constraints
			fl := false
			if c.Columns != nil && len(c.Columns) > 0 {
				for _, cnc := range c.Columns {
					for _, co := range tdst.Columns {
						if *co.Name == cnc {
							fl = true
							break
						}
					}
					if fl {
						break
					}
				}
			} else if *c.Condition != "" {
				for _, co := range tdst.Columns {
					if strings.Contains(*c.Condition, *co.Name) {
						fl = true
						break
					}
				}
			}

			if fl {
				updtConstraints = append(updtConstraints, c)
			}
		}

		tdst.Constraints = updtConstraints
	}

	defer wg.Done()
}

func objectifyColumns(tdb *pb2.Db, action *c.SchemaDiffAction, ch chan int) {
	st := time.Now()
	start := 0
	size := 10000
	columnCount := 0
	qtd := time.Duration(0).Nanoseconds()

	db := getConn(action)
	for {
		count := 0

		cntxt := make([]interface{}, 0)
		cntxt = append(cntxt, start)
		cntxt = append(cntxt, size)
		cntxt = append(cntxt, *tdb.Name)

		qt1 := time.Now()
		rows, err := db.Query(action.Column.Query(cntxt))
		if err != nil {
			c.Fatal("Error querying Columns", err)
		}
		qt2 := time.Now()

		qtd += qt2.Sub(qt1).Nanoseconds()

		for rows.Next() {
			//columnCount++
			count++
			tsp := action.Column.FromResult(rows, action.Tables)
			if tsp != nil {
				columnCount++
				_, ok := action.Tables[*tsp.Name]
				if !ok {
					//fmt.Printf("Objectifying table %s", *tsp.Name)
					action.Tables[*tsp.Name] = tsp
				}
			}
		}
		err = rows.Err()
		if err != nil {
			c.Fatal("Error fetching Column details from database", err)
		}

		if count == 0 || count < size {
			break
		}
		start += count

		rows.Close()
	}

	defer db.Close()
	ch <- 1
	fmt.Printf("Total table columns = %d, Query time = %ds, Total time = %ds\n", columnCount, qtd/fc, time.Since(st).Nanoseconds()/fc)
}

func objectifyConstraints(tdb *pb2.Db, action *c.SchemaDiffAction, ch chan int) {
	st := time.Now()
	start := 0
	size := 2000
	columnCount := 0
	qtd := time.Duration(0).Nanoseconds()

	db := getConn(action)
	for {
		count := 0

		cntxt := make([]interface{}, 0)
		cntxt = append(cntxt, start)
		cntxt = append(cntxt, size)
		cntxt = append(cntxt, *tdb.Name)

		qt1 := time.Now()
		rows, err := db.Query(action.Constraint.Query(cntxt))
		if err != nil {
			c.Fatal("Error querying Constraints", err)
		}
		qt2 := time.Now()

		qtd += qt2.Sub(qt1).Nanoseconds()

		cntxt1 := make([]interface{}, 0)
		cntxt1 = append(cntxt1, action.Tables)
		cntxt1 = append(cntxt1, action.Constraints)

		for rows.Next() {
			//columnCount++
			count++
			tsp := action.Constraint.FromResult(rows, cntxt1)
			if tsp != nil {
				columnCount++
				//fmt.Printf("Objectifying constraint %s\n", *tsp.Name)
			}
		}
		err = rows.Err()
		if err != nil {
			c.Fatal("Error fetching Constraint details from database", err)
		}

		if count == 0 || count < size {
			break
		}
		start += count

		rows.Close()
	}

	defer db.Close()

	ch <- 1
	fmt.Printf("Total table constraints = %d, Query time = %ds, Total time = %ds\n", columnCount, qtd/fc, time.Since(st).Nanoseconds()/fc)
}

func objectifyIndexes(tdb *pb2.Db, action *c.SchemaDiffAction, ch chan int) {
	st := time.Now()
	indxCount := 0
	qtd := time.Duration(0).Nanoseconds()

	db := getConn(action)

	cntxt := make([]interface{}, 0)

	qt1 := time.Now()
	rows, err := db.Query(action.Index.Query(cntxt))
	if err != nil {
		c.Fatal("Error querying Indexes", err)
	}
	qt2 := time.Now()

	qtd += qt2.Sub(qt1).Nanoseconds()

	for rows.Next() {
		indxCount++

		cntxt1 := make([]interface{}, 0)
		cntxt1 = append(cntxt1, action.Indexess)

		tvw := action.Index.FromResult(rows, cntxt1)
		if tvw.Definition != nil && tdb.SchemaName != nil {
			*tvw.Definition = strings.Replace(*tvw.Definition, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tvw.Definition = strings.Replace(*tvw.Definition, *tdb.SchemaName+".", "", -1)
		}
		//fmt.Printf("Objectifying index %s\n", *tvw.Name)
		_, ok := action.Indexes[*tvw.TableName]
		if !ok {
			action.Indexes[*tvw.TableName] = make([]*pb2.Index, 0)
		}
		action.Indexes[*tvw.TableName] = append(action.Indexes[*tvw.TableName], tvw)

		_, ok = action.Indexess[*tvw.Name]
		if !ok {
			action.Indexess[*tvw.Name] = tvw
		}
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Index details from database", err)
	}

	rows.Close()

	defer db.Close()

	fmt.Printf("Total table indexes %d, Query time = %ds, Total time = %ds\n", indxCount, qtd/fc, time.Since(st).Nanoseconds()/fc)
	ch <- 1
}

func objectifyTriggers(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	st := time.Now()

	db := getConn(action)

	rows, err := db.Query(action.Trigger.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying Triggers", err)
	}
	qtd := time.Since(st).Nanoseconds()

	trgCount := 0
	for rows.Next() {
		trgCount++
		tvw := action.Trigger.FromResult(rows, nil)
		if tdb.SchemaName != nil {
			*tvw.Definition = strings.Replace(*tvw.Definition, "\""+*tdb.SchemaName+"\".", "", -1)
			//*tvw.Definition = strings.Replace(*tvw.Definition, *tdb.SchemaName+".", "", -1)
			if tvw.FunctionDef != nil {
				*tvw.FunctionDef = strings.Replace(*tvw.FunctionDef, "\""+*tdb.SchemaName+"\".", "", -1)
				//*tvw.FunctionDef = strings.Replace(*tvw.FunctionDef, *tdb.SchemaName+".", "", -1)
			}
		}

		//fmt.Printf("Objectifying trigger %s\n", *tvw.Name)

		spl := action.Trigger.DefineQuery(*tvw.Function)
		if spl != "" {
			qt1 := time.Now()
			rows, err := db.Query(spl)
			if err != nil {
				c.Fatal("Error querying Trigger definition", err)
			}
			qt2 := time.Now()

			qtd += qt2.Sub(qt1).Nanoseconds()

			for rows.Next() {
				def := action.Trigger.GetDefinition(rows)
				tvw.FunctionDef = proto.String(def)
				break
			}
			err = rows.Err()
			if err != nil {
				c.Fatal("Error fetching Trigger definition details from database", err)
			}

			rows.Close()
		}

		if tvw.TableName != nil {
			_, ok := action.Triggers[*tvw.TableName]
			if !ok {
				action.Triggers[*tvw.TableName] = make([]*pb2.Trigger, 0)
			}
			action.Triggers[*tvw.TableName] = append(action.Triggers[*tvw.TableName], tvw)
		}
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Trigger details from database", err)
	}

	fmt.Printf("Total table triggers %d, Query time = %ds, Total time = %ds\n", trgCount, qtd/fc, time.Since(st).Nanoseconds()/fc)

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func mergeTriggersIndexesWithTables(tdb *pb2.Db, action *c.SchemaDiffAction) {
	for _, tb := range tdb.Tables {
		triggers, ok := action.Triggers[*tb.Name]
		if ok {
			if *tdb.Type != pb2.Db_POSTGRESQL {
				tb.Triggers = triggers
			} else {
				tb.Triggers = c.MergeDuplicates(triggers)
			}
		}
		indexes, ok := action.Indexes[*tb.Name]
		if ok {
			tb.Indexes = indexes
		}
	}
}
