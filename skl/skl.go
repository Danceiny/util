package skl

import (
	"math"

	"github.com/zhiqiangxu/util"
)

const (
	// DefaultMaxLevel for skl
	DefaultMaxLevel int = 18
	// DefaultProbability for skl
	DefaultProbability float64 = 1 / math.E
)

type link struct {
	next *element
}

type links []link

type element struct {
	links
	key   int64
	value interface{}
}

type skl struct {
	links
	maxLevel       int
	length         int
	probability    float64
	probTable      []uint32
	prevLinksCache []links
}

// NewSkipList creates a new SkipList
func NewSkipList() SkipList {
	return NewSkipListWithMaxLevel(DefaultMaxLevel)
}

// NewSkipListWithMaxLevel creates a new SkipList with specified maxLevel
func NewSkipListWithMaxLevel(maxLevel int) SkipList {
	return &skl{
		links:          make(links, maxLevel),
		maxLevel:       maxLevel,
		probability:    DefaultProbability,
		probTable:      probabilityTable(DefaultProbability, maxLevel),
		prevLinksCache: make([]links, maxLevel),
	}
}

func probabilityTable(probability float64, maxLevel int) (table []uint32) {
	for i := 0; i < maxLevel; i++ {
		prob := math.Pow(probability, float64(i))
		table = append(table, uint32(prob*math.MaxUint32))
	}
	return table
}

func (s *skl) Add(key int64, value interface{}) {
	prevs := s.getPrevLinks(key)
	ele := prevs[0][0].next
	if ele != nil && ele.key <= key {
		ele.value = value
		return
	}

	ele = &element{
		links: make(links, s.randLevel()),
		key:   key,
		value: value,
	}

	for i := range ele.links {
		ele.links[i].next = prevs[i][i].next
		prevs[i][i].next = ele
	}

	s.length++
}

func (s *skl) randLevel() (level int) {

	r := util.FastRand()

	level = 1
	for level < s.maxLevel && r < s.probTable[level] {
		level++
	}
	return
}

// 找到每一层上毗邻于该key对应元素之前的links
func (s *skl) getPrevLinks(key int64) []links {
	var prev = s.links
	var current *element

	prevs := s.prevLinksCache
	for i := s.maxLevel - 1; i >= 0; i-- {
		current = prev[i].next

		for current != nil && current.key < key {
			prev = current.links
			current = current.links[i].next
		}

		prevs[i] = prev
	}

	return prevs
}

func (s *skl) Get(key int64) (value interface{}, ok bool) {
	prev := s.links
	var current *element
	for i := s.maxLevel - 1; i >= 0; i-- {
		current = prev[i].next
		for current != nil && current.key < key {
			prev = current.links
			current = current.links[i].next
		}
	}

	if current != nil && current.key <= key {
		return current.value, true
	}

	return nil, false
}

func (s *skl) Remove(key int64) {
	prevs := s.getPrevLinks(key)
	if ele := prevs[0][0].next; ele != nil && ele.key <= key {

		for i, l := range ele.links {
			prevs[i][i].next = l.next
		}
		s.length--
	}
}

func (s *skl) Head() (key int64, value interface{}, ok bool) {
	nele := s.links[0].next
	if nele != nil {
		key = nele.key
		value = nele.value
		ok = true
	}
	return
}

func (s *skl) NewIterator() SkipListIterator {
	return &sklIter{s: s}
}