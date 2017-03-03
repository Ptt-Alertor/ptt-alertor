package article

type Article interface {
	ContainKeyword(keyword string) bool
}
