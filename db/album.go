package db

type Album struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func FetchAlbums() ([]Album, error) {
	var albums []Album
	rs := db.Find(&albums)
	return albums, rs.Error
}

func FetchAlbumById(id uint) (Album, error) {
	var album Album
	rs := db.Where("id = ?", id).Find(&album)
	return album, rs.Error
}

func AddOrUpdateAlbum(newAlbum *Album) error {
	var existingAlbum Album

	// Try to find the existing album by ID
	if err := db.First(&existingAlbum, newAlbum.ID).Error; err != nil {
		// If not found, create a new record
		rs := db.Create(newAlbum)
		return rs.Error
	}

	// If found, update the existing record
	rs := db.Model(&existingAlbum).Updates(newAlbum)
	return rs.Error
}

func DeleteAlbum(id uint) (Album, error) {
	var album Album
	rs := db.Where("id = ?", id).Delete(&album)
	return album, rs.Error
}
