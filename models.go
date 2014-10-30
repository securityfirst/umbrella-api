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
	Id         int64  `json:"id" db:"id"`
	Title      string `json:"title" db:"title" binding:"required"`
	Text       string `json:"text" db:"text" binding:"required"`
	Value      int64  `json:"value" db:"value" binding:"required"`
	Parent     int64  `json:"parent" db:"parent" binding:"required"`
	Category   int64  `json:"category" db:"category" binding:"required"`
	Status     string `json:"-" db:"status"`
	CreatedAt  int64  `json:"-" db:"created_at"`
	ApprovedAt int64  `json:"-" db:"approved_at"`
	ApprovedBy int64  `json:"-" db:"approved_by"`
	Author     int64  `json:"-" db:"author"`
}

type JSONCheckItem struct {
	Title    string `json:"title" db:"title"`
	Text     string `json:"text" db:"text"`
	Value    int64  `json:"value" db:"value"`
	Parent   int64  `json:"parent" db:"parent"`
	Category int64  `json:"category" db:"category"`
	Status   string `json:"status" db:"status"`
}

type Segment struct {
	Id         int64  `json:"id" db:"id"`
	Title      string `json:"title" db:"title" binding:"required"`
	SubTitle   string `json:"subtitle" db:"subtitle" binding:"required"`
	Body       string `json:"body" db:"body" binding:"required"`
	Category   int64  `json:"category" db:"category" binding:"required"`
	Status     string `json:"-" db:"status"`
	CreatedAt  int64  `json:"-" db:"created_at"`
	ApprovedAt int64  `json:"-" db:"approved_at"`
	ApprovedBy int64  `json:"-" db:"approved_by"`
	Author     int64  `json:"-" db:"author"`
}

type JSONSegment struct {
	Title    string `json:"title" db:"title"`
	SubTitle string `json:"subtitle" db:"subtitle"`
	Body     string `json:"body" db:"body"`
	Category int64  `json:"category" db:"category"`
	Status   string `json:"status" db:"status"`
}
