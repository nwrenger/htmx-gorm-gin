package db

type Permission int64

const (
	None Permission = iota
	ReadOnly
	Write
)

type Permissions struct {
	AccessAdmin Permission `json:"access_admin"`
	AccessBook  Permission `json:"access_book"`
}

type Login struct {
	User       string      `gorm:"primaryKey" json:"user"`
	Password   string      `json:"password"`
	Permission Permissions `json:"permission" gorm:"embedded"`
}

func FetchLogin(user string) (Login, error) {
	var login Login
	rs := db.Where("user = ?", user).Find(&login)
	return login, rs.Error
}

func AddOrUpdateLogin(newLogin *Login) error {
	var existingLogin Login

	// Try to find the existing Login by ID
	if err := db.First(&existingLogin, newLogin.User).Error; err != nil {
		// If not found, create a new record
		rs := db.Create(newLogin)
		return rs.Error
	}

	// If found, update the existing record
	rs := db.Model(&existingLogin).Updates(newLogin)
	return rs.Error
}

func DeleteLogin(user string) (Login, error) {
	var login Login
	rs := db.Where("user = ?", user).Delete(&login)
	return login, rs.Error
}
