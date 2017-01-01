package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	log "github.com/Sirupsen/logrus"
	"sort"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

type Stats struct {
	ID string
	Uses int
}

type PairList []Stats

func (a PairList) Len() int           { return len(a) }
func (a PairList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PairList) Less(i, j int) bool { return a[i].Uses > a[j].Uses }

func addStatServer(serverID string){
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Batch(func(tx *bolt.Tx) error {
		
		b := tx.Bucket([]byte("serverStats"))
		v := b.Get([]byte(serverID))
		
		tmp, _ := strconv.Atoi(string(v))
		buf := strconv.Itoa(tmp+1)
		err := b.Put([]byte(serverID), []byte(buf))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func addStatUser(serverID string, UserID string){
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Batch(func(tx *bolt.Tx) error {
		
		b := tx.Bucket([]byte(serverID))
		v := b.Get([]byte(UserID))
		tmp, _ := strconv.Atoi(string(v))
		buf := strconv.Itoa(tmp+1)
		err := b.Put([]byte(UserID), []byte(buf))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getServerStats(serverID string) string{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:Stats3")
	}
	defer db.Close()
	tmp := ""
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("serverStats"))
		v := b.Get([]byte(serverID))
		tmp = string(v)
		return nil
	})
	return "Total Air Horns: " + tmp;
}

func getUserStats(serverID string, UserID string) string{

	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:Stats4")
	}
	defer db.Close()
	tmp := ""
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(serverID))
		v := b.Get([]byte(UserID))
		tmp = string(v)
		return nil
	})
	return "Total Air Horns: " + tmp;
}

func getLeaderboard(serverID string, s *discordgo.Session) string{

	users := []Stats{}

	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:Stats4")
	}
	defer db.Close()
	
	db.View(func(tx *bolt.Tx) error {
	
		b := tx.Bucket([]byte(serverID))
		
		i := 0
		b.ForEach(func(k, v []byte) error {
			tmp, _ := strconv.Atoi(string(v))
			users = append(users, Stats{string(k), tmp})
			i++
			return nil
		})
		return nil

	})
	sort.Sort(PairList(users))
	buffer := bytes.NewBufferString("")
	buffer.WriteString("** Most Horny List **")
	buffer.WriteString("```")
	for i := 0;i < len(users); i++ {
		user, err := s.User(users[i].ID)
		if err != nil {
			log.Fatal(err)
			return ("Err contact administrator. Code:Stats5")
		}
		index := strconv.Itoa(i+1)
		buffer.WriteString(index + " ➤  # " + user.Username + " \n")
		buffer.WriteString("\t " + strconv.Itoa(users[i].Uses) + " \n")
		if (i==10){
			i=len(users)
		}
	}
	
	buffer.WriteString("```")
	return buffer.String();
}

func statTrack(play *Play){
	createTable(play.GuildID)
	addStatServer(play.GuildID)
	addStatUser(play.GuildID, play.UserID)
}