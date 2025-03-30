package dfa

import (
	"strings"
	"sync"
)

const (
	// default Invalid Words, e.g: some bad w..o.r,,d
	defaultInvalidWords = " ,~,!,@,#,$,%,^,&,*,(,),_,-,+,=,?,<,>,.,—,，,。,/,\\,|,《,》,？,;,:,：,',‘,；,“,¥,·"
	// default ReplaceStr
	defaultReplaceStr = "****"
)

// DFA
type DFA struct {
	l            sync.Mutex
	trie         *Trie
	replaceStr   string
	invalidWords map[string]struct{}
}

// New
func New() *DFA {
	f := &DFA{
		trie:         NewTrie(),
		replaceStr:   defaultReplaceStr,
		invalidWords: make(map[string]struct{}),
	}
	for _, s := range defaultInvalidWords {
		f.invalidWords[string(s)] = struct{}{}
	}
	return f
}

// AddBadWords
func (f *DFA) AddBadWords(words []string) {
	f.l.Lock()
	defer f.l.Unlock()
	if len(words) > 0 {
		for _, s := range words {
			f.trie.Insert(s)
		}
	}
}

// SetInvalidChar
func (f *DFA) SetInvalidChar(chars string) {
	f.l.Lock()
	defer f.l.Unlock()
	f.invalidWords = make(map[string]struct{})
	for _, s := range chars {
		f.invalidWords[string(s)] = struct{}{}
	}
}

// SetReplaceStr
func (f *DFA) SetReplaceStr(str string) {
	f.l.Lock()
	defer f.l.Unlock()

	f.replaceStr = str
}

// Check the strings
func (f *DFA) Check(txt string) ([]string, []string, bool) {
	_, found, target, b := f.check(txt, false)
	return found, target, b
}

// CheckAndReplace
func (f *DFA) CheckAndReplace(txt string) (string, []string, []string, bool) {
	return f.check(txt, true)
}

// FilterInvalidChar
func (f *DFA) FilterInvalidChar(txt ...string) []string {
	res := make([]string, 0, len(txt))
	for _, s := range txt {
		str := []rune(s)
		for i, c := range str {
			if _, ok := f.invalidWords[string(c)]; ok {
				str = append(str[:i], str[i+1:]...)
			}
		}
		res = append(res, string(str))
	}
	return res
}

// check and replace the strings
func (f *DFA) check(txt string, replace bool) (dist string, found []string, target []string, b bool) {
	var (
		str        = []rune(txt)
		ok         bool
		node       *Node
		nodeMap    map[rune]*Node
		start, tag = -1, -1
		result     string
		tmp        = ""
	)
	target = make([]string, 0, 0)
	f.l.Lock()
	defer f.l.Unlock()

	for i, val := range str {
		if _, ok = f.invalidWords[string(val)]; ok {
			continue
		}

		if nodeMap == nil {
			node = f.trie.Child(string(val))
			if node != nil {
				tag++
				if tag == 0 {
					start = i
				}
				tmp = node.Value
				if !node.IsEnd {
					nodeMap = node.Child
				} else {
					target = append(target, tmp)
					tmp = ""
					found = append(found, string(str[start:i+1]))
					if replace {
						result = strings.Replace(result, string(str[start:i+1]), f.replaceStr, 1)
						if result == "" {
							result = strings.Replace(txt, string(str[start:i+1]), f.replaceStr, 1)
						}
					}
					tag = -1
					start = -1
					nodeMap = nil
				}
			} else {
				nodeMap = nil
				start = -1
				tag = -1
			}
		} else {
			if node, ok = nodeMap[val]; ok {
				tmp += node.Value
				if !node.IsEnd {
					nodeMap = node.Child
				} else {
					target = append(target, tmp)
					tmp = ""
					found = append(found, string(str[start:i+1]))
					if replace {
						result = strings.Replace(result, string(str[start:i+1]), f.replaceStr, 1)
						if result == "" {
							result = strings.Replace(txt, string(str[start:i+1]), f.replaceStr, 1)
						}
					}
					tag = -1
					start = -1
					nodeMap = nil
				}
			}
		}
	}
	b = len(found) > 0
	return
}
