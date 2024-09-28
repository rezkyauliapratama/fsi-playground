package main

import (
	"database/sql"
	"net/smtp"

	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Struct to represent the transaction data in the "after" field
type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      string    `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Struct to represent the Kafka message containing the transaction data
type KafkaMessage struct {
	After Transaction `json:"after"`
}

// Function to query user data from the MySQL database
func getUserByID(db *sql.DB, userID string) (string, error) {
	var email string
	query := "Select email from dbz_user_management_users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&email)
	return email, err
}

func main() {
	// SMTP server configuration for MailCatcher
	smtpHost := "localhost"
	smtpPort := "1025" // Default MailCatcher SMTP port
	from := "transaction@bank.com"

	// MySQL connection setup
	dsn := "root:root_password@tcp(localhost:4002)/mail"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Consumer configuration for default semantics (At-Least-Once)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:19092,localhost:29092,localhost:39092", // Kafka broker address
		"group.id":                "mail-system",                                     // Consumer group ID
		"auto.offset.reset":       "earliest",                                        // Start reading from the earliest offset
		"enable.auto.commit":      true,                                              // Automatically commit offsets (default)
		"auto.commit.interval.ms": 1000,                                              // Commit offsets every 5 seconds (default)
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	topic := "trx.transaction.transactions"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	// Main loop to consume messages from Kafka
	log.Println("Consumer started, waiting for messages...")
	for {
		// Read a message from the Kafka topic
		msg, err := consumer.ReadMessage(-1) // Blocking call to read messages
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
		}

		// Create an instance of KafkaMessage to hold the deserialized data
		var kafkaMessage KafkaMessage

		// Unmarshal the JSON message into the KafkaMessage struct
		err = json.Unmarshal([]byte(msg.Value), &kafkaMessage)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		// Query the user from the MySQL database
		email, err := getUserByID(db, kafkaMessage.After.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("user with ID %s not found", kafkaMessage.After.UserID)
			} else {
				log.Printf("failed to query user: %v", err)
			}
			continue
		}

		to := []string{email}

		// Email content
		subject := fmt.Sprintf("Subject: Transaction IDR %s in Bank xx \n", kafkaMessage.After.Amount)
		body := fmt.Sprintf("You do transaction IDR %s ", kafkaMessage.After.Amount)
		message := []byte(subject + "\n" + body)

		// Send email
		err = smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, message)
		if err != nil {
			log.Fatal("Failed to send email:", err)
		}

		// Simulate processing delay
		time.Sleep(2 * time.Second)
	}
}
