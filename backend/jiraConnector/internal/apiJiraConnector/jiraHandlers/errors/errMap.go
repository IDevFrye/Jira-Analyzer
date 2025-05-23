package errors

import (
	"errors"
	"net/http"
)

type errMap map[error]int

var (
	ErrorsUpdate = errMap{
		ErrNoProject:    http.StatusNotFound,
		ErrParamProject: http.StatusBadRequest,
		ErrUpdProject:   http.StatusInternalServerError,
		ErrPushProject:  http.StatusInternalServerError,
	}

	ErrorsProject = errMap{
		ErrParamLimitPage: http.StatusBadRequest,
		ErrGetProjectPage: http.StatusInternalServerError,
		ErrEncodeAns:      http.StatusInternalServerError,
	}
)

func GetStatusCode(m errMap, err error) int {
	if err == nil {
		return http.StatusOK
	}

	for e, c := range m {
		if errors.Is(err, e) {
			return c
		}
	}

	return http.StatusInternalServerError
}
