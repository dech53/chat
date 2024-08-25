package dao

import "socket/model"

func AddUser(user *model.User) bool {
	row := DB.Create(user).RowsAffected
	if row == 0 {
		return false
	}
	return true
}
func GetUserByPattern(pattern, value string) (user *model.User) {
	user = &model.User{}
	// 使用参数化查询来防止 SQL 注入
	query := pattern + " LIKE ?"
	DB.Where(query, value).First(&user)
	return user
}
