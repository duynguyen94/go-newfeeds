package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	mathrand "math/rand"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	totalUsers             = 10_000_000
	avgFollows             = 200
	topUsers               = 10
	topFollows             = 100_000
	batchSize              = 1_000
	postsPerUser           = 5
	postingUsersPercentage = 0.1
	maxConcurrentInserts   = 5
	maxOpenConns           = 10
	maxIdleConns           = 5
)

var (
	db           *sql.DB
	randomGen    *mathrand.Rand
	insertWg     sync.WaitGroup
	postInsertWg sync.WaitGroup
	insertSem    chan struct{}
)

type User struct {
	HashedPassword string
	Salt           string
	FirstName      string
	LastName       string
	Dob            string
	Email          string
	UserName       string
}

func init() {
	randomGen = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	insertSem = make(chan struct{}, maxConcurrentInserts)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, n)
	for i := range result {
		num := randomGen.Intn(len(letters))
		result[i] = letters[num]
	}
	return string(result)
}

func hashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func randomDate() string {
	min := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC).Unix()

	sec := randomGen.Int63n(max-min) + min
	return time.Unix(sec, 0).Format("2006-01-02")
}

func generateUser(id int) User {
	firstName := randomString(6)
	lastName := randomString(8)
	dob := randomDate()
	email := fmt.Sprintf("%s.%s%d@example.com", firstName, lastName, id)
	userName := fmt.Sprintf("%s%d", randomString(10), id)
	salt := randomString(16)
	password := randomString(12)
	hashedPassword := hashPassword(password, salt)

	return User{
		HashedPassword: hashedPassword,
		Salt:           salt,
		FirstName:      firstName,
		LastName:       lastName,
		Dob:            dob,
		Email:          email,
		UserName:       userName,
	}
}

func generatePost(userID int) (string, string) {
	contentText := randomString(100)
	contentImagePath := fmt.Sprintf("/images/%d/%s.jpg", userID, randomString(10))
	return contentText, contentImagePath
}

func insertUsers(users []User) {
	defer insertWg.Done()
	insertSem <- struct{}{} // Acquire semaphore

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO `user` (hashed_password, salt, first_name, last_name, dob, email, user_name) VALUES ")

	vals := make([]interface{}, 0, batchSize*7)

	for _, user := range users {
		queryBuilder.WriteString("(?, ?, ?, ?, ?, ?, ?),")
		vals = append(vals, user.HashedPassword, user.Salt, user.FirstName, user.LastName, user.Dob, user.Email, user.UserName)
	}

	query := strings.TrimSuffix(queryBuilder.String(), ",")
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Fatal(err)
	}

	<-insertSem // Release semaphore
	fmt.Printf("Inserted %d users\n", len(users))
}

func insertPosts(posts [][]interface{}) {
	defer postInsertWg.Done()
	insertSem <- struct{}{} // Acquire semaphore

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO `post` (fk_user_id, content_text, content_image_path) VALUES ")

	vals := make([]interface{}, 0, len(posts)*3)

	for _, post := range posts {
		queryBuilder.WriteString("(?, ?, ?),")
		vals = append(vals, post...)
	}

	query := strings.TrimSuffix(queryBuilder.String(), ",")
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Fatal(err)
	}

	<-insertSem // Release semaphore
	fmt.Printf("Inserted %d posts\n", len(posts))
}

func generateFollowers() {
	mathrand.Seed(time.Now().UnixNano())

	// Insert followers in batches
	batch := []string{}
	vals := []interface{}{}

	// Assign followers to the 10 most-followed users
	for i := 1; i <= topUsers; i++ {
		userID := i
		followers := mathrand.Perm(totalUsers)[:topFollows]
		for _, followerID := range followers {
			if followerID+1 != userID { // Ensure a user does not follow themselves
				batch = append(batch, "(?, ?)")
				vals = append(vals, userID, followerID+1)
				if len(batch) >= batchSize {
					insertBatchFollowers(batch, vals)
					batch = []string{}
					vals = []interface{}{}
				}
			}
		}
		fmt.Printf("User %d assigned %d followers\n", userID, topFollows)
	}

	// Insert any remaining followers
	if len(batch) > 0 {
		insertBatchFollowers(batch, vals)
	}

	// Assign an average of 200 followers to the rest of the users
	for userID := topUsers + 1; userID <= totalUsers; userID++ {
		numFollows := randomGen.Intn(avgFollows * 2) // To vary the number of follows per user around the average
		followers := mathrand.Perm(totalUsers)[:numFollows]
		for _, followerID := range followers {
			if followerID+1 != userID { // Ensure a user does not follow themselves
				batch = append(batch, "(?, ?)")
				vals = append(vals, userID, followerID+1)
				if len(batch) >= batchSize {
					insertBatchFollowers(batch, vals)
					batch = []string{}
					vals = []interface{}{}
				}
			}
		}
		if userID%1000 == 0 {
			fmt.Printf("Processed %d users for followers\n", userID)
		}
	}

	// Insert any remaining followers
	if len(batch) > 0 {
		insertBatchFollowers(batch, vals)
	}
}

func insertBatchFollowers(batch []string, vals []interface{}) {
	query := "INSERT IGNORE INTO `user_user` (fk_user_id, fk_follower_id) VALUES " + strings.Join(batch, ",")
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %d follower pairs\n", len(batch))
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve MySQL connection details from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set maximum open and idle connections to avoid "too many connections" error
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	// Step 1: Generate and insert users
	var usersToPost []int // Track users that will have posts
	for i := 0; i < totalUsers/batchSize; i++ {
		users := make([]User, 0, batchSize)
		for j := 0; j < batchSize; j++ {
			id := i*batchSize + j
			users = append(users, generateUser(id))
			// Randomly select 10% of users to have posts
			if randomGen.Float64() < postingUsersPercentage {
				usersToPost = append(usersToPost, id)
			}
		}
		insertWg.Add(1)
		go insertUsers(users)
	}
	insertWg.Wait()

	// Step 2: Generate and insert posts for selected users
	for _, userID := range usersToPost {
		posts := make([][]interface{}, 0, postsPerUser)
		for k := 0; k < postsPerUser; k++ {
			contentText, contentImagePath := generatePost(userID)
			posts = append(posts, []interface{}{userID, contentText, contentImagePath})
		}
		postInsertWg.Add(1)
		go insertPosts(posts)
	}
	postInsertWg.Wait()

	// Step 3: Generate followers after users have been seeded
	generateFollowers()
}
