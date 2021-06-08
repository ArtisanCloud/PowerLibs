package http

import "github.com/ArtisanCloud/go-libs/http/contract"

type HttpResponse struct {

}

func (response *HttpResponse) CastResponseToType(res *contract.ResponseContract, responseType string) interface{} {

	return nil
}
