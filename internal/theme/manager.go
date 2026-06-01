package theme

import "sync"

type ThemeManager struct {
	mu      sync.RWMutex
	current Scheme
	listener func(Palette)
}

var (
	global *ThemeManager
	once   sync.Once
)

func Global() *ThemeManager {
	once.Do(func() {
		global = &ThemeManager{
			current: Schemes[0],
		}
	})
	return global
}

func (tm *ThemeManager) Set(name string) {
	for _, s := range Schemes {
		if s.Name == name {
			tm.mu.Lock()
			tm.current = s
			tm.mu.Unlock()

			if tm.listener != nil {
				tm.listener(s.Palette())
			}
			return
		}
	}
}

func (tm *ThemeManager) Get() Scheme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current
}

func (tm *ThemeManager) Palette() Palette {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current.Palette()
}

func (tm *ThemeManager) OnChange(fn func(Palette)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.listener = fn
}

func (tm *ThemeManager) IsDark() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current.IsDark
}

func ThemeNames() []string {
	names := make([]string, len(Schemes))
	for i, s := range Schemes {
		names[i] = s.Name
	}
	return names
}

func ThemeDisplayNames() []string {
	names := make([]string, len(Schemes))
	for i, s := range Schemes {
		names[i] = s.DisplayName
	}
	return names
}
