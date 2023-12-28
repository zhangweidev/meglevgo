package meglevgo

import (
	"strconv"
	"testing"
)

func Test_hashs(t *testing.T) {
	m := NewMeglev([]*Node{
		{Name: "am", Width: 1},
		{Name: "es", Width: 1},
		{Name: "c", Width: 1},
		{Name: "d", Width: 1},
	})
	res := m.hashs(0x12345678, "nam11e")
	t.Log(res)

}

func Test_Population(t *testing.T) {
	meglev := NewMeglev([]*Node{
		{Name: "am", Width: 1},
		{Name: "es", Width: 1},
		{Name: "c", Width: 1},
		{Name: "d", Width: 2},
	})
	for i, p := range meglev.permutaions {
		t.Log(i, p)
	}
	t.Log(meglev.lookups)

}

func Test_Meglev(t *testing.T) {

	meglev := NewMeglev([]*Node{
		{Name: "am", Width: 1},
		{Name: "es", Width: 1},
		// {Name: "c", Width: 1},
		// {Name: "d", Width: 2},
	})

	m := make(map[string]int)
	for i := 0; i < 10000000; i++ {
		name := meglev.Get(strconv.Itoa(i))
		if val, ok := m[name]; ok {
			m[name] = val + 1
		} else {
			m[name] = 1
		}
	}
	t.Log(m)

}

func Test_gcd(t *testing.T) {

	t.Log(gcd(3, 5))
	t.Log(gcd(5, 10))
	t.Log(gcd(10, 15))
	t.Log(gcd(10, gcd(100, gcd(90, 5))))

	t.Log(gcd(-1, 2))
}

func Test_meglev_gcd(t *testing.T) {
	meglev := NewMeglev([]*Node{
		{Name: "am", Width: 100},
		{Name: "es", Width: 200},
		{Name: "c", Width: 10},
		{Name: "d", Width: 30},
		{Name: "ea", Width: 50},
	})

	for _, node := range meglev.Nodes {
		t.Log(node.Name, node.gcdWidth)
	}

}
