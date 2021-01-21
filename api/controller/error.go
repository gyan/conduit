package controller

import (
	"net/http"

	"github.com/arpit32/conduit/api/constant"
	"github.com/arpit32/conduit/api/constant/codes"
	"github.com/arpit32/conduit/api/errors"
)

//HTTPErrorController ...
type HTTPErrorController struct {
	BaseController
}

//ResourceNotFound ...
func (c *HTTPErrorController) ResourceNotFound(w http.ResponseWriter, r *http.Request) {
	err := errors.New(codes.NotFound, constant.ResourceNotFound)
	c.WriteError(r, w, err)
}
