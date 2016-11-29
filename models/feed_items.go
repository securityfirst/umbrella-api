package models

import (
	"database/sql"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/go-gorp/gorp"
)

type FeedItem struct {
	Id          int64  `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Country     string `json:"-"`
	Source      int64  `json:"-"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
}

func (f *FeedItem) UpdateRelief(db *gorp.DbMap) {
	var alreadyExists FeedItem
	trans, err := db.Begin()
	utils.CheckErr(err)
	err = trans.SelectOne(&alreadyExists, "select * from feed_items where country= ? and source = ? and url = ? order by updated_at desc", f.Country, f.Source, f.URL)
	utils.CheckErr(err)
	if err == nil {
		if alreadyExists.UpdatedAt < f.UpdatedAt { // only if it has been updated since
			alreadyExists.UpdatedAt = f.UpdatedAt
			alreadyExists.Title = f.Title
			alreadyExists.Description = f.Description
			_, err := trans.Update(&alreadyExists)
			utils.CheckErr(err)
		}
	} else {
		utils.CheckErr(trans.Insert(f))
	}
	trans.Commit()
}

func (f *FeedItem) UpdateOthers(db *gorp.DbMap) {
	var alreadyExists FeedItem
	trans, err := db.Begin()
	utils.CheckErr(err)
	err = trans.SelectOne(&alreadyExists, "select * from feed_items where country= ? and source = ? and url = ? order by updated_at desc", f.Country, f.Source, f.URL)
	utils.CheckErr(err)
	if err == sql.ErrNoRows {
		utils.CheckErr(trans.Insert(f))
	}
	trans.Commit()
}
