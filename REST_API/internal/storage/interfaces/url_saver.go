//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver --output=../../../mocks


package interfaces

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

