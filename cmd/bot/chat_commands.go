package main

import (
    "bytes"
	"strings"

    log "github.com/Sirupsen/logrus"
    "github.com/bwmarrin/discordgo"
)

type TextFunction func(*discordgo.Session, *discordgo.MessageCreate, []string, *discordgo.Guild) string

type TextCollection struct {
    Commands    []string
    Function    TextFunction
	HelpText	string
}

var LISTSOUNDS *TextCollection = &TextCollection{
    Commands: []string{
        "sounds",
    },
    Function: listSFX,
	HelpText: "Used to get list of sound commands",
}

var HELP *TextCollection = &TextCollection{
    Commands: []string{
        "help",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		return getHelp()
    },
	HelpText: "",
}

var IGNORE *TextCollection = &TextCollection{
    Commands: []string{
        "ignore",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
        return ignoreChannel(m.ChannelID, channel.Name, channel.GuildID)
    },
	HelpText: "Tells the bot to ignore a text channel \n Usage: !dicbot ignore",
}

var LISTEN *TextCollection = &TextCollection{
    Commands: []string{
        "listen",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
        return listenChannel(m.ChannelID, channel.Name, channel.GuildID)
    },
	HelpText: "Tells the bot to listen to a text channel \n Usage: !dicbot listen",
}

var VIGNORE *TextCollection = &TextCollection{
    Commands: []string{
        "vignore",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel := getCurrentVoiceChannel(m.Author, g)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to find channel to play sound in")
		}
        return ignoreChannel(channel.ID, channel.Name, g.ID)
    },
	HelpText: "Tells the bot to avoid a voice channel \n Usage: !dicbot vignore",
}
var VLISTEN *TextCollection = &TextCollection{
    Commands: []string{
        "vlisten",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel := getCurrentVoiceChannel(m.Author, g)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to find channel to play sound in")
		}
        return listenChannel(channel.ID, channel.Name, g.ID)
    },
	HelpText: "Tells the bot to use a voice channel \n Usage: !dicbot vlisten",
}

var MOD *TextCollection = &TextCollection{
    Commands: []string{
        "mod",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		output := ""
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
		if len(m.Mentions) >= 1{
			for _, mention := range m.Mentions {
				log.Info("Adding mod " + mention.Username)
				output = addMod(mention.ID, mention.Username, channel.GuildID)
				break
			}
		}else{
			output = "Please use mentions for this command"
		}
		return output
    },
	HelpText: "Adds a mod able to use the admin commands \n Usage: !dicbot mod @user",
}

var UNMOD *TextCollection = &TextCollection{
    Commands: []string{
        "unmod",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		output := ""
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
		if len(m.Mentions) >= 1{
			for _, mention := range m.Mentions {
				output = removeMod(mention.ID, mention.Username, channel.GuildID)
				break
			}
		}else{
			output = "Please use mentions for this command"
		}
		return output
    },
	HelpText: "Removes a mod from being able to use the admin commands \n Usage: !dicbot unmod @user",
}

var BLOCK *TextCollection = &TextCollection{
    Commands: []string{
        "block",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		output := ""
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
		if len(m.Mentions) >= 1{
			for _, mention := range m.Mentions {
				output = block(mention.ID, mention.Username, channel.GuildID)
				break
			}
		}else{
			output = "Please use mentions for this command"
		}
		return output
    },
	HelpText: "Blocks a user from being able to use the sound commands \n Usage: !dicbot block @user",
}

var UNBLOCK *TextCollection = &TextCollection{
    Commands: []string{
        "unblock",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		output := ""
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
		if len(m.Mentions) >= 1{
			for _, mention := range m.Mentions {
				output = unblock(mention.ID, mention.Username, channel.GuildID)
				break
			}
		}else{
			output = "Please use mentions for this command"
		}
		return output
    },
	HelpText: "Removes a blocked user from the list \n Usage: !dicbot unblock @user",
}

var MUTE *TextCollection = &TextCollection{
    Commands: []string{
        "mute",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
        return mute(channel.GuildID)
    },
	HelpText: "Global mutes the bot on the server \n Usage: !dicbot mute",
}

var UNMUTE *TextCollection = &TextCollection{
    Commands: []string{
        "unmute",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		channel, _ := discord.State.Channel(m.ChannelID)
		if channel == nil {
			log.WithFields(log.Fields{
				"channel": m.ChannelID,
				"message": m.ID,
			}).Warning("Failed to grab channel")
		}
        return unmute(channel.GuildID)
    },
	HelpText: "Unmutes the bot for the server \n Usage: !dicbot unmute",
}

var MODLIST *TextCollection = &TextCollection{
    Commands: []string{
        "modlist",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
        return getModList(g.ID)
    },
	HelpText: "Gets the list of bot mods for the server \n Usage: !dicbot modlist",
}

var BLOCKLIST *TextCollection = &TextCollection{
    Commands: []string{
        "blocklist",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
        return getBlockList(g.ID)
    },
	HelpText: "Gets the list of blocked users for the server \n Usage: !dicbot blocklist",
}

var BG *TextCollection = &TextCollection{
    Commands: []string{
        "bg",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
        if(m.Author.Username=="bgsteiner"){
			//bgsteiner ,_ := discord.User("150017863090962432")
			return "Even if they dont love you I still love you <@150017863090962432>"
		}
		return ""
    },
	HelpText: "Nope.avi",
}

var STATUS *TextCollection = &TextCollection{
    Commands: []string{
        "status",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
        displayBotStats(m.ChannelID)
		return ""
    },
	HelpText: "Displays status of the bot \n Usage: !dicbot status",
}

var STATS *TextCollection = &TextCollection{
    Commands: []string{
        "stats",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
        if len(m.Mentions) >= 1 {
			displayUserStats(m.ChannelID, utilGetMentioned(s, m).ID, s)
		} else if len(parts) >= 3 {
			displayUserStats(m.ChannelID, parts[2], s)
		} else {
			displayServerStats(m.ChannelID, g.ID, s)
		}
		return ""
    },
	HelpText: "Displays status of the bot \n Usage: !dicbot status",
}

var SETSTATUS *TextCollection = &TextCollection{
    Commands: []string{
        "setstatus",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		msg := ""
		for i := 2; i < len(parts); i++ {
			msg += parts[i] + " "
		}
        s.UpdateStatus(0, strings.Title(msg))
		return ""
    },
	HelpText: "Sets the Playing status of the bot \n Usage: !dicbot status",
}

var LEADERBOARD *TextCollection = &TextCollection{
    Commands: []string{
        "leaderboard",
    },
    Function: func (s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string {
		return getLeaderboard(g.ID, s)
    },
	HelpText: "Gets Leaderboard \n Usage: !dicbot leaderboard",
}

var CHATCOMMANDS []*TextCollection = []*TextCollection{
    LISTSOUNDS,
	BG,
}

var SUBCOMMANDS []*TextCollection = []*TextCollection{
	IGNORE,
	LISTEN,
	MOD,
	UNMOD,
	BLOCK,
	UNBLOCK,
	MUTE,
	UNMUTE,
	VIGNORE,
	VLISTEN,
	MODLIST,
	BLOCKLIST,
	STATUS,
	STATS,
	SETSTATUS,
	LEADERBOARD,
}

func listSFX(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) string{

	if(len(parts) == 1){
		buffer := bytes.NewBufferString("")
		buffer.WriteString("** Dichorn Bot sound commands **")
		buffer.WriteString("```")
		for _, coll := range COLLECTIONS {
			for _, cmds := range coll.Commands {
				buffer.WriteString("!"+cmds)
				buffer.WriteString(", ")
			}
			buffer.WriteString("\n")
		}
		buffer.WriteString("```")
		return buffer.String();
	}
	
	if(len(parts) >= 1){
		soundCmd := strings.Replace(parts[1], "!", "", 1)
		exists := false
		for _, coll := range COLLECTIONS {
			for _, cmds := range coll.Commands {
				if(soundCmd==cmds){
					exists = true
				}
			}
		}
		if(exists){
			buffer := bytes.NewBufferString("")
			buffer.WriteString("** Sounds for !"+soundCmd + "**")
			buffer.WriteString("``` \n")
			for _, coll := range COLLECTIONS {
				if scontains(soundCmd, coll.Commands...) {
					for _, snd := range coll.Sounds {
						buffer.WriteString(snd.Name)
						buffer.WriteString("\n")
					}
				}
			}
			buffer.WriteString("```")
			return buffer.String();
		}else{
			return "Sorry that is not a command"
		}
	}
	return ""
}
func getHelp() string{
	buffer := bytes.NewBufferString("")
	buffer.WriteString("** !dicbot Bot Commands **")
	buffer.WriteString("```")
	for _, coll := range SUBCOMMANDS {
		buffer.WriteString(coll.Commands[0] + " - ")
		buffer.WriteString(coll.HelpText)
		buffer.WriteString(" \n")
	}
	buffer.WriteString("```")
	return buffer.String();
}