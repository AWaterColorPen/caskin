package caskin

type baseService struct {
	Enforcer IEnforcer
	DB       MetaDB
}
