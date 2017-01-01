package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	redis "gopkg.in/redis.v3"
)

var (
	// discordgo session
	discord *discordgo.Session

	// Redis client connection (used for stats)
	rcli *redis.Client

	// Sound encoding settings
	BITRATE        = 128
	MAX_QUEUE_SIZE = 6

	// Owner
	OWNER string
	
	//command prefix
	PREFIX = "!"
	
	CHATCOMMAND = "dicbot"
)

// Attempts to find the current users voice channel inside a given guild
func getCurrentVoiceChannel(user *discordgo.User, guild *discordgo.Guild) *discordgo.Channel {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			channel, _ := discord.State.Channel(vs.ChannelID)
			return channel
		}
	}
	return nil
}

// Returns a random integer between min and max
func randomRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "Airhorny Bot")
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	//event.Guild.Unavailable != nil
	log.Info(event.Guild.Unavailable)
	if event.Guild.Unavailable != nil {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			s.ChannelMessageSend(channel.ID, "**AIRHORN BOT READY FOR HORNING. TYPE `!AIRHORN` WHILE IN A VOICE CHANNEL TO ACTIVATE**")
			return
		}
	}
}

func utilSumRedisKeys(keys []string) int {
	results := make([]*redis.StringCmd, 0)

	rcli.Pipelined(func(pipe *redis.Pipeline) error {
		for _, key := range keys {
			results = append(results, pipe.Get(key))
		}
		return nil
	})

	var total int
	for _, i := range results {
		t, _ := strconv.Atoi(i.Val())
		total += t
	}

	return total
}

func utilGetMentioned(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.User {
	for _, mention := range m.Mentions {
		if mention.ID != s.State.Ready.User.ID {
			return mention
		}
	}
	return nil
}

func airhornBomb(cid string, guild *discordgo.Guild, user *discordgo.User, cs string, s *discordgo.Session) {
	count, _ := strconv.Atoi(cs)
	_, _ = s.ChannelMessageSend(cid, ":ok_hand:"+strings.Repeat(":trumpet:", count))

	// Cap it at something
	if count > 100 {
		return
	}

	play := createPlay(user, guild, AIRHORN, nil)
	vc, err := discord.ChannelVoiceJoin(play.GuildID, play.ChannelID, true, true)
	if err != nil {
		return
	}

	for i := 0; i < count; i++ {
		AIRHORN.Random().Play(vc)
	}

	vc.Disconnect()
}

// Handles bot operator messages, should be refactored (lmao)
func handleBotControlMessages(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) {
	if scontains(parts[1], "status") {
		displayBotStats(m.ChannelID)
	} else if scontains(parts[1], "stats") {
		if len(m.Mentions) >= 2 {
			displayUserStats(m.ChannelID, utilGetMentioned(s, m).ID, s)
		} else if len(parts) >= 3 {
			displayUserStats(m.ChannelID, parts[2], s)
		} else {
			displayServerStats(m.ChannelID, g.ID, s)
		}
	} else if scontains(parts[1], "aps") {
		_, _ = s.ChannelMessageSend(m.ChannelID, ":ok_hand: give me a sec m8")
		go calculateAirhornsPerSecond(m)
	} else if scontains(parts[1], "sudo") {
		input := ""
		for i, v := range parts {
			if i > 1 {
				input += v + " "
			}
		}
		log.Info("Sudo: " + input)
		_, _ = s.ChannelMessageSend(m.ChannelID, input)
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if len(m.Content) <= 0 || (m.Content[0] != '!' && len(m.Mentions) < 1) {
		return
	}

	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := discord.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}
	
	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(strings.ToLower(msg), " ")

	// If this is a mention, it should come from the owner (otherwise we don't care)
	if len(m.Mentions) >= 1{
		mentioned := false
		for _, mention := range m.Mentions {
			mentioned = (mention.ID == s.State.Ready.User.ID)
			if mentioned {
				break
			}
		}
		if mentioned {
			if  m.Author.ID == OWNER && len(parts) > 1{
				handleBotControlMessages(s, m, parts, guild)
				return
			}
		}
	}
	
	if !strings.HasPrefix(m.Content, "!") && len(m.Mentions) < 1 {
		return
	}
	
	baseCommand := strings.Replace(parts[0], PREFIX, "", 1)
	
	if(baseCommand == CHATCOMMAND){
		if(isMod(m.Author.ID, channel.GuildID) || m.Author.ID == OWNER || m.Author.ID == guild.OwnerID){
			if(len(parts)==1 || parts[1]=="help"){
				_, _ = s.ChannelMessageSend(m.ChannelID,getHelp())
				return
			}
		
			log.Info("Processing Command " + parts[1])
			subCommand := parts[1]
			for _, acoll := range SUBCOMMANDS {
				if scontains(subCommand, acoll.Commands...) {
					_, _ = s.ChannelMessageSend(m.ChannelID,
						acoll.Function(s, m, parts, guild))
				}
			}
		}
	}else{
		for _, ccoll := range CHATCOMMANDS {
			if scontains(baseCommand, ccoll.Commands...) {
				_, _ = s.ChannelMessageSend(m.ChannelID,
					ccoll.Function(s, m, parts, guild))
			}
		}
	}
	
	if(isIgnored(m.ChannelID, channel.GuildID)){
		log.Info(m.ChannelID + " Is ignored")
		return
	}
	if(isBlocked(m.Author.ID, channel.GuildID)){
		log.Info(m.Author.ID + " Is blocked")
		return
	}
	if(isMuted(channel.GuildID)){
		log.Info(channel.GuildID + " Is muted")
		return
	}

	//Sound Commands
	for _, coll := range COLLECTIONS {
		if scontains(baseCommand, coll.Commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var sound *Sound
			if len(parts) > 1 {
				for _, s := range coll.Sounds {
					if parts[1] == s.Name {
						sound = s
					}
				}

				if sound == nil {
					if baseCommand == "rek" {
						_, _ = s.ChannelMessageSend(m.ChannelID, "** Banned "+parts[1]+" **")
						go enqueuePlay(m.Author, guild, coll, sound)
					}
					return
				}
			}

			go enqueuePlay(m.Author, guild, coll, sound)
			return
		}
	}
}

func main() {

	var (
		Token      = flag.String("t", "", "Discord Authentication Token")
		Redis      = flag.String("r", "", "Redis Connection String")
		Shard      = flag.String("s", "", "Shard ID")
		ShardCount = flag.String("c", "", "Number of shards")
		Owner      = flag.String("o", "", "Owner ID")
		err        error
	)
	flag.Parse()
	
	f, err := os.OpenFile("airhornbot.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Info("Log Started")
	log.Info("Connecting to DB")
	createTable("ignoreList");
	createTable("modList");
	createTable("muteList");
	createTable("blockList");
	createTable("serverStats");
	log.Info("DB Setup")

	if *Owner != "" {
		OWNER = *Owner
	}

	// Preload all the sounds
	log.Info("Preloading sounds...")
	for _, coll := range COLLECTIONS {
		coll.Load()
	}

	// If we got passed a redis server, try to connect
	if *Redis != "" {
		log.Info("Connecting to redis...")
		rcli = redis.NewClient(&redis.Options{Addr: *Redis, DB: 0})
		_, err = rcli.Ping().Result()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Failed to connect to redis")
			return
		}
	}

	// Create a discord session
	log.Info("Starting discord session...")
	discord, err = discordgo.New(*Token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord session")
		return
	}

	// Set sharding info
	discord.ShardID, _ = strconv.Atoi(*Shard)
	discord.ShardCount, _ = strconv.Atoi(*ShardCount)

	if discord.ShardCount <= 0 {
		discord.ShardCount = 1
	}

	discord.AddHandler(onReady)
	discord.AddHandler(onGuildCreate)
	discord.AddHandler(onMessageCreate)

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord websocket connection")
		return
	}

	// We're running!
	log.Info("AIRHORNBOT is ready to horn it up.")
	
	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
