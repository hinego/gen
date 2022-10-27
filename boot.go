package gen

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"strings"
)

// DBType database type
type DBType string

const (
	// dbMySQL Gorm Drivers mysql || postgres || sqlite || sqlserver
	DbMySQL     DBType = "mysql"
	DbPostgres  DBType = "postgres"
	DbSQLite    DBType = "sqlite"
	DbSQLServer DBType = "sqlserver"
)
const (
	// DefaultOutPath default path
	DefaultOutPath = "./dao/query"
)

// CmdParams is command line parameters
type CmdParams struct {
	DSN               string   `yaml:"dsn"`               // consult[https://gorm.io/docs/connecting_to_the_database.html]"
	DB                string   `yaml:"db"`                // input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
	Tables            []string `yaml:"tables"`            // enter the required data table or leave it blank
	OnlyModel         bool     `yaml:"onlyModel"`         // only generate model
	OutPath           string   `yaml:"outPath"`           // specify a directory for output
	OutFile           string   `yaml:"outFile"`           // query code file name, default: gen.go
	WithUnitTest      bool     `yaml:"withUnitTest"`      // generate unit test for query code
	ModelPkgName      string   `yaml:"modelPkgName"`      // generated model code's package name
	FieldNullable     bool     `yaml:"fieldNullable"`     // generate with pointer when field is nullable
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
	FieldSignable     bool     `yaml:"fieldSignable"`     // detect integer field's unsigned type, adjust generated data type
}

// YamlConfig is yaml config struct
type YamlConfig struct {
	Version  string     `yaml:"version"`  //
	Database *CmdParams `yaml:"database"` //
}

func Connect(t DBType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}
	switch t {
	case DbMySQL:
		return gorm.Open(mysql.Open(dsn))
	case DbPostgres:
		return gorm.Open(postgres.Open(dsn))
	case DbSQLite:
		return gorm.Open(sqlite.Open(dsn))
	case DbSQLServer:
		return gorm.Open(sqlserver.Open(dsn))
	default:
		return nil, fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", t)
	}
}

// genModels is gorm/gen generated models
func genModels(g *Generator, db *gorm.DB, tables []string) (models []interface{}, err error) {
	var tablesList []string
	if len(tables) == 0 {
		// Execute tasks for all tables in the database
		tablesList, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	} else {
		tablesList = tables
	}

	// Execute some data table tasks
	models = make([]interface{}, len(tablesList))
	for i, tableName := range tablesList {
		if opt := g.GetModel(tableName); opt != nil {
			models[i] = g.GenerateModel(tableName, WithMethod(opt))
		} else {
			models[i] = g.GenerateModel(tableName)
		}
	}
	return models, nil
}

// loadConfigFile load config file from path
func loadConfigFile(path string) (*CmdParams, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint
	var yamlConfig YamlConfig
	if cmdErr := yaml.NewDecoder(file).Decode(&yamlConfig); cmdErr != nil {
		return nil, cmdErr
	}
	return yamlConfig.Database, nil
}

func ArgParse() *CmdParams {
	// choose is file or flag
	genPath := flag.String("c", "", "is path for gen.yml")
	dsn := flag.String("dsn", "", "consult[https://gorm.io/docs/connecting_to_the_database.html]")
	db := flag.String("db", "mysql", "input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]")
	tableList := flag.String("tables", "", "enter the required data table or leave it blank")
	onlyModel := flag.Bool("onlyModel", false, "only generate models (without query file)")
	outPath := flag.String("outPath", "./dao/query", "specify a directory for output")
	outFile := flag.String("outFile", "", "query code file name, default: gen.go")
	withUnitTest := flag.Bool("withUnitTest", false, "generate unit test for query code")
	modelPkgName := flag.String("modelPkgName", "", "generated model code's package name")
	fieldNullable := flag.Bool("fieldNullable", false, "generate with pointer when field is nullable")
	fieldWithIndexTag := flag.Bool("fieldWithIndexTag", false, "generate field with gorm index tag")
	fieldWithTypeTag := flag.Bool("fieldWithTypeTag", false, "generate field with gorm column type tag")
	fieldSignable := flag.Bool("fieldSignable", false, "detect integer field's unsigned type, adjust generated data type")
	flag.Parse()
	var cmdParse CmdParams
	if *genPath != "" {
		if configFileParams, err := loadConfigFile(*genPath); err == nil && configFileParams != nil {
			cmdParse = *configFileParams
		}
	}
	// cmd first
	if *dsn != "" {
		cmdParse.DSN = *dsn
	}
	if *db != "" {
		cmdParse.DB = *db
	}
	if *tableList != "" {
		cmdParse.Tables = strings.Split(*tableList, ",")
	}
	if *onlyModel {
		cmdParse.OnlyModel = true
	}
	if *outPath != DefaultOutPath {
		cmdParse.OutPath = *outPath
	}
	if *outFile != "" {
		cmdParse.OutFile = *outFile
	}
	if *withUnitTest {
		cmdParse.WithUnitTest = *withUnitTest
	}
	if *modelPkgName != "" {
		cmdParse.ModelPkgName = *modelPkgName
	}
	if *fieldNullable {
		cmdParse.FieldNullable = *fieldNullable
	}
	if *fieldWithIndexTag {
		cmdParse.FieldWithIndexTag = *fieldWithIndexTag
	}
	if *fieldWithTypeTag {
		cmdParse.FieldWithTypeTag = *fieldWithTypeTag
	}
	if *fieldSignable {
		cmdParse.FieldSignable = *fieldSignable
	}
	return &cmdParse
}
