package main

import (
    "log"
    "os"

    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

func main() {
    log.Println("Connecting to Gmail server")
    connection, err := client.DialTLS("imap.gmail.com:993", nil)
    if err != nil {
        log.Fatal(err)
    }

    defer connection.Logout()
    log.Println("Connected!")
    emailAddress := os.Getenv("EMAIL_ADDRESS")
    emailPasswd := os.Getenv("EMAIL_PASSWD")
    err = connection.Login(emailAddress, emailPasswd)
    if err != nil {
        log.Fatalf("Unable to login to email account. Reason: %s", err)
    }
    log.Println("Logged in!")

    inbox, err := connection.Select("INBOX", true)
    if err != nil {
        log.Fatalf("Unable to select INBOX. Reason: %s\n", err)
    }
    log.Println("Inbox selected")
    messageCount := inbox.Messages
    log.Printf("%d messages currently in INBOX\n", messageCount)
    to := inbox.Messages
    from := to - 10
    seqSet := new(imap.SeqSet)
    seqSet.AddRange(from, to)
    messages := make(chan *imap.Message, 10)
    done := make(chan error, 1)
    go func() {
        done <- connection.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages)
    }()
    for msg := range messages {
        log.Printf("Subject: %s, Date :%s\n", msg.Envelope.Subject, msg.Envelope.Date)
    }
    err = <-done
    if err != nil {
        log.Fatalf("Unable to fetch more messages. Reason: %s\n", err)
    }
    log.Println("Done")
}
