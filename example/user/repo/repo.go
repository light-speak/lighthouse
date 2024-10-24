package repo

import (
	"user/models"

	"github.com/light-speak/lighthouse/graphql/excute"
	"github.com/light-speak/lighthouse/graphql/model"
	"gorm.io/gorm"
)

func BuildQuery_Post(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	return model.GetDB().Model(&models.Post{}).Scopes(scopes...)
}

func BuildQuery_User(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	return model.GetDB().Model(&models.User{}).Scopes(scopes...)
}

func List_Post() ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	var datas []*models.Post
	err := BuildQuery_Post().Find(&datas).Error
	if err != nil {
		return nil, err
	}
	for _, data := range datas {
		item := make(map[string]interface{})
		post := data
		item["id"] = post.ID
		item["title"] = post.Title
		item["content"] = post.Content
		item["user"] = map[string]interface{}{
			"id":   post.User.ID,
			"name": post.User.Name,
		}
		res = append(res, item)
	}
	return res, nil
}

func First__User() (map[string]interface{}, error) {
	res := make(map[string]interface{})
	user := &models.User{}
	err := BuildQuery_User().First(user).Error
	if err != nil {
		return nil, err
	}
	res["id"] = user.ID
	res["name"] = user.Name
	return res, nil
}

func Load__Post(id int64) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	loader := model.GetLoader[*models.Post](model.GetDB(), nil)
	post, err := loader.Load(id)
	if err != nil {
		return nil, err
	}
	res["id"] = post.ID
	res["title"] = post.Title
	res["content"] = post.Content
	return res, nil
}

func Load__User(id int64) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	loader := model.GetLoader[*models.User](model.GetDB(), nil)
	user, err := loader.Load(id)
	if err != nil {
		return nil, err
	}
	res["id"] = user.ID
	res["name"] = user.Name
	return res, nil
}

func GetCount(buildCountQuery func(db *gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	err := buildCountQuery(model.GetDB()).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func init() {
	excute.AddQuickList("posts", List_Post)
	excute.AddQuickFirst("user", First__User)
}
