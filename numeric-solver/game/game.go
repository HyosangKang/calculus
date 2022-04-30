package game

import (
	"fmt"
	"sync"
	"time"
)

type Game struct {
	verbose bool
	n       int
	minions []*minion
}

// nss is the array of minion's info [num, hp, dps]
func NewGame(v bool, nss ...[]int) *Game {
	g := &Game{
		verbose: v,
		n:       len(nss),
	}

	id := 0
	for i, ns := range nss {
		for j := 0; j < ns[0]; j++ {
			m := newMinion(g, i, id, ns[1], ns[2])
			g.minions = append(g.minions, m)
			id++
		}
	}
	return g
}

func (g *Game) Start() {
	var wg sync.WaitGroup
	for _, m := range g.minions {
		m := m
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.deploy()
		}()
	}
	wg.Wait()
}

func (g *Game) Report() {
	r := make([]int, g.n)
	for i := 0; i < g.n; i++ {
		for _, m := range g.minions {
			if m.tag == i && m.hp > 0 {
				r[i]++
			}
		}
	}
	fmt.Println(r)
}

type minion struct {
	game    *Game
	tget    *minion
	tag, id int
	hp, dps int
}

func newMinion(g *Game, tag, id, hp, dps int) *minion {
	return &minion{
		game: g,
		tget: nil,
		id:   id,
		tag:  tag,
		hp:   hp,
		dps:  dps,
	}
}

func (m *minion) deploy() {
	for {
		time.Sleep(1 * time.Millisecond)
		if m.hp <= 0 {
			if m.game.verbose {
				fmt.Printf("minion [%d(%d)] hp is %d\n", m.id, m.tag, m.hp)
			}
			break
		}
		if !m.attack() {
			if m.game.verbose {
				fmt.Printf("minion [%d(%d)] has no target\n", m.id, m.tag)
			}
			break
		}
	}
}

func (m *minion) attack() bool {
	if m.tget == nil {
		m.target()
	} else if m.tget.hp <= 0 {
		m.target()
	}
	if m.tget == nil {
		return false
	}
	m.tget.hp -= m.dps
	if m.game.verbose {
		fmt.Printf("minion [%d(%d)] deals dmg to minion [%d(%d)]\n", m.id, m.tag, m.tget.id, m.tget.tag)
	}
	return true
}

func (m *minion) target() {
	for _, n := range m.game.minions {
		if n == m {
			continue
		}
		if n.hp <= 0 {
			continue
		}
		if n.tag == m.tag {
			continue
		}
		m.tget = n
		return
	}
	m.tget = nil
}
