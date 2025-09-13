package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/5HT2C/doses-logger-pw-journal/types"
)

var (
	fp      = flag.String("f", "", "file path")
	loadUrl = flag.String("url", "http://localhost:6010/media/doses.json", "URL for doses.json")
)

func main() {
	flag.Parse()

	var doses []types.Dose
	if err := getJsonFromUrl(&doses, *loadUrl); err != nil {
		panic(err)
	}

	var j types.Journal
	if len(*fp) > 0 {
		if b, err := os.ReadFile(*fp); err != nil {
			panic(err)
		} else {

			if err := json.Unmarshal(b, &j); err != nil {
				j = types.Journal{}
			}
		}
	}

	d := types.Doses(doses)
	j.Experiences = append(
		j.Experiences,
		d.ToExperienceGroups()...,
	)
	j.JournalCompanions = append(
		j.JournalCompanions,
		d.ToJournalCompanions(j.UniqueSubstances())...,
	)

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
