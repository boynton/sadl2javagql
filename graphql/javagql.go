package graphql

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/boynton/sadl"
)

func (gen *Generator) CreateGraphqlServer() {
	if gen.Err != nil {
		return
	}
	gen.CreateServerDataAndFuncMap(gen.SourceDir, gen.ResourceDir)
	gen.ServerData.ExtraResources = gen.GraphqlResourceAsString()
	gen.CreateServer()
	if gen.Err != nil {
		return
	}
	rezDir := filepath.Join(gen.OutDir, gen.ResourceDir)
	gen.CreateGraphqlSchema(rezDir)
	gen.CreateGraphqlHandler()
	gen.CreateGraphqlRequestPojo()
	gen.CreateGraphqlResponsePojo()
}

const graphqlResource = `    @POST
    @Path("/graphql")
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public GraphqlResponse query(GraphqlRequest req) throws Exception {
        try {
            GraphqlResponse res = {{graphqlClass}}.execute(req, {{implName .Name}});
            return res;
        } catch (Exception e) {
            e.printStackTrace();
            throw e;
        }
    }
`

func (gen *Generator) GraphqlClass() string {
	return gen.Capitalize(gen.Model.Name) + "Graphql"
}

func (gen *Generator) GraphqlResourceAsString() string {
	if gen.Graphql != nil {
		s := strings.Replace(graphqlResource, "{{graphqlClass}}", gen.GraphqlClass(), -1)
		s = strings.Replace(s, "{{implName .Name}}", gen.Uncapitalize(gen.Name)+"Controller", -1)
		return s
	}
	return ""
}

func (gen *Generator) CreateGraphqlHandler() {
	if gen.Err != nil {
		return
	}
	serviceName := gen.Capitalize(gen.Model.Name)
	className := serviceName + "Graphql"
	funcMap := template.FuncMap{
		"openBrace": func() string { return "{" },
		"graphqlFetchers": func() string {
			return gen.GraphqlFetchers()
		},
		"graphqlClass": func() string {
			return className
		},
	}
	gen.CreateJavaFileFromTemplate(className, graphqlHandlerTemplate, gen.ServerData, funcMap, gen.ModelPackage)
}

const graphqlHandlerTemplate = `
import java.util.Map;
import java.util.List;
import java.util.ArrayList;

import graphql.ExecutionInput;
import graphql.ExecutionResult;
import graphql.GraphQL;
import graphql.GraphQLError;
import graphql.schema.GraphQLSchema;
import graphql.schema.DataFetcher;
import graphql.schema.DataFetchingEnvironment;
import graphql.schema.idl.RuntimeWiring;
import graphql.schema.idl.SchemaGenerator;
import graphql.schema.idl.SchemaParser;
import graphql.schema.idl.TypeDefinitionRegistry;

import static graphql.schema.idl.RuntimeWiring.newRuntimeWiring;

public class {{graphqlClass}} {

    static String getSchema(String name) {
        try {
            ClassLoader classLoader = new {{graphqlClass}}().getClass().getClassLoader();
            java.net.URL url = classLoader.getResource(name);
            if (url != null) {
                java.io.File file = new java.io.File(url.getFile());
                return new String(java.nio.file.Files.readAllBytes(java.nio.file.Paths.get(file.getPath())), java.nio.charset.StandardCharsets.UTF_8);
            } else {
                throw new Exception("Cannot find resource: " + name);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return "";
    }

    public static GraphqlResponse execute(GraphqlRequest req, {{.InterfaceClass}} impl) {
        String query = req.getQuery();
        Map<String,Object> variables = req.getVariables();

        String schema = getSchema("schema.gql");
        
        SchemaParser schemaParser = new SchemaParser();
        TypeDefinitionRegistry typeDefinitionRegistry = schemaParser.parse(schema);
        
        RuntimeWiring runtimeWiring = newRuntimeWiring()
{{graphqlFetchers}}            .build();
        
        SchemaGenerator schemaGenerator = new SchemaGenerator();
        GraphQLSchema graphQLSchema = schemaGenerator.makeExecutableSchema(typeDefinitionRegistry, runtimeWiring);
        
        GraphQL build = GraphQL.newGraphQL(graphQLSchema).build();
        
        ExecutionInput.Builder input = new ExecutionInput.Builder().query(query);
        if (variables != null) {
            input.variables(variables);
        }
        ExecutionResult executionResult = build.execute(input);

        ArrayList<Object> lstErrors = new ArrayList<Object>();
        for (GraphQLError err : executionResult.getErrors()) {
            lstErrors.add(err.toSpecification());
        }
        GraphqlResponse res = GraphqlResponse.builder().data(executionResult.getData()).errors(lstErrors).build();
        System.out.println("=>\n" + Util.pretty(res));
        return res;
    }

    //Resolvers could be implemented as inner classes here

}
`

func (gen *Generator) GraphqlFetchers() string {
	gql := gen.Graphql
	indent := "            "
	result := ""
	for _, op := range gql.Actions {
		rType, _, _ := gen.TypeName(op.Return, op.Return.Type, true)
		indent2 := indent + "            "
		q := indent + `.type("Query", builder -> builder.dataFetcher("` + op.Name + `", new DataFetcher<` + rType + ">() {\n"
		q = q + indent2 + "public " + rType + " get(DataFetchingEnvironment env) throws Exception {\n"
		if len(op.Params) > 0 {
			for _, param := range op.Params {
				q = q + indent2 + "    " + param.Type + " " + param.Name + ` = env.getArgument("` + param.Name + "\");\n"
			}
		}
		q = q + indent2 + "    " + op.Provider + "Response res = impl." + gen.Uncapitalize(op.Provider) + "(" + op.Provider + "Request.builder()"
		if len(op.Params) > 0 {
			for _, param := range op.Params {
				q = q + "." + param.Name + "(" + param.Name + ")"
			}
		}
		q = q + ".build());\n"
		q = q + indent2 + "    return res.get" + gen.Capitalize(op.Name) + "();\n"
		q = q + indent2 + "}\n"
		result = result + q + indent + "        }))\n"
	}
	return result
}

func (gen *Generator) CreateGraphqlRequestPojo() {
	if gen.Err != nil {
		return
	}
	ts := &sadl.TypeSpec{
		Type: "Struct",
		Fields: []*sadl.StructFieldDef{
			&sadl.StructFieldDef{
				Name:     "query",
				Required: true,
				TypeSpec: sadl.TypeSpec{
					Type: "String",
				},
			},
			&sadl.StructFieldDef{
				Name: "operationName",
				TypeSpec: sadl.TypeSpec{
					Type: "String",
				},
			},
			&sadl.StructFieldDef{
				Name: "variables",
				TypeSpec: sadl.TypeSpec{
					Type:  "Map",
					Keys:  "String",
					Items: "Any",
				},
			},
		},
	}
	className := "GraphqlRequest"
	gen.CreatePojo(ts, className, "", nil)
}

func (gen *Generator) CreateGraphqlResponsePojo() {
	if gen.Err != nil {
		return
	}
	ts := &sadl.TypeSpec{
		Type: "Struct",
		Fields: []*sadl.StructFieldDef{
			&sadl.StructFieldDef{
				Name: "data",
				TypeSpec: sadl.TypeSpec{
					Type:  "Map",
					Keys:  "String",
					Items: "Any",
				},
			},
			&sadl.StructFieldDef{
				Name: "errors",
				TypeSpec: sadl.TypeSpec{
					Type:  "Array",
					Items: "Any",
				},
			},
		},
	}
	className := "GraphqlResponse"
	gen.CreatePojo(ts, className, "", nil)
}
