package types

var (
	JournalColors      = []string{"RED", "ORANGE", "YELLOW", "GREEN", "MINT", "TEAL", "CYAN", "BLUE", "INDIGO", "PURPLE", "PINK", "BROWN", "FIRE_ENGINE_RED", "CORAL", "TOMATO", "CINNABAR", "RUST", "ORANGE_RED", "AUBURN", "SADDLE_BROWN", "DARK_ORANGE", "DARK_GOLD", "KHAKI", "BRONZE", "GOLD", "OLIVE", "OLIVE_DRAB", "DARK_OLIVE_GREEN", "MOSS_GREEN", "LIME_GREEN", "LIME", "FOREST_GREEN", "SEA_GREEN", "JUNGLE_GREEN", "LIGHT_SEA_GREEN", "DARK_TURQUOISE", "DODGER_BLUE", "ROYAL_BLUE", "DEEP_LAVENDER", "BLUE_VIOLET", "DARK_VIOLET", "HELIOTROPE", "BYZANTIUM", "MAGENTA", "DARK_MAGENTA", "FUCHSIA", "DEEP_PINK", "GRAYISH_MAGENTA", "HOT_PINK", "JAZZBERRY_JAM", "MAROON"}
	JournalClassColors = map[string]string{}
)

type Journal struct {
	Experiences       []Experience       `json:"experiences"`
	JournalCompanions []JournalCompanion `json:"substanceCompanions"`
}

type JournalCompanion struct {
	SubstanceName string       `json:"substanceName"`
	Color         JournalColor `json:"color"`
}

type JournalCompanions []JournalCompanions

type Experience struct {
	Title        string      `json:"title"`
	Text         string      `json:"text"`
	CreationDate int64       `json:"creationDate"`
	SortDate     int64       `json:"sortDate"`
	Ingestions   []Ingestion `json:"ingestions"`
}

type Ingestion struct {
	SubstanceName       string  `json:"substanceName"`
	Time                *int64  `json:"time"`
	EndTime             *int64  `json:"endTime"`
	CreationDate        int64   `json:"creationDate"`
	AdministrationRoute *string `json:"administrationRoute"`
	Dose                float64 `json:"dose"`
	IsDoseAnEstimate    bool    `json:"isDoseAnEstimate"`
	Units               string  `json:"units"`
	Notes               string  `json:"notes"`
}

type JournalColor string

func (j Journal) UniqueSubstances() UniqueSubstances {
	m := make(map[string]int64)

	for _, c := range j.JournalCompanions {
		if n, ok := m[c.SubstanceName]; !ok {
			m[c.SubstanceName] = 1
		} else {
			m[c.SubstanceName] = n + 1
		}
	}

	return m
}

func (j Journal) Append(doses []Dose) Journal {
	d := Doses(doses)

	j.Experiences = append(
		j.Experiences,
		d.ToExperienceGroups()...,
	)
	j.JournalCompanions = append(
		j.JournalCompanions,
		d.ToJournalCompanions(j.UniqueSubstances())...,
	)

	return j
}
