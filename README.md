# sadl2javagql
An example of a SADL extension to handle GraphQL queries.

Use the sadl2javagql command just like the sadl2java command. It additionally supports the "graphql"
parser extension (see the end of the swapi.sadl example file), and generates server glue to support its operation
as an additional /graphql endpoint in addition to the normal HTTP endpoints.

See examples/javaswapiimpl for more details on running it and how to query against it.

