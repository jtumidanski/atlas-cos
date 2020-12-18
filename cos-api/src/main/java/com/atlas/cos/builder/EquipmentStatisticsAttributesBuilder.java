package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.EquipmentStatisticsAttributes;

import builder.AttributeResultBuilder;

public class EquipmentStatisticsAttributesBuilder extends RecordBuilder<EquipmentStatisticsAttributes, EquipmentStatisticsAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer itemId;

   private Integer strength;

   private Integer dexterity;

   private Integer intelligence;

   private Integer luck;

   private Integer hp;

   private Integer mp;

   private Integer weaponAttack;

   private Integer magicAttack;

   private Integer weaponDefense;

   private Integer magicDefense;

   private Integer accuracy;

   private Integer avoidability;

   private Integer hands;

   private Integer speed;

   private Integer jump;

   private Integer slots;

   @Override
   public EquipmentStatisticsAttributes construct() {
      return new EquipmentStatisticsAttributes(itemId, strength, dexterity, intelligence, luck, hp, mp, weaponAttack, magicAttack,
            weaponDefense, magicDefense, accuracy, avoidability, hands, speed, jump, slots);
   }

   @Override
   public EquipmentStatisticsAttributesBuilder getThis() {
      return this;
   }

   public EquipmentStatisticsAttributesBuilder setItemId(Integer itemId) {
      this.itemId = itemId;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setStrength(Integer strength) {
      this.strength = strength;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setDexterity(Integer dexterity) {
      this.dexterity = dexterity;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setLuck(Integer luck) {
      this.luck = luck;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setHp(Integer hp) {
      this.hp = hp;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setMp(Integer mp) {
      this.mp = mp;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setWeaponAttack(Integer weaponAttack) {
      this.weaponAttack = weaponAttack;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setMagicAttack(Integer magicAttack) {
      this.magicAttack = magicAttack;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setWeaponDefense(Integer weaponDefense) {
      this.weaponDefense = weaponDefense;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setMagicDefense(Integer magicDefense) {
      this.magicDefense = magicDefense;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setAccuracy(Integer accuracy) {
      this.accuracy = accuracy;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setAvoidability(Integer avoidability) {
      this.avoidability = avoidability;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setHands(Integer hands) {
      this.hands = hands;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setSpeed(Integer speed) {
      this.speed = speed;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setJump(Integer jump) {
      this.jump = jump;
      return getThis();
   }

   public EquipmentStatisticsAttributesBuilder setSlots(Integer slots) {
      this.slots = slots;
      return getThis();
   }
}
