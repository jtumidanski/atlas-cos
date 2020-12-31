package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.Skill;
import com.atlas.cos.model.SkillData;

import transformer.SqlTransformer;

public class SkillTransformer implements SqlTransformer<SkillData, Skill> {
   @Override
   public SkillData transform(Skill resultSet) {
      return new SkillData(resultSet.getId(),
            resultSet.getSkillId(),
            resultSet.getSkillLevel(),
            resultSet.getMasterLevel(),
            resultSet.getExpiration());
   }
}
