package main

import (
	sql "database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	_ "gopkg.in/goracle.v2"
	"os"
	c "sqldiffer/common"
	db "sqldiffer/db"
	pb2 "sqldiffer/protos"
	sq "sqldiffer/sequence"
	sp "sqldiffer/storedproc"
	spp "sqldiffer/storedproc/storedprocparam"
	tb "sqldiffer/table"
	co "sqldiffer/table/column"
	cn "sqldiffer/table/constraint"
	in "sqldiffer/table/index"
	tr "sqldiffer/table/trigger"
	vw "sqldiffer/view"
	"strings"
	"sync"
	"time"
)

var opts struct {
	DatabaseType string `short:"t" long:"type" description:"The database type" required:"true" choice:"postgres" choice:"oracle" choice:"mysql" choice:"sqlserver"`
	DatabaseName string `short:"n" long:"name" description:"The database name" required:"true"`
	Host         string `short:"i" long:"host" description:"The database host" required:"true"`
	Port         int32  `short:"p" long:"port" description:"The database port" required:"true"`
	User         string `short:"u" long:"user" description:"The database user" required:"true"`
	Password     string `short:"w" long:"password" description:"The database password" required:"true"`
	SchemaName   string `short:"s" long:"sch-nam" description:"The database schema name"`
	FileName     string `short:"f" long:"fil-nam" description:"The generated schema/diff file name"`
}

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
		Host:         &opts.Host,
		Port:         opts.Port,
		DatabaseType: &opts.DatabaseType,
		DatabaseName: &opts.DatabaseName,
		User:         &opts.User,
		Password:     &opts.Password,
		SchemaName:   &opts.SchemaName,
		FileName:     &opts.FileName,
		Parallel:     true,
		IsDiffNeeded: false,
	}
	//fmt.Println("%!", action)
	generateSchema(&action)
}

func generateSchema(action *c.SchemaDiffAction) {
	action.Procs = make(map[string]*pb2.StoredProcedure)
	action.Tables = make(map[string]*pb2.Table)
	action.TablesP = make(map[string]bool)
	action.Indexes = make(map[string][]*pb2.Index)
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
		*action.DatabaseType = "goracle"
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
		ps := pb2.Db_SQLSERVER
		tdb.Type = &ps
	}

	var wg sync.WaitGroup
	wg.Add(5)

	ch := make(chan int)

	go objectifyStoredProcedures(action, &tdb, &wg)
	go objectifyViews(action, &tdb, &wg)
	go objectifyTriggers(action, &tdb, &wg)
	go objectifyIndexes(action, &tdb, &wg)
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
	if *action.DatabaseType == "goracle" {
		*action.DatabaseType = "oracle"
	}

	data, err := proto.Marshal(&tdb)
	if err != nil {
		c.Fatal("Marshalling error: ", err)
	}

	f.Write(data)
	f.Close()
}

func objectifyStoredProcedures(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	url := action.Db.GenerateURL(action)

	//fmt.Println(url)
	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	//fmt.Println(action.StoredProcedure.Query(*tdb.Name))
	rows, err := db.Query(action.StoredProcedure.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying StoredProcedures", err)
	}

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

	fmt.Printf("Total stored procedures %d\n", len(tdb.StoredProcs))

	defer rows.Close()

	spParamCount := 0
	if action.DuplicateProcNamesAllowed {
		for _, tsp := range tdb.StoredProcs {
			//fmt.Printf("Objectifying procedure %s\n", *tsp.Name)

			cntxt := make([]interface{}, 0)
			cntxt = append(cntxt, *tsp.Name)
			cntxt = append(cntxt, *tsp.NumParams)
			cntxt = append(cntxt, *tdb.Name)

			rows, err := db.Query(action.StoredProcedureParam.Query(cntxt))
			if err != nil {
				c.Fatal("Error querying StoredProcedureParams", err)
			}

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
		size := 1000
		for {
			count := 0

			cntxt := make([]interface{}, 0)
			cntxt = append(cntxt, start)
			cntxt = append(cntxt, size)
			cntxt = append(cntxt, *tdb.Name)

			rows, err := db.Query(action.StoredProcedureParam.Query(cntxt))
			if err != nil {
				c.Fatal("Error querying StoredProcedureParams", err)
			}

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

			if count == 0 {
				break
			}
			start += count

			defer rows.Close()
		}
	}
	fmt.Printf("Total stored procedure parameters = %d\n", spParamCount)

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifyViews(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	url := action.Db.GenerateURL(action)

	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	rows, err := db.Query(action.View.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying Views", err)
	}

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

	fmt.Printf("Total views %d\n", len(tdb.Views))

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifySequences(action *c.SchemaDiffAction, tdb *pb2.Db, ch chan int) {
	url := action.Db.GenerateURL(action)

	query := action.Sequence.Query(*tdb.Name)
	if query == "" {
		ch <- 1
		return
	}

	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	rows, err := db.Query(query)
	if err != nil {
		c.Fatal("Error querying Sequences", err)
	}

	uniqViewMap := make(map[string]*pb2.Sequence)
	for rows.Next() {
		tvw := action.Sequence.FromResult(rows, &uniqViewMap)
		if tvw == nil {
			continue
		}

		cacheValQ := action.Sequence.ExQuery(*tvw.Name)
		if cacheValQ != "" {
			rows, err = db.Query(cacheValQ)
			if err != nil {
				c.Fatal("Error querying Sequence cache value", err)
			}

			for rows.Next() {
				action.Sequence.UpdateSequence(rows, tvw)
			}
			err = rows.Err()
			if err != nil {
				c.Fatal("Error fetching Sequence cache value details from database", err)
			}
		}

		tdb.Sequences = append(tdb.Sequences, tvw)
		if tvw.Name != nil {
			//fmt.Printf("Objectifying sequence %s\n", *tvw.Name)
		}
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Sequence details from database", err)
	}

	fmt.Printf("Total sequences %d\n", len(tdb.Sequences))

	defer rows.Close()
	defer db.Close()

	ch <- 1
}

func objectifyTables(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	url := action.Db.GenerateURL(action)

	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	cntxt := make([]interface{}, 0)
	cntxt = append(cntxt, action.SchemaName)
	cntxt = append(cntxt, *tdb.Name)

	rows, err := db.Query(action.Table.Query(cntxt))
	if err != nil {
		c.Fatal("Error querying Tables", err)
	}

	for rows.Next() {
		tbl := action.Table.FromResult(rows, nil)
		tdb.Tables = append(tdb.Tables, tbl)
		action.Tables[*tbl.Name] = tbl
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Table details from database", err)
	}

	fmt.Printf("Total tables %d\n", len(tdb.Tables))

	objectifyColumns(db, tdb, action)
	objectifyConstraints(db, tdb, action)

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
					if strings.Index(*c.Condition, *co.Name) != -1 {
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

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifyColumns(db *sql.DB, tdb *pb2.Db, action *c.SchemaDiffAction) {
	start := 0
	size := 1000
	columnCount := 0
	for {
		count := 0

		cntxt := make([]interface{}, 0)
		cntxt = append(cntxt, start)
		cntxt = append(cntxt, size)
		cntxt = append(cntxt, *tdb.Name)

		rows, err := db.Query(action.Column.Query(cntxt))
		if err != nil {
			c.Fatal("Error querying Columns", err)
		}

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

		if count == 0 {
			break
		}
		start += count

		rows.Close()
	}

	fmt.Printf("Total table columns = %d\n", columnCount)
}

func objectifyConstraints(db *sql.DB, tdb *pb2.Db, action *c.SchemaDiffAction) {
	start := 0
	size := 1000
	columnCount := 0
	for {
		count := 0

		cntxt := make([]interface{}, 0)
		cntxt = append(cntxt, start)
		cntxt = append(cntxt, size)
		cntxt = append(cntxt, *tdb.Name)

		//fmt.Println(action.Constraint.Query(cntxt))
		rows, err := db.Query(action.Constraint.Query(cntxt))
		if err != nil {
			c.Fatal("Error querying Constraints", err)
		}

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

		if count == 0 {
			break
		}
		start += count

		rows.Close()
	}

	fmt.Printf("Total table constraints = %d\n", columnCount)
}

func objectifyIndexes(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	url := action.Db.GenerateURL(action)

	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	rows, err := db.Query(action.Index.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying Indexes", err)
	}

	indxCount := 0
	for rows.Next() {
		indxCount++
		tvw := action.Index.FromResult(rows, nil)
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
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Index details from database", err)
	}

	fmt.Printf("Total indexes %d\n", indxCount)

	defer rows.Close()
	defer db.Close()

	defer wg.Done()
}

func objectifyTriggers(action *c.SchemaDiffAction, tdb *pb2.Db, wg *sync.WaitGroup) {
	url := action.Db.GenerateURL(action)

	db, err := sql.Open(*action.DatabaseType, url)
	if err != nil {
		c.Fatal("Error connecting to the database", err)
	}

	rows, err := db.Query(action.Trigger.Query(*tdb.Name))
	if err != nil {
		c.Fatal("Error querying Triggers", err)
	}

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
			rows, err := db.Query(spl)
			if err != nil {
				c.Fatal("Error querying Trigger definition", err)
			}
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

		_, ok := action.Triggers[*tvw.TableName]
		if !ok {
			action.Triggers[*tvw.TableName] = make([]*pb2.Trigger, 0)
		}
		action.Triggers[*tvw.TableName] = append(action.Triggers[*tvw.TableName], tvw)
	}
	err = rows.Err()
	if err != nil {
		c.Fatal("Error fetching Trigger details from database", err)
	}

	fmt.Printf("Total triggers %d\n", trgCount)

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
