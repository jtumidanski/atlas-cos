package com.atlas.cos.builder;

import com.atlas.cos.attribute.CharacterAttributes;

import builder.AttributeResultBuilder;
import builder.Builder;

public class CharacterAttributesBuilder extends Builder<CharacterAttributes, CharacterAttributesBuilder>
      implements AttributeResultBuilder {
   @Override
   public CharacterAttributes construct() {
      return new CharacterAttributes();
   }

   @Override
   public CharacterAttributesBuilder getThis() {
      return this;
   }

   public CharacterAttributesBuilder setAccountId(Integer accountId) {
      return add(attr -> attr.setAccountId(accountId));
   }

   public CharacterAttributesBuilder setWorldId(Integer worldId) {
      return add(attr -> attr.setWorldId(worldId));
   }

   public CharacterAttributesBuilder setName(String name) {
      return add(attr -> attr.setName(name));
   }

   public CharacterAttributesBuilder setLevel(Integer level) {
      return add(attr -> attr.setLevel(level));
   }

   public CharacterAttributesBuilder setExperience(Integer experience) {
      return add(attr -> attr.setExperience(experience));
   }

   public CharacterAttributesBuilder setGachaponExperience(Integer gachaponExperience) {
      return add(attr -> attr.setGachaponExperience(gachaponExperience));
   }

   public CharacterAttributesBuilder setStrength(Integer strength) {
      return add(attr -> attr.setStrength(strength));
   }

   public CharacterAttributesBuilder setDexterity(Integer dexterity) {
      return add(attr -> attr.setDexterity(dexterity));
   }

   public CharacterAttributesBuilder setLuck(Integer luck) {
      return add(attr -> attr.setLuck(luck));
   }

   public CharacterAttributesBuilder setIntelligence(Integer intelligence) {
      return add(attr -> attr.setIntelligence(intelligence));
   }

   public CharacterAttributesBuilder setHp(Integer hp) {
      return add(attr -> attr.setHp(hp));
   }

   public CharacterAttributesBuilder setMp(Integer mp) {
      return add(attr -> attr.setMp(mp));
   }

   public CharacterAttributesBuilder setMaxHp(Integer maxHp) {
      return add(attr -> attr.setMaxHp(maxHp));
   }

   public CharacterAttributesBuilder setMaxMp(Integer maxMp) {
      return add(attr -> attr.setMaxMp(maxMp));
   }

   public CharacterAttributesBuilder setMeso(Integer meso) {
      return add(attr -> attr.setMeso(meso));
   }

   public CharacterAttributesBuilder setHpMpUsed(Integer hpMpUsed) {
      return add(attr -> attr.setHpMpUsed(hpMpUsed));
   }

   public CharacterAttributesBuilder setJobId(Integer jobId) {
      return add(attr -> attr.setJobId(jobId));
   }

   public CharacterAttributesBuilder setSkinColor(Integer skinColor) {
      return add(attr -> attr.setSkinColor(skinColor));
   }

   public CharacterAttributesBuilder setGender(Integer gender) {
      return add(attr -> attr.setGender(gender));
   }

   public CharacterAttributesBuilder setFame(Integer fame) {
      return add(attr -> attr.setFame(fame));
   }

   public CharacterAttributesBuilder setHair(Integer hair) {
      return add(attr -> attr.setHair(hair));
   }

   public CharacterAttributesBuilder setFace(Integer face) {
      return add(attr -> attr.setFace(face));
   }

   public CharacterAttributesBuilder setAp(Integer ap) {
      return add(attr -> attr.setAp(ap));
   }

   public CharacterAttributesBuilder setSp(String sp) {
      return add(attr -> attr.setSp(sp));
   }

   public CharacterAttributesBuilder setMapId(Integer mapId) {
      return add(attr -> attr.setMapId(mapId));
   }

   public CharacterAttributesBuilder setSpawnPoint(Integer spawnPoint) {
      return add(attr -> attr.setSpawnPoint(spawnPoint));
   }

   public CharacterAttributesBuilder setGm(Integer gm) {
      return add(attr -> attr.setGm(gm));
   }
}
