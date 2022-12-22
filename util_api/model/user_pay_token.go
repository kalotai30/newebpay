package model

import (
	"time"
	. "util_api/database/mysql"
	"util_api/util/log"
)

type UserPayToken struct {
	Id             int64     `table:"id"`
	UserId         int64     `table:"user_id"`
	Token          string    `table:"token"`
	TokenLife      time.Time `table:"tokenLife"`
	TokenUpdatedAt time.Time `table:"token_updated_at"`
	TokenType      int64     `table:"token_type"`
}

func (model *UserPayToken) SetUserId(userId int64) *UserPayToken {
	model.UserId = userId
	return model
}

func (model *UserPayToken) SetToken(token string) *UserPayToken {
	model.Token = token
	return model
}

func (model *UserPayToken) SetTokenLife(tokenLife time.Time) *UserPayToken {
	model.TokenLife = tokenLife
	return model
}

func (model *UserPayToken) SetTokenUpdatedAt(tokenUpdatedAt time.Time) *UserPayToken {
	model.TokenUpdatedAt = tokenUpdatedAt
	return model
}

func (model *UserPayToken) SetTokenType(tokenType int64) *UserPayToken {
	model.TokenType = tokenType
	return model
}

func (model *UserPayToken) QueryOne() *UserPayToken {
	table := Model(model)

	if model.UserId > 0 {
		table.Where("user_id", "=", model.UserId)
	}

	table.Where("token_type", "=", model.TokenType)

	log.Error(table.Select([]string{"id", "user_id", "token", "tokenLife", "token_updated_at", "token_type"}).
		Find().Scan(&model.Id, &model.UserId, &model.Token, &model.TokenLife, &model.TokenUpdatedAt, &model.TokenType))
	return model
}

func (model *UserPayToken) Insert() (int64, error) {
	return Model(model).Insert()
}

func (model *UserPayToken) Update(columns []string) error {
	table := Model(model)

	if model.UserId != 0 {
		table.Where("user_id", "=", model.UserId)
	}

	return table.Update(columns)
}
