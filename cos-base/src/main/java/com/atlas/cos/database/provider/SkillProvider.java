package com.atlas.cos.database.provider;

import java.util.List;
import java.util.Optional;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;

import com.app.database.util.QueryProviderUtil;
import com.atlas.cos.database.transformer.SkillTransformer;
import com.atlas.cos.entity.Skill;
import com.atlas.cos.model.SkillData;

public final class SkillProvider {
   private SkillProvider() {
   }

   public static List<SkillData> getSkills(EntityManager entityManager, int characterId) {
      TypedQuery<Skill> query = entityManager.createQuery("FROM Skill s WHERE s.characterId = :characterId", Skill.class);
      query.setParameter("characterId", characterId);
      return QueryProviderUtil.list(query, new SkillTransformer());
   }

   public static Optional<SkillData> getSkill(EntityManager entityManager, int characterId, int skillId) {
      TypedQuery<Skill> query = entityManager.createQuery("FROM Skill s WHERE s.characterId = :characterId AND s.skillId = "
            + ":skillId", Skill.class);
      query.setParameter("characterId", characterId);
      query.setParameter("skillId", skillId);
      return QueryProviderUtil.optionalElement(query, new SkillTransformer());
   }
}