package data

import (
	"errors"
	"log"
	"math"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mcwiet/go-test/pkg/model"
)

type UserPoolClient interface {
	AdminGetUser(*cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error)
	ListUsers(*cognito.ListUsersInput) (*cognito.ListUsersOutput, error)
	DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
}

type UserDao struct {
	client     UserPoolClient
	userPoolId string
}

// Create a new DAO object for accessing users
func NewUserDao(client UserPoolClient, userPoolId string) UserDao {
	return UserDao{
		client:     client,
		userPoolId: userPoolId,
	}
}

// Get a user given a username
func (u *UserDao) GetByUsername(username string) (model.User, error) {
	ret, err := u.client.AdminGetUser(&cognito.AdminGetUserInput{
		UserPoolId: &u.userPoolId,
		Username:   &username,
	})

	if err != nil {
		log.Println(err)
		return model.User{}, errors.New("error retrieving user")
	}

	return convertAttributesToUser(username, ret.UserAttributes), nil
}

// List users
func (u *UserDao) Query(first int, after string) ([]model.User, string, error) {
	remaining := first
	users := []model.User{}
	paginationToken := after

	for remaining > 0 {
		limit := int64(math.Min(50, float64(remaining)))

		tempToken := &paginationToken
		if *tempToken == "" {
			tempToken = nil
		}
		ret, err := u.client.ListUsers(&cognito.ListUsersInput{
			UserPoolId:      &u.userPoolId,
			Limit:           &limit,
			PaginationToken: tempToken,
		})

		if err != nil {
			log.Println(err)
			return []model.User{}, "", errors.New("error retrieving users")
		}

		for _, userData := range ret.Users {
			user := convertAttributesToUser(*userData.Username, userData.Attributes)
			users = append(users, user)
		}

		if ret.PaginationToken == nil {
			paginationToken = ""
			remaining = 0
		} else {
			paginationToken = *ret.PaginationToken
			remaining -= 50
		}
	}

	return users, paginationToken, nil
}

// Get the (estimated) total count of users in the user pool
func (u *UserDao) GetTotalCount() (int, error) {
	ret, err := u.client.DescribeUserPool(&cognito.DescribeUserPoolInput{
		UserPoolId: &u.userPoolId,
	})

	if err != nil {
		log.Println(err)
		return 0, errors.New("error getting total number of users")
	}

	return int(*ret.UserPool.EstimatedNumberOfUsers), nil
}

// Convert a set of attributes into a user object
func convertAttributesToUser(username string, attrs []*cognito.AttributeType) model.User {
	attrMap := map[string]string{}
	for _, attr := range attrs {
		attrMap[*attr.Name] = *attr.Value
	}

	return model.User{
		Username: username,
		Email:    attrMap["email"],
		Name:     attrMap["name"],
	}
}
