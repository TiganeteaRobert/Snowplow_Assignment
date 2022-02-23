# Snowplow JSON Schema Validator

Snowplow is a JSON schema validator capable of storing JSON schemas, retrieving them, and validating JSON input using JSON schemas.

## Requirements

```docker-compose```

## Installation

Simply run the following command inside the main directory

```bash
docker-compose up
```

# API Endpoints

## Create schema
Method: ```POST /schema/SCHEMAID```
### Examples
#### - Successful schema creation
#### Request
```POST /schema/config_schema_1```

##### Body
```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "source": {
      "type": "string"
    },
    "destination": {
      "type": "string"
    },
    "timeout": {
      "type": "integer",
      "minimum": 0,
      "maximum": 32767
    },
    "chunks": {
      "type": "object",
      "properties": {
        "size": {
          "type": "integer"
        },
        "number": {
          "type": "integer"
        }
      },
      "required": ["size"]
    }
  },
  "required": ["source", "destination"]
}
```
#### Response
```Status 200 OK```
```json
{
    "action": "uploadSchema",
    "id": "2",
    "status": "success"
}
```

#### - Invalid JSON schema input
#### Request
```POST /schema/config_schema_1```

##### Body
```json
{
  "$schema": "http://json-schema.org/draft-04/schema#"
  "type": "object",
  "properties": {
    "source": {
      "type": "string"
    },
    "destination": {
      "type": "string"
    },
    "timeout": {
      "type": "integer",
      "minimum": 0,
      "maximum": 32767
    },
    "chunks": {
      "type": "object",
      "properties": {
        "size": {
          "type": "integer"
        },
        "number": {
          "type": "integer"
        }
      },
      "required": ["size"]
    }
  },
  "required": ["source", "destination"]
}
```
#### Response
```Status 400 Bad Request```
```json
{
  "action": "uploadSchema",
  "id": "11",
  "status": "error",
  "message": "Invalid JSON"
}
```

#### - Schema ID already exists
#### Request
```POST /schema/config_schema_1```

##### Body
```json
{
  "$schema": "http://json-schema.org/draft-04/schema#"
  "type": "object",
  "properties": {
    "source": {
      "type": "string"
    },
    "destination": {
      "type": "string"
    },
    "timeout": {
      "type": "integer",
      "minimum": 0,
      "maximum": 32767
    },
    "chunks": {
      "type": "object",
      "properties": {
        "size": {
          "type": "integer"
        },
        "number": {
          "type": "integer"
        }
      },
      "required": ["size"]
    }
  },
  "required": ["source", "destination"]
}
```
#### Response
```Status 409 Conflict```
```json
{
  "action": "uploadSchema",
  "id": "11",
  "status": "error",
  "message": "Provided ID already exists"
}
```

## Get schema by ID
Method: ```GET /schema/SCHEMAID```
### Examples
#### - Found schema by ID
#### Request
```GET /schema/config_schema_1```
#### Response
```Content-Type: application/json; charset=utf-8```

```Status 200 OK```
```
"{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"properties\":{\"chunks\":{\"properties\":{\"number\":{\"type\":\"integer\"},\"size\":{\"type\":\"integer\"}},\"required\":[\"size\"],\"type\":\"object\"},\"destination\":{\"type\":\"string\"},\"source\":{\"type\":\"string\"},\"timeout\":{\"maximum\":32767,\"minimum\":0,\"type\":\"integer\"}},\"required\":[\"source\",\"destination\"],\"type\":\"object\"}"
```

#### - Schema not found
#### Request
```GET /schema/config_schema_1```
#### Response
```Status 404 Not Found```

## Validate JSON input using schema
Method: ```POST /schema/validate/SCHEMAID```

### Examples
#### - Successful validation of JSON input
#### Request
```POST /schema/validate/config_schema_1```

##### Body
```json
{
  "source": "/home/alice/image.iso",
  "destination": "/mnt/storage",
  "timeout": 100,
  "chunks": {
    "size": 1024,
    "number": 1
  }
}
```
#### Response
```Status 200 OK```
```json
{
    "action": "validateDocument",
    "id": "2",
    "status": "success"
}
```
#### - Validation errors 
#### Request
```POST /schema/validate/config_schema_1```

##### Body
```json
{
  "timeout": null,
  "chunks": {
    "size": "string value",
    "number": null
  }
}
```
#### Response
```Status 200 OK```
```json
{
    "action": "validateDocument",
    "id": "1",
    "status": "error",
    "message": {
        "/chunks/size": "expected integer, but got string",
        "missing properties": [
            "source",
            "destination"
        ]
    }
}
```
###### Empty properties are removed from the object

#### - Invalid JSON input
#### Request
```POST /schema/validate/config_schema_1```

##### Body
```json
{
  "source": "/home/alice/image.iso"
  "destination": "/mnt/storage",
  "timeout": null,
  "chunks": {
    "size": 1024,
    "number": "some string"
  }
}
```
#### Response
```Status 400 Bad Request```
```json
{
    "action": "uploadSchema",
    "id": "2",
    "status": "error",
    "message": "Invalid JSON"
}
```

#### - JSON schema not found
#### Request
```POST /schema/validate/config_schema_1```

##### Body
```json
{
  "source": "/home/alice/image.iso"
  "destination": "/mnt/storage",
  "timeout": null,
  "chunks": {
    "size": 1024,
    "number": "some string"
  }
}
```
#### Response
```Status 404 Not Found```

## License
[MIT](https://choosealicense.com/licenses/mit/)