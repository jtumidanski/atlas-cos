package com.atlas.cos.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.rest.attribute.CharacterAttributes;

import builder.AttributeResultBuilder;

public class CharacterAttributesBuilder extends RecordBuilder<CharacterAttributes, CharacterAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer accountId;

   private Integer worldId;

   private String name;

   private Integer level;

   private Integer experience;

   private Integer gachaponExperience;

   private Integer strength;

   private Integer dexterity;

   private Integer luck;

   private Integer intelligence;

   private Integer hp;

   private Integer mp;

   private Integer maxHp;

   private Integer maxMp;

   private Integer meso;

   private Integer hpMpUsed;

   private Integer jobId;

   private Integer skinColor;

   private Byte gender;

   private Integer fame;

   private Integer hair;

   private Integer face;

   private Integer ap;

   private String sp;

   private Integer mapId;

   private Integer spawnPoint;

   private Integer gm;

   private Integer x;

   private Integer y;

   private Byte stance;

   @Override
   public CharacterAttributes construct() {
      return new CharacterAttributes(accountId, worldId, name, level, experience, gachaponExperience, strength, dexterity, luck,
            intelligence, hp, mp, maxHp, maxMp, meso, hpMpUsed, jobId, skinColor, gender, fame, hair, face, ap, sp, mapId,
            spawnPoint, gm, x, y, stance);
   }

   @Override
   public CharacterAttributesBuilder getThis() {
      return this;
   }

   public CharacterAttributesBuilder setAccountId(Integer accountId) {
      this.accountId = accountId;
      return getThis();
   }

   public CharacterAttributesBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public CharacterAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }

   public CharacterAttributesBuilder setLevel(Integer level) {
      this.level = level;
      return getThis();
   }

   public CharacterAttributesBuilder setExperience(Integer experience) {
      this.experience = experience;
      return getThis();
   }

   public CharacterAttributesBuilder setGachaponExperience(Integer gachaponExperience) {
      this.gachaponExperience = gachaponExperience;
      return getThis();
   }

   public CharacterAttributesBuilder setStrength(Integer strength) {
      this.strength = strength;
      return getThis();
   }

   public CharacterAttributesBuilder setDexterity(Integer dexterity) {
      this.dexterity = dexterity;
      return getThis();
   }

   public CharacterAttributesBuilder setLuck(Integer luck) {
      this.luck = luck;
      return getThis();
   }

   public CharacterAttributesBuilder setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
      return getThis();
   }

   public CharacterAttributesBuilder setHp(Integer hp) {
      this.hp = hp;
      return getThis();
   }

   public CharacterAttributesBuilder setMp(Integer mp) {
      this.mp = mp;
      return getThis();
   }

   public CharacterAttributesBuilder setMaxHp(Integer maxHp) {
      this.maxHp = maxHp;
      return getThis();
   }

   public CharacterAttributesBuilder setMaxMp(Integer maxMp) {
      this.maxMp = maxMp;
      return getThis();
   }

   public CharacterAttributesBuilder setMeso(Integer meso) {
      this.meso = meso;
      return getThis();
   }

   public CharacterAttributesBuilder setHpMpUsed(Integer hpMpUsed) {
      this.hpMpUsed = hpMpUsed;
      return getThis();
   }

   public CharacterAttributesBuilder setJobId(Integer jobId) {
      this.jobId = jobId;
      return getThis();
   }

   public CharacterAttributesBuilder setSkinColor(Integer skinColor) {
      this.skinColor = skinColor;
      return getThis();
   }

   public CharacterAttributesBuilder setGender(Byte gender) {
      this.gender = gender;
      return getThis();
   }

   public CharacterAttributesBuilder setFame(Integer fame) {
      this.fame = fame;
      return getThis();
   }

   public CharacterAttributesBuilder setHair(Integer hair) {
      this.hair = hair;
      return getThis();
   }

   public CharacterAttributesBuilder setFace(Integer face) {
      this.face = face;
      return getThis();
   }

   public CharacterAttributesBuilder setAp(Integer ap) {
      this.ap = ap;
      return getThis();
   }

   public CharacterAttributesBuilder setSp(String sp) {
      this.sp = sp;
      return getThis();
   }

   public CharacterAttributesBuilder setMapId(Integer mapId) {
      this.mapId = mapId;
      return getThis();
   }

   public CharacterAttributesBuilder setSpawnPoint(Integer spawnPoint) {
      this.spawnPoint = spawnPoint;
      return getThis();
   }

   public CharacterAttributesBuilder setGm(Integer gm) {
      this.gm = gm;
      return getThis();
   }

   public CharacterAttributesBuilder setX(Integer x) {
      this.x = x;
      return getThis();
   }

   public CharacterAttributesBuilder setY(Integer y) {
      this.y = y;
      return getThis();
   }

   public CharacterAttributesBuilder setStance(Byte stance) {
      this.stance = stance;
      return getThis();
   }
}
