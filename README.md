# GORM Gen So Simple

```go
package main

import (
	"gorm.io/gen"
	"log"
)

var g *gen.Generator
var config *gen.CmdParams

func main() {
	m := []any{
		//table.Token{},
		//table.User{},
		//table.Node{},
		//table.Bot{},
		//table.Chat{},
		//table.Link{},
	}
	if err := g.LinkModel(m...); err != nil {
		log.Println(err)
		return
	}
	if !config.OnlyModel {
		g.ApplyBasic(g.GenerateAllTable()...)
	}
	g.Execute()
}
func init() {
	log.SetFlags(log.Llongfile)
	config = gen.ArgParse()
	if config == nil {
		log.Fatalf("parse config fail")
	}
	//schema.RegisterSerializer("auto", database.AutoSerializer{})
	if db, err := gen.Connect(gen.DBType(config.DB), config.DSN); err != nil {
		log.Fatalf("connect db server fail: %v", err)
	} else {
		g = gen.NewGenerator(gen.Config{
			Mode:              gen.WithDefaultQuery | gen.WithoutContext,
			OutPath:           config.OutPath,
			OutFile:           config.OutFile,
			ModelPkgPath:      config.ModelPkgName,
			WithUnitTest:      config.WithUnitTest,
			FieldNullable:     false,
			FieldWithIndexTag: true,
			FieldWithTypeTag:  true,
			FieldSignable:     config.FieldSignable,
		})
		g.UseDB(db)
	}
}
```