package context

func (l *LuxContext) SetStatus(status int) {
	l.Response.StatusCode = status
}

func (l *LuxContext) SetOK() {
	l.Response.StatusCode = 200
}

func (l *LuxContext) SetAccepted() {
	l.Response.StatusCode = 202
}

func (l *LuxContext) SetNoContent() {
	l.Response.StatusCode = 204
}

func (l *LuxContext) SetResetContent() {
	l.Response.StatusCode = 205
}

func (l *LuxContext) SetFound() {
	l.Response.StatusCode = 302
}

func (l *LuxContext) SetBadRequest() {
	l.Response.StatusCode = 400
}

func (l *LuxContext) SetUnauthorized() {
	l.Response.StatusCode = 401
}

func (l *LuxContext) SetForbidden() {
	l.Response.StatusCode = 403
}

func (l *LuxContext) SetNotFound() {
	l.Response.StatusCode = 404
}

func (l *LuxContext) SetInternalServerError() {
	l.Response.StatusCode = 500
}

func (l *LuxContext) SetNotImplemented() {
	l.Response.StatusCode = 501
}

func (l *LuxContext) SetServiceUnavailable() {
	l.Response.StatusCode = 503
}

func (l *LuxContext) SetConflict() {
	l.Response.StatusCode = 409
}

func (l *LuxContext) SetUnsupportedMediaType() {
	l.Response.StatusCode = 415
}
