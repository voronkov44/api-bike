package users

import (
	"bike/pkg/db"
)

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindByID(id uint) (*User, error) {
	var user User
	result := repo.database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// ListAll - получает всех пользователей с пагинацией и сортировкой
func (repo *UserRepository) ListAll(limit, offset int, sortBy, nameFilter, emailFilter string) ([]User, error) {
	var list []User

	database := repo.database.DB

	// ФИЛЬТРАЦИЯ по имени
	if nameFilter != "" {
		database = database.Where("name LIKE ?", "%"+nameFilter+"%")
	}

	// ФИЛЬТРАЦИЯ по email
	if emailFilter != "" {
		database = database.Where("email LIKE ?", "%"+emailFilter+"%")
	}

	// СОРТИРОВКА
	switch sortBy {
	case "name":
		database = database.Order("name asc")
	case "name_desc":
		database = database.Order("name desc")
	case "email":
		database = database.Order("email asc")
	case "email_desc":
		database = database.Order("email desc")
	case "created_at":
		database = database.Order("created_at asc")
	case "created_at_desc":
		database = database.Order("created_at desc")
	default:
		database = database.Order("created_at desc") // по умолчанию
	}

	// Пагинация
	if limit > 0 {
		database = database.Limit(limit).Offset(offset)
	}

	result := database.Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// Count - возвращает общее количество пользователей
func (repo *UserRepository) Count(nameFilter, emailFilter string) (int64, error) {
	var count int64
	database := repo.database.DB.Model(&User{})

	if nameFilter != "" {
		database = database.Where("name LIKE ?", "%"+nameFilter+"%")
	}
	if emailFilter != "" {
		database = database.Where("email LIKE ?", "%"+emailFilter+"%")
	}

	result := database.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// SearchByEmail - поиск пользователей по email (частичное совпадение)
func (repo *UserRepository) SearchByEmail(email string) ([]User, error) {
	var users []User

	if email == "" {
		return users, nil
	}

	result := repo.database.DB.
		Where("email LIKE ?", "%"+email+"%").
		Order("email asc").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
