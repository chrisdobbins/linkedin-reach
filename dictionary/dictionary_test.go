package dictionary

import "testing"

const maxUniqueChars = 4

var validWords = map[string]struct{}{"abc": struct{}{},
	"abcd":                   struct{}{},
	"test":                   struct{}{},
	"aaaaaaaabbbbbbbbcccccc": struct{}{}}

var invalidWords = map[string]struct{}{
	"abcdefghi": struct{}{},
	"mkprlni":   struct{}{},
        "xyzwqpolimnf": struct{}{},
        "qwertyuiuop": struct{}{},
        "zdcvxvbnmlkjhgfdsa": struct{}{},
        "yqwhgfdmalkopzxc": struct{}{},
        "vxvbnmlkjhgfd": struct{}{},
}

type mockGetter struct{}

func (mg mockGetter) get() ([]string, error) {
	output := []string{}
	for word, _ := range invalidWords {
		output = append(output, word)
	}
	for word, _ := range validWords {
		output = append(output, word)
	}
	return output, nil
}

func TestGetOne(t *testing.T) {
       wc := WordCriteria{}
       d := &Dict{}
       _, err := d.GetOne(wc) 
       if err == nil {
          t.Logf("expected \"no words available\" error, got %v", err)
         t.Fail()
       }
        
	wc = WordCriteria{
		MaxUniqueChars: maxUniqueChars,
	}
	mg := mockGetter{}
	if err := d.populate(mg); err != nil {
		t.Fatal(err)
	}
	testWord, err := d.GetOne(wc)
	if err != nil {
		t.Logf("expcted nil error, got %s", err.Error())
		t.Fail()
	}
        if len(testWord) == 0 {
            t.Log("no word returned")
            t.Fail()
        }
	if _, ok := invalidWords[testWord]; ok {
		t.Logf("invalid word %s returned from Dict.GetOne", testWord)
		t.Fail()
	}
}
