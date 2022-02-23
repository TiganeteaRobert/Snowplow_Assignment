package schema

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
	dataAlreadyExists = `Provided ID already exists`
	invalidJSON       = `Invalid JSON`
	internalError     = `Internal server error`
	propertyMissing   = `required property is missing`

	// logging
	divider = `--------------`
)
