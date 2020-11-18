package com.atlas.cos.builder;

import builder.RecordBuilder;
import com.atlas.cos.attribute.CharacterAttributes;

import builder.AttributeResultBuilder;

public class CharacterAttributesBuilder extends RecordBuilder<CharacterAttributes, CharacterAttributesBuilder> implements AttributeResultBuilder {
   private static final String ACCOUNT_ID = "ACCOUNT_ID";

   private static final String WORLD_ID = "WORLD_ID";

   private static final String NAME = "NAME";

   private static final String LEVEL = "LEVEL";

   private static final String EXPERIENCE = "EXPERIENCE";

   private static final String GACHAPON_EXPERIENCE = "GACHAPON_EXPERIENCE";

   private static final String STRENGTH = "STRENGTH";

   private static final String DEXTERITY = "DEXTERITY";

   private static final String LUCK = "LUCK";

   private static final String INTELLIGENCE = "INTELLIGENCE";

   private static final String HP = "HP";

   private static final String MP = "MP";

   private static final String MAX_HP = "MAX_HP";

   private static final String MAX_MP = "MAX_MP";

   private static final String MESO = "MESO";

   private static final String HP_MP_USED = "HP_MP_USED";

   private static final String JOB_ID = "JOB_ID";

   private static final String SKIN_COLOR = "SKIN_COLOR";

   private static final String GENDER = "GENDER";

   private static final String FAME = "FAME";

   private static final String HAIR = "HAIR";

   private static final String FACE = "FACE";

   private static final String AP = "AP";

   private static final String SP = "SP";

   private static final String MAP_ID = "MAP_ID";

   private static final String SPAWN_POINT = "SPAWN_POINT";

   private static final String GM = "GM";

   @Override
   public CharacterAttributes construct() {
      return new CharacterAttributes(get(ACCOUNT_ID), get(WORLD_ID), get(NAME), get(LEVEL), get(EXPERIENCE), get(GACHAPON_EXPERIENCE), get(STRENGTH), get(DEXTERITY), get(LUCK), get(INTELLIGENCE), get(HP), get(MP), get(MAX_HP), get(MAX_MP), get(MESO), get(HP_MP_USED), get(JOB_ID), get(SKIN_COLOR), get(GENDER), get(FAME), get(HAIR), get(FACE), get(AP), get(SP), get(MAP_ID), get(SPAWN_POINT), get(GM));
   }

   @Override
   public CharacterAttributesBuilder getThis() {
      return this;
   }

   public CharacterAttributesBuilder setAccountId(int accountId) {
      return set(ACCOUNT_ID, accountId);
   }

   public CharacterAttributesBuilder setWorldId(int worldId) {
      return set(WORLD_ID, worldId);
   }

   public CharacterAttributesBuilder setName(String name) {
      return set(NAME, name);
   }

   public CharacterAttributesBuilder setLevel(int level) {
      return set(LEVEL, level);
   }

   public CharacterAttributesBuilder setExperience(int experience) {
      return set(EXPERIENCE, experience);
   }

   public CharacterAttributesBuilder setGachaponExperience(int gachaponExperience) {
      return set(GACHAPON_EXPERIENCE, gachaponExperience);
   }

   public CharacterAttributesBuilder setStrength(int strength) {
      return set(STRENGTH, strength);
   }

   public CharacterAttributesBuilder setDexterity(int dexterity) {
      return set(DEXTERITY, dexterity);
   }

   public CharacterAttributesBuilder setLuck(int luck) {
      return set(LUCK, luck);
   }

   public CharacterAttributesBuilder setIntelligence(int intelligence) {
      return set(INTELLIGENCE, intelligence);
   }

   public CharacterAttributesBuilder setHp(int hp) {
      return set(HP, hp);
   }

   public CharacterAttributesBuilder setMp(int mp) {
      return set(MP, mp);
   }

   public CharacterAttributesBuilder setMaxHp(int maxHp) {
      return set(MAX_HP, maxHp);
   }

   public CharacterAttributesBuilder setMaxMp(int maxMp) {
      return set(MAX_MP, maxMp);
   }

   public CharacterAttributesBuilder setMeso(int meso) {
      return set(MESO, meso);
   }

   public CharacterAttributesBuilder setHpMpUsed(int hpMpUsed) {
      return set(HP_MP_USED, hpMpUsed);
   }

   public CharacterAttributesBuilder setJobId(int jobId) {
      return set(JOB_ID, jobId);
   }

   public CharacterAttributesBuilder setSkinColor(int skinColor) {
      return set(SKIN_COLOR, skinColor);
   }

   public CharacterAttributesBuilder setGender(byte gender) {
      return set(GENDER, gender);
   }

   public CharacterAttributesBuilder setFame(int fame) {
      return set(FAME, fame);
   }

   public CharacterAttributesBuilder setHair(int hair) {
      return set(HAIR, hair);
   }

   public CharacterAttributesBuilder setFace(int face) {
      return set(FACE, face);
   }

   public CharacterAttributesBuilder setAp(int ap) {
      return set(AP, ap);
   }

   public CharacterAttributesBuilder setSp(String sp) {
      return set(SP, sp);
   }

   public CharacterAttributesBuilder setMapId(int mapId) {
      return set(MAP_ID, mapId);
   }

   public CharacterAttributesBuilder setSpawnPoint(int spawnPoint) {
      return set(SPAWN_POINT, spawnPoint);
   }

   public CharacterAttributesBuilder setGm(int gm) {
      return set(GM, gm);
   }
}
