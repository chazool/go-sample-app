package dto

type HostPost struct {
	Host string `query:"host" validate:"Host"`
	Port int    `query:"port" validate:"Port"`
}
