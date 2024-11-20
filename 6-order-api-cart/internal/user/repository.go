package user

import "order-api/pkg/db"

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.database.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	existingUser, err := repo.FindByPhone(user.Phone)

	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if user.Name != "" {
		updates["name"] = user.Name
	}

	if user.SessionId != "" {
		updates["phone"] = user.Phone
	}
	if user.SessionId != "" {
		updates["session_id"] = user.SessionId
	}

	if len(updates) == 0 {
		return existingUser, nil
	}

	result := repo.database.DB.Model(&User{}).
		Where("phone = ?", user.Phone).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	return repo.FindByPhone(user.Phone)
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {
	var user User
	result := repo.database.DB.First(&user, "Phone = ?", phone)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) FindBySessionId(sessionId string) (*User, error) {
	var user User
	result := repo.database.DB.First(&user, "session_id = ?", sessionId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
