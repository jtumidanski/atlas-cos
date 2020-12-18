package com.atlas.cos.attribute;

import rest.AttributeResult;

public record EquipmentStatistics(Integer itemId, Integer strength, Integer dexterity, Integer intelligence, Integer luck,
                                  Integer hp, Integer mp, Integer weaponAttack, Integer magicAttack, Integer weaponDefense,
                                  Integer magicDefense, Integer accuracy, Integer avoidability, Integer hands, Integer speed,
                                  Integer jump, Integer slots) implements AttributeResult {
}
