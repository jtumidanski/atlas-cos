package character

type builderConfiguration struct {
	useStarting4AP          bool
	useAutoAssignStartersAP bool
}

func (b *builderConfiguration) UseStarting4AP() bool {
	return b.useStarting4AP
}

func (b *builderConfiguration) UseAutoAssignStartersAP() bool {
	return b.useAutoAssignStartersAP
}

type builder struct {
	accountId    uint32
	worldId      byte
	name         string
	level        byte
	strength     uint16
	dexterity    uint16
	intelligence uint16
	luck         uint16
	maxHp        uint16
	maxMp        uint16
	jobId        uint16
	skinColor    byte
	gender       byte
	hair         uint32
	face         uint32
	ap           uint16
	mapId        uint32
}

func (b *builder) setJobId(jobId uint16) {
	b.jobId = jobId
}

func (b *builder) setMapId(mapId uint32) {
	b.mapId = mapId
}

func (b *builder) Build() *Model {
	return &Model{
		accountId:          b.accountId,
		worldId:            b.worldId,
		name:               b.name,
		level:              b.level,
		experience:         0,
		gachaponExperience: 0,
		strength:           b.strength,
		dexterity:          b.dexterity,
		intelligence:       b.intelligence,
		luck:               b.luck,
		hp:                 0,
		mp:                 0,
		maxHp:              b.maxHp,
		maxMp:              b.maxMp,
		meso:               0,
		hpMpUsed:           0,
		jobId:              b.jobId,
		skinColor:          b.skinColor,
		gender:             b.gender,
		fame:               0,
		hair:               b.hair,
		face:               b.face,
		ap:                 b.ap,
		sp:                 "",
		mapId:              b.mapId,
		spawnPoint:         0,
		gm:                 0,
	}
}

func NewBuilder(c builderConfiguration, accountId uint32, worldId byte, name string, skinColor byte, gender byte, hair uint32, face uint32) *builder {
	b := &builder{
		accountId: accountId,
		worldId:   worldId,
		name:      name,
		level:     1,
		jobId:     0,
		skinColor: skinColor,
		gender:    gender,
		hair:      hair,
		face:      face,
	}

	if !c.UseStarting4AP() {
		if c.UseAutoAssignStartersAP() {
			b.strength = 12
			b.dexterity = 5
			b.intelligence = 4
			b.luck = 4
			b.ap = 0
		} else {
			b.strength = 4
			b.dexterity = 4
			b.intelligence = 4
			b.luck = 4
			b.ap = 9
		}
	} else {
		b.strength = 4
		b.dexterity = 4
		b.intelligence = 4
		b.luck = 4
		b.ap = 0
	}

	b.maxHp = 50
	b.maxMp = 5
	return b
}

func NewBuilder2(c builderConfiguration, accountId uint32, worldId byte, name string, skinColor byte, gender byte, hair uint32, face uint32, level byte, mapId uint32) *builder {
	b := NewBuilder(c, accountId, worldId, name, skinColor, gender, hair, face)
	b.level = level
	b.mapId = mapId
	return b
}

func NewBuilder3(c builderConfiguration, accountId uint32, worldId byte, name string, skinColor byte, gender byte, hair uint32, face uint32, level byte, mapId uint32, jobId uint16) *builder {
	b := NewBuilder2(c, accountId, worldId, name, skinColor, gender, hair, face, level, mapId)
	b.jobId = jobId
	return b
}
