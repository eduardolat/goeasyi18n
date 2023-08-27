package goeasyi18n

type TranslateString struct {
	Key     string
	Default string

	// For pluralization
	Zero string // Optional
	One  string // Optional
	Two  string // Optional
	Few  string // Optional
	Many string // Optional

	// For genders
	Male      string // Optional
	Female    string // Optional
	NonBinary string // Optional

	// For pluralization with male gender
	ZeroMale string // Optional
	OneMale  string // Optional
	TwoMale  string // Optional
	FewMale  string // Optional
	ManyMale string // Optional

	// For pluralization with female gender
	ZeroFemale string // Optional
	OneFemale  string // Optional
	TwoFemale  string // Optional
	FewFemale  string // Optional
	ManyFemale string // Optional

	// For pluralization with non binary gender
	ZeroNonBinary string // Optional
	OneNonBinary  string // Optional
	TwoNonBinary  string // Optional
	FewNonBinary  string // Optional
	ManyNonBinary string // Optional
}

type TranslateStrings []TranslateString

type PluralizationFunc func(count int) string

type Data map[string]any
