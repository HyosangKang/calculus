package game

import (
	"fmt"
	"time"
)

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
			m.kia()
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

func (m *minion) kia() {
	k := kia{
		time: int(time.Since(m.game.start).Milliseconds()),
		tag:  m.tag,
		id:   m.id,
	}
	m.game.report = append(m.game.report, k)
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
