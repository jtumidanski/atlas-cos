package character

import (
	"atlas-cos/inventory"
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameAdjustHealth     = "adjust_health_command"
	consumerNameAdjustMana       = "adjust_mana_command"
	consumerNameAdjustMeso       = "adjust_meso_command"
	consumerNameAssignAP         = "assign_ap_command"
	consumerNameAssignSP         = "assign_sp_command"
	consumerNameChangeMap        = "change_map_command"
	consumerNameAdjustJob        = "adjust_job_command"
	consumerNameDropItem         = "drop_item_command"
	consumerNameResetAP          = "reset_ap_command"
	consumerNameMoveItem         = "move_item_command"
	consumerNameAdjustExperience = "character_experience_event" // TODO change this
	consumerNameGainLevelEvent   = "character_level_event"      // TODO change this
	consumerNameMovementEvent    = "character_movement_event"   // TODO change this
	consumerNameStatusEvent      = "character_status_event"
	consumerNameEquipItem        = "equip_item_command"
	consumerNameUnequipItem      = "unequip_item_command"
	topicTokenAdjustHealth       = "TOPIC_ADJUST_HEALTH"
	topicTokenAdjustMana         = "TOPIC_ADJUST_MANA"
	topicTokenAdjustMeso         = "TOPIC_ADJUST_MESO"
	topicTokenAssignAP           = "TOPIC_ASSIGN_AP_COMMAND"
	topicTokenAssignSP           = "TOPIC_ASSIGN_SP_COMMAND"
	topicTokenChangeMap          = "TOPIC_CHANGE_MAP_COMMAND"
	topicTokenAdjustJob          = "TOPIC_CHARACTER_ADJUST_JOB"
	topicTokenDropItem           = "TOPIC_CHARACTER_DROP_ITEM"
	topicTokenResetAP            = "TOPIC_CHARACTER_RESET_AP"
	topicTokenMoveItem           = "TOPIC_MOVE_ITEM"
	topicTokenAdjustExperience   = "TOPIC_CHARACTER_EXPERIENCE_EVENT"
	topicTokenGainLevel          = "TOPIC_CHARACTER_LEVEL_EVENT"
	topicTokenMovement           = "TOPIC_CHARACTER_MOVEMENT"
	topicTokenStatusEvent        = "TOPIC_CHARACTER_STATUS"
	topicTokenEquipItem          = "TOPIC_EQUIP_ITEM"
	topicTokenUnequipItem        = "TOPIC_UNEQUIP_ITEM"
)

func AdjustHealthCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[adjustHealthCommand](consumerNameAdjustHealth, topicTokenAdjustHealth, groupId, handleAdjustHealthCommand(db))
	}
}

type adjustHealthCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func handleAdjustHealthCommand(db *gorm.DB) kafka.HandlerFunc[adjustHealthCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command adjustHealthCommand) {
		if command.Amount == 0 {
			l.Infoln("Received erroneous command to adjust by 0. This should be cleaned up.")
			return
		}
		AdjustHealth(l, db, span)(command.CharacterId, command.Amount)
	}
}

func AdjustManaCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[adjustManaCommand](consumerNameAdjustMana, topicTokenAdjustMana, groupId, handleAdjustManaCommand(db))
	}
}

type adjustManaCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func handleAdjustManaCommand(db *gorm.DB) kafka.HandlerFunc[adjustManaCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command adjustManaCommand) {
		if command.Amount == 0 {
			l.Infoln("Received erroneous command to adjust by 0. This should be cleaned up.")
			return
		}
		AdjustMana(l, db, span)(command.CharacterId, command.Amount)
	}
}

func AdjustMesoCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[adjustMesoCommand](consumerNameAdjustMeso, topicTokenAdjustMeso, groupId, handleAdjustMesoCommand(db))
	}
}

type adjustMesoCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int32  `json:"amount"`
	Show        bool   `json:"show"`
}

func handleAdjustMesoCommand(db *gorm.DB) kafka.HandlerFunc[adjustMesoCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command adjustMesoCommand) {
		AdjustMeso(l, db, span)(command.CharacterId, command.Amount, command.Show)
	}
}

func AssignAPCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[assignAPCommand](consumerNameAssignAP, topicTokenAssignAP, groupId, handleAssignAPCommand(db))
	}
}

type assignAPCommand struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func handleAssignAPCommand(db *gorm.DB) kafka.HandlerFunc[assignAPCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command assignAPCommand) {
		switch command.Type {
		case "STRENGTH":
			AssignStrength(l, db, span)(command.CharacterId)
			break
		case "DEXTERITY":
			AssignDexterity(l, db, span)(command.CharacterId)
			break
		case "INTELLIGENCE":
			AssignIntelligence(l, db, span)(command.CharacterId)
			break
		case "LUCK":
			AssignLuck(l, db, span)(command.CharacterId)
			break
		case "HP":
			AssignHp(l, db, span)(command.CharacterId)
			break
		case "MP":
			AssignMp(l, db, span)(command.CharacterId)
			break
		}
	}
}

func AssignSPCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[assignSPCommand](consumerNameAssignSP, topicTokenAssignSP, groupId, handleAssignSPCommand(db))
	}
}

type assignSPCommand struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func handleAssignSPCommand(db *gorm.DB) kafka.HandlerFunc[assignSPCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command assignSPCommand) {
		AssignSP(l, db, span)(command.CharacterId, command.SkillId)
	}
}

func ChangeMapCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[changeMapCommand](consumerNameChangeMap, topicTokenChangeMap, groupId, handleChangeMapCommand(db))
	}
}

type changeMapCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
}

func handleChangeMapCommand(db *gorm.DB) kafka.HandlerFunc[changeMapCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command changeMapCommand) {
		ChangeMap(l, db, span)(command.CharacterId, command.WorldId, command.ChannelId, command.MapId, command.PortalId)
	}
}

func AdjustJobCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[adjustJobCommand](consumerNameAdjustJob, topicTokenAdjustJob, groupId, handleAdjustJobCommand(db))
	}
}

type adjustJobCommand struct {
	CharacterId uint32 `json:"characterId"`
	JobId       uint16 `json:"jobId"`
}

func handleAdjustJobCommand(db *gorm.DB) kafka.HandlerFunc[adjustJobCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command adjustJobCommand) {
		err := AdjustJob(l, db, span)(command.CharacterId, command.JobId)
		if err != nil {
			l.WithError(err).Errorf("Unable to adjust the job of character %d.", command.CharacterId)
		}
	}
}

func DropItemCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[dropItemCommand](consumerNameDropItem, topicTokenDropItem, groupId, handleDropItemCommand(db))
	}
}

type dropItemCommand struct {
	WorldId       byte   `json:"worldId"`
	ChannelId     byte   `json:"channelId"`
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Quantity      int16  `json:"quantity"`
}

func handleDropItemCommand(db *gorm.DB) kafka.HandlerFunc[dropItemCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command dropItemCommand) {
		l.Debugf("Request to drop %d item in slot %d for character %d.", command.Quantity, command.Source, command.CharacterId)
		_ = DropItem(l, db, span)(command.WorldId, command.ChannelId, command.CharacterId, command.InventoryType, command.Source, command.Quantity)
	}
}

func ResetAPCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[resetAPCommand](consumerNameResetAP, topicTokenResetAP, groupId, handleResetAPCommand(db))
	}
}

type resetAPCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func handleResetAPCommand(db *gorm.DB) kafka.HandlerFunc[resetAPCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command resetAPCommand) {
		err := ResetAP(l, db, span)(command.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to reset AP of character %d.", command.CharacterId)
		}
	}
}

func MoveItemCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[moveItemCommand](consumerNameMoveItem, topicTokenMoveItem, groupId, handleMoveItemCommand(db))
	}
}

type moveItemCommand struct {
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Destination   int16  `json:"destination"`
}

func handleMoveItemCommand(db *gorm.DB) kafka.HandlerFunc[moveItemCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command moveItemCommand) {
		_ = MoveItem(l, db, span)(command.CharacterId, command.InventoryType, command.Source, command.Destination)
	}
}

func AdjustExperienceCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[gainExperienceEvent](consumerNameAdjustExperience, topicTokenAdjustExperience, groupId, handleAdjustExperienceCommand(db))
	}
}

type gainExperienceEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

func handleAdjustExperienceCommand(db *gorm.DB) kafka.HandlerFunc[gainExperienceEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event gainExperienceEvent) {
		GainExperience(l, db, span)(event.CharacterId, event.PersonalGain+event.PartyGain)
	}
}

func GainLevelCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[gainLevelEvent](consumerNameGainLevelEvent, topicTokenGainLevel, groupId, handleGainLevelEvent(db))
	}
}

type gainLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func handleGainLevelEvent(db *gorm.DB) kafka.HandlerFunc[gainLevelEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event gainLevelEvent) {
		GainLevel(l, db, span)(event.CharacterId)
	}
}

func MovementCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[movementCommand](consumerNameMovementEvent, topicTokenMovement, groupId, handleMovementEvent(db))
	}
}

type movementCommand struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement rawMovement `json:"rawMovement"`
}

type rawMovement []byte

func handleMovementEvent(db *gorm.DB) kafka.HandlerFunc[movementCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event movementCommand) {
		if event.X != 0 || event.Y != 0 {
			MoveCharacter(l, db, span)(event.CharacterId, event.X, event.Y, event.Stance)
		} else if event.Stance != 0 {
			UpdateStance(l, db)(event.CharacterId, event.Stance)
		}
	}
}

func StatusEventConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterStatusEvent](consumerNameStatusEvent, topicTokenStatusEvent, groupId, handleStatusEvent(db))
	}
}

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func handleStatusEvent(db *gorm.DB) kafka.HandlerFunc[characterStatusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterStatusEvent) {
		if event.Type == "LOGIN" {
			UpdateLoginPosition(l, db, span)(event.CharacterId)
		}
	}
}

func EquipItemConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterEquipItemCommand](consumerNameEquipItem, topicTokenEquipItem, groupId, handleEquipItemCommand(db))
	}
}

type characterEquipItemCommand struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func handleEquipItemCommand(db *gorm.DB) kafka.HandlerFunc[characterEquipItemCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command characterEquipItemCommand) {
		l.Debugf("CharacterId = %d, Source = %d, Destination = %d.", command.CharacterId, command.Source, command.CharacterId)
		e, err := inventory.GetEquippedItemBySlot(l, db)(command.CharacterId, command.Source)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item to equip for character %d in slot %d.", command.CharacterId, command.Source)
			return
		}
		inventory.EquipItemForCharacter(l, db, span)(command.CharacterId, e.EquipmentId())
	}
}

func UnEquipItemConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterUnequipItem](consumerNameUnequipItem, topicTokenUnequipItem, groupId, handleUnEquipItemCommand(db))
	}
}

type characterUnequipItem struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func handleUnEquipItemCommand(db *gorm.DB) kafka.HandlerFunc[characterUnequipItem] {
	return func(l logrus.FieldLogger, span opentracing.Span, command characterUnequipItem) {
		l.Debugf("CharacterId = %d, Source = %d, Destination = %d.", command.CharacterId, command.Source, command.CharacterId)
		e, err := inventory.GetEquippedItemBySlot(l, db)(command.CharacterId, command.Source)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item to equip for character %d in slot %d.", command.CharacterId, command.Source)
			return
		}
		inventory.UnequipItemForCharacter(l, db, span)(command.CharacterId, e.EquipmentId(), command.Source)
	}
}
