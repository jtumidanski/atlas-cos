package com.atlas.cos.database.administrator;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.SkillTransformer;
import com.atlas.cos.entity.Skill;
import com.atlas.cos.model.SkillData;

import javax.persistence.EntityManager;

public final class SkillAdministrator {
   private SkillAdministrator() {
   }

   public static SkillData createSkill(EntityManager entityManager, int characterId, int skillId, int level, int masterLevel, long expiration) {
      Skill skill = new Skill();
      skill.setCharacterId(characterId);
      skill.setSkillId(skillId);
      skill.setSkillLevel(level);
      skill.setMasterLevel(masterLevel);
      skill.setExpiration(expiration);
      QueryAdministratorUtil.insert(entityManager, skill);
      return new SkillTransformer().transform(skill);
   }

   public static void updateSkill(EntityManager entityManager, int id, byte level, int masterLevel, long expiration) {
      QueryAdministratorUtil.update(entityManager, Skill.class, id, skill -> {
         skill.setSkillLevel((int) level);
         skill.setMasterLevel(masterLevel);
         skill.setExpiration(expiration);
      });
   }
}