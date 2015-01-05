package main

type StringSet struct {
	set map[string]struct{}
}

func initStringSet(set *StringSet) {
	if set.set == nil {
		set.set = make(map[string]struct{})
	}
}

func (set *StringSet) Add(value string) {
	initStringSet(set)
	set.set[value] = struct{}{}
}

func (set *StringSet) AddAll(values []string) {
	initStringSet(set)
	for _, value := range values {
		set.set[value] = struct{}{}
	}
}

func (set *StringSet) Contains(value string) bool {
	initStringSet(set)
	_, ok := set.set[value]
	return ok
}
