package main

import "github.com/go-gorp/gorp"

type Umbrella struct {
	Db *gorp.DbMap
}

type LoginJSON struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id       int64  `db:"id" json:"-"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `json:"-"`
	Token    string `json:"-"`
	Role     int    `db:"role" json:"-"`
}

type CheckItem struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Value      int64  `json:"value"`
	Parent     int64  `json:"parent"`
	Category   int64  `json:"category"`
	Difficulty int64  `json:"difficulty"`
	Custom     int64  `json:"custom"`
	Disabled   int64  `json:"disabled"`
	NoCheck    int64  `json:"no_check" db:"no_check"`
	Status     string `json:"-" db:"status"`
	CreatedAt  int64  `json:"-" db:"created_at"`
	ApprovedAt int64  `json:"-" db:"approved_at"`
	ApprovedBy int64  `json:"-" db:"approved_by"`
	Author     int64  `json:"-" db:"author"`
}

// type JSONCheckItem struct {
// 	Title    string `json:"title" db:"title"`
// 	Text     string `json:"text" db:"text"`
// 	Value    int64  `json:"value" db:"value"`
// 	Parent   int64  `json:"parent" db:"parent"`
// 	Category int64  `json:"category" db:"category"`
// 	Status   string `json:"status" db:"status"`
// }

type Segment struct {
	Id         int64  `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	SubTitle   string `json:"subtitle" db:"subtitle"`
	Body       string `json:"body" db:"body"`
	Category   int64  `json:"category" db:"category"`
	Difficulty int64  `json:"difficulty" db:"difficulty"`
	Status     string `json:"-" db:"status"`
	CreatedAt  int64  `json:"-" db:"created_at"`
	ApprovedAt int64  `json:"-" db:"approved_at"`
	ApprovedBy int64  `json:"-" db:"approved_by"`
	Author     int64  `json:"-" db:"author"`
}

// type JSONSegment struct {
// 	Title    string `json:"title" db:"title"`
// 	SubTitle string `json:"subtitle" db:"subtitle"`
// 	Body     string `json:"body" db:"body"`
// 	Category int64  `json:"category" db:"category"`
// 	Status   string `json:"status" db:"status"`
// }

type Category struct {
	Id            int64  `json:"id" db:"id"`
	Parent        int64  `json:"parent" db:"parent"`
	ParentName    string `json:"parentName" db:"parent_name"`
	Category      string `json:"category" db:"category"`
	HasDifficulty int64  `json:"has_difficulty" db:"has_difficulty"`
	DiffBeginner  int64  `json:"difficulty_beginner" db:"diff_beginner"`
	DiffAdvanced  int64  `json:"difficulty_advanced" db:"diff_advanced"`
	DiffExpert    int64  `json:"difficulty_expert" db:"diff_expert"`
	TextBeginner  string `json:"text_beginner" db:"text_beginner"`
	TextAdvanced  string `json:"text_advanced" db:"text_advanced"`
	TextExpert    string `json:"text_expert" db:"text_expert"`
	Status        string `json:"-" db:"status"`
	CreatedAt     int64  `json:"-" db:"created_at"`
	ApprovedAt    int64  `json:"-" db:"approved_at"`
	ApprovedBy    int64  `json:"-" db:"approved_by"`
	Author        int64  `json:"-" db:"author"`
	SortOrder     int64  `json:"-" db:"sort_order"`
}

// type CategoryInsert struct {
// 	Id         int64  `json:"id" db:"id"`
// 	Parent     int64  `json:"parent" db:"parent" binding:"required"`
// 	Category   string `json:"category" db:"category" binding:"required"`
// 	Status     string `json:"-" db:"status"`
// 	CreatedAt  int64  `json:"-" db:"created_at"`
// 	ApprovedAt int64  `json:"-" db:"approved_at"`
// 	ApprovedBy int64  `json:"-" db:"approved_by"`
// 	Author     int64  `json:"-" db:"author"`
// 	SortOrder  int64  `json:"-" db:"sort_order"`
// }

// type JSONCategory struct {
// 	Parent   int64  `json:"parent" db:"parent"`
// 	Category string `json:"category" db:"category"`
// 	Status   string `json:"status" db:"status"`
// }
