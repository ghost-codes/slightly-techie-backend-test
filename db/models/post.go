package db

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	ID     int64 `gorm:"primaryKey"`
	Text   string
	UserID int64
	User   User
}

func (store *Store) CreatePost(post *Post) error {
	if err := store.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (store *Store) UpdatePost(post *Post) error {
	if err := store.Save(post).Error; err != nil {
		return err
	}
	return nil
}

func (store *Store) DeletePost(id int64, uid int64) error {
	if err := store.Delete(&Post{UserID: uid, ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (store *Store) ViewOnePost(id int64, uid int64) (*Post, error) {
	post := &Post{
		ID:     id,
		UserID: uid,
	}
	if err := store.First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (store *Store) ViewAll(uid int64) ([]Post, error) {
	post := []Post{}
	if err := store.Where("user_id = ?", uid).Find(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}
