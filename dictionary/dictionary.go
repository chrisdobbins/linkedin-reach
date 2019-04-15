package dictionary

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Dict []string

type dictionaryGetter struct{}

type WordCriteria struct {
	MaxUniqueChars int
}

type getter interface {
	get() ([]string, error)
}

func New() (Dict, error) {
	dg := dictionaryGetter{}
	d := &Dict{}
	if err := d.populate(dg); err != nil {
		return *d, fmt.Errorf("unable to populate dictionary: %s", err.Error())
	}
	return *d, nil
}

func (d Dict) GetOne(wc WordCriteria) (string, error) {
	if len(d) == 0 {
		return "", errors.New("no words available")
	}
	wordIdx := rand.Intn(len(d))
	word := d[wordIdx]

	if wc.MaxUniqueChars == 0 {
		return word, nil
	}

	chars := map[rune]struct{}{}
	for _, ch := range word {
		chars[ch] = struct{}{}
	}
	if len(chars) > wc.MaxUniqueChars || len(word) == 0 {
		getOne := Dict.GetOne
		dCopy := []string(d)
		dCopy = append(append([]string{}, d[:wordIdx]...), d[wordIdx+1:]...)
		return getOne(Dict(dCopy), wc)
	}
	return word, nil
}

func (d *Dict) populate(g getter) error {
	newDict, err := g.get()
	if err != nil {
		return fmt.Errorf("unable to populate dictionary: %s", err.Error())
	}
	*d = Dict(newDict)
	return nil
}

func (d dictionaryGetter) get() ([]string, error) {
	errorTemplate := "unable to retrieve word list: %s"
	url := "http://app.linkedin-reach.io/words"
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf(errorTemplate, fmt.Sprintf("error getting words from %s", url))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, err.Error())
	}
	defer resp.Body.Close()
	return strings.Split(string(body), "\n"), nil
}
