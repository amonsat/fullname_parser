package fullname_parser

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ParsedName struct {
	Title  string
	First  string
	Middle string
	Last   string
	Nick   string
	Suffix string
}

var (
	suffixList = []string{"esq", "esquire", "jr", "jnr", "sr", "snr", "2", "ii", "iii", "iv",
		"v", "clu", "chfc", "cfp", "md", "phd", "j.d.", "ll.m.", "m.d.", "d.o.", "d.c.",
		"p.c.", "ph.d."}

	prefixList = []string{"a", "ab", "antune", "ap", "abu", "al", "alm", "alt", "bab", "bäck",
		"bar", "bath", "bat", "beau", "beck", "ben", "berg", "bet", "bin", "bint", "birch",
		"björk", "björn", "bjur", "da", "dahl", "dal", "de", "degli", "dele", "del",
		"della", "der", "di", "dos", "du", "e", "ek", "el", "escob", "esch", "fleisch",
		"fitz", "fors", "gott", "griff", "haj", "haug", "holm", "ibn", "kauf", "kil",
		"koop", "kvarn", "la", "le", "lind", "lönn", "lund", "mac", "mhic", "mic", "mir",
		"na", "naka", "neder", "nic", "ni", "nin", "nord", "norr", "ny", "o", "ua", `ui\'`,
		"öfver", "ost", "över", "öz", "papa", "pour", "quarn", "skog", "skoog", "sten",
		"stor", "ström", "söder", "ter", "ter", "tre", "türk", "van", "väst", "väster",
		"vest", "von"}

	titleList = []string{"mr", "mrs", "ms", "miss", "dr", "herr", "monsieur", "hr", "frau",
		"a v m", "admiraal", "admiral", "air cdre", "air commodore", "air marshal",
		"air vice marshal", "alderman", "alhaji", "ambassador", "baron", "barones",
		"brig", "brig gen", "brig general", "brigadier", "brigadier general",
		"brother", "canon", "capt", "captain", "cardinal", "cdr", "chief", "cik", "cmdr",
		"coach", "col", "colonel", "commandant", "commander", "commissioner",
		"commodore", "comte", "comtessa", "congressman", "conseiller", "consul",
		"conte", "contessa", "corporal", "councillor", "count", "countess",
		"crown prince", "crown princess", "dame", "datin", "dato", "datuk",
		"datuk seri", "deacon", "deaconess", "dean", "dhr", "dipl ing", "doctor",
		"dott", "dott sa", "dr ing", "dra", "drs", "embajador", "embajadora", "en",
		"encik", "eng", "eur ing", "exma sra", "exmo sr", "f o", "father",
		"first lieutient", "first officer", "flt lieut", "flying officer", "fr",
		"frau", "fraulein", "fru", "gen", "generaal", "general", "governor", "graaf",
		"gravin", "group captain", "grp capt", "h e dr", "h h", "h m", "h r h", "hajah",
		"haji", "hajim", "her highness", "her majesty", "herr", "high chief",
		"his highness", "his holiness", "his majesty", "hon", "hr", "hra", "ing", "ir",
		"jonkheer", "judge", "justice", "khun ying", "kolonel", "lady", "lcda", "lic",
		"lieut", "lieut cdr", "lieut col", "lieut gen", "lord", "m", "m l", "m r",
		"madame", "mademoiselle", "maj gen", "major", "master", "mevrouw", "miss",
		"mlle", "mme", "monsieur", "monsignor", "mstr", "nti", "pastor",
		"president", "prince", "princess", "princesse", "prinses", "prof",
		"prof sir", "professor", "puan", "puan sri", "rabbi", "rear admiral", "rev",
		"rev canon", "rev dr", "rev mother", "reverend", "rva", "senator", "sergeant",
		"sheikh", "sheikha", "sig", "sig na", "sig ra", "sir", "sister", "sqn ldr", "sr",
		"sr d", "sra", "srta", "sultan", "tan sri", "tan sri dato", "tengku", "teuku",
		"than puying", "the hon dr", "the hon justice", "the hon miss", "the hon mr",
		"the hon mrs", "the hon ms", "the hon sir", "the very rev", "toh puan", "tun",
		"vice admiral", "viscount", "viscountess", "wg cdr"}

	conjunctionList = []string{"&", "and", "et", "e", "of", "the", "und", "y"}
)

var (
	nameParts  []string
	nameCommas []bool
)

func ParseFullname(fullname string) (parsedName ParsedName) {
	log.Debug("Start parsing fullname: ", fullname)

	//nicknames: remove and store
	nicknames := findNicknames(&fullname)
	parsedName.Nick = strings.Join(nicknames, ",")

	//split name to parts and store commas
	splitName(fullname)

	//suffix: remove and store
	if len(nameParts) > 1 {
		suffixes := findSuffixes()
		parsedName.Suffix = strings.Join(suffixes, ", ")
	}

	//titles: remove and store
	if len(nameParts) > 1 {
		titles := findTitles()
		parsedName.Title = strings.Join(titles, ", ")
	}

	// Join name prefixes to following names
	if len(nameParts) > 1 {
		joinPrefixes()
	}

	// Join conjunctions to surrounding names
	if len(nameParts) > 1 {
		joinConjunctions()
	}

	// Suffix: remove and store items after extra commas as suffixes
	if len(nameParts) > 1 {
		extraSuffixes := findExtraSuffixes()
		if len(extraSuffixes) > 0 {
			if parsedName.Suffix != "" {
				parsedName.Suffix += ", " + strings.Join(extraSuffixes, ", ")
			} else {
				parsedName.Suffix = strings.Join(extraSuffixes, ", ")
			}
		}
	}

	// Last name: remove and store last name
	if len(nameParts) > 0 {
		parsedName.Last = findLastname()
	}

	// First name: remove and store first part as first name
	if len(nameParts) > 0 {
		parsedName.First = findFirstname()
	}

	// Middle name: store all remaining parts as middle name
	if len(nameParts) > 0 {
		parsedName.Middle = findMiddlename()
	}

	log.Debugf("Parsing complete: %+v", parsedName)
	return
}

func findNicknames(fullname *string) []string {
	var re = regexp.MustCompile(`\s?[\'\"\(\[]([^\[\]\)\)\'\"]+)[\'\"\)\]]`)
	var partsFound []string
	tempString := *fullname

	matches := re.FindAllStringSubmatch(tempString, -1)
	for _, v := range matches {
		partsFound = append(partsFound, v[1])
	}

	log.Debugf("Founded %v nickname(s): %v", len(partsFound), partsFound)
	log.Debug("Clearing")

	for _, v := range matches {
		tempString = strings.Replace(tempString, v[0], "", -1)
	}
	*fullname = tempString

	log.Debug("Cleared fullname: ", *fullname)
	return partsFound
}

func findSuffixes() []string {
	log.Debug("Searching suffixes")
	return findParts(suffixList)
}

func findTitles() []string {
	log.Debug("Searching titles")
	return findParts(titleList)
}

func splitName(fullname string) {
	log.Debug("Spliting fullname")
	re := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	fullname = re.ReplaceAllLiteralString(fullname, " ")
	nameParts = strings.Split(strings.TrimSpace(fullname), " ")
	for i, v := range nameParts {
		nameParts[i] = strings.TrimSpace(v)
		nameCommas = append(nameCommas, false)
		if strings.HasSuffix(v, ",") {
			nameParts[i] = strings.TrimSuffix(v, ",")
			nameCommas[i] = true
		}
	}

	log.Debug("Splitted parts: ", nameParts)
	log.Debug("Splitted commas: ", nameCommas)
}

func findParts(list []string) []string {
	var partsFound []string

	for _, namePart := range nameParts {
		if namePart == "" {
			continue
		}

		partToCheck := strings.ToLower(namePart)
		partToCheck = strings.TrimSuffix(partToCheck, ".")

		for _, suf := range list {
			if suf == partToCheck {
				partsFound = append(partsFound, namePart)
			}
		}
	}

	log.Debugf("Founded %v parts: %v", len(partsFound), partsFound)
	log.Debug("Clearing")

	for _, partFound := range partsFound {
		foundIndex := -1
		for i, namePart := range nameParts {
			if partFound == namePart {
				foundIndex = i
				break
			}
		}
		if foundIndex > -1 {
			nameParts = append(nameParts[:foundIndex], nameParts[foundIndex+1:]...)
			if nameCommas[foundIndex] && foundIndex != len(nameCommas)-1 {
				nameCommas = append(nameCommas[:foundIndex+1], nameCommas[foundIndex+1+1:]...)
			} else {
				nameCommas = append(nameCommas[:foundIndex], nameCommas[foundIndex+1:]...)
			}
		}
	}

	log.Debug("Cleared parts: ", nameParts)
	log.Debug("Cleared commas: ", nameCommas)

	return partsFound
}

func joinPrefixes() {
	log.Debug("Join prefixes")

	if len(nameParts) > 1 {
		for i := len(nameParts) - 2; i >= 0; i-- {
			for _, pref := range prefixList {
				if pref == nameParts[i] {
					nameParts[i] = nameParts[i] + " " + nameParts[i+1]
					nameParts = append(nameParts[:i+1], nameParts[i+2:]...)
					nameCommas = append(nameCommas[:i], nameCommas[i+1:]...)
				}
			}
		}
	}

	log.Debug("Prefixes joined: ", strings.Join(nameParts, ","))
	log.Debug("Cleared commas: ", nameCommas)
}

func joinConjunctions() {
	log.Debug("Join conjunctions")

	if len(nameParts) > 2 {
		for i := len(nameParts) - 3; i >= 0; i-- {
			for _, conj := range conjunctionList {
				if conj == nameParts[i+1] {
					nameParts[i] = nameParts[i] + " " + nameParts[i+1] + " " + nameParts[i+2]
					nameParts = append(nameParts[:i+1], nameParts[i+3:]...)
					nameCommas = append(nameCommas[:i], nameCommas[i+2:]...)
					i--
				}
			}
		}
	}
	log.Debug("Conjunctions joined: ", strings.Join(nameParts, ","))
	log.Debug("Cleared commas: ", nameCommas)
}

func findExtraSuffixes() (extraSuffixes []string) {
	commasCount := 0
	for _, v := range nameCommas {
		if v {
			commasCount++
		}
	}
	if commasCount > 1 {
		for i := len(nameParts) - 1; i >= 2; i-- {
			if nameCommas[i] {
				extraSuffixes = append(extraSuffixes, nameParts[i])
				nameParts = append(nameParts[:i], nameParts[i+1:]...)
				nameCommas = append(nameCommas[:i], nameCommas[i+1:]...)
			}
		}
	}

	log.Debugf("Founded %v extra suffixes: %v", len(extraSuffixes), extraSuffixes)
	log.Debug("Cleared commas: ", nameCommas)

	return
}

func findLastname() (lastname string) {
	log.Debug("Searching lastname")

	commaIndex := -1
	for i, v := range nameCommas {
		if v {
			commaIndex = i
		}
	}

	if commaIndex == -1 {
		commaIndex = len(nameParts) - 1
	}

	lastname = nameParts[commaIndex]
	nameParts = append(nameParts[:commaIndex], nameParts[commaIndex+1:]...)
	nameCommas = nameCommas[:0]

	log.Debug("Founded lastname: ", lastname)
	log.Debug("Cleared parts: ", nameParts)
	log.Debug("Cleared commas: ", nameCommas)
	return
}

func findFirstname() (firstname string) {
	log.Debug("Searching firstname")
	firstname = nameParts[0]
	nameParts = nameParts[1:]
	log.Debug("Founded firstname: ", firstname)
	log.Debug("Cleared parts: ", nameParts)
	return
}

func findMiddlename() (middlename string) {
	log.Debug("Searching middlename")
	middlename = strings.Join(nameParts, " ")
	nameParts = nameParts[:0]
	log.Debug("Founded middlename(s): ", middlename)
	log.Debug("Cleared parts: ", nameParts)
	return
}
