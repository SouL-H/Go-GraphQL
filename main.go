package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

func populate() []Tutorial {
	author := &Author{Name: "TestData", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "TestData",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "TestData"},
		},
	}
	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials

}
func main() {
	tutorials := populate()
	fmt.Println("Hello GraphQL")
	var commnetType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)
	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commnetType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get Tutorial By Id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Full Tutorial List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL Schema,err %v", err)
	}
	query := `
		{
			tutorial(id:1){
				title
				author{
					Name
					Tutorials
				}
			}
		}
	`

	
	// oldQuery := `
	// 	{
	// 		list{
	// 			id
	// 			title
	// 		}
	// 	}
	// `
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation,errors:%+v", r.Errors)
	}
	rJSON, err := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
