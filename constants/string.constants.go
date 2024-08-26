package constants

import "time"

const InternalServerError = "Internal Server Error"
const BadRequest = "Invalid Request"
const Unauthorized = "Unauthorized Access Detected"
const ValidationError = "Validation Error Occured"
const IdCannotBeZero = "ID cannot be zero"
const UUIDInvalid = "UUID is invalid"
const IdEmpty = "ID field is empty"
const ShortTimeout = time.Second * 5
const LongTimeout = time.Second * 7
const NotFound = "Could not find requested object"
