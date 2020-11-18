package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.CharacterAttributes;

import builder.AttributeResultBuilder;

public class CharacterAttributesBuilder extends RecordBuilder<CharacterAttributes, CharacterAttributesBuilder>
      implements AttributeResultBuilder {
   private int accountId;

   private int worldId;

   private String name;

   private int level;

   private int experience;

   private int gachaponExperience;

   private int strength;

   private int dexterity;

   private int luck;

   private int intelligence;

   private int hp;

   private int mp;

   private int maxHp;

   private int maxMp;

   private int meso;

   private int hpMpUsed;

   private int jobId;

   private int skinColor;

   private byte gender;

   private int fame;

   private int hair;

   private int face;

   private int ap;

   private String sp;

   private int mapId;

   private int spawnPoint;

   private int gm;

   @Override
   public CharacterAttributes construct() {
      return new CharacterAttributes(accountId, worldId, name, level, experience, gachaponExperience, strength, dexterity, luck,
            intelligence, hp, mp, maxHp, maxMp, meso, hpMpUsed, jobId, skinColor, gender, fame, hair, face, ap, sp, mapId,
            spawnPoint, gm);
   }

   @Override
   public CharacterAttributesBuilder getThis() {
      return this;
   }

   public CharacterAttributesBuilder setAccountId(int accountId) {
      this.accountId = accountId;
      return getThis();
   }

   public CharacterAttributesBuilder setWorldId(int worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public CharacterAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }

   public CharacterAttributesBuilder setLevel(int level) {
      this.level = level;
      return getThis();
   }

   public CharacterAttributesBuilder setExperience(int experience) {
      this.experience = experience;
      return getThis();
   }

   public CharacterAttributesBuilder setGachaponExperience(int gachaponExperience) {
      this.gachaponExperience = gachaponExperience;
      return getThis();
   }

   public CharacterAttributesBuilder setStrength(int strength) {
      this.strength = strength;
      return getThis();
   }

   public CharacterAttributesBuilder setDexterity(int dexterity) {
      this.dexterity = dexterity;
      return getThis();
   }

   public CharacterAttributesBuilder setLuck(int luck) {
      this.luck = luck;
      return getThis();
   }

   public CharacterAttributesBuilder setIntelligence(int intelligence) {
      this.intelligence = intelligence;
      return getThis();
   }

   public CharacterAttributesBuilder setHp(int hp) {
      this.hp = hp;
      return getThis();
   }

   public CharacterAttributesBuilder setMp(int mp) {
      this.mp = mp;
      return getThis();
   }

   public CharacterAttributesBuilder setMaxHp(int maxHp) {
      this.maxHp = maxHp;
      return getThis();
   }

   public CharacterAttributesBuilder setMaxMp(int maxMp) {
      this.maxMp = maxMp;
      return getThis();
   }

   public CharacterAttributesBuilder setMeso(int meso) {
      this.meso = meso;
      return getThis();
   }

   public CharacterAttributesBuilder setHpMpUsed(int hpMpUsed) {
      this.hpMpUsed = hpMpUsed;
      return getThis();
   }

   public CharacterAttributesBuilder setJobId(int jobId) {
      this.jobId = jobId;
      return getThis();
   }

   public CharacterAttributesBuilder setSkinColor(int skinColor) {
      this.skinColor = skinColor;
      return getThis();
   }

   public CharacterAttributesBuilder setGender(byte gender) {
      this.gender = gender;
      return getThis();
   }

   public CharacterAttributesBuilder setFame(int fame) {
      this.fame = fame;
      return getThis();
   }

   public CharacterAttributesBuilder setHair(int hair) {
      this.hair = hair;
      return getThis();
   }

   public CharacterAttributesBuilder setFace(int face) {
      this.face = face;
      return getThis();
   }

   public CharacterAttributesBuilder setAp(int ap) {
      this.ap = ap;
      return getThis();
   }

   public CharacterAttributesBuilder setSp(String sp) {
      this.sp = sp;
      return getThis();
   }

   public CharacterAttributesBuilder setMapId(int mapId) {
      this.mapId = mapId;
      return getThis();
   }

   public CharacterAttributesBuilder setSpawnPoint(int spawnPoint) {
      this.spawnPoint = spawnPoint;
      return getThis();
   }

   public CharacterAttributesBuilder setGm(int gm) {
      this.gm = gm;
      return getThis();
   }
}
