package voyager

//go:generate counterfeiter . Source

type Source interface {
	AssetNames() []string
	Asset(name string) ([]byte, error)
}
