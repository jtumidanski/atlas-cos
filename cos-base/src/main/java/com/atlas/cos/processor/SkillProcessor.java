package com.atlas.cos.processor;

import com.atlas.cos.database.administrator.SkillAdministrator;
import com.atlas.cos.database.provider.SkillProvider;
import com.atlas.cos.model.SkillData;
import com.atlas.cos.model.SkillInformation;
import database.Connection;

import java.util.Optional;

public final class SkillProcessor {
   private SkillProcessor() {
   }

   public static Optional<SkillData> getSkill(int characterId, int skillId) {
      return Connection.instance().element(entityManager -> SkillProvider.getSkill(entityManager, characterId, skillId));
   }

   public static void awardSkill(int characterId, int skillId) {
      SkillInformation skillInformation = SkillInformationProcessor.getSkillInformation(skillId).join();
      int maxLevel = skillInformation.effects().size();
      Connection.instance().with(entityManager -> SkillAdministrator
            .createSkill(entityManager, characterId, skillId, 0, maxLevel, -1));
   }

   public static void updateSkill(int characterId, int skillId, int level, int masterLevel, long expiration) {
      Connection.instance().with(entityManager -> {
         SkillData skillData = SkillProvider
               .getSkill(entityManager, characterId, skillId)
               .orElseGet(() -> SkillAdministrator.createSkill(entityManager, characterId, skillId, 0, masterLevel, -1));
         SkillAdministrator.updateSkill(entityManager, skillData.id(), (byte) level, masterLevel, expiration);
      });
   }
}
