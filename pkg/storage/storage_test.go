package storage

import (
	"Task36a41/pkg/rss"
	"testing"
	"time"
)

func setupDatabase(db *Storage) error {
	// Drop table если создано, иначе создать
	_, err := db.db.Exec(`DROP TABLE IF EXISTS posts;
        CREATE TABLE posts (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            pub_time BIGINT NOT NULL,
            link TEXT NOT NULL UNIQUE
        );`)
	return err
}

func TestSaveAndGetPosts(t *testing.T) {
	// Используем строку подключения к тестовой базе данных
	db, err := New("user=postgres password=vlad5043 dbname=News sslmode=disable")
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Настраиваем базу данных
	if err := setupDatabase(db); err != nil {
		t.Fatalf("Error setting up database: %v", err)
	}

	// Очищаем таблицу перед началом тестов
	_, err = db.db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatalf("Error cleaning table: %v", err)
	}

	// Пример поста для тестирования
	post := rss.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		PubDate: time.Now().Format(time.RFC1123Z), // Используем текущее время
		Link:    "http://example.com/test",
	}

	// Сохраняем пост
	err = db.SavePost(post)
	if err != nil {
		t.Fatalf("Error saving post: %v", err)
	}

	// Получаем посты
	posts, err := db.GetLastNPosts(1)
	if err != nil {
		t.Fatalf("Error getting posts: %v", err)
	}

	// Проверяем результат
	if len(posts) != 1 || posts[0].Title != post.Title {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}
}
