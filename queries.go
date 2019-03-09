package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/securityfirst/umbrella-api/feed"
	"github.com/securityfirst/umbrella-api/models"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func (um *Umbrella) checkUser(c *gin.Context) (user models.User, err error) {
	token := c.Request.Header.Get("token")
	if utf8.RuneCountInString(token) > 0 {
		err = um.Db.SelectOne(&user, "select id, name, email, password, token, role from users where token=?", token)
	}
	return user, err
}

func (um *Umbrella) checkWebUser(token string) (user models.User, err error) {
	if utf8.RuneCountInString(token) > 0 {
		err = um.Db.SelectOne(&user, "select id, name, email, password, token, role from users where token=?", token)
	}
	return user, err
}

func (um *Umbrella) getAllPublishedSegments(c *gin.Context) (segments []models.Segment, err error) {
	_, err = um.Db.Select(&segments, "select s1.id, s1.title, s1.subtitle, s1.body, s1.category, s1.difficulty from segments s1 where status=:status order by id asc", map[string]interface{}{
		"status": "published",
	})
	return segments, err
}

func (um *Umbrella) getCountry(urlCountry string) string {
	country, err := um.Db.SelectStr("select iso2 from countries_index where iso2 = :iso2 order by id asc limit 1", map[string]interface{}{
		"iso2": strings.ToLower(strings.TrimSpace(urlCountry)),
	})
	utils.CheckErr(err)
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
	utils.CheckErr(err)
	if err == nil {
		lastChecked = []int64{checked.Relief, checked.FCO, checked.UN, checked.CDC, checked.GDASC, checked.CADATA}
	}
	if len(lastChecked) != feed.SourceCount {
		lastChecked = make([]int64, feed.SourceCount)
	}
	return lastChecked
}

func (um *Umbrella) getCountryInfo(urlCountry string) (country models.Country, err error) {
	err = um.Db.SelectOne(&country, "select id, name, iso3, iso2, reliefweb_int, search from countries_index where iso2 = :iso2 order by id asc limit 1", map[string]interface{}{
		"iso2": strings.ToLower(strings.TrimSpace(urlCountry)),
	})
	utils.CheckErr(err)
	return country, err
}

func (um *Umbrella) getDbFeedItems(sources []string, country string, since int64) (feedItems []models.FeedItem, err error) {
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
	utils.CheckErr(err)
}

func (um *Umbrella) getSegmentById(c *gin.Context, segmentId int64) (segment models.Segment, err error) {
	err = um.Db.SelectOne(&segment, "select id, title, subtitle, difficulty, body, category, status, created_at, author, approved_at, approved_by from segments WHERE id=:segment_id ORDER BY id ASC", map[string]interface{}{
		"status":     "published",
		"segment_id": segmentId,
	})
	return segment, err
}

func (um *Umbrella) getSegmentByCatIdAndDifficulty(categoryId int64, difficulty string) (segments []models.Segment, err error) {
	var diffInt int = 1
	switch difficulty {
	case "advanced":
		diffInt = 2
	case "expert":
		diffInt = 3
	}
	_, err = um.Db.Select(&segments, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE status = :status and category=:category_id and difficulty=:difficulty ORDER BY id ASC", map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
		"difficulty":  diffInt,
	})
	return segments, err
}

func (um *Umbrella) getCheckItemsByCatIdAndDifficulty(categoryId int64, difficulty string) (checkItems []models.CheckItem, err error) {
	var diffInt int = 1
	switch difficulty {
	case "advanced":
		diffInt = 2
	case "expert":
		diffInt = 3
	}
	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category, EXISTS(SELECT * FROM check_items ci WHERE check_items.parent = 0 AND ci.parent = check_items.id LIMIT 1) as has_subitems from check_items WHERE status=:status AND category=:category_id ORDER BY sort_order ASC, id ASC", map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
		"difficulty":  diffInt,
	})
	return checkItems, err
}

func (um *Umbrella) getAllPublishedCheckItems(c *gin.Context) (checkItems []models.CheckItem, err error) {
	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category, difficulty, custom, disabled, no_check from check_items WHERE status=:status ORDER BY sort_order ASC, id ASC", map[string]interface{}{
		"status": "published",
	})
	return checkItems, err
}

func (um *Umbrella) getAllPublishedCategories(c *gin.Context) (categories []models.Category, err error) {
	_, err = um.Db.Select(&categories, "select c.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, EXISTS(SELECT * FROM categories c2 WHERE c2.parent = c.id LIMIT 1) as has_subcategories, c.category, c.parent, c.has_difficulty, c.diff_beginner, c.diff_advanced, c.diff_expert, COALESCE(c.text_beginner, '') as text_beginner, COALESCE(c.text_advanced, '') as text_advanced, COALESCE(c.text_expert, '') as text_expert, c.`status`,c. created_at, c.author, c.approved_at, c.approved_by FROM categories as c LEFT JOIN categories as cat ON cat.id = c.parent WHERE c.status=:status ORDER BY id ASC, c.sort_order ASC", map[string]interface{}{
		"status": "published",
	})
	return categories, err
}

func (um *Umbrella) getCategoryById(c *gin.Context, categoryId int64) (category models.Category, err error) {
	err = um.Db.SelectOne(&category, `
		select
			categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name,
			categories.category, categories.parent, categories.status, categories.created_at, categories.author,
			categories.approved_at, categories.approved_by
		from
			categories LEFT JOIN categories as cat ON cat.id = categories.parent
		where
			categories.status=:status AND categories.id=:category_id ORDER BY id ASC
		`, map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
	})
	return category, err
}
