package dictionary

import "testing"

type mockGetter struct{}

func (mg mockGetter) get() ([]string, error) {
	return []string{"abcd", "test", "aaaaaaaabbbbbbbbcccccc", "abcdefghi", "mkprlni", "xyzwqpolimnf", "qwertyuiuop", "zdcvxvbnmlkjhgfdsa", "yqwhgfdmalkopzxc", "vxvbnmlkjhgfd"}, nil
}

func TestGetOneEmptyDictionary(t *testing.T) {
	wc := WordCriteria{}
	d := &Dict{}
	w, err := d.GetOne(wc)
	if err == nil {
		t.Errorf("err is %v", err)
	}
	if len(w) > 0 {
		t.Errorf("word is %s", w)
	}
}

func TestGetOneMaxUniqueChars(t *testing.T) {
	invalidWords := map[string]struct{}{
		"abcdefghi":          struct{}{},
		"mkprlni":            struct{}{},
		"xyzwqpolimnf":       struct{}{},
		"qwertyuiuop":        struct{}{},
		"zdcvxvbnmlkjhgfdsa": struct{}{},
		"yqwhgfdmalkopzxc":   struct{}{},
		"vxvbnmlkjhgfd":      struct{}{},
	}
	d := &Dict{}
	maxUniqueChars := 4
	wc := WordCriteria{
		MaxUniqueChars: maxUniqueChars,
	}
	mg := mockGetter{}
	if err := d.populate(mg); err != nil {
		t.Fatal(err)
	}
	testWord, err := d.GetOne(wc)
	if err != nil {
		t.Errorf("expcted nil error, got %s", err.Error())
	}
	if len(testWord) == 0 {
		t.Error("no word returned")
	}
	if _, ok := invalidWords[testWord]; ok {
		t.Errorf("invalid word %s returned from Dict.GetOne", testWord)
	}
}

func TestGetOneHappyPath(t *testing.T) {
	words := map[string]struct{}{"abc": struct{}{},
		"abcd":                   struct{}{},
		"test":                   struct{}{},
		"aaaaaaaabbbbbbbbcccccc": struct{}{},
		"abcdefghi":              struct{}{},
		"mkprlni":                struct{}{},
		"xyzwqpolimnf":           struct{}{},
		"qwertyuiuop":            struct{}{},
		"zdcvxvbnmlkjhgfdsa":     struct{}{},
		"yqwhgfdmalkopzxc":       struct{}{},
		"vxvbnmlkjhgfd":          struct{}{},
	}
	d := &Dict{}
	mg := mockGetter{}
	wc := WordCriteria{}
	if err := d.populate(mg); err != nil {
		t.Fatal(err)
	}
	w, err := d.GetOne(wc)
	if err != nil {
		t.Errorf("d.GetOne returned error  %s", err.Error())
	}
	if _, ok := words[w]; !ok {
		t.Errorf("d.GetOne(wc) returned word %s", w)
	}
}
