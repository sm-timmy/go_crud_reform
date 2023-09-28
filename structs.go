package main

type user struct {
	Id   int    `param:"id" query:"id" form:"id" json:"id" `
	Name string `param:"name" query:"name" form:"name" json:"name"`
	Age  int    `param:"age" query:"age" form:"age" json:"age" `
}
