package com.atlas.cos.rest;

import com.atlas.cos.attribute.BlockedNameAttributes;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.BlockedNameAttributesBuilder;
import com.atlas.cos.builder.CharacterAttributesBuilder;
import com.atlas.cos.model.BlockedNameData;
import com.atlas.cos.model.CharacterData;

import builder.ResultObjectBuilder;

public class ResultObjectFactory {
   public static ResultObjectBuilder create(BlockedNameData blockedNameData) {
      return new ResultObjectBuilder(BlockedNameAttributes.class, blockedNameData.name())
            .setAttribute(new BlockedNameAttributesBuilder().setName(blockedNameData.name()));
   }

   public static ResultObjectBuilder create(CharacterData data) {
      return new ResultObjectBuilder(CharacterAttributes.class, data.id())
            .setAttribute(new CharacterAttributesBuilder()
                  .setAccountId(data.accountId())
                  .setWorldId(data.worldId())
                  .setName(data.name())
                  .setLevel(data.level())
                  .setExperience(data.experience())
                  .setGachaponExperience(data.gachaponExperience())
                  .setStrength(data.strength())
                  .setDexterity(data.dexterity())
                  .setLuck(data.luck())
                  .setIntelligence(data.intelligence())
                  .setHp(data.hp())
                  .setMp(data.mp())
                  .setMaxHp(data.maxHp())
                  .setMaxMp(data.maxMp())
                  .setMeso(data.meso())
                  .setHpMpUsed(data.hpMpUsed())
                  .setJobId(data.jobId())
                  .setSkinColor(data.skinColor())
                  .setGender(data.gender())
                  .setFame(data.fame())
                  .setHair(data.hair())
                  .setFace(data.face())
                  .setAp(data.ap())
                  .setSp(data.sp())
                  .setMapId(data.mapId())
                  .setSpawnPoint(data.spawnPoint())
                  .setGm(data.gm())
            );
   }
}
