package chat

import (
	"net"

	"github.com/zeusito/gochat/internal/models"
)

type IService interface {
	MemberJoin(claims models.MyClaims, conn net.Conn)
}
