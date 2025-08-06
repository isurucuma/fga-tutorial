package main

import (
	"context"
	"encoding/json"
	"fmt"
	openfga "github.com/openfga/go-sdk"

	. "github.com/openfga/go-sdk/client"
)

// for the sake of simplicity, I have not handled errors in this example.
func main() {
	// Initialize OpenFGA client
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl: "http://localhost:8080",
	})
	if err != nil {
		panic(fmt.Sprintf("Error initializing OpenFGA client: %v", err))
	}

	// Create a new store
	resp, err := fgaClient.CreateStore(context.Background()).
		Body(ClientCreateStoreRequest{
			Name: "Document Management System",
		}).Execute()

	// Store the store ID for future use
	storeId := resp.Id

	// here we need to update the client configuration with the new store ID
	err = fgaClient.SetStoreId(storeId)

	// Configure our authorization model
	var authModelString = `{
		  "schema_version": "1.1",
		  "type_definitions": [
			{
			  "type": "user",
			  "relations": {},
			  "metadata": null
			},
			{
			  "type": "team",
			  "relations": {
				"member": {
				  "this": {}
				},
				"department": {
				  "this": {}
				}
			  },
			  "metadata": {
				"relations": {
				  "member": {
					"directly_related_user_types": [
					  {
						"type": "user"
					  }
					]
				  },
				  "department": {
					"directly_related_user_types": [
					  {
						"type": "department"
					  }
					]
				  }
				}
			  }
			},
			{
			  "type": "department",
			  "relations": {
				"member": {
				  "tupleToUserset": {
					"computedUserset": {
					  "relation": "member"
					},
					"tupleset": {
					  "relation": "team"
					}
				  }
				},
				"team": {
				  "this": {}
				}
			  },
			  "metadata": {
				"relations": {
				  "member": {
					"directly_related_user_types": []
				  },
				  "team": {
					"directly_related_user_types": [
					  {
						"type": "team"
					  }
					]
				  }
				}
			  }
			},
			{
			  "type": "document",
			  "relations": {
				"owner": {
				  "this": {}
				},
				"editor": {
				  "union": {
					"child": [
					  {
						"this": {}
					  },
					  {
						"computedUserset": {
						  "relation": "owner"
						}
					  }
					]
				  }
				}
			  },
			  "metadata": {
				"relations": {
				  "owner": {
					"directly_related_user_types": [
					  {
						"type": "user"
					  }
					]
				  },
				  "editor": {
					"directly_related_user_types": [
					  {
						"type": "user"
					  },
					  {
						"type": "team",
						"relation": "member"
					  },
					  {
						"type": "department",
						"relation": "member"
					  }
					]
				  }
				}
			  }
			}
		  ]
}`

	var body openfga.WriteAuthorizationModelRequest
	err = json.Unmarshal([]byte(authModelString), &body)

	data, err := fgaClient.WriteAuthorizationModel(context.Background()).
		Body(body).
		Execute()

	// Save the model ID for future use
	modelId := data.AuthorizationModelId

	// Update the client configuration with the new model ID
	err = fgaClient.SetAuthorizationModelId(modelId)

	options := ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}

	scenarioTuples := ClientWriteRequest{
		Writes: []ClientTupleKey{
			{
				User:     "user:alice",
				Relation: "owner",
				Object:   "document:doc-001",
			},
			{
				User:     "team:engineering#member",
				Relation: "editor",
				Object:   "document:doc-001",
			},
			{
				User:     "user:bob",
				Relation: "member",
				Object:   "team:engineering",
			},
			{
				User:     "team:engineering",
				Relation: "team",
				Object:   "department:product",
			},
			{
				User:     "department:product#member",
				Relation: "editor",
				Object:   "document:doc-002",
			},
			{
				User:     "team:technical-support#member",
				Relation: "editor",
				Object:   "document:doc-003",
			},
		},
	}
	_, err = fgaClient.Write(context.Background()).
		Body(scenarioTuples).
		Options(options).
		Execute()

	// Can Bob edit doc-001?
	checkOptions := ClientCheckOptions{
		AuthorizationModelId: &modelId,
	}

	checkBody := ClientCheckRequest{
		User:     "user:bob",
		Relation: "editor",
		Object:   "document:doc-001",
	}

	checkData, err := fgaClient.Check(context.Background()).
		Body(checkBody).
		Options(checkOptions).
		Execute()
	fmt.Printf("Can Bob edit doc-001? %v\n", *checkData.Allowed) // true

	// Can Bob edit doc-002?
	checkBody.Object = "document:doc-002"
	checkData, err = fgaClient.Check(context.Background()).
		Body(checkBody).
		Options(checkOptions).
		Execute()
	fmt.Printf("Can Bob edit doc-002? %v\n", *checkData.Allowed) // true

	// Can Bob edit doc-003?
	checkBody.Object = "document:doc-003"
	checkData, err = fgaClient.Check(context.Background()).
		Body(checkBody).
		Options(checkOptions).
		Execute()
	fmt.Printf("Can Bob edit doc-003? %v\n", *checkData.Allowed) // false

	// Can Alice edit doc-001?
	checkBody.User = "user:alice"
	checkBody.Object = "document:doc-001"
	checkData, err = fgaClient.Check(context.Background()).
		Body(checkBody).
		Options(checkOptions).
		Execute()
	fmt.Printf("Can Alice edit doc-001? %v\n", *checkData.Allowed) // true

	// other than check you can issue query operations like the following
	// to get the list of documents that Bob can edit
	// to get the list of teams that can edit doc-001

}
