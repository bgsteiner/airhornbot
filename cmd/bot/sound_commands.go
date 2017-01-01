package main

var AIRHORN *SoundCollection = &SoundCollection{
	Prefix: "airhorn",
	Commands: []string{
		"airhorn",
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
		"anotha",
		"anothaone",
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
		"johncena",
		"cena",
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
		"ethan",
		"eb",
		"ethanbradberry",
		"h3h3",
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
		"stan",
		"stanislav",
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
		"birthday",
		"bday",
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
		"wowthatscool",
		"wtc",
	},
	Sounds: []*Sound{
		createSound("thatscool", 50, 250),
	},
}

var MOAN *SoundCollection = &SoundCollection{
	Prefix: "moan",
	Commands: []string{
		"moan",
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
		"ss",
		"sandstorm",
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
		"iflanzy",
		"flanzy",
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
		"rek",
	},
	Sounds: []*Sound{
		createSound("law", 1000, 250),
	},
}
var FITNESS *SoundCollection = &SoundCollection{
	Prefix: "fitness",
	Commands: []string{
		"fit",
		"fitness",
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

var CANADA *SoundCollection = &SoundCollection{
	Prefix:    "canada",
	Commands: []string{
		"canada",
	},
	Sounds: []*Sound{
		createSound("anthem", 10, 250),
		createSound("moose1", 10, 250),
		createSound("moose2", 10, 250),
	},
}

var MEMES *SoundCollection = &SoundCollection{
    Prefix: "meme",
    Commands: []string{
        "meme",
        "maymay",
        "memes",
    },
    Sounds: []*Sound{
        createSound("headshot", 3, 250),
        createSound("wombo", 3, 250),
        createSound("triple", 3, 250),
        createSound("camera", 3, 250),
        createSound("gandalf", 3, 250),
        createSound("mad", 50, 0),
        createSound("ateam", 50, 0),
        createSound("bennyhill", 50, 0),
        createSound("tuba", 50, 0),
        createSound("donethis", 50, 0),
        createSound("leeroy", 50, 0),
        createSound("slam", 50, 0),
        createSound("nerd", 50, 0),
        createSound("kappa", 50, 0),
        createSound("digitalsports", 50, 0),
        createSound("csi", 50, 0),
        createSound("nogod", 50, 0),
        createSound("welcomebdc", 50, 0),
    },
}

var DICTION *SoundCollection = &SoundCollection{
    Prefix: "diction",
    Commands: []string{
        "diction",
        "emd",
    },
    Sounds: []*Sound{
		createSound("chandelier", 3, 250),
		createSound("fly", 3, 250),
		createSound("123drink", 3, 250),
		createSound("shame", 3, 250),
    },
}

var BATMAN *SoundCollection = &SoundCollection{
	Prefix: "bat",
	Commands: []string{
		"batman",
		"bat",
	},
	Sounds: []*Sound{
		createSound("batman", 3, 250),
		createSound("imbatman", 3, 250),
	},
}

var RICKROLL *SoundCollection = &SoundCollection{
	Prefix: "rickroll",
	Commands: []string{
		"rickroll",
		"rr",
	},
	Sounds: []*Sound{
		createSound("1", 3, 250),
		createSound("2", 3, 250),
		createSound("3", 3, 250),
		createSound("4", 3, 250),
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
	CANADA,
	MEMES,
	DICTION,
	BATMAN,
	RICKROLL,
}