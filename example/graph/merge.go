// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.

package graph

import "context"


func (r *Resolver) mergeUsers(ctx context.Context, users []*User) ([]*User, error) {
	var err error
	for _, user := range users {
		user, err = r.mergeUser(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (r *Resolver) mergeUser(ctx context.Context, user *User) (*User, error) {
	return user, nil
}



