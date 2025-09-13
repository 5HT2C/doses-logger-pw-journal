package types

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	timeFmt = "Mon Jan 02 2006"
)

var (
	timeLoc = make(map[string]time.Location)

	dosageRegex = regexp.MustCompile(`([0-9]+)(\.[0-9]+)?([\p{Greek}A-Za-z]+)`)
	doseStim    = regexp.MustCompile(`(am(ph|f)etamine|Bromantane|phenidate|-[A-Z]PH|-[FCM]?M[CA]|Pi?[HV]P|finil|drone|ylamine|MDM?A|APB)`)
	doseAmph    = regexp.MustCompile(`(am(ph|f)etamine|-[FCM]?M[CA])`)
	doseExcl    = regexp.MustCompile(`(alcohol|cigarette)`)
	doseOpiate  = regexp.MustCompile(`(Heroin|O-DSMT|Tramadol|AP-237|Kratom|codine)`)
	doseDisso   = regexp.MustCompile(`(PCE|PCP|PCM|DXM|Ketamine|Memantine|phenidine|Nitrous)`)
	doseBenzo   = regexp.MustCompile(`(epam|olam)`)
	doseGABA    = regexp.MustCompile(`(Phenibut|Gabapentin|Pregabalin|GHB|GBL|BDO|Alcohol)`)
	doseTrypt   = regexp.MustCompile(`[45]-[A-Za-z]{2,3}-[A-Za-z]{3,4}`)
)

type TimeData struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Timezone  string    `json:"timezone,omitempty"`
}

func (t *TimeData) TimeInLocation() time.Time {
	loc, ok := timeLoc[t.Timezone]

	if !ok {
		if l, err := time.LoadLocation(t.Timezone); err != nil {
			return t.Timestamp
		} else {
			loc = *l
		}
	}

	return t.Timestamp.In(&loc)
}

type Dose struct { // timezone,date,time,dosage,drug,roa,note
	Position int `json:"position"` // order added, at the top so that it's marshaled as at the top
	TimeData
	Created *TimeData `json:"created,omitempty"`
	Date    string    `json:"date,omitempty"`
	Time    string    `json:"time,omitempty"`
	Dosage  Dosage    `json:"dosage,omitempty"`
	Drug    string    `json:"drug,omitempty"`
	RoA     RoA       `json:"roa,omitempty"`
	Note    string    `json:"note,omitempty"`
}

type Doses []Dose

type RoA string

// ToIngestionRoA converts the RoA to Journal format
// Docs: https://github.com/pwarchive/psychonautwiki-journal-android/blob/73f013752cea2f05558c1ed091cdccd3dfcde62b/app/src/main/java/com/isaakhanimann/journal/data/substances/AdministrationRoute.kt#L27
func (r RoA) ToIngestionRoA() *string {
	roa := ""

	switch r {
	case
		"Buccal",
		"Inhaled",
		"Insufflated",
		"Intramuscular",
		"Oral",
		"Rectal",
		"Smoked",
		"Subcutaneous",
		"Sublingual":
		roa = strings.ToUpper(string(r))

	case "Intrabuccal":
		roa = "INTRAMUSCULAR"

	case "Intranasal":
		roa = "INSUFFLATED"

	case "Topical":
		roa = "TRANSDERMAL"

	case "Vaporized":
		roa = "SMOKED"
	}

	return &roa
}

type Dosage string

func (d Dosage) ToIngestionDosage() (float64, string) {
	//log.Printf("%#v\n", dosageRegex.FindStringSubmatch(string(d)))

	if len(d) == 0 {
		return 0, string(d)
	}

	m := dosageRegex.FindStringSubmatch(string(d))
	dose, err := strconv.ParseFloat(m[1]+m[2], 64)
	unit := m[3]

	if err != nil {
		if doseInt, err := strconv.ParseInt(m[1], 10, 64); err != nil {
			panic(err)
		} else {
			dose = float64(doseInt)
		}
	}

	return dose, unit
}

func (d Dose) ToIngestion() Ingestion {
	dTime := d.Timestamp.UnixMilli()
	dCreated := dTime

	if d.Created != nil {
		dCreated = d.Created.Timestamp.UnixMilli()
	}

	dDose, dUnit := d.Dosage.ToIngestionDosage()

	return Ingestion{
		SubstanceName:       d.Drug,
		Time:                &dTime,
		EndTime:             nil, // TODO: impl
		CreationDate:        dCreated,
		AdministrationRoute: d.RoA.ToIngestionRoA(),
		Dose:                dDose,
		IsDoseAnEstimate:    false,
		Units:               dUnit,
		Notes:               d.Note,
	}
}

func (doses Doses) ToExperienceGroups() []Experience {
	lastDate := ""
	lastTitle := ""
	lastGroup := Experience{}

	groups := make([]Experience, 0)

	for _, d := range doses {
		// If the current dose is on a new day, split into a new group
		if d.Date != lastDate || len(lastTitle) == 0 {
			// Add the previous group if we have a valid lastGroup (i.e., no default value)
			if lastGroup.CreationDate != 0 {
				groups = append(groups, lastGroup)
			}

			lastDate = d.Date
			lastTitle = d.TimeInLocation().Format(timeFmt)
			lastGroup = Experience{
				Title:        lastTitle,
				Text:         "",
				CreationDate: d.Timestamp.UnixMilli(),
				SortDate:     d.Timestamp.UnixMilli(),
				Ingestions:   make([]Ingestion, 0),
			}
		}

		// We have a group to append to, or we just made one
		lastGroup.Ingestions = append(lastGroup.Ingestions, d.ToIngestion())
	}

	return groups
}

func (doses Doses) UniqueSubstances() UniqueSubstances {
	m := make(map[string]int64)

	for _, d := range doses {
		if n, ok := m[d.Drug]; !ok {
			m[d.Drug] = 1
		} else {
			m[d.Drug] = n + 1
		}
	}

	return m
}

func (doses Doses) ToJournalCompanions(jUniq UniqueSubstances) []JournalCompanion {
	dUniq := doses.UniqueSubstances()
	jComp := make([]JournalCompanion, 0)

	for substance, _ := range dUniq {
		// Create companion info for Journal
		if _, ok := jUniq[substance]; !ok {
			companion := JournalCompanion{
				SubstanceName: substance,
				Color:         "CYAN",
			}

			jComp = append(jComp, companion)
		}
	}

	return jComp
}

func (d Dose) ParsedTime() (time.Time, error) {
	timeZero := time.Unix(0, 0)
	loc, err := time.LoadLocation(d.Timezone)
	if err != nil {
		return timeZero, err
	}

	if pt, err := time.ParseInLocation("2006/01/0215:04", d.Date+d.Time, loc); err == nil {
		return pt, nil
	} else {
		return timeZero, err
	}
}
