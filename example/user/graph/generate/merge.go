// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.

package generate

import (
	"context"
	"user/graph/models"
)

func MergeUsers(ctx context.Context, users []*models.User) ([]*models.User, error) {
	var err error
	for _, user := range users {
		user, err = MergeUser(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func  MergeUser(ctx context.Context, user *models.User) (*models.User, error) {
	return user, nil
}
