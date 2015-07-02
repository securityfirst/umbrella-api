package main

import (
	"errors"

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

func (um *Umbrella) getAllPublishedSegments(c *gin.Context) (segments []Segment, err error) {
	_, err = um.Db.Select(&segments, "select s1.id, s1.title, s1.subtitle, s1.body, s1.category, s1.difficulty from segments s1 where status=:status AND (category, id) in (select category, max(id) from segments s2 where status=:status group by category)", map[string]interface{}{
		"status": "published",
	})
	return segments, err
}

func (um *Umbrella) getAllPublishedSegmentsByCat(c *gin.Context, category int64) (segments []Segment, err error) {
	_, err = um.Db.Select(&segments, "SELECT s.id, s.title, s.subtitle, s.body, s.category FROM segments s INNER JOIN (SELECT `category`, MAX(`created_at`) as max_date FROM segments GROUP BY category) b ON s.category = b.category AND s.created_at = b.max_date WHERE s.category=:category", map[string]interface{}{
		"status":   "published",
		"category": category,
	})
	return segments, err
}

func (um *Umbrella) getSegmentById(c *gin.Context, segmentId int64) (segment Segment, err error) {
	err = um.Db.SelectOne(&segment, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE id=:segment_id ORDER BY id ASC", map[string]interface{}{
		"status":     "published",
		"segment_id": segmentId,
	})
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.Fail(404, errors.New("The resource could not be found"))
	}
	return segment, err
}

func (um *Umbrella) getSegmentByCatId(c *gin.Context, categoryId int64) (segment Segment, err error) {
	err = um.Db.SelectOne(&segment, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE category=:category_id ORDER BY id DESC LIMIT 1", map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
	})
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.Fail(404, errors.New("The resource could not be found"))
	}
	return segment, err
}

func (um *Umbrella) getAllPublishedCheckItems(c *gin.Context) (checkItems []CheckItem, err error) {
	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category from check_items WHERE status=:status ORDER BY sort_order ASC, id ASC", map[string]interface{}{
		"status": "published",
	})
	return checkItems, err
}

func (um *Umbrella) getAllPublishedCheckItemsByCat(c *gin.Context, categoryId int64) (checkItems []CheckItem, err error) {
	_, err = um.Db.Select(&checkItems, "select id, title, text, value, parent, category from check_items WHERE status=:status AND category=:category_id ORDER BY sort_order ASC, id ASC", map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
	})
	return checkItems, err
}

func (um *Umbrella) getCheckItemById(c *gin.Context, checkItemId int64) (checkItem CheckItem, err error) {
	err = um.Db.SelectOne(&checkItem, "select id, title, text, value, parent, category, status, created_at, author, approved_at, approved_by from check_items WHERE id=:check_item_id ORDER BY id ASC", map[string]interface{}{
		"status":        "published",
		"check_item_id": checkItemId,
	})
	return checkItem, err
}

func (um *Umbrella) getAllPublishedCategories(c *gin.Context) (categories []Category, err error) {
	_, err = um.Db.Select(&categories, "select categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, categories.category, categories.parent, categories.`status`,categories. created_at, categories.author, categories.approved_at, categories.approved_by from categories LEFT JOIN categories as cat ON cat.id = categories.parent WHERE categories.status=:status ORDER BY id ASC", map[string]interface{}{
		"status": "published",
	})
	return categories, err
}

func (um *Umbrella) getAllPublishedCategoriesByParent(c *gin.Context, parent int64) (categories []Category, err error) {
	_, err = um.Db.Select(&categories, "select categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, categories.category, categories.parent, categories.`status`,categories. created_at, categories.author, categories.approved_at, categories.approved_by from categories LEFT JOIN categories as cat ON cat.id = categories.parent WHERE categories.status=:status AND categories.id=:parent_id ORDER BY id ASC", map[string]interface{}{
		"status":    "published",
		"parent_id": parent,
	})
	return categories, err
}

func (um *Umbrella) getCategoryById(c *gin.Context, categoryId int64) (category Category, err error) {
	err = um.Db.SelectOne(&category, "select categories.id, (case when cat.category IS NOT NULL then cat.category else '' end) as parent_name, categories.category, categories.parent, categories.`status`,categories. created_at, categories.author, categories.approved_at, categories.approved_by from categories LEFT JOIN categories as cat ON cat.id = categories.parent WHERE categories.status=:status AND categories.id=:category_id ORDER BY id ASC", map[string]interface{}{
		"status":      "published",
		"category_id": categoryId,
	})
	return category, err
}
