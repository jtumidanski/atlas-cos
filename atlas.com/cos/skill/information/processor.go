package information

import (
	"atlas-cos/model"
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(skillId uint32) model.Provider[Model] {
	return func(skillId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestSkill(skillId), makeSkill)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(skillId uint32) (Model, error) {
	return func(skillId uint32) (Model, error) {
		return ByIdModelProvider(l, span)(skillId)()
	}
}

func makeSkill(data requests.DataBody[attributes]) (Model, error) {
	attr := data.Attributes
	return Model{
		action:        attr.Action,
		element:       attr.Element,
		animationTime: attr.AnimationTime,
		effects:       makeEffects(attr.Effects),
	}, nil
}

func makeEffects(es []effects) []Effect {
	var results = make([]Effect, 0)
	for _, e := range es {
		results = append(results, makeEffect(e))
	}
	return results
}

func makeEffect(e effects) Effect {
	return Effect{
		weaponAttack:         e.WeaponAttack,
		magicAttack:          e.MagicAttack,
		weaponDefense:        e.WeaponDefense,
		magicDefense:         e.MagicDefense,
		accuracy:             e.Accuracy,
		avoidability:         e.Avoidability,
		speed:                e.Speed,
		jump:                 e.Jump,
		hp:                   e.HP,
		mp:                   e.MP,
		hpr:                  e.HPR,
		mpr:                  e.MPR,
		mhprRate:             e.MHPRRate,
		mmprRate:             e.MMPRRate,
		mobSkill:             e.MobSkill,
		mobSkillLevel:        e.MobSkillLevel,
		mhpR:                 e.MHPR,
		mmpR:                 e.MMPR,
		hpCon:                e.HPCon,
		mpCon:                e.MPCon,
		duration:             e.Duration,
		target:               e.Target,
		barrier:              e.Barrier,
		mob:                  e.Mob,
		overtime:             e.OverTime,
		repeatEffect:         e.RepeatEffect,
		moveTo:               e.MoveTo,
		cp:                   e.CP,
		nuffSkill:            e.NuffSkill,
		skill:                e.Skill,
		x:                    e.X,
		y:                    e.Y,
		mobCount:             e.MobCount,
		moneyCon:             e.MoneyCon,
		cooldown:             e.Cooldown,
		morphId:              e.MorphId,
		ghost:                e.Ghost,
		fatigue:              e.Fatigue,
		berserk:              e.Berserk,
		booster:              e.Booster,
		prop:                 e.Prop,
		itemCon:              e.ItemCon,
		itemConNo:            e.ItemConNo,
		damage:               e.Damage,
		attackCount:          e.AttackCount,
		fixDamage:            e.FixDamage,
		bulletCount:          e.BulletCount,
		bulletConsume:        e.BulletConsume,
		mapProtection:        e.MapProtection,
		cureAbnormalStatuses: e.CureAbnormalStatuses,
		statups:              makeStatups(e.Statups),
	}
}

func makeStatups(statups []buffStatAmount) []Statup {
	results := make([]Statup, 0)
	for _, s := range statups {
		results = append(results, makeStatup(s))
	}
	return results
}

func makeStatup(s buffStatAmount) Statup {
	return Statup{
		buff:   s.Buff,
		amount: s.Amount,
	}
}
