package model

import (
	. "app_api/database/mysql"
	"app_api/util/log"
	"database/sql"
	"time"
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

	log.Error(table.Select([]string{"id", "user_id", "token", "tokenLife", "token_updated_at", "token_type"}).
		Where("user_id", "=", model.UserId).
		Where("token_type", "=", model.TokenType).
		Find().Scan(&model.Id, &model.UserId, &model.Token, &model.TokenLife, &model.TokenUpdatedAt, &model.TokenType))
	return model
}

func (model *UserPayToken) QueryOneByTokeLife() *UserPayToken {
	table := Model(model)

	log.Error(table.Select([]string{"id", "user_id", "token", "tokenLife", "token_updated_at", "token_type"}).
		Where("user_id", "=", model.UserId).
		Where("token_type", "=", model.TokenType).
		Where("tokenLife", ">=", time.Now().Format("2006-01-02")).
		Find().Scan(&model.Id, &model.UserId, &model.Token, &model.TokenLife, &model.TokenUpdatedAt, &model.TokenType))
	return model
}

func (model *UserPayToken) QueryAll(option func(*UserPayToken)) {
	table := Model(model)

	table.Select([]string{"id", "user_id", "token", "tokenLife", "token_updated_at", "token_type"}).
		Where("user_id", "=", model.UserId).
		OrderBy([]string{"token_type"}, []string{"asc"}).
		Get(func(rows *sql.Rows) (isBreak bool) {
			err := rows.Scan(&model.Id, &model.UserId, &model.Token, &model.TokenLife, &model.TokenUpdatedAt, &model.TokenType)
			log.Error(err)
			if err == nil && option != nil {
				option(model)
			}
			return
		})
}
