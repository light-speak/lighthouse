package repo

import (
	"encoding/json"
	"fmt"
	"user/models"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/excute"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

type SelectRelation struct {
	Relation      *ast.Relation
	selectColumns map[string]interface{}
}

func GetUserProvide() map[string]*ast.Relation {
	return map[string]*ast.Relation{"created_at": {}, "id": {}, "name": {}, "posts": {}, "updated_at": {}}
}

func GetPostProvide() map[string]*ast.Relation {
	return map[string]*ast.Relation{"content": {}, "created_at": {}, "deleted_at": {}, "id": {}, "title": {}, "updated_at": {}, "user": {Name: "user", RelationType: ast.RelationTypeBelongsTo, ForeignKey: "user_id", Reference: "id"}, "user_id": {}}
}

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

func First__User(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	selectColumns := make([]string, 0)
	selectRelations := make(map[string]*ast.Relation, 0)
	for key, value := range columns {
		if value != nil && len(value.(map[string]interface{})) > 0 {
			relation := GetUserProvide()[key]
			selectRelations[key] = relation
			selectColumns = append(selectColumns, relation.ForeignKey)
		} else {
			selectColumns = append(selectColumns, key)
		}
	}
	user := &models.User{}
	err := BuildQuery_User().Scopes(scopes...).Select(selectColumns).First(user).Error
	if err != nil {
		return nil, err
	}
	return StructToMap(user)
}

func First__Post(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	selectColumns := make([]string, 0)
	selectRelations := make(map[string]*SelectRelation, 0)
	for key, value := range columns {
		if value != nil && len(value.(map[string]interface{})) > 0 {
			relation := GetPostProvide()[key]
			selectRelations[key] = &SelectRelation{relation, value.(map[string]interface{})}
			selectColumns = append(selectColumns, relation.ForeignKey)
		} else {
			selectColumns = append(selectColumns, key)
		}
	}
	post := &models.Post{}
	err := BuildQuery_Post().Scopes(scopes...).Select(selectColumns).First(post).Error
	if err != nil {
		return nil, err
	}
	res, err := StructToMap(post)
	if err != nil {
		return nil, err
	}
	for _, r := range selectRelations {
		fieldValue, err := Field_Post(post, r.Relation.ForeignKey)
		if err != nil {
			return nil, err
		}
		data, err := excute.GetQuickFirst(utils.UcFirst(r.Relation.Name))(r.selectColumns, func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", r.Relation.Reference), fieldValue)
		})
		if err != nil {
			return nil, err
		}
		res[r.Relation.Name] = data
	}

	return res, nil
}

func Field_Post(post *models.Post, key string) (interface{}, error) {
	switch key {
	case "id":
		return post.ID, nil
	case "title":
		return post.Title, nil
	case "content":
		return post.Content, nil
	case "user_id":
		return post.User_id, nil
	case "created_at":
		return post.Created_at, nil
	case "updated_at":
		return post.Updated_at, nil
	case "deleted_at":
		return post.Deleted_at, nil
	default:
		return nil, fmt.Errorf("field %s not found", key)
	}
}

func GetField_User(user *models.User, key string) (interface{}, error) {
	switch key {
	case "id":
		return user.ID, nil
	case "name":
		return user.Name, nil
	case "created_at":
		return user.Created_at, nil
	case "updated_at":
		return user.Updated_at, nil
	default:
		return nil, fmt.Errorf("field %s not found", key)
	}
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
	excute.AddQuickList("Posts", List_Post)
	excute.AddQuickFirst("Post", First__Post)
	excute.AddQuickFirst("User", First__User)
}

func StructToMap(m model.ModelInterface) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
