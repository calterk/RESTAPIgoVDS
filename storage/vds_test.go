package storage

import (
	"app/entity"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupVds() {
	s = &VdsStorage{
		vdsMap: make(map[uint64]entity.Vds),
		next:   1,
		mutx:   sync.RWMutex{},
	}
}

func TestCreateVds(t *testing.T) {
	setupVds()
	vds := entity.Vds{
		Uid:    1,
		Login:  "vdslogin",
		Pass:   "password",
		Name:   "Test VDS",
		Status: true,
	}
	createdVds := CreateVds(vds)

	assert.Equal(t, uint64(1), createdVds.Vid)
	assert.Equal(t, vds.Uid, createdVds.Uid)
	assert.Equal(t, vds.Login, createdVds.Login)
	assert.Equal(t, vds.Pass, createdVds.Pass)
	assert.Equal(t, vds.Name, createdVds.Name)
	assert.Equal(t, vds.Status, createdVds.Status)

	// Проверяем, что VDS была добавлена в карту
	assert.Equal(t, createdVds, s.vdsMap[createdVds.Vid])
}
