package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.EquipmentAttributes;

import builder.AttributeResultBuilder;

public class EquipmentAttributesBuilder extends RecordBuilder<EquipmentAttributes, EquipmentAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer itemId;

   private Short slot;

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
   public EquipmentAttributes construct() {
      return new EquipmentAttributes(itemId, slot, strength, dexterity, intelligence, luck, hp, mp, weaponAttack, magicAttack,
            weaponDefense, magicDefense, accuracy, avoidability, hands, speed, jump, slots);
   }

   @Override
   public EquipmentAttributesBuilder getThis() {
      return this;
   }

   public EquipmentAttributesBuilder setItemId(Integer itemId) {
      this.itemId = itemId;
      return getThis();
   }

   public EquipmentAttributesBuilder setSlot(Short slot) {
      this.slot = slot;
      return getThis();
   }

   public EquipmentAttributesBuilder setStrength(Integer strength) {
      this.strength = strength;
      return getThis();
   }

   public EquipmentAttributesBuilder setDexterity(Integer dexterity) {
      this.dexterity = dexterity;
      return getThis();
   }

   public EquipmentAttributesBuilder setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
      return getThis();
   }

   public EquipmentAttributesBuilder setLuck(Integer luck) {
      this.luck = luck;
      return getThis();
   }

   public EquipmentAttributesBuilder setHp(Integer hp) {
      this.hp = hp;
      return getThis();
   }

   public EquipmentAttributesBuilder setMp(Integer mp) {
      this.mp = mp;
      return getThis();
   }

   public EquipmentAttributesBuilder setWeaponAttack(Integer weaponAttack) {
      this.weaponAttack = weaponAttack;
      return getThis();
   }

   public EquipmentAttributesBuilder setMagicAttack(Integer magicAttack) {
      this.magicAttack = magicAttack;
      return getThis();
   }

   public EquipmentAttributesBuilder setWeaponDefense(Integer weaponDefense) {
      this.weaponDefense = weaponDefense;
      return getThis();
   }

   public EquipmentAttributesBuilder setMagicDefense(Integer magicDefense) {
      this.magicDefense = magicDefense;
      return getThis();
   }

   public EquipmentAttributesBuilder setAccuracy(Integer accuracy) {
      this.accuracy = accuracy;
      return getThis();
   }

   public EquipmentAttributesBuilder setAvoidability(Integer avoidability) {
      this.avoidability = avoidability;
      return getThis();
   }

   public EquipmentAttributesBuilder setHands(Integer hands) {
      this.hands = hands;
      return getThis();
   }

   public EquipmentAttributesBuilder setSpeed(Integer speed) {
      this.speed = speed;
      return getThis();
   }

   public EquipmentAttributesBuilder setJump(Integer jump) {
      this.jump = jump;
      return getThis();
   }

   public EquipmentAttributesBuilder setSlots(Integer slots) {
      this.slots = slots;
      return getThis();
   }
}
