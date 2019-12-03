package sola

// OriginContext Middleware change ctx to it's origin
func OriginContext(next Handler) Handler {
	return func(c Context) error {
		if origin := c.Origin(); origin != nil {
			return next(origin)
		}
		return next(c)
	}
}
