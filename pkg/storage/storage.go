package storage

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
	return &Storage{path: path}
}

// Checks if user has a record. If hasn't, creates a zero record.
// Otherwise, responds with a relevant message.
func (s *Storage) Start(chatID int64) string {

	// open file
	file, err := os.OpenFile(s.path, os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// check if user has a record
	for scanner.Scan() {

		pair := strings.Split(scanner.Text(), "=")

		// user has a record
		if pair[0] == strconv.FormatInt(chatID, 10) {

			// check user status

			if pair[1] == "1" {
				// subbed user
				return "Welcome back! You're already subbed. Type /unsub to stop getting memes."
			} else if pair[1] == "0" {
				// unsubbed user
				return "Welcome back! You aren't subbed. Type /sub to get memes."
			} else {
				// invalid user
				log.Println("ERROR: unknown user status.")
				return "Your status is invalid. Please, contact the admin about this! \n @januarycoming admin tg"
			}
		}
	}

	// user hasn't record
	writer := bufio.NewWriter(file)
	bs, err := writer.Write([]byte(strconv.Itoa(int(chatID)) + "=" + "0" + "\n"))
	log.Println("Message len: ", len([]byte(strconv.Itoa(int(chatID))+"="+"0"+"\n")))
	log.Println("Bytes wrote: ", bs)
	if err != nil {
		log.Println(err)
	}

	writer.Flush()

	return "Welcome! Type /sub to subscribe memes."
}

// looks for the user record and puts it to 1.
// If user hasn't record (unreachable scenario), asks to type /start.
func (s *Storage) Subscribe(chatID int64) string {

	// open file
	file, err := os.OpenFile(s.path, os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	editPos := 0

	// look for the record
	for scanner.Scan() {
		// edit position is the end if the line
		editPos += len(scanner.Text())
		pair := strings.Split(scanner.Text(), "=")

		// check if user has his record
		if pair[0] == strconv.FormatInt(chatID, 10) {

			// check if he is subbed
			if pair[1] == "1" {
				return "You are already subbed."
			} else {

				// sub user

				if editPos == 11 {
					editPos--
				} else {
					shift := editPos/10 - 2
					editPos += shift
				}

				_, err := file.Seek(int64(editPos), io.SeekStart)
				if err != nil {
					log.Fatal(err)
				}

				_, err = file.Write([]byte("1\n"))
				if err != nil {
					log.Fatal(err)
				}

				return "Now you're getting memes. Type /unsub to stop this."
			}
		}
	}

	return "Oh, you skipped the /start function, so i can't edit you in data base. \nPlease, type /start."
}

// looks for the user record and puts it to 0.
// If user hasn't record (unreachable scenario), asks to type /start.
func (s *Storage) Unsubscribe(chatID int64) string {

	// open file
	file, err := os.OpenFile(s.path, os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	editPos := 0

	// look for the record
	for scanner.Scan() {
		// edit position is the end if the line
		editPos += len(scanner.Text())
		pair := strings.Split(scanner.Text(), "=")

		// check if user has his record
		if pair[0] == strconv.FormatInt(chatID, 10) {

			// check if he is subbed
			if pair[1] == "1" {

				// unsub user

				if editPos == 11 {
					editPos--
				} else {
					shift := editPos/10 - 2
					editPos += shift
				}

				_, err := file.Seek(int64(editPos), io.SeekStart)
				if err != nil {
					log.Fatal(err)
				}

				_, err = file.Write([]byte("0\n"))
				if err != nil {
					log.Fatal(err)
				}

				return "You aren't getting memes anymore."

			} else {
				return "You are not subbed yet."
			}
		}
	}

	return "Oh, you skipped the /start function, so i can't edit you in data base. \nPlease, type /start"
}
