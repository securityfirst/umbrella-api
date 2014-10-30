package main

type LoginJSON struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

type CheckItem struct {
	Id       int64  `json:"id" db:"id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Value    int    `json:"value"`
	Parent   int    `json:"parent"`
	Category int    `json:"category"`
}

type Segment struct {
	Id       int64  `json:"id" db:"id"`
	Title    string `json:"title" db:"title" binding:"required"`
	SubTitle string `json:"subtitle" db:"subtitle" binding:"required"`
	Body     string `json:"body" db:"body" binding:"required"`
	Category int    `json:"category" db:"category" binding:"required"`
}

type EditSegment struct {
	Id       int64  `json:"id" db:"id" binding:"required"`
	Title    string `json:"title" db:"title"`
	SubTitle string `json:"subtitle" db:"subtitle"`
	Body     string `json:"body" db:"body"`
	Category int    `json:"category" db:"category"`
}
