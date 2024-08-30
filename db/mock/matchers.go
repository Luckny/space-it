package mockdb

import (
	"fmt"
	reflect "reflect"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/util"
	gomock "go.uber.org/mock/gomock"
)

// ----> Register user matchers

type eqRegisterUserParamsMatcher struct {
	arg db.RegisterUserParams
}

func (e eqRegisterUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.RegisterUserParams)
	if !ok {
		return false
	}

	if err := util.CheckPassword(e.arg.Password, arg.Password); err != nil {
		return false
	}

	return reflect.DeepEqual(e.arg.Email, arg.Email)
}

func (e eqRegisterUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqRegisterUserParams(arg db.RegisterUserParams) gomock.Matcher {
	return eqRegisterUserParamsMatcher{arg}
}

// ----> Audit loggin matchers

// Unauthenticated Request matchers

type eqUnauthenticatedLogParam struct {
	arg db.CreateUnauthenticatedRequestLogParams
}

func (e eqUnauthenticatedLogParam) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUnauthenticatedRequestLogParams)
	if !ok {
		return false
	}

	return reflect.DeepEqual(e.arg.Method, arg.Method) &&
		reflect.DeepEqual(e.arg.Path, arg.Path)
}

func (e eqUnauthenticatedLogParam) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUnauthenticatedLogParam(
	arg db.CreateUnauthenticatedRequestLogParams,
) gomock.Matcher {
	return eqUnauthenticatedLogParam{arg}
}

// Authenticated request matcher
type eqAuthenticatedLogParam struct {
	arg db.CreateAuthenticatedRequestLogParams
}

func (e eqAuthenticatedLogParam) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAuthenticatedRequestLogParams)
	if !ok {
		return false
	}

	return reflect.DeepEqual(e.arg.Method, arg.Method) &&
		reflect.DeepEqual(e.arg.Path, arg.Path)
}

func (e eqAuthenticatedLogParam) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqAuthenticatedLogParam(
	arg db.CreateAuthenticatedRequestLogParams,
) gomock.Matcher {
	return eqAuthenticatedLogParam{arg}
}

// Response Log matcher

type eqResponseLogParam struct {
	arg db.CreateResponseLogParams
}

func (e eqResponseLogParam) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateResponseLogParams)
	if !ok {
		return false
	}

	return reflect.DeepEqual(e.arg.Status, arg.Status) &&
		reflect.DeepEqual(e.arg.ID, arg.ID)
}

func (e eqResponseLogParam) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqResponseLogParam(
	arg db.CreateResponseLogParams,
) gomock.Matcher {
	return eqResponseLogParam{arg}
}
