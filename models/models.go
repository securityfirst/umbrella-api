package models

type Country struct {
	Id          int64    `db:"id"`
	Name        string   `db:"name"`
	Iso2        string   `db:"iso2"`
	Iso3        string   `db:"iso3"`
	ReliefWeb   int64    `db:"reliefweb_int"`
	Search      string   `db:"search"`
	SearchSlice []string `db:"-"`
}

type CheckItem struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Value       int64  `json:"value"`
	Parent      int64  `json:"parent"`
	HasSubItems bool   `json:"-" db:"has_subitems"`
	Category    int64  `json:"category"`
	Difficulty  int64  `json:"difficulty"`
	Custom      int64  `json:"custom"`
	Disabled    int64  `json:"disabled"`
	NoCheck     int64  `json:"no_check" db:"no_check"`
	Status      string `json:"-" db:"status"`
	CreatedAt   int64  `json:"-" db:"created_at"`
	ApprovedAt  int64  `json:"-" db:"approved_at"`
	ApprovedBy  int64  `json:"-" db:"approved_by"`
	Author      int64  `json:"-" db:"author"`
}

type FeedLastChecked struct {
	Id          int64
	Country     string
	Source      string
	LastChecked int64 `db:"last_checked"`
}

type FeedRequestLog struct {
	Id        int64
	Country   string
	Sources   string
	Status    int64
	CheckedAt int64 `db:"checked_at"`
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
	Id               int64  `json:"id" db:"id"`
	Title            string `json:"title" db:"title"`
	SubTitle         string `json:"subtitle" db:"subtitle"`
	Body             string `json:"body" db:"body"`
	Category         int64  `json:"category" db:"category"`
	Difficulty       int64  `json:"difficulty" db:"difficulty"`
	DifficultyString string `json:"difficulty_string,omitempty" db:"-"`
	Status           string `json:"-" db:"status"`
	CreatedAt        int64  `json:"-" db:"created_at"`
	ApprovedAt       int64  `json:"-" db:"approved_at"`
	ApprovedBy       int64  `json:"-" db:"approved_by"`
	Author           int64  `json:"-" db:"author"`
}

// type JSONSegment struct {
// 	Title    string `json:"title" db:"title"`
// 	SubTitle string `json:"subtitle" db:"subtitle"`
// 	Body     string `json:"body" db:"body"`
// 	Category int64  `json:"category" db:"category"`
// 	Status   string `json:"status" db:"status"`
// }

type Category struct {
	Id               int64  `json:"id" db:"id"`
	Parent           int64  `json:"parent" db:"parent"`
	ParentName       string `json:"parentName" db:"parent_name"`
	Category         string `json:"category" db:"category"`
	HasSubcategories bool   `json:"-" db:"has_subcategories"`
	HasDifficulty    int64  `json:"has_difficulty" db:"has_difficulty"`
	DiffBeginner     int64  `json:"difficulty_beginner" db:"diff_beginner"`
	DiffAdvanced     int64  `json:"difficulty_advanced" db:"diff_advanced"`
	DiffExpert       int64  `json:"difficulty_expert" db:"diff_expert"`
	TextBeginner     string `json:"text_beginner" db:"text_beginner"`
	TextAdvanced     string `json:"text_advanced" db:"text_advanced"`
	TextExpert       string `json:"text_expert" db:"text_expert"`
	Status           string `json:"-" db:"status"`
	CreatedAt        int64  `json:"-" db:"created_at"`
	ApprovedAt       int64  `json:"-" db:"approved_at"`
	ApprovedBy       int64  `json:"-" db:"approved_by"`
	Author           int64  `json:"-" db:"author"`
	SortOrder        int64  `json:"-" db:"sort_order"`
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

type Language struct {
	Name  string `json:"name" db:"name"`
	Label string `json:"label" db:"label"`
}
