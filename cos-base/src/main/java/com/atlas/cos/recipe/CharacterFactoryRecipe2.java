package com.atlas.cos.recipe;

import java.util.ArrayList;
import java.util.List;

import com.atlas.cos.model.MapleJob;
import com.atlas.cos.model.SkillData;

public class CharacterFactoryRecipe2 {
   private Integer strength;

   private Integer dexterity;

   private Integer intelligence;

   private Integer luck;

   private Integer maxHp;

   private Integer maxMp;

   private Integer remainingAp;

   private Integer remainingSp;

   private Integer meso;

   private MapleJob job;

   private Integer level;

   private Integer map;

   private List<SkillData> skills;

   protected CharacterFactoryRecipe2() {
      if (!YamlConfig.config.server.USE_STARTING_AP_4) {
         if (YamlConfig.config.server.USE_AUTOASSIGN_STARTERS_AP) {
            strength = 12;
            dexterity = 5;
            intelligence = 4;
            luck = 4;
            remainingAp = 0;
         } else {
            strength = 4;
            dexterity = 4;
            intelligence = 4;
            luck = 4;
            remainingAp = 9;
         }
      } else {
         strength = 4;
         dexterity = 4;
         intelligence = 4;
         luck = 4;
         remainingAp = 0;
      }

      maxHp = 50;
      maxMp = 5;
      remainingSp = 0;
      meso = 0;
      skills = new ArrayList<>();
   }

   public CharacterFactoryRecipe2(MapleJob job, Integer level, Integer map) {
      this();
      this.job = job;
      this.level = level;
      this.map = map;
   }

   public void addStartingSkillLevel(Integer skillId, Integer level) {
      this.skills.add(new SkillData(skillId, level));
   }

   public List<SkillData> getStartingSkillLevel() {
      return skills;
   }

   public CharacterFactoryRecipe2 setStrength(Integer strength) {
      this.strength = strength;
      return this;
   }

   public CharacterFactoryRecipe2 setDexterity(Integer dexterity) {
      this.dexterity = dexterity;
      return this;
   }

   public CharacterFactoryRecipe2 setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
      return this;
   }

   public CharacterFactoryRecipe2 setLuck(Integer luck) {
      this.luck = luck;
      return this;
   }

   public CharacterFactoryRecipe2 setMaxHp(Integer maxHp) {
      this.maxHp = maxHp;
      return this;
   }

   public CharacterFactoryRecipe2 setMaxMp(Integer maxMp) {
      this.maxMp = maxMp;
      return this;
   }

   public CharacterFactoryRecipe2 setRemainingAp(Integer remainingAp) {
      this.remainingAp = remainingAp;
      return this;
   }

   public CharacterFactoryRecipe2 setRemainingSp(Integer remainingSp) {
      this.remainingSp = remainingSp;
      return this;
   }

   public CharacterFactoryRecipe2 increaseRemainingSp(Integer amount) {
      this.remainingSp += amount;
      return this;
   }

   public CharacterFactoryRecipe2 setMeso(Integer meso) {
      this.meso = meso;
      return this;
   }

//   public void apply(MapleCharacter character) {
//
//      character.init(str, dex, intelligence, luk, maxHp, maxMp, meso);
//      character.setLevel(level);
//      character.setJob(job);
//      character.setMapId(map);
//      character.setMaxHp(maxHp);
//      character.setMaxMp(maxMp);
//      character.setLevel(level);
//      character.setRemainingAp(remainingAp);
//      character.setRemainingSp(GameConstants.getSkillBook(character.getJob().getId()), remainingSp);
//
//
//      skills.forEach(entry -> SkillFactory.getSkill(entry.getLeft())
//            .ifPresent(skill -> character.changeSkillLevel(skill, entry.getRight().byteValue(), skill.getMaxLevel(), -1)));
//   }

   public Integer getStrength() {
      return strength;
   }

   public Integer getDexterity() {
      return dexterity;
   }

   public Integer getIntelligence() {
      return intelligence;
   }

   public Integer getLuck() {
      return luck;
   }

   public Integer getMaxHp() {
      return maxHp;
   }

   public Integer getMaxMp() {
      return maxMp;
   }

   public Integer getRemainingAp() {
      return remainingAp;
   }

   public Integer getRemainingSp() {
      return remainingSp;
   }

   public Integer getMeso() {
      return meso;
   }

   public MapleJob getJob() {
      return job;
   }

   public Integer getLevel() {
      return level;
   }

   public Integer getMap() {
      return map;
   }

   public List<SkillData> getSkills() {
      return skills;
   }
}
