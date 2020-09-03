package openapi

// As go-swagger already returns a lot of different errors (like auth or
// validation) on it's own in specific format (usual HTTP Status Code plus
// body with JSON like {"code":600,"message":"some error"}) and already
// uses 2 error codes (HTTP Status and value of "code" field in JSON),
// which may be the same for some errors (like 404) or differ for others
// (like 422) - we should mimic this behaviour and also provide 2 codes
// for each of our own errors.
type errCode struct {
	status int   // HTTP Status Code.
	extra  int32 // Code for use in JSON body, may differ from HTTP Status Code.
}

// NewErrCode _MUST_ be used to create all used error codes, because it
// also registers each statusCode as a label for metrics.
//
// If extraCode is 0 then it'll be set to statusCode.
//
// As go-swagger already uses 6xx codes it's recommended to set extraCode
// to either 0 or >=700 to avoid conflicts.
func newErrCode(statusCode int, extraCode int32) errCode {
	codeLabels = append(codeLabels, statusCode)
	if extraCode == 0 {
		extraCode = int32(statusCode)
	}
	return errCode{status: statusCode, extra: extraCode}
}

// All error codes used by handlers should be declared here.
//
//nolint:gochecknoglobals,gomnd // Const.
var (
	codeInternal      = newErrCode(500, 0)
	codeContactExists = newErrCode(409, 1000)
)
