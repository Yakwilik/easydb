// Code generated by mockery v2.53.4. DO NOT EDIT.

package easydb

import (
	context "context"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"
)

// mockPgxquerier is an autogenerated mock type for the pgxquerier type
type mockPgxquerier struct {
	mock.Mock
}

type mockPgxquerier_Expecter struct {
	mock *mock.Mock
}

func (_m *mockPgxquerier) EXPECT() *mockPgxquerier_Expecter {
	return &mockPgxquerier_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, sql, args
func (_m *mockPgxquerier) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 pgconn.CommandTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)); ok {
		return rf(ctx, sql, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgconn.CommandTag); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		r0 = ret.Get(0).(pgconn.CommandTag)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockPgxquerier_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type mockPgxquerier_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *mockPgxquerier_Expecter) Exec(ctx interface{}, sql interface{}, args ...interface{}) *mockPgxquerier_Exec_Call {
	return &mockPgxquerier_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *mockPgxquerier_Exec_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *mockPgxquerier_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockPgxquerier_Exec_Call) Return(_a0 pgconn.CommandTag, _a1 error) *mockPgxquerier_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockPgxquerier_Exec_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)) *mockPgxquerier_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, sql, args
func (_m *mockPgxquerier) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 pgx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgx.Rows, error)); ok {
		return rf(ctx, sql, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Rows); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockPgxquerier_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type mockPgxquerier_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *mockPgxquerier_Expecter) Query(ctx interface{}, sql interface{}, args ...interface{}) *mockPgxquerier_Query_Call {
	return &mockPgxquerier_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *mockPgxquerier_Query_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *mockPgxquerier_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockPgxquerier_Query_Call) Return(_a0 pgx.Rows, _a1 error) *mockPgxquerier_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockPgxquerier_Query_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgx.Rows, error)) *mockPgxquerier_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: ctx, sql, args
func (_m *mockPgxquerier) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRow")
	}

	var r0 pgx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Row); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Row)
		}
	}

	return r0
}

// mockPgxquerier_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type mockPgxquerier_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *mockPgxquerier_Expecter) QueryRow(ctx interface{}, sql interface{}, args ...interface{}) *mockPgxquerier_QueryRow_Call {
	return &mockPgxquerier_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *mockPgxquerier_QueryRow_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *mockPgxquerier_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockPgxquerier_QueryRow_Call) Return(_a0 pgx.Row) *mockPgxquerier_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockPgxquerier_QueryRow_Call) RunAndReturn(run func(context.Context, string, ...interface{}) pgx.Row) *mockPgxquerier_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

// newMockPgxquerier creates a new instance of mockPgxquerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockPgxquerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockPgxquerier {
	mock := &mockPgxquerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
