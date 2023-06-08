package types

type ImgWOptions struct {
	Engine    string
	Namespace string
	Address   string

	//Registry Address (docker hub or redii.net)
	ServerAddress string
	UserId        string
	UserPasswd    string
}
