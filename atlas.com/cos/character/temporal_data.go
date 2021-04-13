package character

import "sync"

type temporalData struct {
	x      int16
	y      int16
	stance byte
}

func (d *temporalData) UpdatePosition(x int16, y int16) *temporalData {
	return &temporalData{
		x:      x,
		y:      y,
		stance: d.stance,
	}
}

func (d *temporalData) Update(x int16, y int16, stance byte) *temporalData {
	return &temporalData{
		x:      x,
		y:      y,
		stance: stance,
	}
}

func (d *temporalData) UpdateStance(stance byte) *temporalData {
	return &temporalData{
		x:      d.x,
		y:      d.y,
		stance: stance,
	}
}

type temporalRegistry struct {
	data           map[uint32]*temporalData
	mutex          *sync.RWMutex
	characterLocks map[uint32]*sync.RWMutex
}

func (r temporalRegistry) UpdatePosition(characterId uint32, x int16, y int16) {
	r.lockCharacter(characterId)
	if val, ok := r.data[characterId]; ok {
		r.data[characterId] = val.UpdatePosition(x, y)
	} else {
		r.data[characterId] = &temporalData{
			x:      x,
			y:      y,
			stance: 0,
		}
	}
	r.unlockCharacter(characterId)
}

func (r temporalRegistry) lockCharacter(characterId uint32) {
	if val, ok := r.characterLocks[characterId]; ok {
		val.Lock()
	} else {
		r.mutex.Lock()
		lock := sync.RWMutex{}
		r.characterLocks[characterId] = &lock
		r.mutex.Unlock()
		lock.Lock()
	}
}

func (r temporalRegistry) unlockCharacter(characterId uint32) {
	if val, ok := r.characterLocks[characterId]; ok {
		val.Unlock()
	}
}

func (r temporalRegistry) Update(characterId uint32, x int16, y int16, stance byte) {
	r.lockCharacter(characterId)
	if val, ok := r.data[characterId]; ok {
		r.data[characterId] = val.Update(x, y, stance)
	} else {
		r.data[characterId] = &temporalData{
			x:      x,
			y:      y,
			stance: stance,
		}
	}
	r.unlockCharacter(characterId)
}

func (r temporalRegistry) UpdateStance(characterId uint32, stance byte) {
	r.lockCharacter(characterId)
	if val, ok := r.data[characterId]; ok {
		r.data[characterId] = val.UpdateStance(stance)
	} else {
		r.data[characterId] = &temporalData{
			x:      0,
			y:      0,
			stance: stance,
		}
	}
	r.unlockCharacter(characterId)
}

var t *temporalRegistry
var once sync.Once

func GetTemporalRegistry() *temporalRegistry {
	once.Do(func() {
		t = &temporalRegistry{
			data:           make(map[uint32]*temporalData),
			mutex:          &sync.RWMutex{},
			characterLocks: make(map[uint32]*sync.RWMutex),
		}
	})
	return t
}
