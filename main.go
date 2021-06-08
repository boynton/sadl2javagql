package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/boynton/sadl"
	"github.com/boynton/sadl2javagql/graphql"
)

func main() {
	pDir := flag.String("dir", ".", "output directory for generated artifacts")
	pSrc := flag.String("src", "src/main/java", "output directory for generated source tree, relative to dir")
	pRez := flag.String("rez", "src/main/resources", "output directory for generated resources, relative to dir")
	pPackage := flag.String("package", "", "Java package for generated source")
	pServer := flag.Bool("server", false, "generate server code")
	pLombok := flag.Bool("lombok", false, "generate Lombok annotations")
	pGetters := flag.Bool("getters", false, "generate setters/getters instead of the default fluent style")
	pInstants := flag.Bool("instants", false, "Use java.time.Instant. By default, use generated Timestamp class")
	pPom := flag.Bool("pom", false, "Create Maven pom.xml file to build the project")
	flag.Parse()
	argv := flag.Args()
	argc := len(argv)
	if argc == 0 {
		fmt.Fprintf(os.Stderr, "usage: sadl2java -dir projdir -src relative_source_dir -rez relative_resource_dir -package java.package.name -pom -server -getters -lombok some_model.sadl\n")
		os.Exit(1)
	}
	config := sadl.NewData()
	config.Put("model", true)
	config.Put("server", true)
	config.Put("example-implementation", false)
	if *pSrc != "" {
		config.Put("source", *pSrc)
	}
	if *pRez != "" {
		config.Put("resource", *pRez)
	}
	if *pPackage != "" {
		config.Put("package", *pPackage)
	}
	if *pLombok {
		config.Put("lombok", true)
	}
	if *pGetters {
		config.Put("getters", true)
	}
	if *pInstants {
		config.Put("instants", true)
	}
	model, err := sadl.ParseSadlFile(argv[0], config, graphql.NewExtension())
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** %v\n", err)
		os.Exit(1)
	}
	gen := graphql.NewGenerator(model, *pDir, config)
	gen.CreateModel()
	/*	for _, td := range model.Types {
			gen.CreatePojoFromDef(td, nil)
		}
		if gen.NeedTimestamp {
			gen.CreateTimestamp()
		}
		gen.CreateUtil()
	*/
	if gen.Err != nil {
		fmt.Fprintf(os.Stderr, "*** %v\n", err)
		os.Exit(1)
	}
	if *pServer {
		gen.CreateGraphqlServer()
	}
	if *pPom {
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			domain = "my.domain"
		}
		gen.CreatePom(graphqlDepends)
	}
	if gen.Err != nil {
		fmt.Fprintf(os.Stderr, "*** %v\n", gen.Err)
		os.Exit(1)
	}
}

const graphqlDepends = `      <dependency>
        <groupId>com.graphql-java</groupId>
        <artifactId>graphql-java</artifactId>
        <version>2019-02-20T00-59-31-9356c3d</version>
      </dependency>
`
