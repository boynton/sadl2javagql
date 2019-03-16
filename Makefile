all:: bin/sadl2javagql
	./bin/sadl2javagql -dir gen -package model -pom -server examples/swapi.sadl
	cp examples/javaswapiimpl/Main.java gen/src/main/java
	(cd gen; mvn compile && mvn exec:java)

bin/sadl2javagql::
	mkdir -p bin
	go build -o bin/sadl2javagql github.com/boynton/sadl-gql/cmd/sadl2javagql

clean::
	rm -rf gen bin

proper::
	go fmt github.com/boynton/sadl-gql/graphql
	go vet github.com/boynton/sadl-gql/graphql
	go fmt github.com/boynton/sadl-gql/cmd/sadl2javagql
	go vet github.com/boynton/sadl-gql/cmd/sadl2javagql
