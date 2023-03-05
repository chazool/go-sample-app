package dto

type HostPost struct {
	Host string `query:"host" validate:"Host"`
	Post int    `query:"port" validate:"Port"`
}
