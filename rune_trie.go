package trie

// runeTrie is a rune-wise trie implementation.
type runeTrie struct {
	root *runeNode
	size int
}

type runeNode struct {
	children map[rune]*runeNode
	last     bool
}

// NewRuneTrie creates a rune-wise new trie.
func NewRuneTrie() Trie {
	return &runeTrie{root: &runeNode{make(map[rune]*runeNode), false}}
}

// Insert inserts a word into the trie and returns true
// if the word already exists.
func (t *runeTrie) Insert(word string) bool {
	exists := true
	current := t.root
	for _, letter := range word {
		n, ok := current.children[letter]
		if !ok {
			exists = false

			n = &runeNode{make(map[rune]*runeNode), false}
			current.children[letter] = n
		}

		current = n
	}
	current.last = true

	if !exists {
		t.size++
	}

	return exists
}

// Contains returns true if an exact match of the word is found, otherwise false.
func (t *runeTrie) Contains(word string) bool {
	n, _ := t.nodeByPrefix(word)

	return n != nil && n.last
}

// Size returns the number of keys in the tree.
func (t *runeTrie) Size() int {
	return t.size
}

// StartsWith returns true if there is any word in the trie that starts
// with the given prefix.
func (t *runeTrie) StartsWith(prefix string) bool {
	node, _ := t.nodeByPrefix(prefix)

	return node != nil
}

// Finds and returns words by prefix.
func (t *runeTrie) SearchByPrefix(prefix string) []string {
	node, r := t.nodeByPrefix(prefix)

	return search(node, r, []rune(prefix[:len(prefix)-1]))
}

func (t *runeTrie) nodeByPrefix(prefix string) (*runeNode, rune) {
	current := t.root
	var r rune
	for _, letter := range prefix {
		n, ok := current.children[letter]
		if !ok {
			return nil, 0
		}

		current = n
		r = letter
	}

	return current, r
}

// Finds the longest unique suffix of a given word.
// It iterates each character of the word and tries to find
// if a character is shared with other words. The suffix
// of the given word that is not shared with other words is returned.
func (t *runeTrie) FindLongestUniqueSuffix(word string) (string, bool) {
	currentSuffix := ""
	current := t.root

	for i := 0; i < len(word)-1; i++ {
		letter := rune(word[i])
		n, ok := current.children[letter]
		if !ok {
			return "", false
		}

		if len(n.children) > 1 || n.last {
			currentSuffix = ""
		} else {
			currentSuffix += string(letter)
		}
		current = n
	}

	lastLetter := rune(word[len(word)-1])
	current, ok := current.children[lastLetter]
	if !ok {
		return "", false
	}
	if !current.last {
		return "", false
	}
	if len(current.children) > 0 {
		return "", true
	}

	return currentSuffix + string(lastLetter), true
}

func search(currentNode *runeNode, currentRune rune, prefix []rune) []string {
	words := []string{}
	if currentNode == nil {
		return words
	}

	newPrefix := append(prefix, currentRune)
	if currentNode.last {
		words = append(words, string(newPrefix))
	}

	for letter, node := range currentNode.children {
		newWords := search(node, letter, newPrefix)
		words = append(words, newWords...)
	}

	return words
}
