package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record CharacterAttributes(Integer accountId, Integer worldId, String name, Integer level, Integer experience,
                                  Integer gachaponExperience, Integer strength, Integer dexterity, Integer luck,
                                  Integer intelligence, Integer hp, Integer mp, Integer maxHp, Integer maxMp, Integer meso,
                                  Integer hpMpUsed, Integer jobId, Integer skinColor, Byte gender, Integer fame, Integer hair,
                                  Integer face, Integer ap, String sp, Integer mapId, Integer spawnPoint, Integer gm, Integer x,
                                  Integer y, Byte stance)
      implements AttributeResult {
}
