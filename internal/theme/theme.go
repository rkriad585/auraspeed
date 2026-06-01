package theme

import "github.com/gdamore/tcell/v2"

type Scheme struct {
	Name        string
	DisplayName string
	IsDark      bool
	Colors      []string
}

type Palette struct {
	Background tcell.Color
	Surface    tcell.Color
	Primary    tcell.Color
	Secondary  tcell.Color
	Accent     tcell.Color
	Foreground tcell.Color
	Success    tcell.Color
	Warning    tcell.Color
	Error      tcell.Color
	Header     tcell.Color
	Muted      tcell.Color
	Border     tcell.Color

	BackgroundHex string
	SurfaceHex    string
	PrimaryHex    string
	SecondaryHex  string
	AccentHex     string
	ForegroundHex string
	SuccessHex    string
	WarningHex    string
	ErrorHex      string
	HeaderHex     string
	MutedHex      string
	BorderHex     string
}

var Schemes = []Scheme{
	{
		Name: "dark", DisplayName: "Dark Theme", IsDark: true,
		Colors: []string{"#0f172a", "#111827", "#1e293b", "#334155", "#64748b", "#e2e8f0"},
	},
	{
		Name: "light", DisplayName: "Light Theme", IsDark: false,
		Colors: []string{"#ffffff", "#f8fafc", "#e2e8f0", "#cbd5e1", "#475569", "#0f172a"},
	},
	{
		Name: "sunny-beach", DisplayName: "Sunny Beach Day", IsDark: true,
		Colors: []string{"#264653", "#2a9d8f", "#e9c46a", "#f4a261", "#e76f51"},
	},
	{
		Name: "olive-garden", DisplayName: "Olive Garden Feast", IsDark: true,
		Colors: []string{"#606c38", "#283618", "#fefae0", "#dda15e", "#bc6c25"},
	},
	{
		Name: "ocean-breeze", DisplayName: "Summer Ocean Breeze", IsDark: true,
		Colors: []string{"#e63946", "#f1faee", "#a8dadc", "#457b9d", "#1d3557"},
	},
	{
		Name: "summer-fun", DisplayName: "Refreshing Summer Fun", IsDark: true,
		Colors: []string{"#8ecae6", "#219ebc", "#023047", "#ffb703", "#fb8500"},
	},
	{
		Name: "black-gold", DisplayName: "Black & Gold Elegance", IsDark: true,
		Colors: []string{"#000000", "#14213d", "#fca311", "#e5e5e5", "#ffffff"},
	},
	{
		Name: "vibrant-fiesta", DisplayName: "Vibrant Color Fiesta", IsDark: true,
		Colors: []string{"#ffbe0b", "#fb5607", "#ff006e", "#8338ec", "#3a86ff"},
	},
	{
		Name: "light-steel", DisplayName: "Light Steel", IsDark: false,
		Colors: []string{"#f8f9fa", "#e9ecef", "#dee2e6", "#ced4da", "#adb5bd", "#6c757d", "#495057", "#343a40", "#212529"},
	},
	{
		Name: "golden-twilight", DisplayName: "Golden Twilight", IsDark: true,
		Colors: []string{"#000814", "#001d3d", "#003566", "#ffc300", "#ffd60a"},
	},
	{
		Name: "deep-sea", DisplayName: "Deep Sea", IsDark: true,
		Colors: []string{"#0d1b2a", "#1b263b", "#415a77", "#778da9", "#e0e1dd"},
	},
	{
		Name: "bright-green", DisplayName: "Bright Green", IsDark: true,
		Colors: []string{"#004b23", "#006400", "#007200", "#008000", "#38b000", "#70e000", "#9ef01a", "#ccff33"},
	},
	{
		Name: "vivid-nightfall", DisplayName: "Vivid Nightfall", IsDark: true,
		Colors: []string{"#10002b", "#240046", "#3c096c", "#5a189a", "#7b2cbf", "#9d4edd", "#c77dff", "#e0aaff"},
	},
}

func hexToColor(hex string) tcell.Color {
	return tcell.GetColor(hex)
}

func MustParse(name string) Palette {
	for _, s := range Schemes {
		if s.Name == name {
			return s.Palette()
		}
	}
	return Schemes[0].Palette()
}

func (s Scheme) Palette() Palette {
	c := s.Colors
	p := Palette{}

	get := func(idx int) string {
		if idx < len(c) {
			return c[idx]
		}
		return c[len(c)-1]
	}

	p.BackgroundHex = get(0)
	p.SurfaceHex = get(1)
	p.PrimaryHex = get(2)
	p.SecondaryHex = get(3)
	p.AccentHex = get(4)
	p.ForegroundHex = get(len(c) - 1)

	if len(c) >= 6 {
		if s.IsDark {
			p.SuccessHex = c[4]
			p.WarningHex = c[3]
			p.ErrorHex = c[2]
		} else {
			p.SuccessHex = c[3]
			p.WarningHex = c[4]
			p.ErrorHex = c[len(c)-2]
		}
	} else {
		p.SuccessHex = p.SecondaryHex
		p.WarningHex = p.AccentHex
		p.ErrorHex = p.PrimaryHex
	}

	p.Background = hexToColor(p.BackgroundHex)
	p.Surface = hexToColor(p.SurfaceHex)
	p.Primary = hexToColor(p.PrimaryHex)
	p.Secondary = hexToColor(p.SecondaryHex)
	p.Accent = hexToColor(p.AccentHex)
	p.Foreground = hexToColor(p.ForegroundHex)
	p.Success = hexToColor(p.SuccessHex)
	p.Warning = hexToColor(p.WarningHex)
	p.Error = hexToColor(p.ErrorHex)

	p.HeaderHex = p.PrimaryHex
	p.MutedHex = p.ForegroundHex
	p.BorderHex = p.SecondaryHex
	p.Header = p.Primary
	p.Muted = p.Foreground
	p.Border = p.Secondary

	return p
}
