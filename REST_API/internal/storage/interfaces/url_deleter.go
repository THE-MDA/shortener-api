package interfaces

type URLDeleter interface {
	DeleteURL(alias string) error
}
