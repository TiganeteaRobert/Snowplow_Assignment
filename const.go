package main

const (
	// schema
	schemaID       = `SCHEMAID`
	schemaEndpoint = `/:SCHEMAID`

	// Response actions
	uploadSchema = `uploadSchema`
	validateDoc  = `validateDocument`

	// Response statuses
	successStatus = `success`
	errorStatus   = `error`

	// Response error messages
	dataAlreadyExists  = `Provided ID already exists`
	invalidJSON        = `Invalid JSON`
	internalError      = `Internal server error`
	propertiesRequired = "The following required properties are missing: %s"
	wrongTypes         = "\nThe following properties' types are incorrect: %s"

	// logging
	divider = `--------------`
)
