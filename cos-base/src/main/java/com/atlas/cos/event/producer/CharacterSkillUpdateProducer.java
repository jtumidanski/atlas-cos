package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterSkillUpdateEvent;

public final class CharacterSkillUpdateProducer {
   public static void updateSkill(int characterId, int skillId, int level, int masterLevel, long expiration) {
      EventProducerRegistry.getInstance().send(CharacterSkillUpdateEvent.class,
            EventConstants.TOPIC_CHARACTER_SKILL_UPDATE_EVENT, characterId,
            new CharacterSkillUpdateEvent(characterId, skillId, level, masterLevel, expiration));
   }
}
