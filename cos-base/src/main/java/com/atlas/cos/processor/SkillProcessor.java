package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.database.provider.SkillProvider;
import com.atlas.cos.model.SkillData;

import database.Connection;

public final class SkillProcessor {
   private SkillProcessor() {
   }

   public static Optional<SkillData> getSkill(int characterId, int skillId) {
      return Connection.instance().element(entityManager -> SkillProvider.getSkill(entityManager, characterId, skillId));
   }
}
