package src

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func LoadUsersInfoFromCSV(filename string) map[int32]User {

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error to open filename: %s", err)
		return nil
	}
	defer f.Close()

	var membersInfo = make(map[int32]User)
	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		i, err := strconv.ParseInt(line[0], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		id := int32(i)

		i, err = strconv.ParseInt(line[6], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		timestamp := int32(i)

		i, err = strconv.ParseInt(line[7], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		joinedChatDate := int32(i)

		membersInfo[id] = User{Id: id, FirstName: line[1], LastName: line[2], Username: line[3], PhoneNumber: line[4], Status: line[5], Timestamp: timestamp, JoinChatDate: joinedChatDate}

	}
	return membersInfo
}

func SaveToCSV(data map[int32]User, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()
	var isFirstIteration = true

	var columnNames = []string{"ID", "FirstName", "LastName", "Username", "PhoneNumber", "Status", "Timestamp", "JoinedChatDate"}
	for _, value := range data {
		if isFirstIteration {
			if err := w.Write(columnNames); err != nil {
				return err
			}
		}
		isFirstIteration = false

		if err := w.Write([]string{fmt.Sprint(value.Id), value.FirstName, value.LastName, value.Username, value.PhoneNumber, value.Status, fmt.Sprint(value.Timestamp), fmt.Sprint(value.JoinChatDate)}); err != nil {
			return err
		}
	}
	return nil
}


func SaveMsgToCSV(data map[int64]Message, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()
	var isFirstIteration = true

	var columnNames = []string{"MessageID", "ContentType", "ContentText", "SenderUserID", "Date", "ReplyMessageID"}
	for _, value := range data {
		if isFirstIteration {
			if err := w.Write(columnNames); err != nil {
				return err
			}
		}
		isFirstIteration = false

		if err := w.Write([]string{fmt.Sprint(value.MsgId), value.ContentType, value.ContentText, fmt.Sprint(value.SenderUserId), fmt.Sprint(value.Date), fmt.Sprint(value.ReplyToMessageId)}); err != nil {
			return err
		}
	}
	return nil
}

func LoadMessagesFromCSV(filename string) map[int64]Message {

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error to open filename: %s", err)
		return nil
	}
	defer f.Close()

	var messages  = make(map[int64]Message)
	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		i, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		id := int64(i)

		i, err = strconv.ParseInt(line[3], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		senderUserId := int32(i)

		i, err = strconv.ParseInt(line[4], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		date := int32(i)

		i, err = strconv.ParseInt(line[5], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		replyMessageId := int64(i)

		messages[id] = Message{MsgId: id, ContentType: line[1], ContentText: line[2], SenderUserId: senderUserId, Date: date, ReplyToMessageId: replyMessageId}

	}
	return messages
}
