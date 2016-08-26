package main

import (
	"bytes"
	"encoding/binary"
	//"encoding/json"
	"flag"
	"fmt"
	"io"
	//"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/boltdb/bolt"
	redis "gopkg.in/redis.v3"
)

var (
	// discordgo session
	discord *discordgo.Session

	// Redis client connection (used for stats)
	rcli *redis.Client

	// Map of Guild id's to *Play channels, used for queuing and rate-limiting guilds
	queues map[string]chan *Play = make(map[string]chan *Play)

	// Sound encoding settings
	BITRATE        = 128
	MAX_QUEUE_SIZE = 6

	// Owner
	OWNER string
)

// Play represents an individual use of the !airhorn command
type Play struct {
	GuildID   string
	ChannelID string
	UserID    string
	Sound     *Sound

	// The next play to occur after this, only used for chaining sounds like anotha
	Next *Play

	// If true, this was a forced play using a specific airhorn sound name
	Forced bool
}

type SoundCollection struct {
	Prefix    string
	Commands  []string
	Sounds    []*Sound
	ChainWith *SoundCollection

	soundRange int
}

// Sound represents a sound clip
type Sound struct {
	Name string

	// Weight adjust how likely it is this song will play, higher = more likely
	Weight int

	// Delay (in milliseconds) for the bot to wait before sending the disconnect request
	PartDelay int

	// Buffer to store encoded PCM packets
	buffer [][]byte
}

// Array of all the sounds we have
var AIRHORN *SoundCollection = &SoundCollection{
	Prefix: "airhorn",
	Commands: []string{
		"!airhorn",
	},
	Sounds: []*Sound{
		createSound("default", 1000, 250),
		createSound("reverb", 800, 250),
		createSound("spam", 800, 0),
		createSound("tripletap", 800, 250),
		createSound("fourtap", 800, 250),
		createSound("distant", 500, 250),
		createSound("echo", 500, 250),
		createSound("clownfull", 250, 250),
		createSound("clownshort", 250, 250),
		createSound("clownspam", 250, 0),
		createSound("highfartlong", 200, 250),
		createSound("highfartshort", 200, 250),
		createSound("midshort", 100, 250),
		createSound("truck", 10, 250),
	},
}

var KHALED *SoundCollection = &SoundCollection{
	Prefix:    "another",
	ChainWith: AIRHORN,
	Commands: []string{
		"!anotha",
		"!anothaone",
	},
	Sounds: []*Sound{
		createSound("one", 1, 250),
		createSound("one_classic", 1, 250),
		createSound("one_echo", 1, 250),
	},
}

var CENA *SoundCollection = &SoundCollection{
	Prefix: "jc",
	Commands: []string{
		"!johncena",
		"!cena",
	},
	Sounds: []*Sound{
		createSound("airhorn", 1, 250),
		createSound("echo", 1, 250),
		createSound("full", 1, 250),
		createSound("jc", 1, 250),
		createSound("nameis", 1, 250),
		createSound("spam", 1, 250),
	},
}

var ETHAN *SoundCollection = &SoundCollection{
	Prefix: "ethan",
	Commands: []string{
		"!ethan",
		"!eb",
		"!ethanbradberry",
		"!h3h3",
	},
	Sounds: []*Sound{
		createSound("areyou_classic", 100, 250),
		createSound("areyou_condensed", 100, 250),
		createSound("areyou_crazy", 100, 250),
		createSound("areyou_ethan", 100, 250),
		createSound("classic", 100, 250),
		createSound("echo", 100, 250),
		createSound("high", 100, 250),
		createSound("slowandlow", 100, 250),
		createSound("cuts", 30, 250),
		createSound("beat", 30, 250),
		createSound("sodiepop", 1, 250),
	},
}

var COW *SoundCollection = &SoundCollection{
	Prefix: "cow",
	Commands: []string{
		"!stan",
		"!stanislav",
	},
	Sounds: []*Sound{
		createSound("herd", 10, 250),
		createSound("moo", 10, 250),
		createSound("x3", 1, 250),
	},
}

var BIRTHDAY *SoundCollection = &SoundCollection{
	Prefix: "birthday",
	Commands: []string{
		"!birthday",
		"!bday",
	},
	Sounds: []*Sound{
		createSound("horn", 50, 250),
		createSound("horn3", 30, 250),
		createSound("sadhorn", 25, 250),
		createSound("weakhorn", 25, 250),
	},
}

var WOW *SoundCollection = &SoundCollection{
	Prefix: "wow",
	Commands: []string{
		"!wowthatscool",
		"!wtc",
	},
	Sounds: []*Sound{
		createSound("thatscool", 50, 250),
	},
}

var MOAN *SoundCollection = &SoundCollection{
	Prefix: "moan",
	Commands: []string{
		"!moan",
	},
	Sounds: []*Sound{
		createSound("1", 1000, 250),
		createSound("2", 1000, 250),
		createSound("3", 1000, 250),
		createSound("4", 1000, 250),
		createSound("5", 1000, 250),
	},
}

var SANDSTORM *SoundCollection = &SoundCollection{
	Prefix: "sandstorm",
	Commands: []string{
		"!ss",
		"!sandstorm",
	},
	Sounds: []*Sound{
		createSound("toy1", 1000, 250),
		createSound("toy2", 1000, 250),
		createSound("toy3", 1000, 250),
		createSound("toy4", 1000, 250),
		createSound("toy5", 1000, 250),
		createSound("1", 1000, 250),
		createSound("2", 1000, 250),
		createSound("3", 1000, 250),
		createSound("4", 1000, 250),
		createSound("5", 1000, 250),
		createSound("6", 1000, 250),
		createSound("7", 1000, 250),
	},
}

var FLANZY *SoundCollection = &SoundCollection{
	Prefix: "flanzy",
	Commands: []string{
		"!iflanzy",
		"!flanzy",
	},
	Sounds: []*Sound{
		createSound("burp", 1000, 250),
		createSound("hole", 1000, 250),
		createSound("moist1", 10, 250),
		createSound("moist2", 10, 250),
		createSound("slurp1", 1000, 250),
		createSound("slurp2", 1000, 250),
		createSound("slurp3", 1000, 250),
		createSound("slurp4", 1000, 250),
		createSound("slurp5", 1000, 250),
		createSound("slurp6", 1000, 250),
		createSound("slurp7", 1000, 250),
		createSound("spooning", 1000, 250),
	},
}

var REK *SoundCollection = &SoundCollection{
	Prefix: "rek",
	Commands: []string{
		"!rek",
	},
	Sounds: []*Sound{
		createSound("law", 1000, 250),
	},
}
var FITNESS *SoundCollection = &SoundCollection{
	Prefix: "fitness",
	Commands: []string{
		"!fit",
		"!fitness",
	},
	Sounds: []*Sound{
		createSound("1", 1000, 250),
		createSound("begin", 1000, 250),
		createSound("bing", 1000, 250),
		createSound("bring", 1000, 250),
		createSound("curlup", 1000, 250),
		createSound("over", 1000, 250),
		createSound("pacer", 1000, 250),
		createSound("run", 1000, 250),
		createSound("running", 1000, 250),
		createSound("start", 1000, 250),
		createSound("updown", 1000, 250),
	},
}

var COLLECTIONS []*SoundCollection = []*SoundCollection{
	AIRHORN,
	KHALED,
	CENA,
	ETHAN,
	COW,
	BIRTHDAY,
	WOW,
	MOAN,
	SANDSTORM,
	FLANZY,
	REK,
	FITNESS,
}

var db *bolt.DB

func CreateBucket(name []byte){
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // store some data
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
    })

    if err != nil {
        log.Fatal(err)
    }
	
	
}

func checkList(m *discordgo.MessageCreate) bool{
	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.Warning("Failed to grab channel with id " + m.ChannelID)
		return false
	}
	var channelID = channel.ID
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	is := false
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("ignoreList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == channelID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func check(m *discordgo.MessageCreate) bool{
	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.Warning("Failed to grab channel with id " + m.ChannelID)
		return false
	}
	var channelID = channel.ID
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	is := false
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("ignoreList"))

		b.ForEach(func(k, v []byte) error {
			if string(k) == channelID{
				is = true
			}
			return nil
		})
		return nil
	})
	return is
}

func addToList(m *discordgo.MessageCreate){
	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.Warning("Failed to grab channel with id " + m.ChannelID)
		return
	}
	var channelName = channel.Name
	var channelID = channel.ID
	log.Info("Ignoring " + channelName)
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ignoreList"))
		err := b.Put([]byte(channelID), []byte(channelName))
		return err
	})
}

func removeFromList(m *discordgo.MessageCreate){
	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.Warning("Failed to grab channel with id " + m.ChannelID)
		return
	}
	var channelName = channel.Name
	var channelID = channel.ID
	log.Info("Ignoring " + channelName)
	
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ignoreList"))
		err := b.Delete([]byte(channelID))
		return err
	})
}

// Create a Sound struct
func createSound(Name string, Weight int, PartDelay int) *Sound {
	return &Sound{
		Name:      Name,
		Weight:    Weight,
		PartDelay: PartDelay,
		buffer:    make([][]byte, 0),
	}
}

func (sc *SoundCollection) Load() {
	for _, sound := range sc.Sounds {
		sc.soundRange += sound.Weight
		sound.Load(sc)
	}
}

func (s *SoundCollection) Random() *Sound {
	var (
		i      int
		number int = randomRange(0, s.soundRange)
	)

	for _, sound := range s.Sounds {
		i += sound.Weight

		if number < i {
			return sound
		}
	}
	return nil
}

// Load attempts to load an encoded sound file from disk
// DCA files are pre-computed sound files that are easy to send to Discord.
// If you would like to create your own DCA files, please use:
// https://github.com/nstafie/dca-rs
// eg: dca-rs --raw -i <input wav file> > <output file>
func (s *Sound) Load(c *SoundCollection) error {
	path := fmt.Sprintf("audio/%v_%v.dca", c.Prefix, s.Name)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// read opus frame length from dca file
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			fmt.Println("error reading from dca file :", err)
			return err
		}

		// read encoded pcm from dca file
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("error reading from dca file :", err)
			return err
		}

		// append encoded pcm data to the buffer
		s.buffer = append(s.buffer, InBuf)
	}
}

// Plays this sound over the specified VoiceConnection
func (s *Sound) Play(vc *discordgo.VoiceConnection) {
	vc.Speaking(true)
	defer vc.Speaking(false)

	for _, buff := range s.buffer {
		vc.OpusSend <- buff
	}
}

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

// Prepares a play
func createPlay(user *discordgo.User, guild *discordgo.Guild, coll *SoundCollection, sound *Sound) *Play {
	// Grab the users voice channel
	channel := getCurrentVoiceChannel(user, guild)
	if channel == nil {
		log.WithFields(log.Fields{
			"user":  user.ID,
			"guild": guild.ID,
		}).Warning("Failed to find channel to play sound in")
		return nil
	}

	// Create the play
	play := &Play{
		GuildID:   guild.ID,
		ChannelID: channel.ID,
		UserID:    user.ID,
		Sound:     sound,
		Forced:    true,
	}

	// If we didn't get passed a manual sound, generate a random one
	if play.Sound == nil {
		play.Sound = coll.Random()
		play.Forced = false
	}

	// If the collection is a chained one, set the next sound
	if coll.ChainWith != nil {
		play.Next = &Play{
			GuildID:   play.GuildID,
			ChannelID: play.ChannelID,
			UserID:    play.UserID,
			Sound:     coll.ChainWith.Random(),
			Forced:    play.Forced,
		}
	}

	return play
}

// Prepares and enqueues a play into the ratelimit/buffer guild queue
func enqueuePlay(user *discordgo.User, guild *discordgo.Guild, coll *SoundCollection, sound *Sound) {
	play := createPlay(user, guild, coll, sound)
	if play == nil {
		return
	}

	// Check if we already have a connection to this guild
	//   yes, this isn't threadsafe, but its "OK" 99% of the time
	_, exists := queues[guild.ID]

	if exists {
		if len(queues[guild.ID]) < MAX_QUEUE_SIZE {
			queues[guild.ID] <- play
		}
	} else {
		queues[guild.ID] = make(chan *Play, MAX_QUEUE_SIZE)
		playSound(play, nil)
	}
}

func trackSoundStats(play *Play) {
	if rcli == nil {
		return
	}

	_, err := rcli.Pipelined(func(pipe *redis.Pipeline) error {
		var baseChar string

		if play.Forced {
			baseChar = "f"
		} else {
			baseChar = "a"
		}

		base := fmt.Sprintf("airhorn:%s", baseChar)
		pipe.Incr("airhorn:total")
		pipe.Incr(fmt.Sprintf("%s:total", base))
		pipe.Incr(fmt.Sprintf("%s:sound:%s", base, play.Sound.Name))
		pipe.Incr(fmt.Sprintf("%s:user:%s:sound:%s", base, play.UserID, play.Sound.Name))
		pipe.Incr(fmt.Sprintf("%s:guild:%s:sound:%s", base, play.GuildID, play.Sound.Name))
		pipe.Incr(fmt.Sprintf("%s:guild:%s:chan:%s:sound:%s", base, play.GuildID, play.ChannelID, play.Sound.Name))
		pipe.SAdd(fmt.Sprintf("%s:users", base), play.UserID)
		pipe.SAdd(fmt.Sprintf("%s:guilds", base), play.GuildID)
		pipe.SAdd(fmt.Sprintf("%s:channels", base), play.ChannelID)
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Failed to track stats in redis")
	}
}

// Play a sound
func playSound(play *Play, vc *discordgo.VoiceConnection) (err error) {
	log.WithFields(log.Fields{
		"play": play,
	}).Info("Playing sound")

	if vc == nil {
		vc, err = discord.ChannelVoiceJoin(play.GuildID, play.ChannelID, false, false)
		// vc.Receive = false
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to play sound")
			delete(queues, play.GuildID)
			return err
		}
	}

	// If we need to change channels, do that now
	if vc.ChannelID != play.ChannelID {
		vc.ChangeChannel(play.ChannelID, false, false)
		time.Sleep(time.Millisecond * 125)
	}

	// Track stats for this play in redis
	go trackSoundStats(play)

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(time.Millisecond * 32)

	// Play the sound
	play.Sound.Play(vc)

	// If this is chained, play the chained sound
	if play.Next != nil {
		playSound(play.Next, vc)
	}

	// If there is another song in the queue, recurse and play that
	if len(queues[play.GuildID]) > 0 {
		play := <-queues[play.GuildID]
		playSound(play, vc)
		return nil
	}

	// If the queue is empty, delete it
	time.Sleep(time.Millisecond * time.Duration(play.Sound.PartDelay))
	delete(queues, play.GuildID)
	vc.Disconnect()
	return nil
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "Airhorny Bot")
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
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

func scontains(key string, options ...string) bool {
	for _, item := range options {
		if item == key {
			return true
		}
	}
	return false
}

func calculateAirhornsPerSecond( m *discordgo.MessageCreate) {
	current, _ := strconv.Atoi(rcli.Get("airhorn:a:total").Val())
	time.Sleep(time.Second * 10)
	latest, _ := strconv.Atoi(rcli.Get("airhorn:a:total").Val())

	discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current APS: %v", (float64(latest-current))/10.0))
}

func displayBotStats(cid string) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	users := 0
	for _, guild := range discord.State.Ready.Guilds {
		users += len(guild.Members)
	}

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "Discordgo: \t%s\n", discordgo.VERSION)
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	fmt.Fprintf(w, "Memory: \t%s / %s (%s total allocated)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Tasks: \t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "Servers: \t%d\n", len(discord.State.Ready.Guilds))
	fmt.Fprintf(w, "Users: \t%d\n", users)
	fmt.Fprintf(w, "```\n")
	w.Flush()
	discord.ChannelMessageSend(cid, buf.String())
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

func displayUserStats(cid, uid string, s *discordgo.Session) {
	keys, err := rcli.Keys(fmt.Sprintf("airhorn:*:user:%s:sound:*", uid)).Result()
	if err != nil {
		return
	}

	totalAirhorns := utilSumRedisKeys(keys)
	_, _ = s.ChannelMessageSend(cid, fmt.Sprintf("Total Airhorns: %v", totalAirhorns))
}

func displayServerStats(cid, sid string, s *discordgo.Session) {
	keys, err := rcli.Keys(fmt.Sprintf("airhorn:*:guild:%s:sound:*", sid)).Result()
	if err != nil {
		return
	}

	totalAirhorns := utilSumRedisKeys(keys)
	_, _ = s.ChannelMessageSend(cid, fmt.Sprintf("Total Airhorns: %v", totalAirhorns))
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
	} else if scontains(parts[1], "bomb") && len(parts) >= 3 {
		airhornBomb(m.ChannelID, g, utilGetMentioned(s, m), parts[2], s)
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
		_, _ = s.ChannelMessageSend(m.ChannelID, input)
	} else if scontains(parts[1], "ignore") {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.Warning("Failed to grab channel with id " + m.ChannelID)
			return
		}
		addToList(m)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Added #" + channel.Name + " to the ignore list")
	} else if scontains(parts[1], "listen") {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.Warning("Failed to grab channel with id " + m.ChannelID)
			return
		}
		removeFromList(m)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Removed #" + channel.Name + " to the ignore list")
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
	cmds := "``` **Airhorn Bot Commands** \n !airhorn \n !anotha , !anothaone \n !johncena , !cena \n !ethan , !eb , !ethanbradberry , !h3h3 \n !stanislav , !stan \n !birthday , !bday \n !wowthatscool , !wtc \n !moan \n !sandstorm , !ss \n !flanzy , !iflanzy \n !fit , !fitness \n !rek ```"

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
			}else if len(parts) == 1{
				_, _ = s.ChannelMessageSend(m.ChannelID, cmds)
			}
		}
		return
	}

	// Find the collection for the command we got
	for _, coll := range COLLECTIONS {
		if scontains(parts[0], coll.Commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var sound *Sound
			if len(parts) > 1 {
				for _, s := range coll.Sounds {
					if parts[1] == s.Name {
						sound = s
					}
				}

				if sound == nil {
					if parts[0] == "!rek" {
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
	
	// reply with formatted sound list
	
	if parts[0] == "!sounds" {
		log.Info("Listing Commands")
		if len(parts) > 1 {
			for _, coll := range COLLECTIONS {
				if scontains("!"+parts[1], coll.Commands...) {
					soundList := "``` **Sounds for " + parts[1] + "** \n"
					for _, s := range coll.Sounds {
						soundList = soundList + s.Name + "\n"
					}
					soundList = soundList + "```"
					_, _ = s.ChannelMessageSend(m.ChannelID, soundList)
					return
				}
			}
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, cmds)
		return
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
	dbName := []byte("ignoreList")
	CreateBucket(dbName);
	dbName = []byte("modList")
	CreateBucket(dbName);
	dbName = []byte("muteList")
	CreateBucket(dbName);
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
