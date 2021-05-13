package information

import "github.com/sirupsen/logrus"

func GetById(l logrus.FieldLogger) func(skillId uint32) (*Model, error) {
	return func(skillId uint32) (*Model, error) {
		s, err := requestSkill(skillId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve skill %d information.", skillId)
			return nil, err
		}
		return makeSkill(s.Data), nil
	}
}

func makeSkill(data dataBody) *Model {
	attr := data.Attributes

	return &Model{
		action:        attr.Action,
		element:       attr.Element,
		animationTime: attr.AnimationTime,
		effects:       makeEffects(attr.Effects),
	}
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
