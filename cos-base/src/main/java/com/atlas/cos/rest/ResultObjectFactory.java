package com.atlas.cos.rest;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.attribute.BlockedNameAttributes;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.attribute.LocationAttributes;
import com.atlas.cos.attribute.LocationType;
import com.atlas.cos.builder.BlockedNameAttributesBuilder;
import com.atlas.cos.builder.CharacterAttributesBuilder;
import com.atlas.cos.builder.EquipmentAttributesBuilder;
import com.atlas.cos.builder.LocationAttributesBuilder;
import com.atlas.cos.model.BlockedNameData;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.CharacterTemporalData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.SavedLocationData;

import builder.ResultObjectBuilder;

public class ResultObjectFactory {
   public static ResultObjectBuilder create(BlockedNameData blockedNameData) {
      return new ResultObjectBuilder(BlockedNameAttributes.class, blockedNameData.name())
            .setAttribute(new BlockedNameAttributesBuilder().setName(blockedNameData.name()));
   }

   public static ResultObjectBuilder create(CharacterData data) {

      CharacterTemporalData temporalData = CharacterTemporalRegistry.getInstance()
            .getTemporalData(data.id());

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
                  .setX(temporalData.x())
                  .setY(temporalData.y())
                  .setStance(temporalData.stance())
            );
   }

   public static ResultObjectBuilder create(EquipmentData data) {
      return new ResultObjectBuilder(EquipmentAttributes.class, data.id())
            .setAttribute(new EquipmentAttributesBuilder()
                  .setEquipmentId(data.equipmentId())
                  .setSlot(data.slot())
            );
   }

   public static ResultObjectBuilder create(SavedLocationData data) {
      return new ResultObjectBuilder(LocationAttributes.class, data.id())
            .setAttribute(new LocationAttributesBuilder()
                  .setType(LocationType.valueOf(data.type()))
                  .setMapId(data.mapId())
                  .setPortalId(data.portalId())
            );
   }
}
