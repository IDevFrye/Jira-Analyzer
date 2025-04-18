package errors

import "errors"

var (
	ErrParamLimitPage = errors.New("incorrect limit or page param - need integer > 0")
	ErrParamProject   = errors.New("we don't have project with such name")

	ErrUpdProject  = errors.New("something went wrong and i can't update project")
	ErrPushProject = errors.New("something went wrong and i can't save issues")

	ErrGetProjectPage = errors.New("something went wrong and i can't get page of projects")

	ErrEncodeAns = errors.New("something went wrong and i can't encode ans for this request")

	ErrNoProject = errors.New("jira doesn't have such project")
)
