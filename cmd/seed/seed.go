package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"

	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/db"
	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/d4rthvadr/dusky-go/internal/utils"

	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/joho/godotenv"
)

var fakeUsernames = []string{
	"alice",
	"bob",
	"charlie",
	"dave",
	"eve",
	"frank",
}

func main() {

	err := godotenv.Load()
	logger := utils.NewLogger()
	defer logger.Sync()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	config, err := config.InitializeConfig()
	if err != nil {
		logger.Fatal("Error initializing config:", err)
	}

	db, err := db.New(config.Db.Addr, config.Db.MaxOpenConns, config.Db.MaxIdleConns, config.Db.MaxIdleTime)
	if err != nil {
		logger.Panic("Error connecting to the database:", err)
	}

	store := store.NewStorage(db)
	if err = Seed(store); err != nil {
		logger.Fatal("Error seeding data:", err)
	}
}

func Seed(store store.Storage) error {
	logger := utils.NewLogger()
	defer logger.Sync()
	logger.Info("seeding...")
	ctx := context.Background()

	users := generateUsers(3)

	userIds := make([]int64, len(users))

	for index, user := range users {
		if err := store.Users.Create(ctx, &user); err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
		userIds[index] = user.ID
	}

	posts := generatePosts(200, userIds)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, &post); err != nil {
			return fmt.Errorf("error creating post: %w", err)
		}
	}

	fmt.Println("done")

	return nil
}

func getFakeUsername(i int) string {
	return fakeUsernames[i%len(fakeUsernames)] + fmt.Sprintf("%d", i+1)
}

// generateRandomHash generates a random 4-character hexadecimal string
func generateRandomHash() string {
	bytes := make([]byte, 10) // 10 bytes will give us a 20-character hex string
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func generateUsers(count int) []models.User {
	users := make([]models.User, count)

	for i := 0; i < count; i++ {
		username := getFakeUsername(i) + "_" + generateRandomHash()
		users[i] = models.User{
			Username: username,
			Email:    fmt.Sprintf("user_%s@example.com", username),
			Password: "password123",
		}
	}
	return users
}

func generatePosts(count int, userIDs []int64) []models.Post {
	posts := make([]models.Post, count)
	for i := 0; i < count; i++ {
		userID := userIDs[mrand.Intn(len(userIDs))]
		posts[i] = models.Post{
			Title:   fmt.Sprintf("Post Title %d", i+1),
			Content: fmt.Sprintf("This is the content of post %d.", i+1),
			UserID:  userID,
			Tags:    []string{"tag1", "tag2"},
		}
	}
	return posts
}
