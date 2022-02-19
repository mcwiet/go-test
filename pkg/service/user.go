package service

import (
	"github.com/mcwiet/go-test/pkg/model"
)

type UserDao interface {
	GetByUsername(id string) (model.User, error)
	GetTotalCount() (int, error)
	List(first int, after string) ([]model.User, string, error)
}

type UserService struct {
	userDao UserDao
	encoder CursorEncoder
}

func NewUserService(userDao UserDao, encoder CursorEncoder) UserService {
	return UserService{
		userDao: userDao,
		encoder: encoder,
	}
}

// Get a user from their username
func (u *UserService) GetByUsername(username string) (model.User, error) {
	user, err := u.userDao.GetByUsername(username)
	return user, err
}

// Get the first N users after the provided token
func (u *UserService) List(first int, after string) (model.UserConnection, error) {
	decodedToken, err := u.encoder.Decode(after)
	if err != nil {
		return model.UserConnection{}, err
	}

	users, token, err := u.userDao.List(first, decodedToken)
	if err != nil {
		return model.UserConnection{}, err
	}

	token = u.encoder.Encode(token)
	totalCount, err := u.userDao.GetTotalCount()

	connection := model.UserConnection{
		TotalCount: totalCount,
		Edges:      []model.UserEdge{},
		PageInfo: model.PageInfo{
			EndCursor:   token,
			HasNextPage: token != "",
		},
	}

	for _, user := range users {
		connection.Edges = append(connection.Edges, model.UserEdge{
			Node: user,
		})
	}

	return connection, err
}
