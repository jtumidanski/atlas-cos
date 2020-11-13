package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

import transformer.SqlTransformer;

public class CharacterDataTransformer implements SqlTransformer<CharacterData, Character> {
   @Override
   public CharacterData transform(Character character) {
      return new CharacterData(character.getId(), character.getAccountId(), character.getWorld(), character.getName(),
            character.getLevel(), character.getExp(), character.getGachaponExp(), character.getStrength(),
            character.getDexterity(), character.getLuck(), character.getIntelligence(), character.getHp(), character.getMp(),
            character.getMaxHp(), character.getMaxMp(), character.getMeso(), character.getHpMpUsed(), character.getJob(),
            character.getSkinColor(), character.getGender(), character.getFame(), character.getHair(), character.getFace(),
            character.getAp(), character.getSp(), character.getMap(), character.getSpawnPoint(), character.getGm());
   }
}
