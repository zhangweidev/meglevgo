package meglev

import (
	"crypto/md5"
	"encoding/binary"
	"sort"
	"strings"
)

var (
	LOOKUPSNUM int = 65537
)

type Node struct {
	Name  string
	Width int
}

type Meglev struct {
	Nodes       []*Node
	lookups     []int
	permutaions [][]int
}

func NewMeglev(nodes []*Node) *Meglev {

	sort.SliceStable(nodes, func(i, j int) bool {
		return strings.Compare(nodes[i].Name, nodes[j].Name) < 0
	})

	meglev := &Meglev{
		Nodes:       nodes,
		lookups:     make([]int, LOOKUPSNUM),
		permutaions: make([][]int, len(nodes)),
	}
	meglev.generatePopulation()
	meglev.lookuptable()
	return meglev
}

func (m *Meglev) Get(key string) string {
	i := hashs("emhhbmd3ZWlkZXY=", key)
	index := m.lookups[i%int64(LOOKUPSNUM)]
	return m.Nodes[index].Name
}

// func (m *Meglev) AddNode(name string, width int) error {
// 	// 将 node 添加到链表中按 name 排序
// 	node := &Node{
// 		Name:  name,
// 		Width: width,
// 	}
// 	m.Nodes = append(m.Nodes, node)
// 	return nil
// }

func hashs(token, name string) int64 {
	key := token + name

	hasher := md5.New()
	hasher.Write([]byte(key))
	encode := hasher.Sum(nil)
	res := int64(binary.BigEndian.Uint64(encode))
	if res < 0 {
		res = -res
	}
	return res
}

// 生成查找表
func (m *Meglev) generatePopulation() {

	if len(m.Nodes) == 0 {
		return
	}

	for i := 0; i < len(m.Nodes); i++ {
		node := m.Nodes[i]

		offset := int(hashs("emhhbmd3ZWlkZXY=", node.Name) % int64(LOOKUPSNUM))
		skip := int(hashs("aGVsbG93b3Jk", node.Name)%int64(LOOKUPSNUM-1) + 1)

		row := make([]int, LOOKUPSNUM)
		for j := 0; j < int(LOOKUPSNUM); j++ {
			row[j] = (offset + i*skip) % (LOOKUPSNUM)
		}
		m.permutaions[i] = row
	}
}

// 生成填充表
func (m *Meglev) lookuptable() {

	if len(m.Nodes) == 0 {
		return
	}
	next := make([]int, len(m.Nodes))
	sum := 0
	for {
		for i := 0; i < len(m.Nodes); i++ {
			for w := 0; w < m.Nodes[i].Width; w++ { // 选择 width 次
				c := m.permutaions[i][next[i]]
				for m.lookups[c] > 0 {
					next[i] = next[i] + 1
					c = m.permutaions[i][next[i]]
				}
				m.lookups[c] = i
				next[i] = next[i] + 1
				sum++
				if sum == len(m.Nodes) {
					return
				}
			}
		}
	}

}
