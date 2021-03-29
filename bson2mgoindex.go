package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/fatih/structtag"
)

type StructFields []*ast.Field

var mgoCreateIndexLayout string = `db.getCollection("%s").createIndex(%s,{background: true});`

var bsonPath = flag.String("f", "", "file path")

// go run mian.go -f ./models/bson_test.go

func main() {

	flag.Parse()

	log.Printf("start create %s index file\n", *bsonPath)

	var indexs []string

	fieldsMap, err := ParseStruct(*bsonPath, nil, "mgo")
	if err != nil {
		fmt.Println(err)
	}

	funcMap, err := ParseFunc(*bsonPath, nil, "CollectName")
	if err != nil {
		fmt.Println(err)
	}

	for structName, fields := range fieldsMap {

		// fmt.Printf("StructName:%s\n", structName)
		// fmt.Printf("CollectName:%s\n", funcMap[structName])
		for _, field := range fields {

			var bsonName string
			var isIndex bool
			var indexType string
			var indexSort int

			mgoIndexMap := make(map[string]int, 0)

			// fmt.Printf("	FieldName:%s\n", field.Names[0].Name)
			// fmt.Printf("	FieldType:%s\n", field.Type)
			// fmt.Printf("	FieldTag:%s\n", field.Tag.Value)
			tags, _ := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
			bsonTag, _ := tags.Get("bson")
			// fmt.Printf("		BsonTagKey:%s\n", bsonTag.Key)
			// fmt.Printf("		BsonTagName:%s\n", bsonTag.Name)
			// fmt.Printf("		BsonTagOptions:%s\n", bsonTag.Options)
			mgoTag, _ := tags.Get("mgo")
			// fmt.Printf("		MgoTagKey:%s\n", mgoTag.Key)
			// fmt.Printf("		MgoTagName:%s\n", mgoTag.Name)
			// fmt.Printf("		MgoTagOptions:%s\n", mgoTag.Options)

			bsonName = bsonTag.Name
			mtns := strings.Split(mgoTag.Name, ";")
			//TODO: len(mtns) = 1
			index := strings.Split(mtns[0], ":")
			if len(index) > 1 {
				if index[0] == "index" {
					isIndex = true
					indexType = index[0]
					indexSort, _ = strconv.Atoi(index[1])
				}
			}
			if isIndex && indexType == "index" {
				mgoIndexMap[bsonName] = indexSort
			}
			indexStr, _ := jsoniter.MarshalToString(mgoIndexMap)
			in := fmt.Sprintf(mgoCreateIndexLayout, funcMap[structName], indexStr)
			indexs = append(indexs, in)
		}

	}
	si := strings.LastIndex(*bsonPath, "/")
	fileHead := (*bsonPath)[si+1:]
	fileName := "mgo_index_" + fileHead + "_" + time.Now().Format("20060102") + ".js"
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	for _, v := range indexs {
		fmt.Println(v)
		file.WriteString(v + "\n")
	}
}

func ParseStruct(filename string, src []byte, tagName string) (structMap map[string]StructFields, err error) {
	structMap = make(map[string]StructFields)

	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			return structMap, err
		}
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return structMap, err
	}

	collectStructs := func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		// 获取结构体名称
		structName := ts.Name.Name

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			tag := field.Tag.Value
			tag = strings.Trim(tag, "`")
			tags, err := structtag.Parse(string(tag))
			if err != nil {
				return true
			}
			_, err = tags.Get(tagName)
			if err == nil {
				structMap[structName] = append(structMap[structName], field)
			}
		}
		return false
	}

	ast.Inspect(file, collectStructs)

	return structMap, nil
}

func ParseFunc(filename string, src []byte, FuncName string) (funcMap map[string]string, err error) {
	funcMap = make(map[string]string)

	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			return funcMap, err
		}
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return funcMap, err
	}

	collectStructs := func(x ast.Node) bool {
		fs, ok := x.(*ast.FuncDecl)
		if !ok || fs.Type == nil {
			return true
		}

		// 获取方法名称
		funcName := fs.Name.Name

		if funcName != FuncName {
			return true
		}

		// 获取接受者名称
		recvName := ""
		// 获取函数体表名
		collectName := ""

		for _, field := range fs.Recv.List {
			recvName = field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
		}

		for _, s := range fs.Body.List {
			collectName = s.(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
			collectName = strings.Trim(collectName, `"`)
		}
		funcMap[recvName] = collectName
		return false
	}

	ast.Inspect(file, collectStructs)

	return funcMap, nil
}
