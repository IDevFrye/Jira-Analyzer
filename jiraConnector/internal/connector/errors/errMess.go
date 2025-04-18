package errors

import "errors"

var (
	ErrMakeRequest    = errors.New("error make request")
	ErrMaxTimeRequest = errors.New("unsucsess request - the maximum request execution time has been reached")

	ErrReadResponseBody = errors.New("can't read responce body")
	ErrUnmarshalAns     = errors.New("can't unmarshal responce body")

	ErrGetIssues   = errors.New("can't get issues")
	ErrGetProjects = errors.New("can't get project")
)
