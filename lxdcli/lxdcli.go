package lxdcli
import (
	"github.com/lxc/lxd"
)

type MClient struct {
	Client lxd.Client
}