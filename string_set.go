package main

type StringSet map[string]byte

func NewStringSet() StringSet {
	return StringSet(make(map[string]byte))
}

func (ss StringSet) Add(s string) bool {
	if _, ok := ss[s]; ok {
		return false
	} else {
		ss[s] = 1
		return true
	}
}

func (ss StringSet) Remove(s string) bool {
	if _, ok := ss[s]; ok {
		delete(ss, s)
		return true
	} else {
		return false
	}
}

func (ss StringSet) Contains(s string) bool {
	if _, ok := ss[s]; ok {
		return true
	} else {
		return false
	}
}

func (ss StringSet) Each(f func(string)) {
	for s, _ := range ss {
		f(s)
	}
}
func (ss StringSet) EachUntil(f func(string) bool) {
	for s, _ := range ss {
		stop := f(s)
		if stop {
			break
		}
	}
}
