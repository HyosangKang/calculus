package game

import (
	"sort"
	"sync"
	"time"
)

type Game struct {
	verbose bool
	n       int
	tagn    int
	minions []*minion
	report  []kia
	start   time.Time
}

// nss is the array of minion's info [num, hp, dps]
func NewGame(v bool, nss ...[]int) *Game {
	g := &Game{
		verbose: v,
		n:       len(nss),
		tagn:    len(nss),
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
	g.start = time.Now()
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

func (g *Game) Report(t float64) ([][]float64, [][]float64) {
	sort.Slice(g.report, func(i, j int) bool {
		return g.report[i].time < g.report[j].time
	})
	duration := float64(g.report[len(g.report)-1].time)
	count := make([]float64, g.tagn)
	for i := 0; i < g.tagn; i++ {
		for _, m := range g.minions {
			if m.tag == i {
				count[i]++
			}
		}
	}

	xss := [][]float64{}
	yss := [][]float64{}
	for i := 0; i < g.tagn; i++ {
		xs := []float64{0}
		ys := []float64{count[i]}
		for _, k := range g.report {
			if k.tag == i {
				xs = append(xs, float64(k.time)/duration*t)
				count[i]--
				ys = append(ys, count[i])
			}
		}
		xss = append(xss, xs)
		yss = append(yss, ys)
	}
	return xss, yss
}
