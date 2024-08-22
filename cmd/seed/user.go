package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	mathrand "math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	totalUsers = 10000
	avgFollows = 200
	topUsers   = 10
	topFollows = 1000
	batchSize  = 1000
)

// User represents a user in the database
type User struct {
	HashedPassword string
	Salt           string
	FirstName      string
	LastName       string
	Dob            string
	Email          string
	UserName       string
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
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

	sec := randInt(min, max)
	return time.Unix(sec, 0).Format("2006-01-02")
}

func randInt(min, max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return min + n.Int64()
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

func insertUsers(db *sql.DB, users []User) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO `user` (hashed_password, salt, first_name, last_name, dob, email, user_name) VALUES ")

	vals := []interface{}{}

	for _, user := range users {
		queryBuilder.WriteString("(?, ?, ?, ?, ?, ?, ?),")
		vals = append(vals, user.HashedPassword, user.Salt, user.FirstName, user.LastName, user.Dob, user.Email, user.UserName)
	}

	// Trim the last comma
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
}

func generateFollowers(db *sql.DB) {
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
					insertBatchFollowers(db, batch, vals)
					batch = []string{}
					vals = []interface{}{}
				}
			}
		}
	}

	// Insert any remaining followers
	if len(batch) > 0 {
		insertBatchFollowers(db, batch, vals)
	}

	// Assign an average of 200 followers to the rest of the users
	for userID := topUsers + 1; userID <= totalUsers; userID++ {
		numFollows := mathrand.Intn(avgFollows * 2) // To vary the number of follows per user around the average
		followers := mathrand.Perm(totalUsers)[:numFollows]
		for _, followerID := range followers {
			if followerID+1 != userID { // Ensure a user does not follow themselves
				batch = append(batch, "(?, ?)")
				vals = append(vals, userID, followerID+1)
				if len(batch) >= batchSize {
					insertBatchFollowers(db, batch, vals)
					batch = []string{}
					vals = []interface{}{}
				}
			}
		}
		if userID%1000 == 0 {
			fmt.Printf("Processed %d users\n", userID)
		}
	}

	// Insert any remaining followers
	if len(batch) > 0 {
		insertBatchFollowers(db, batch, vals)
	}
}

func insertBatchFollowers(db *sql.DB, batch []string, vals []interface{}) {
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

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Step 1: Generate and insert users
	for i := 0; i < totalUsers/batchSize; i++ {
		var users []User
		for j := 0; j < batchSize; j++ {
			id := i*batchSize + j
			users = append(users, generateUser(id))
		}
		insertUsers(db, users)
		fmt.Printf("Inserted %d users\n", (i+1)*batchSize)
	}

	// Step 2: Generate followers after users have been seeded
	generateFollowers(db)
}
