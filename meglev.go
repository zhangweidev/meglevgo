package meglevgo

import (
	"sort"
	"strings"

	"github.com/dchest/siphash"
)

var (
	LOOKUPSNUM     int = 65537 //
	LOOKUPSNUM_MAX int = 655373
)

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

type Node struct {
	Name     string
	Width    int
	gcdWidth int
}

type Meglev struct {
	Nodes       []*Node
	lookups     []int
	permutaions [][]int

	lookupsnum int
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
	width_sum := 0
	if len(meglev.Nodes) >= 2 {
		g := gcd(meglev.Nodes[0].Width, meglev.Nodes[1].Width)
		for i := 2; i < len(meglev.Nodes); i++ {
			g = gcd(g, meglev.Nodes[i].Width)
		}
		for i := 0; i < len(meglev.Nodes); i++ {
			meglev.Nodes[i].gcdWidth = meglev.Nodes[i].Width / g
			width_sum += meglev.Nodes[i].gcdWidth
		}
	}
	if width_sum < 10000 {
		meglev.lookupsnum = LOOKUPSNUM
	} else {
		meglev.lookupsnum = LOOKUPSNUM_MAX
	}

	meglev.generatePopulation()
	meglev.lookuptable()
	return meglev
}

func (m *Meglev) hashs(token uint64, key string) uint64 {
	return siphash.Hash(token, 0, []byte(key))
}

// 生成哈希
// func hashs(token, name string) int64 {
// 	key := token + name

// 	hasher := md5.New()
// 	hasher.Write([]byte(key))
// 	encode := hasher.Sum(nil)
// 	res := int64(binary.BigEndian.Uint64(encode))
// 	if res < 0 {
// 		res = -res
// 	}
// 	return res
// }

func (m *Meglev) Get(key string) string {
	i := m.hashs(0xabcd1234, key) % uint64(m.lookupsnum)
	index := m.lookups[i]
	return m.Nodes[index-1].Name
}

// 生成查找表
func (m *Meglev) generatePopulation() {

	if len(m.Nodes) == 0 {
		return
	}

	for i := 0; i < len(m.Nodes); i++ {
		node := m.Nodes[i]

		offset := int(m.hashs(0xabcdeeff, node.Name) % uint64(m.lookupsnum))
		skip := int(m.hashs(0xffeedcba, node.Name)%uint64(m.lookupsnum-1) + 1)

		row := make([]int, m.lookupsnum)
		for j := 0; j < int(m.lookupsnum); j++ {
			row[j] = (offset + j*skip) % (m.lookupsnum)
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
			for w := 0; w < m.Nodes[i].gcdWidth; w++ { // 选择 width 次

				if next[i] >= m.lookupsnum {
					continue
				}
				c := m.permutaions[i][next[i]]
				for m.lookups[c] > 0 {
					next[i] = next[i] + 1
					if next[i] >= m.lookupsnum {
						break
					}
					c = m.permutaions[i][next[i]]
				}
				if m.lookups[c] == 0 {
					m.lookups[c] = i + 1
					next[i] = next[i] + 1
					sum++
				}
				if sum == m.lookupsnum {
					return
				}
			}
		}
	}
}
