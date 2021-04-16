package attributes

type SkillDataContainer struct {
	Data SkillData `json:"data"`
}

type SkillDataListContainer struct {
	Data []SkillData `json:"data"`
}

type SkillData struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes SkillAttributes `json:"attributes"`
}

type SkillAttributes struct {
	Action        bool                 `json:"action"`
	Element       string               `json:"element"`
	AnimationTime uint32               `json:"animationTime"`
	Effects       []StatEffectAttributes `json:"effects"`
}

type StatEffectAttributes struct {
	WeaponAttack  uint16  `json:"weaponAttack"`
	MagicAttack   uint16  `json:"magicAttack"`
	WeaponDefense uint16  `json:"weaponDefense"`
	MagicDefense  uint16  `json:"magicDefense"`
	Accuracy      uint16  `json:"accuracy"`
	Avoidability  uint16  `json:"avoidability"`
	Speed         uint16  `json:"speed"`
	Jump          uint16  `json:"jump"`
	HP            uint16  `json:"hp"`
	MP            uint16  `json:"mp"`
	HPR           float64 `json:"hpR"`
	MPR           float64 `json:"mpR"`
	MHPRRate      uint16  `json:"MHPRRate"`
	MMPRRate      uint16  `json:"MMPRRate"`
	MobSkill      uint16  `json:"mobSkill"`
	MobSkillLevel uint16  `json:"mobSkillLevel"`
	MHPR          byte    `json:"mhpr"`
	MMPR          byte    `json:"mmpr"`
	HPCon         uint16  `json:"HPCon"`
	MPCon         uint16  `json:"MPCon"`
	Duration      uint32  `json:"duration"`
	Target        uint32  `json:"target"`
	Barrier       uint32  `json:"barrier"`
	Mob           uint32  `json:"mob"`
	OverTime      bool    `json:"overTime"`
	RepeatEffect  bool    `json:"repeatEffect"`
	MoveTo        uint32  `json:"moveTo"`
	CP            uint32  `json:"cp"`
	NuffSkill     uint32  `json:"nuffSkill"`
	Skill         bool    `json:"skill"`
	X             int16   `json:"x"`
	Y             int16   `json:"y"`
	MobCount      uint32  `json:"mobCount"`
	MoneyCon      uint32  `json:"moneyCon"`
	Cooldown      uint32  `json:"cooldown"`
	MorphId       uint32  `json:"morphId"`
	Ghost         uint32  `json:"ghost"`
	Fatigue       uint32  `json:"fatigue"`
	Berserk       uint32  `json:"berserk"`
	Booster       uint32  `json:"booster"`
	Prop          float64 `json:"prop"`
	ItemCon       uint32  `json:"itemCon"`
	ItemConNo     uint32  `json:"itemConNo"`
	Damage        uint32  `json:"damage"`
	AttackCount   uint32  `json:"attackCount"`
	FixDamage     uint32  `json:"fixDamage"`
	//LT Point
	//RB Point
	BulletCount          uint16                     `json:"bulletCount"`
	BulletConsume        uint16                     `json:"bulletConsume"`
	MapProtection        byte                       `json:"mapProtection"`
	CureAbnormalStatuses []string                   `json:"cureAbnormalStatuses"`
	Statups              []BuffStatAmountAttributes `json:"statups"`
	MonsterStatus        []MonsterStatusAttributes  `json:"monsterStatus"`
	CardStats            CardItemUpStatsAttributes  `json:"cardStats"`
}

type BuffStatAmountAttributes struct {
	Buff   string `json:"buff"`
	Amount uint32 `json:"amount"`
}

type MonsterStatusAttributes struct {
	Status string `json:"status"`
	Value  uint32 `json:"value"`
}

type CardItemUpStatsAttributes struct {
	ItemCode    uint32           `json:"itemCode"`
	Probability uint32           `json:"probability"`
	Areas       []AreaAttributes `json:"areas"`
	InParty     bool             `json:"inParty"`
}

type AreaAttributes struct {
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}
