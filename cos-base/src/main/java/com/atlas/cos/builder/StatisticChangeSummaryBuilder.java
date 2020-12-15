package com.atlas.cos.builder;

import com.atlas.cos.model.StatisticChangeSummary;

public class StatisticChangeSummaryBuilder {
   private int hp;

   private int maxHp;

   private int mp;

   private int maxMp;

   private int strength;

   private int dexterity;

   private int intelligence;

   private int luck;

   private int ap;

   private int level;

   public StatisticChangeSummaryBuilder() {
   }

   public StatisticChangeSummaryBuilder(StatisticChangeSummary other) {
      this.hp = other.hp();
      this.maxHp = other.maxHp();
      this.mp = other.mp();
      this.maxMp = other.maxMp();
      this.strength = other.strength();
      this.dexterity = other.dexterity();
      this.intelligence = other.intelligence();
      this.luck = other.luck();
      this.ap = other.ap();
   }

   public StatisticChangeSummary build() {
      return new StatisticChangeSummary(hp, maxHp, mp, maxMp, strength, dexterity, intelligence, luck, ap, level);
   }

   public StatisticChangeSummaryBuilder setHp(int hp) {
      this.hp = hp;
      return this;
   }

   public StatisticChangeSummaryBuilder setMaxHp(int maxHp) {
      this.maxHp = maxHp;
      return this;
   }

   public StatisticChangeSummaryBuilder setMp(int mp) {
      this.mp = mp;
      return this;
   }

   public StatisticChangeSummaryBuilder setMaxMp(int maxMp) {
      this.maxMp = maxMp;
      return this;
   }

   public StatisticChangeSummaryBuilder setStrength(int strength) {
      this.strength = strength;
      return this;
   }

   public StatisticChangeSummaryBuilder setDexterity(int dexterity) {
      this.dexterity = dexterity;
      return this;
   }

   public StatisticChangeSummaryBuilder setIntelligence(int intelligence) {
      this.intelligence = intelligence;
      return this;
   }

   public StatisticChangeSummaryBuilder setLuck(int luck) {
      this.luck = luck;
      return this;
   }

   public StatisticChangeSummaryBuilder setAp(int ap) {
      this.ap = ap;
      return this;
   }

   public StatisticChangeSummaryBuilder setLevel(int level) {
      this.level = level;
      return this;
   }
}
