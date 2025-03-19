package build

import "user/internal/jwt"

func (b *Builder) jwtService() *jwt.Service {
	return jwt.NewService(b.config.JWT.Secret)
}
