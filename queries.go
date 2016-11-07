package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func (um *Umbrella) checkUser(c *gin.Context) (user User, err error) {
	token := c.Request.Header.Get("token")
	if utf8.RuneCountInString(token) > 0 {
		err = um.Db.SelectOne(&user, "select id, name, email, password, token, role from users where token=?", token)
	}
	return user, err
}

func (um *Umbrella) checkWebUser(token string) (user User, err error) {
	if utf8.RuneCountInString(token) > 0 {
		err = um.Db.SelectOne(&user, "select id, name, email, password, token, role from users where token=?", token)
	}
	return user, err
}

func (um *Umbrella) getAllPublishedSegments(c *gin.Context) (segments []Segment, err error) {
	_, err = um.Db.Select(&segments, "select s1.id, s1.title, s1.subtitle, s1.body, s1.category, s1.difficulty from segments s1 where status=:status order by id asc", map[string]interface{}{
		"status": "published",
	})
	return segments, err
}

func (um *Umbrella) getCountry(urlCountry string) string {
	country, err := um.Db.SelectStr("select iso2 from countries_index where iso2 = :iso2 order by id asc limit 1", map[string]interface{}{
		"iso2": strings.ToLower(strings.TrimSpace(urlCountry)),
	})
	checkErr(err)
	return country
}

func (um *Umbrella) getLastChecked(urlCountry string) (lastChecked []int64) {
	var checked struct {
		Relief int64
		FCO    int64
		UN     int64
		CDC    int64
		GDASC  int64
		CADATA int64
	}
	err := um.Db.SelectOne(&checked, "select COALESCE((SELECT last_checked FROM feed_last_checked WHERE country = :iso2 AND source = 0 limit 1),0) as relief, COALESCE((SELECT last_checked FROM feed_last_checked WHERE country = :iso2 AND source = 1 limit 1),0) as fco, COALESCE((SELECT last_checked FROM feed_last_checked WHERE country = :iso2 AND source = 2 limit 1),0) as un, COALESCE((SELECT last_checked FROM feed_last_checked WHERE source = 3 limit 1),0) as cdc from feed_last_checked limit 1", map[string]interface{}{
		"iso2": strings.ToLower(strings.TrimSpace(urlCountry)),
	})
	checkErr(err)
	if err == nil {
		lastChecked = []int64{checked.Relief, checked.FCO, checked.UN, checked.CDC, checked.GDASC, checked.CADATA}
	}
	if len(lastChecked) != SourceCount {
		lastChecked = make([]int64, SourceCount)
	}
	return lastChecked
}

func (um *Umbrella) getCountryInfo(urlCountry string) (country Country, err error) {
	err = um.Db.SelectOne(&country, "select id, name, iso3, iso2, reliefweb_int, search from countries_index where iso2 = :iso2 order by id asc limit 1", map[string]interface{}{
		"iso2": strings.ToLower(strings.TrimSpace(urlCountry)),
	})
	checkErr(err)
	return country, err
}

func (um *Umbrella) getDbFeedItems(sources []string, country string, since int64) (feedItems []FeedItem, err error) {
	if len(sources) < 1 {
		err = errors.New("No valid sources selected")
	} else if country == "" || len(country) != 2 {
		err = errors.New("Selected country is not valid")
	} else {
		_, err = um.Db.Select(&feedItems, fmt.Sprintf("select * from feed_items where country=:country and updated_at>:since and source in (%v) order by updated_at desc", strings.Join(sources, ",")), map[string]interface{}{
			"country": country,
			"since":   since,
		})
	}
	return feedItems, err
}

func (um *Umbrella) updateLastChecked(country string, source int, updatedAt int64) {
	_, err := um.Db.Exec("INSERT INTO feed_last_checked (country, source, last_checked) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE last_checked = ?", country, source, updatedAt, updatedAt)
	checkErr(err)
}

func (f *FeedItem) updateRelief(um *Umbrella) {
	var alreadyExists FeedItem
	trans, err := um.Db.Begin()
	checkErr(err)
	err = trans.SelectOne(&alreadyExists, "select * from feed_items where country= ? and source = ? and url = ? order by updated_at desc", f.Country, f.Source, f.URL)
	checkErr(err)
	if err == nil {
		if alreadyExists.UpdatedAt < f.UpdatedAt { // only if it has been updated since
			alreadyExists.UpdatedAt = f.UpdatedAt
			alreadyExists.Title = f.Title
			alreadyExists.Description = f.Description
			_, err := trans.Update(&alreadyExists)
			checkErr(err)
		}
	} else {
		checkErr(trans.Insert(f))
	}
	trans.Commit()
}

func (f *FeedItem) updateOthers(um *Umbrella) {
	var alreadyExists FeedItem
	trans, err := um.Db.Begin()
	checkErr(err)
	err = trans.SelectOne(&alreadyExists, "select * from feed_items where country= ? and source = ? and url = ? order by updated_at desc", f.Country, f.Source, f.URL)
	checkErr(err)
	if err == sql.ErrNoRows {
		checkErr(trans.Insert(f))
	}
	trans.Commit()
}

// func (um *Umbrella) getAllPublishedSegmentsByCat(c *gin.Context, category int64) (segments []Segment, err error) {
// 	_, err = um.Db.Select(&segments, "SELECT s.id, s.title, s.subtitle, s.body, s.category, s.difficulty FROM segments s INNER JOIN (SELECT `category`, MAX(`created_at`) as max_date FROM segments GROUP BY category) b ON s.category = b.category AND s.created_at = b.max_date WHERE s.category=:category", map[string]interface{}{
// 		"status":   "published",
// 		"category": category,
// 	})
// 	return segments, err
// }

// func (um *Umbrella) getSegmentById(c *gin.Context, segmentId int64) (segment Segment, err error) {
// 	err = um.Db.SelectOne(&segment, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE id=:segment_id ORDER BY id ASC", map[string]interface{}{
// 		"status":     "published",
// 		"segment_id": segmentId,
// 	})
// 	if err != nil && err.Error() == "sql: no rows in result set" {
// 		c.Fail(404, errors.New("The resource could not be found"))
// 	}
// 	return segment, err
// }

// func (um *Umbrella) getSegmentByCatId(c *gin.Context, categoryId int64) (segment Segment, err error) {
// 	err = um.Db.SelectOne(&segment, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE category=:category_id ORDER BY id DESC LIMIT 1", map[string]interface{}{
// 		"status":      "published",
// 		"category_id": categoryId,
// 	})
// 	if err != nil && err.Error() == "sql: no rows in result set" {
// 		c.Fail(404, errors.New("The resource could not be found"))
// 	}
// 	return segment, err
// }

func (um *Umbrella) getAllPublishedCheckItems(c *gin.Context) (checkItems []CheckItem, err error) {
	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category, difficulty, custom, disabled, no_check from check_items WHERE status=:status ORDER BY sort_order ASC, id ASC", map[string]interface{}{
		"status": "published",
	})
	return checkItems, err
}

// func (um *Umbrella) getAllPublishedCheckItemsByCat(c *gin.Context, categoryId int64) (checkItems []CheckItem, err error) {
// 	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category from check_items WHERE status=:status AND category=:category_id ORDER BY sort_order ASC, id ASC", map[string]interface{}{
// 		"status":      "published",
// 		"category_id": categoryId,
// 	})
// 	return checkItems, err
// }

// func (um *Umbrella) getCheckItemById(c *gin.Context, checkItemId int64) (checkItem CheckItem, err error) {
// 	err = um.Db.SelectOne(&checkItem, "select id, title, text, value, parent, category, status, created_at, author, approved_at, approved_by from check_items WHERE id=:check_item_id ORDER BY id ASC", map[string]interface{}{
// 		"status":        "published",
// 		"check_item_id": checkItemId,
// 	})
// 	return checkItem, err
// }

func (um *Umbrella) getAllPublishedCategories(c *gin.Context) (categories []Category, err error) {
	_, err = um.Db.Select(&categories, "select c.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, c.category, c.parent, c.has_difficulty, c.diff_beginner, c.diff_advanced, c.diff_expert, COALESCE(c.text_beginner, '') as text_beginner, COALESCE(c.text_advanced, '') as text_advanced, COALESCE(c.text_expert, '') as text_expert, c.`status`,c. created_at, c.author, c.approved_at, c.approved_by FROM categories as c LEFT JOIN categories as cat ON cat.id = c.parent WHERE c.status=:status ORDER BY id ASC", map[string]interface{}{
		"status": "published",
	})
	return categories, err
}

// func (um *Umbrella) getAllPublishedCategoriesByParent(c *gin.Context, parent int64) (categories []Category, err error) {
// 	_, err = um.Db.Select(&categories, "select categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, categories.category, categories.parent, categories.`status`,categories. created_at, categories.author, categories.approved_at, categories.approved_by from categories LEFT JOIN categories as cat ON cat.id = categories.parent WHERE categories.status=:status AND categories.id=:parent_id ORDER BY id ASC", map[string]interface{}{
// 		"status":    "published",
// 		"parent_id": parent,
// 	})
// 	return categories, err
// }

// func (um *Umbrella) getCategoryById(c *gin.Context, categoryId int64) (category Category, err error) {
// 	err = um.Db.SelectOne(&category, "select categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, categories.category, categories.parent, categories.`status`,categories. created_at, categories.author, categories.approved_at, categories.approved_by from categories LEFT JOIN categories as cat ON cat.id = categories.parent WHERE categories.status=:status AND categories.id=:category_id ORDER BY id ASC", map[string]interface{}{
// 		"status":      "published",
// 		"category_id": categoryId,
// 	})
// 	return category, err
// }
