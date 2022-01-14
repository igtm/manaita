package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

const (
	DefaultScaffoldFileName           = "SCAFFOLD.md"
	DeprecatedDefaultScaffoldFileName = "CODEGEN.md" // deprecated
	YamlMetaStartCode                 = "---"
	DestFileNameStartCode             = "# "
	CodeTemplateStartCode             = "```"
)

var (
	c = flag.String("c", DefaultScaffoldFileName, "scaffold markdown file path")
	p = flag.String("p", "", "template params")

	funcMap = template.FuncMap{
		"ToUpper":          strings.ToUpper,
		"ToLower":          strings.ToLower,
		"ToSnake":          strcase.ToSnake,
		"ToScreamingSnake": strcase.ToScreamingSnake,
		"ToKebab":          strcase.ToKebab,
		"ToScreamingKebab": strcase.ToScreamingKebab,
		"ToCamel":          strcase.ToCamel,
		"ToLowerCamel":     strcase.ToLowerCamel,
	}
)

func main() {
	flag.Parse()

	errlog := log.New(os.Stderr, "ERROR: ", 0)

	envMap := envToMap()
	currentDir, _ := os.Getwd()
	givenParamMap := paramToMap(*p)

	scaffoldFileName := fmt.Sprintf("%s/%s", currentDir, *c)

	scaffoldFile, err := os.Open(scaffoldFileName)
	if err != nil {
		scaffoldFile, err = os.Open(DeprecatedDefaultScaffoldFileName)
	}
	if err != nil {
		errlog.Println(fmt.Errorf("cannot find '%s'", DefaultScaffoldFileName))
		return
	}
	defer scaffoldFile.Close()

	sc := bufio.NewScanner(scaffoldFile)

	var foundMeta bool
	var endMeta bool
	var metaStr string
	var metaData map[string]interface{}
	paramMap := make(map[string]string)
	var code, dest, destFileName string
	var searchCode bool
	var foundCode bool
	for sc.Scan() {
		loc := sc.Text()
		if !endMeta {
			if strings.HasPrefix(loc, YamlMetaStartCode) {
				if foundMeta {
					endMeta = true
					metaStr += loc + "\n"
					// yaml-meta
					markdown := goldmark.New(
						goldmark.WithExtensions(
							meta.Meta,
						),
					)
					var buf bytes.Buffer
					context := parser.NewContext()
					if err := markdown.Convert([]byte(metaStr), &buf, parser.WithContext(context)); err != nil {
						panic(err)
					}
					metaData = meta.Get(context)
					// validate params
					if metaParams, ok := metaData["Params"]; ok {
						metaParamsSlice := metaParams.([]interface{})
					LOOP:
						for _, m := range metaParamsSlice {
							for pk, pv := range givenParamMap {
								if m == pk {
									paramMap[pk] = pv
									continue LOOP
								}
							}
							panic(fmt.Sprintf("param '%s' is required", m))
						}
					LOOP2:
						for pk, pv := range givenParamMap {
							for _, m := range metaParamsSlice {
								if m == pk {
									paramMap[pk] = pv
									continue LOOP2
								}
							}
							fmt.Printf("param '%s' is not defined on markdown header. so it's not used.\n", pk)
						}
					}
					continue
				}
				foundMeta = true
			}
			if foundMeta {
				metaStr += loc + "\n"
				continue
			}
			// only allowed metadata at first line
			endMeta = true
		}
		if strings.Contains(loc, DestFileNameStartCode) {
			// dest filename
			destFileName = strings.Replace(loc, DestFileNameStartCode, "", -1)
			destFileName = strings.Trim(destFileName, " ")
			// compile filename
			tmpl := template.Must(template.New("").Funcs(funcMap).Parse(destFileName))
			var compiledDest bytes.Buffer
			err = tmpl.Execute(&compiledDest, map[string]interface{}{
				"Env":    envMap,
				"Params": paramMap,
			})
			destFileName = compiledDest.String()
			dest = filepath.Join(currentDir, destFileName)
			searchCode = true
		}
		if searchCode {
			if strings.Contains(loc, CodeTemplateStartCode) {
				if foundCode {
					// end
					// recursively make directories
					os.MkdirAll(filepath.Dir(dest), os.ModePerm)
					destFile, err := os.Create(dest)
					if err != nil {
						errlog.Println(fmt.Errorf("cannot Create '%s'", dest))
						return
					}
					defer destFile.Close()
					fmt.Println(destFileName)

					tmpl := template.Must(template.New("").Funcs(funcMap).Parse(code))
					err = tmpl.Execute(destFile, map[string]interface{}{
						"Env":    envMap,
						"Params": paramMap,
					})
					if err != nil {
						errlog.Println(fmt.Errorf("cannot write code to file '%s'", dest))
						return
					}
					searchCode = false
					foundCode = false
					code = ""
					dest = ""
					destFileName = ""
					continue
				}
				// start
				foundCode = true
				continue
			}
			if foundCode {
				code += sc.Text() + "\n"
			}
		}
	}
}

func envToMap() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		envMap[splitV[0]] = strings.Join(splitV[1:], "=")
	}
	return envMap
}
func paramToMap(p string) map[string]string {
	paramMap := make(map[string]string)

	if p == "" {
		return paramMap
	}

	for _, v := range strings.Split(p, ",") {
		splitV := strings.Split(v, "=")
		if len(splitV) < 2 {
			panic("-p params should be like 'FOO=BAR,AAA=BBB'.")
		}
		paramMap[splitV[0]] = splitV[1]
	}
	return paramMap
}
