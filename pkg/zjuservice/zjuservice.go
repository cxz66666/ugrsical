package zjuservice

type ZjuService interface {
	Login(username, password string) error
}
