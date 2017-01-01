package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	log "github.com/Sirupsen/logrus"
)

var db *bolt.DB

func createTable(name string) bool{
	tableName := []byte(name)
	log.Info("Creating table with name " + name)
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tableName)
		if err != nil {
			return err
		}
		return nil
    })

    if err != nil {
        log.Fatal(err)
		return false
    }
	
	return true
}

func ignoreChannel(channelID string, channelName string, serverID string) string{
	log.Info("Adding channel #"+channelName + " to ignore list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:IC1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ignoreList"))
		err := b.Put([]byte(channelID), []byte(serverID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:IC2")
	}

	return ("Added #" + channelName + " to the ignore List")
}

func listenChannel(channelID string, channelName string, serverID string) string{
	log.Info("Removing " + channelName + " from the ignore list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:LC1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ignoreList"))
		err := b.Delete([]byte(channelID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:LC2")
	}

	return ("Removed #" + channelName + " from the ignore List")
}

func ignoreVoice(channelID string, channelName string, serverID string) string{
	log.Info("Adding channel #"+channelName + " to ignore list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:IVC1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vignoreList"))
		err := b.Put([]byte(channelID), []byte(serverID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:IVC2")
	}

	return ("Added #" + channelName + " to the ignore List")
}

func listenVoice(channelID string, channelName string, serverID string) string{
	log.Info("Removing " + channelName + " from the ignore list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:LVC1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vignoreList"))
		err := b.Delete([]byte(channelID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:LVC2")
	}

	return ("Removed #" + channelName + " from the ignore List")
}

func addMod(userID string, userName string, serverID string) string{
	log.Info("Adding " + userName + " to mod list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:AM1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("modList"))
		err := b.Put([]byte(userID), []byte(serverID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:AM2")
	}

	return ("Added " + userName + " as a mod")
}

func removeMod(userID string, userName string, serverID string) string{
	log.Info("Removing " + userName + " from the mod list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:RM1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("modList"))
		err := b.Delete([]byte(userID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:RM2")
	}

	return ("Removed " + userName + " as a mod")
}

func block(userID string, userName string, serverID string) string{
	log.Info("Adding " + userName + " to block list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:BU1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockList"))
		err := b.Put([]byte(userID), []byte(serverID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:BU2")
	}

	return ("Added " + userName + " to the blocked list")
}

func unblock(userID string, userName string, serverID string) string{
	log.Info("Removing " + userName + " from the block list")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:UB1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockList"))
		err := b.Delete([]byte(userID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:UB2")
	}

	return ("Removed " + userName + " to the blocked list")
}

func mute(serverID string) string{
	log.Info("Muted bot on server")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:M1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("muteList"))
		err := b.Put([]byte(serverID), []byte("true"))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:M2")
	}

	return ("Muted bot on this server")
}

func unmute(serverID string) string{
	log.Info("Unmuted bot on server")
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:UM1")
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("muteList"))
		err := b.Delete([]byte(serverID))
		return err
	})
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:UM2")
	}

	return ("Unmuted bot on this server")
}

func isIgnored(channelID string, serverID string) bool{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()
	
	is := false
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("ignoreList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == channelID && string(v) == serverID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func voiceIsIgnored(channelID string, serverID string) bool{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()
	
	is := false
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("ignoreList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == channelID && string(v) == serverID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func isMod(userID string, serverID string) bool{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()
	
	is := false
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("modList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == userID && string(v) == serverID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func isBlocked(userID string, serverID string) bool{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()
	
	is := false
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == userID && string(v) == serverID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func isMuted(serverID string) bool{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()
	
	is := false
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("muteList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == serverID && string(v) == "true"{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func getModList(serverID string) string{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:ML1")
	}
	defer db.Close()
	
	buffer := bytes.NewBufferString("")
	buffer.WriteString("** Bot Mod List **")
	buffer.WriteString("```")
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("modList"))

		b.ForEach(func(k, v []byte) error {
			if string(v) == serverID{
				name := ""
				member, err := discord.State.Member(string(v), string(k))
				if err != nil {
					user, _ := discord.User(string(k))
					name = user.Username
				}else{
					name = member.Nick
				}
				buffer.WriteString(name + " \n")
			}
			return nil
		})
		return nil
	})
	buffer.WriteString("```")
	return buffer.String();
}

func getBlockList(serverID string) string{
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return ("Err contact administrator. Code:BL1")
	}
	defer db.Close()
	
	buffer := bytes.NewBufferString("")
	buffer.WriteString("** Blocked User List **")
	buffer.WriteString("```")
	
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockList"))

		b.ForEach(func(k, v []byte) error {
			if string(v) == serverID{
				name := ""
				member, err := discord.State.Member(string(v), string(k))
				if err != nil {
					user, _ := discord.User(string(k))
					name = user.Username
				}else{
					name = member.Nick
				}
				buffer.WriteString(name + " \n")
			}
			return nil
		})
		return nil
	})
	buffer.WriteString("```")
	return buffer.String();
}