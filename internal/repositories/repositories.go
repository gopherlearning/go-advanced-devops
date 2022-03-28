package repositories

type Repository interface {
	Update(target, metric string) error
	List(targets ...string) map[string][]string
	ListProm(targets ...string) []byte
}
