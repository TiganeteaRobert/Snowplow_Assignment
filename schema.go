package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// addSchema adds a new schema to Redis
func addSchema(c *gin.Context) {
	InfoLogger.Println(divider + `add schema`)
	schemaID := c.Param(schemaID)
	InfoLogger.Printf(`attempting to add new schema. ID %s`, schemaID)

	var body interface{}
	// move the request body to an interface
	// passing this point means that the body is in correct JSON format
	if err := c.ShouldBindJSON(&body); err != nil {
		ErrorLogger.Printf(`invalid JSON, responding with Bad Request status. error: %v`, err)
		c.IndentedJSON(http.StatusBadRequest, Response{
			ID:      schemaID,
			Action:  uploadSchema,
			Status:  errorStatus,
			Message: invalidJSON,
		})
		return
	}

	// marshal the request body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		ErrorLogger.Printf(`error while marshalling body. error: %v`, err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	// check if the provided schema ID already exists
	exists, err := redisClient.Exists(schemaID).Result()
	if err != nil {
		ErrorLogger.Printf(`error while checking for existing schema in Redis. error: %v`, err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	if exists == 1 {
		InfoLogger.Printf(`found already existing schema with ID %s`, schemaID)
		c.IndentedJSON(http.StatusOK, Response{
			ID:      schemaID,
			Action:  uploadSchema,
			Status:  errorStatus,
			Message: dataAlreadyExists,
		})
		return
	}

	// upload the schema's JSON body to Redis using the schema ID as the key
	if err = redisClient.Set(schemaID, jsonBody, 0).Err(); err != nil {
		ErrorLogger.Printf(`error while adding new schema in Redis. error: %v`, err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	InfoLogger.Printf(`added new schema to Redis. Key: %s. Responding with success status`, schemaID)

	// return success
	c.IndentedJSON(http.StatusOK, Response{
		ID:     schemaID,
		Action: uploadSchema,
		Status: successStatus,
	})
}

// getSchema retrieves a schema from Redis
func getSchema(c *gin.Context) {
	InfoLogger.Println(divider + `get schema`)
	schemaID := c.Param(schemaID)

	InfoLogger.Printf(`attempting to get schema from Redis. ID: %s`, schemaID)

	// get the schema from Redis using the schema ID

	schema := redisClient.Get(schemaID)
	InfoLogger.Println(`requested schema from Redis`)

	// check if there have been errors getting the schema from Redis (except Nil error)
	if schema.Err() != nil && schema.Err().Error() != redis.Nil.Error() {
		ErrorLogger.Printf(`error while getting schema from Redis. error: %v`, schema.Err())
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	// check if the schema has been found
	if schema.Val() == `` {
		InfoLogger.Println(`schema has not been found.`)
		c.IndentedJSON(http.StatusNotFound, schema)
		return
	}

	// respond with schema
	InfoLogger.Println(`retrieved schema from Redis successfully`)
	c.IndentedJSON(http.StatusCreated, schema.Val())
}

// validateJSON validates a JSON body using a JSON schema
func validateJSON(c *gin.Context) {
	InfoLogger.Println(divider + `validate`)
	schemaID := c.Param(schemaID)

	InfoLogger.Printf(`attempting to validate JSON against schema. ID %s`, schemaID)

	// get the schema from Redis using the schema ID
	schema := redisClient.Get(schemaID)
	InfoLogger.Println(`requested schema from Redis`)

	// check if there have been errors getting the schema from Redis (except Nil error)
	if schema.Err() != nil {
		ErrorLogger.Printf(`error while getting schema from Redis. error: %v`, schema.Err())
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	// check if the schema has been found
	if schema.Val() == `` {
		InfoLogger.Println(`schema has not been found.`)
		c.IndentedJSON(http.StatusNotFound, schema)
		return
	}

	InfoLogger.Println(`retrieved schema from Redis successfully, compiling it`)

	// compile the retrieved JSON schema
	val, err := jsonschema.CompileString("schema.json", schema.Val())
	if err != nil {
		ErrorLogger.Printf(`error while compiling schema. error: %v`, err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	// move the request body to an interface
	// passing this point means that the body is in correct JSON format
	var body interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		InfoLogger.Printf(`request body is an invalid JSON, responding with a Bad Request status. validation error: %v`, err)
		c.IndentedJSON(http.StatusBadRequest, Response{
			ID:      schemaID,
			Action:  uploadSchema,
			Status:  errorStatus,
			Message: invalidJSON,
		})
		return
	}

	InfoLogger.Println(`validating request JSON against schema`)

	// validate the JSON request body using the JSON schema
	validationResult := val.Validate(body)

	if validationResult != nil {
		res := validationResult.(*jsonschema.ValidationError)
		errMessages := make(map[string]string, 0)
		extractValidationErrs(res.Causes, &errMessages)
		if len(errMessages) != 0 {
			InfoLogger.Printf(`found %d validation errors agains the schema, sending response`, len(errMessages))
			c.IndentedJSON(http.StatusOK, Response{
				Action:  validateDoc,
				ID:      schemaID,
				Status:  errorStatus,
				Message: errMessages,
			})
			return
		}
	}

	InfoLogger.Printf(`validation successful, sending response`)
	// return success response
	c.IndentedJSON(http.StatusOK, Response{
		Action: validateDoc,
		ID:     schemaID,
		Status: successStatus,
	})
}

// extractValidationErrs iterates through a ValidationError object and recursively extracts
// all validation errors
func extractValidationErrs(causes []*jsonschema.ValidationError, allCauses *map[string]string) {
	for _, cause := range causes {
		if cause.Causes != nil {
			extractValidationErrs(cause.Causes, allCauses)
		}
		if cause.Message != `` {
			c := *allCauses
			if strings.Split(cause.Message, `got `)[1] == `null` {
				InfoLogger.Printf(`found null field %s, removing from object`, cause.InstanceLocation)
			}
			c[cause.InstanceLocation] = cause.Message
			*allCauses = c
		}
	}
}
