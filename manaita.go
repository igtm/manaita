package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/kevinburke/ssh_config"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
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

	GoStyleRawContentDefaultBrand = "master" // NOTE: master always redirects to default branch only not when master exists but other is default branch.
	GoStyleRawContentBaseURL      = "https://raw.githubusercontent.com/%s/%s/%s/%s"
	GitCloneHTTPBaseURL           = "https://github.com/%s/%s.git"
	GitCloneSSHBaseURL            = "git@github.com:%s/%s.git"
)

var (
	c = flag.String("c", DefaultScaffoldFileName, "specify markdown scaffold file path. default name is 'SCAFFOLD.md'")
	p = flag.String("p", "", "specify parameters for scaffold template. these must be defined on markdown  e.g. '-p foo=bar,fizz=buzz'")

	// source patterns
	githubGoGetStyleSourceRegexp, _ = regexp.Compile("^github\\.com/([\\w\\-.]+)/([\\w\\-.]+)/([\\w\\-./]+)(@[\\w\\-./]+)?$")
	httpSourceRegexp, _             = regexp.Compile("^https?://.+")

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

	errlog = log.New(os.Stderr, "ERROR: ", 0)
)

func main() {
	flag.Parse()

	envMap := envToMap()
	currentDir, _ := os.Getwd()
	givenParamMap := paramToMap(*p)

	// detect source file path type
	var reader io.Reader
	if ms := githubGoGetStyleSourceRegexp.FindStringSubmatch(*c); len(ms) >= 4 {
		// go get style github source
		// trying to get from public raw content url
		owner := ms[1]
		repo := ms[2]
		file := ms[3]
		branch := GoStyleRawContentDefaultBrand
		var specifiedBrand bool
		if len(ms) == 5 && ms[4] != "" {
			specifiedBrand = true
			branch = strings.Replace(ms[4], "@", "", 1)
		}
		resp, err := http.Get(fmt.Sprintf(GoStyleRawContentBaseURL, owner, repo, branch, file))
		defer resp.Body.Close()
		reader = resp.Body
		if err != nil || resp.StatusCode != http.StatusOK {
			// fallback to git clone (SSH)
			// Filesystem abstraction based on memory
			fs := memfs.New()
			// Git objects storer based on memory
			storer := memory.NewStorage()

			githubIdentityFile := ssh_config.Get("github.com", "IdentityFile")
			publicKeys, err := ssh.NewPublicKeysFromFile("git", githubIdentityFile, "")
			if err != nil {
				errlog.Println("git clone ssh error: %s", err)
				return
			}

			// Clones the repository into the worktree (fs) and stores all the .git
			sshPath := fmt.Sprintf(GitCloneSSHBaseURL, owner, repo)
			cloneOptions := &git.CloneOptions{
				URL:  sshPath,
				Auth: publicKeys,
				Tags: git.AllTags,
			}
			if specifiedBrand {
				// branch
				cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branch)
				_, err = git.Clone(storer, fs, cloneOptions)
				// FIXME: tag cannot be cloned correctly... maybe library problem..
				//if err != nil {
				//	// tag
				//	cloneOptions.ReferenceName = plumbing.NewTagReferenceName(branch)
				//	_, err = git.Clone(storer, fs, cloneOptions)
				//}
			} else {
				// default branch
				_, err = git.Clone(storer, fs, cloneOptions)
			}
			if err != nil {
				errlog.Printf("cannot git clone '%s': %s\n", sshPath, err)
				return
			}

			// Prints the content of the CHANGELOG file from the cloned repository
			scaffoldFile, err := fs.Open(file)
			if err != nil {
				errlog.Println(fmt.Errorf("cannot open file '%s'", file))
				return
			}

			reader = scaffoldFile
		}
	} else if ms := httpSourceRegexp.FindAllString(*c, -1); len(ms) == 1 {
		// http source
		resp, err := http.Get(ms[0])
		if err != nil {
			errlog.Println(fmt.Errorf("cannot get '%s'", ms[0]))
			return
		}
		if resp.StatusCode != http.StatusOK {
			errlog.Println(fmt.Errorf("cannot get '%s': status code: %d", ms[0], resp.StatusCode))
			return
		}
		defer resp.Body.Close()
		reader = resp.Body
	} else {
		// trying to search local file
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
		reader = scaffoldFile
	}

	sc := bufio.NewScanner(reader)

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
			err := tmpl.Execute(&compiledDest, map[string]interface{}{
				"Env":    envMap,
				"Params": paramMap,
			})
			if err != nil {
				errlog.Println(fmt.Errorf("cannot compile template"))
				return
			}
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
