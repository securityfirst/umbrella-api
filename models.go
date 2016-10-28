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

type FeedItem struct {
	Id          int64  `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Country     string `json:"-"`
	Source      int64  `json:"-"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
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

type Language struct {
	Name  string `json:"name" db:"name"`
	Label string `json:"label" db:"label"`
}

type RWResponse struct {
	Version string `json:"version"`
	Status  int    `json:"status"`
	Time    int    `json:"time"`
	Data    struct {
		Type string `json:"type"`
		ID   int    `json:"id"`
		Item struct {
			ID              int       `json:"id"`
			Name            string    `json:"name"`
			Description     string    `json:"description"`
			Status          string    `json:"status"`
			Iso3            string    `json:"iso3"`
			Featured        bool      `json:"featured"`
			URL             string    `json:"url"`
			DescriptionHTML string    `json:"description-html"`
			Current         bool      `json:"current"`
			Location        []float64 `json:"location"`
		} `json:"item"`
	} `json:"data"`
}

type RWReport struct {
	Version string `json:"version"`
	Status  int    `json:"status"`
	Time    int    `json:"time"`
	Data    struct {
		Type string `json:"type"`
		ID   int    `json:"id"`
		Item struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Body     string `json:"body"`
			Headline struct {
				Title    string `json:"title"`
				Summary  string `json:"summary"`
				Featured bool   `json:"featured"`
			} `json:"headline"`
			File []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				URL         string `json:"url"`
				Filename    string `json:"filename"`
				Filemime    string `json:"filemime"`
				Filesize    string `json:"filesize"`
				Preview     struct {
					URL      string `json:"url"`
					URLLarge string `json:"url-large"`
					URLSmall string `json:"url-small"`
					URLThumb string `json:"url-thumb"`
				} `json:"preview"`
			} `json:"file"`
			PrimaryCountry struct {
				ID       int       `json:"id"`
				Name     string    `json:"name"`
				Iso3     string    `json:"iso3"`
				Location []float64 `json:"location"`
			} `json:"primary_country"`
			Country []struct {
				ID       int       `json:"id"`
				Name     string    `json:"name"`
				Iso3     string    `json:"iso3"`
				Location []float64 `json:"location"`
				Primary  bool      `json:"primary"`
			} `json:"country"`
			Source []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Shortname string `json:"shortname"`
				Longname  string `json:"longname"`
				Homepage  string `json:"homepage"`
				Type      struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"type"`
			} `json:"source"`
			Language []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"language"`
			Theme []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"theme"`
			Format []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"format"`
			OchaProduct []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"ocha_product"`
			Disaster []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Glide string `json:"glide"`
				Type  []struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
					Code string `json:"code"`
				} `json:"type"`
			} `json:"disaster"`
			DisasterType []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"disaster_type"`
			VulnerableGroups []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"vulnerable_groups"`
			URL      string `json:"url"`
			BodyHTML string `json:"body-html"`
			Date     struct {
				Created  int64 `json:"created"`
				Changed  int64 `json:"changed"`
				Original int64 `json:"original"`
			} `json:"date"`
		} `json:"item"`
	} `json:"data"`
}
