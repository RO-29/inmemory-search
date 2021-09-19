package inmemorycache

// trie Implements the trie DS.
type trie struct {
	ChildrenNode [512]*trie
	isTerminal   bool
	Value        interface{}
}

// Set the key:value to trie "cache".
func (t *trie) Set(key string, value interface{}) {
	if len(key) < 1 {
		t.isTerminal = true
		t.Value = value
		return
	}
	index := key[0]
	if t.ChildrenNode[index] == nil {
		t.ChildrenNode[index] = &trie{}
	}
	t.ChildrenNode[index].Set(key[1:], value)
}

// Get retrieve the value for key from trie "cache".
func (t *trie) Get(key string) (interface{}, bool) {
	curr := t
	for i := 0; i < len(key); i++ {
		index := key[i]
		curr = curr.ChildrenNode[index]
		if curr == nil {
			return nil, false
		}
	}
	return curr.Value, curr.isTerminal
}

// PrefixMatch returns the matching *keys* in the trie.
func (t *trie) PrefixMatch(prefix string, total int) []string {
	var result []string
	for i := 0; i < len(prefix); i++ {
		index := prefix[i]
		if t.ChildrenNode[index] == nil {
			return result
		}
		t = t.ChildrenNode[index]
	}
	if t.isTerminal && t.isLastNode() {
		result = append(result, prefix)
		return result
	}
	keys := []string{}
	if t.isTerminal {
		keys = append(keys, prefix)
		if total != -1 {
			total--
		}
	}
	if !t.isLastNode() {
		_, result = t.find(prefix, keys, total)
	}
	return result
}

func (t *trie) isLastNode() bool {
	for i := 0; i < 512; i++ {
		if t.ChildrenNode[i] != nil {
			return false
		}
	}
	return true
}

func (t *trie) find(prefix string, keys []string, repeat int) (int, []string) {
	if t.isLastNode() {
		if t.isTerminal && len(keys) < 1 {
			keys = append(keys, prefix)
		}
		return repeat, keys
	}
	for i := 0; i < 512; i++ {
		if repeat == 0 {
			return repeat, keys
		}
		r := t
		if t.ChildrenNode[i] != nil {
			l := rune(i)
			prefix += string(l)
			r = r.ChildrenNode[i]
			if r.isTerminal {
				keys = append(keys, prefix)
				if repeat != -1 {
					repeat--
				}
			}
			repeat, keys = r.find(prefix, keys, repeat)
			prefix = prefix[0 : len(prefix)-1]
		}
	}
	return repeat, keys
}
