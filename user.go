package main

//go:generate reform

// User represents a row in users table.
//
//reform:users
type User struct {
	ID   int32  `param:"id" query:"id" form:"id" reform:"id,pk"`
	Name string `param:"name" query:"name" form:"name" reform:"name"`
	Age  *int32 `param:"age" query:"age" form:"age" reform:"age"`
}
