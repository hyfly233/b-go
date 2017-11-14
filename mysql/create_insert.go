package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "your_username"
	dbPassword = "your_password"
	dbName     = "your_database"
)

func main() {
	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS test_table (
        id INT AUTO_INCREMENT PRIMARY KEY,
        tinyint_col TINYINT,
        smallint_col SMALLINT,
        mediumint_col MEDIUMINT,
        int_col INT,
        bigint_col BIGINT,
        float_col FLOAT,
        double_col DOUBLE,
        decimal_col DECIMAL(10, 2),
        date_col DATE,
        datetime_col DATETIME,
        timestamp_col TIMESTAMP,
        time_col TIME,
        year_col YEAR,
        char_col CHAR(255),
        varchar_col VARCHAR(255),
        binary_col BINARY(255),
        varbinary_col VARBINARY(255),
        tinyblob_col TINYBLOB,
        blob_col BLOB,
        mediumblob_col MEDIUMBLOB,
        longblob_col LONGBLOB,
        tinytext_col TINYTEXT,
        text_col TEXT,
        mediumtext_col MEDIUMTEXT,
        longtext_col LONGTEXT,
        enum_col ENUM('value1', 'value2', 'value3'),
        set_col SET('value1', 'value2', 'value3')
    ) ENGINE=InnoDB;
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	insertSQL := `
    INSERT INTO test_table (
        tinyint_col, smallint_col, mediumint_col, int_col, bigint_col,
        float_col, double_col, decimal_col, date_col, datetime_col,
        timestamp_col, time_col, year_col, char_col, varchar_col,
        binary_col, varbinary_col, tinyblob_col, blob_col, mediumblob_col,
        longblob_col, tinytext_col, text_col, mediumtext_col, longtext_col,
        enum_col, set_col
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ { // Insert 1000 rows of test data
		_, err = stmt.Exec(
			rand.Intn(128), rand.Intn(256), rand.Intn(512), rand.Intn(1024), rand.Int63(),
			rand.Float32(), rand.Float64(), rand.Float64()*100, "2025-03-04", "2025-03-04 12:34:56",
			"2025-03-04 12:34:56", "12:34:56", 2025, "char_data", "varchar_data",
			[]byte("binary_data"), []byte("varbinary_data"), []byte("tinyblob_data"), []byte("blob_data"), []byte("mediumblob_data"),
			[]byte("longblob_data"), "tinytext_data", "text_data", "mediumtext_data", "longtext_data",
			"value1", "value1,value2",
		)
		if err != nil {
			log.Fatalf("Failed to insert test data: %v", err)
		}
	}

	log.Println("Test data inserted successfully")
}
