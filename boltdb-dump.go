package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	bolt "github.com/coreos/bbolt"
)

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}

func dumpCursor(tx *bolt.Tx, c *bolt.Cursor, indent int) {
	indentStr := strings.Repeat(" ", indent-1)

	for k, v := c.First(); k != nil; k, v = c.Next() {
		var keyStr string
		if isPrintable(k, false) {
			keyStr = "\"" + string(k) + "\""
		} else {
			keyStr = hex.EncodeToString(k)
		}

		if v == nil {
			fmt.Printf(indentStr+"[%s]\n", keyStr)
			newBucket := c.Bucket().Bucket(k)
			if newBucket == nil {
				// from the top-level cursor and not a cursor from a bucket
				newBucket = tx.Bucket(k)
			}
			newCursor := newBucket.Cursor()
			dumpCursor(tx, newCursor, indent+1)
		} else {
			fmt.Printf(indentStr+"%s\n", keyStr)
			var valueStr string
			if isPrintable(v, true) {
				valueStr = string(v)
			} else {
				valueStr = hex.Dump(v)
			}

			// XXX: If this hits a text only value that's larger than
			// bufio.MaxScanTokenSize, bad things will happen.
			s := bufio.NewScanner(bytes.NewReader([]byte(valueStr)))
			for s.Scan() {
				fmt.Printf(indentStr+"  %s\n", s.Text())
			}
			if err := s.Err(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

// Dump everything in the database.
func dump(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		dumpCursor(tx, c, 1)
		return nil
	})
}

func isPrintable(b []byte, allowNL bool) bool {
	for _, v := range b {
		switch {
		case v > 127:
			// Out of the 7 bit ASCII range.
			return false
		case v == 10, v == 13:
			// LF, CR.
			if !allowNL {
				return false
			}
		case v < 32:
			// Other non-printable.
			return false
		default:
		}
	}
	return true
}

func main() {
	// check we have a filename
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename.db>\n", os.Args[0])
		os.Exit(2)
	}

	// the first arg is the database file
	filename := os.Args[1]

	// open this file
	opts := bolt.Options{
		ReadOnly: true,
		Timeout:  1 * time.Second,
	}
	db, err := bolt.Open(filename, 0600, &opts)
	check(err)
	defer db.Close()

	// dump all keys
	dump(db)
}
