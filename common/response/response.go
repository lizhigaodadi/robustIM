package response

/*TODO:Generic http request corresponding package*/

type ImResponse struct {
	body map[string]interface{}
}

func Response() *ImResponse {
	resp := &ImResponse{
		body: make(map[string]interface{}),
	}
	return resp
}

func (resp *ImResponse) Err() *ImResponse {
	resp.body["message"] = "err"
	resp.body["code"] = 400
	return resp
}

func (resp *ImResponse) Ok() *ImResponse {
	resp.body["message"] = "ok"
	resp.body["code"] = 200
	return resp
}

func (resp *ImResponse) Put(key string, val interface{}) *ImResponse {
	resp.body[key] = val
	return resp
}

func (resp *ImResponse) End() map[string]interface{} {
	return resp.body
}
