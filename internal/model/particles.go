package model

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/programmersd21/pyproject-tui/internal/theme"
)

type sandGrain struct {
	char    rune
	gx, gy  int
	x, y    float64
	vx, vy  float64
	phase   float64
	freq    float64
	amp     float64
	settled bool
	color   string
}

type sandField struct {
	targets        []*sandGrain
	active         []*sandGrain
	nextIdx        int
	completed      bool
	colors         []string
	rng            *rand.Rand
	seed           int64
	gradientOffset int
}

func newSandField(_, _ int) *sandField {
	return &sandField{seed: rand.Int63()}
}

func (sf *sandField) setSize(_, _ int) {}

func (sf *sandField) refreshColors() {
	if t := theme.Active(); t != nil && len(t.PreviewColors) > 0 {
		sf.colors = t.PreviewColors
	}
	sf.reset()
}

func (sf *sandField) reset() {
	sf.nextIdx = 0
	sf.active = nil
	sf.completed = false
	sf.gradientOffset = 0
	for _, g := range sf.targets {
		g.settled = false
		g.y = -float64(sf.rng.Intn(8) + 6)
		g.vy = 0.6 + sf.rng.Float64()*0.6
		g.vx = (sf.rng.Float64() - 0.5) * 0.3
		g.phase = 0
		g.freq = 0.03 + sf.rng.Float64()*0.04
		g.amp = 0.3 + sf.rng.Float64()*0.5
		g.color = sf.gradientColor(g.gx, g.gy, 0)
	}
}

func (sf *sandField) gradientColor(gx, gy, offset int) string {
	ci := len(sf.colors)
	if ci == 0 {
		return "#FFFFFF"
	}
	t := (math.Sin(float64(gx)*0.3+float64(offset)*0.12+float64(gy)*0.2) + 1.0) / 2.0
	idx := int(t * float64(ci))
	if idx >= ci {
		idx = ci - 1
	}
	return sf.colors[idx]
}

func (sf *sandField) initFromLogo(logoLines []string) {
	sf.rng = rand.New(rand.NewSource(sf.seed))
	sf.colors = []string{"#88C0D0", "#8FBCBB", "#A3BE8C"}
	if t := theme.Active(); t != nil && len(t.PreviewColors) > 0 {
		sf.colors = t.PreviewColors
	}
	sf.targets = nil
	for gy, line := range logoLines {
		for gx, ch := range []rune(line) {
			if ch == ' ' {
				continue
			}
			sf.targets = append(sf.targets, &sandGrain{
				char: ch,
				gx:   gx,
				gy:   gy,
			})
		}
	}
	sf.reset()
}

func (sf *sandField) tick() {
	sf.gradientOffset++

	if sf.completed || len(sf.targets) == 0 {
		return
	}

	releaseCount := 4 + sf.rng.Intn(12)
	for i := 0; i < releaseCount && sf.nextIdx < len(sf.targets); i++ {
		g := sf.targets[sf.nextIdx]
		g.x = float64(g.gx)
		g.y = -float64(sf.rng.Intn(6) + 4)
		g.vy = 0.6 + sf.rng.Float64()*0.6
		g.vx = (sf.rng.Float64() - 0.5) * 0.3
		g.phase = 0
		g.freq = 0.03 + sf.rng.Float64()*0.04
		g.amp = 0.3 + sf.rng.Float64()*0.5
		g.color = sf.gradientColor(g.gx, g.gy, 0)
		sf.active = append(sf.active, g)
		sf.nextIdx++
	}

	allSettled := true
	for _, g := range sf.active {
		if g.settled {
			continue
		}
		allSettled = false
		g.phase += g.freq
		sway := math.Sin(g.phase) * g.amp
		g.vy += 0.01
		if g.vy > 1.2 {
			g.vy = 1.2
		}

		// Ease out as grain approaches target
		dist := float64(g.gy) - g.y
		if dist < 2.0 && dist > 0 {
			g.vy *= 0.90
			if g.vy < 0.08 {
				g.vy = 0.08
			}
		}

		g.vx *= 0.995
		g.x = float64(g.gx) + sway + g.vx
		g.y += g.vy

		if g.y >= float64(g.gy) {
			g.settled = true
			g.x = float64(g.gx)
			g.y = float64(g.gy)
		}
	}
	if allSettled && sf.nextIdx >= len(sf.targets) {
		sf.active = nil
		sf.completed = true
	}
}

func (sf *sandField) render(logoLines []string) string {
	if len(sf.targets) == 0 {
		sf.initFromLogo(logoLines)
	}
	sf.tick()

	if sf.completed {
		extraTop := 5
		logoW := 0
		for _, l := range logoLines {
			if w := len([]rune(l)); w > logoW {
				logoW = w
			}
		}
		logoH := len(logoLines)
		var sb strings.Builder
		for y := -extraTop; y < logoH; y++ {
			for x := 0; x < logoW; x++ {
				if y >= 0 {
					runes := []rune(logoLines[y])
					if x < len(runes) {
						ch := runes[x]
						if ch == ' ' {
							sb.WriteRune(' ')
						} else {
							style := lipgloss.NewStyle().Foreground(lipgloss.Color(sf.gradientColor(x, y, sf.gradientOffset)))
							sb.WriteString(style.Render(string(ch)))
						}
					} else {
						sb.WriteRune(' ')
					}
				} else {
					sb.WriteRune(' ')
				}
			}
			if y < logoH-1 {
				sb.WriteRune('\n')
			}
		}
		return sb.String()
	}

	settled := make(map[[2]int]*sandGrain)
	for _, g := range sf.targets {
		if g.settled {
			settled[[2]int{g.gx, g.gy}] = g
		}
	}

	logoW := 0
	for _, l := range logoLines {
		if w := len([]rune(l)); w > logoW {
			logoW = w
		}
	}
	logoH := len(logoLines)
	extraTop := 5

	var sb strings.Builder
	for y := -extraTop; y < logoH; y++ {
		for x := 0; x < logoW; x++ {
			var placed bool
			for _, g := range sf.active {
				if g.settled {
					continue
				}
				rx, ry := int(math.Round(g.x)), int(math.Round(g.y))
				if rx == x && ry == y {
					style := lipgloss.NewStyle().Foreground(lipgloss.Color(g.color))
					sb.WriteString(style.Render("·"))
					placed = true
					break
				}
			}
			if placed {
				continue
			}
			if y >= 0 {
				runes := []rune(logoLines[y])
				if x < len(runes) && runes[x] != ' ' {
					if g, ok := settled[[2]int{x, y}]; ok {
						style := lipgloss.NewStyle().Foreground(lipgloss.Color(g.color))
						sb.WriteString(style.Render(string(g.char)))
					} else {
						sb.WriteRune(' ')
					}
				} else {
					sb.WriteRune(' ')
				}
			} else {
				sb.WriteRune(' ')
			}
		}
		if y < logoH-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
