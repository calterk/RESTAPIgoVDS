package storage

import (
	"app/entity"
	"sync"
)

type VdsStorage struct {
	vdsMap map[uint64]entity.Vds
	next   uint64
	mutx   sync.RWMutex
}

var s *VdsStorage

func init() {
	s = new(VdsStorage)
}

func CreateVds(vds entity.Vds) entity.Vds {
	s.mutx.Lock()
	defer s.mutx.Unlock()
	vds.Vid = s.next
	s.next++
	s.vdsMap[vds.Vid] = vds
	return vds
}
