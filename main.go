package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"flag"
	"strings"
)

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
	return fmt.Sprintf(content, toUpperCaseAtOne(fn))
}

var (

	filename = flag.String("n", "app", "-n [name]")
)

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

func builder(Questions map[string]ModuleContent){
	m := Module{
		"",
		"",
	}

	// traverse questions for create file
	for moduleName, moduleContent := range Questions {
		filename := strings.Join([]string{*filename, moduleName, "ts"}, ".")

		m.name = filename
		m.content = moduleContent.content

		if moduleContent.isCreate {
			m.writeString()
		}
	}
}


func main(){
	flag.Parse()

	fn := *filename

	moduleContent := contentFactory(strings.Join([]string{
		"@Module(",
		"imports:[],",
		"providers:[],",
		"controllers:[],",
		")",
		"export class %sModule { }",
	}, "\n"), fn)

	serviceContent := contentFactory(strings.Join([]string{
		"@Injectable()",
		"export class %sService { }",
	},"\n"), fn)

	entityContent := contentFactory(strings.Join([]string{
		"@Entity()",
		"export class %s { }",
	},"\n"), fn)

	controllerContent := contentFactory(strings.Join([]string{
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
	builder(Questions)
}
