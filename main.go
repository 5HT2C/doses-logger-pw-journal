package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	fp      = flag.String("f", "", "file path")
	loadUrl = flag.String("url", "http://localhost:6010/media/doses.json", "URL for doses.json")
)

func main() {
	flag.Parse()

	var doses []Dose
	if err := getJsonFromUrl(&doses, *loadUrl); err != nil {
		panic(err)
	}

	var j Journal
	if len(*fp) > 0 {
		if b, err := os.ReadFile(*fp); err != nil {
			panic(err)
		} else {

			if err := json.Unmarshal(b, &j); err != nil {
				j = Journal{}
			}
		}
	}

	d := Doses(doses)
	j.Experiences = append(j.Experiences, d.ToExperienceGroups()...)

	dUniq := d.UniqueSubstances()
	jUniq := j.UniqueSubstances()

	for substance, _ := range dUniq {
		// Create companion info for Journal
		if _, ok := jUniq[substance]; !ok {
			companion := JournalCompanion{
				SubstanceName: substance,
				Color:         "CYAN",
			}

			j.JournalCompanions = append(j.JournalCompanions, companion)
		}
	}

	if b, err := json.MarshalIndent(j, "", "    "); err != nil {
		panic(err)
	} else {
		fmt.Printf("%s\n", b)
	}
}

func getJsonFromUrl(v any, path string) error {
	response, err := http.Get(path)
	if err != nil {
		fmt.Printf("failed to read json: %v\n", err)
		return err
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read body: %v\n", err)
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		fmt.Printf("failed to unmarshal doses: \n%s\n%v\n", b, err)
		return err
	}

	return nil
}
