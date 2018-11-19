package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"flag"
	"strings"
)

var filename = flag.String("n", "app", "module name")

type Module struct {
	name string
	content string
}

type ModuleContent struct {
	content string
	isCreate bool
}

func defaultModuleContent(content string) ModuleContent{
	return ModuleContent{
		content,
		true,
	}
}

func toUpperCaseAtOne(str string) string {
	return strings.Join([]string{strings.ToUpper(string(str[0])), str[1:]}, "")
}

func contentFactory(content string, fn string) string {
	return fmt.Sprintf(content, processFileName(fn))
}

func processFileName(filename string) string {
	// skewer
	parts := strings.Split(filename, "-")

	for i := range parts {
		parts[i] = toUpperCaseAtOne(parts[i])
	}

	return strings.Join(parts, "")
}

func checkErr(err error){
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (m *Module) writeString(){
	file, err := os.Create( m.name)

	defer file.Close()

	checkErr(err)

	file.WriteString(m.content)
}

func entryQuestion(Questions map[string]ModuleContent){
	input := bufio.NewScanner(os.Stdin)

	for q := range Questions {
		fmt.Printf("%s(true): ", q)
		input.Scan()

		if v := input.Text(); v != "" {
			result, err := strconv.ParseBool(v)

			checkErr(err)

			Questions[q] = ModuleContent{
				"",
				result,
			}
		}
	}
}

func builder(Questions map[string]ModuleContent, filename string){
	m := Module{
		"",
		"",
	}

	// traverse questions for create file
	for moduleName, moduleContent := range Questions {
		filename := strings.Join([]string{filename, moduleName, "ts"}, ".")

		m.name = filename
		m.content = moduleContent.content

		if moduleName == "index" {
			m.name = strings.Join([]string{moduleName, "ts"}, ".")
		}

		if moduleContent.isCreate {
			m.writeString()
		}
	}
}


func main(){
	flag.Parse()

	fn := *filename

	moduleContent := contentFactory(strings.Join([]string{
		"import {Module} from '@nestjs/common';",
		"",
		"@Module({",
		"  imports:[],",
		"  providers:[],",
		"  controllers:[],",
		"})",
		"export class %sModule { }",
	}, "\n"), fn)

	serviceContent := contentFactory(strings.Join([]string{
		"import {Injectable} from '@nestjs/common';",
		"",
		"@Injectable()",
		"export class %sService { }",
	},"\n"), fn)

	entityContent := contentFactory(strings.Join([]string{
		"import {Entity} from 'typeorm';",
		"",
		"@Entity()",
		"export class %s { }",
	},"\n"), fn)

	controllerContent := contentFactory(strings.Join([]string{
		"import {Controller} from '@nestjs/common';",
		"",
		"@Controller('')",
		"export class %sController { }",
	},"\n"), fn)

	Questions := map[string]ModuleContent{
		"module": defaultModuleContent(moduleContent),
		"service": defaultModuleContent(serviceContent),
		"index": defaultModuleContent(""),
		"entity": defaultModuleContent(entityContent),
		"controller": defaultModuleContent(controllerContent),
	}

	entryQuestion(Questions)
	builder(Questions, fn)
}
