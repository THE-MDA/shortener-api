//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter --output=../../../mocks

package interfaces

type URLGetter interface {
	GetURL(alias string) (string, error)
	GetAllURLs() (map[string]string, error)
}
