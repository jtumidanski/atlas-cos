package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.model.EquipmentData;

public class EquipmentDataBuilder extends RecordBuilder<EquipmentData, EquipmentDataBuilder> {
   private int id;

   private int itemId;

   private short slot;

   private int strength;

   private int dexterity;

   private int intelligence;

   private int luck;

   private int weaponAttack;

   private int weaponDefense;

   private int magicAttack;

   private int magicDefense;

   private int accuracy;

   private int avoidability;

   private int speed;

   private int jump;

   private int hp;

   private int mp;

   private int slots;

   private int hands;

   @Override
   public EquipmentData construct() {
      return new EquipmentData(id, itemId, slot, strength, dexterity, intelligence, luck, hp, mp, weaponAttack, weaponDefense,
            magicAttack,
            magicDefense, accuracy, avoidability, hands, speed, jump, slots);
   }

   @Override
   public EquipmentDataBuilder getThis() {
      return this;
   }

   public EquipmentDataBuilder setId(int id) {
      this.id = id;
      return getThis();
   }

   public EquipmentDataBuilder setItemId(int itemId) {
      this.itemId = itemId;
      return getThis();
   }

   public EquipmentDataBuilder setSlot(short slot) {
      this.slot = slot;
      return getThis();
   }

   public EquipmentDataBuilder setStrength(int strength) {
      this.strength = strength;
      return getThis();
   }

   public EquipmentDataBuilder setDexterity(int dexterity) {
      this.dexterity = dexterity;
      return getThis();
   }

   public EquipmentDataBuilder setIntelligence(int intelligence) {
      this.intelligence = intelligence;
      return getThis();
   }

   public EquipmentDataBuilder setLuck(int luck) {
      this.luck = luck;
      return getThis();
   }

   public EquipmentDataBuilder setWeaponAttack(int weaponAttack) {
      this.weaponAttack = weaponAttack;
      return getThis();
   }

   public EquipmentDataBuilder setWeaponDefense(int weaponDefense) {
      this.weaponDefense = weaponDefense;
      return getThis();
   }

   public EquipmentDataBuilder setMagicAttack(int magicAttack) {
      this.magicAttack = magicAttack;
      return getThis();
   }

   public EquipmentDataBuilder setMagicDefense(int magicDefense) {
      this.magicDefense = magicDefense;
      return getThis();
   }

   public EquipmentDataBuilder setAccuracy(int accuracy) {
      this.accuracy = accuracy;
      return getThis();
   }

   public EquipmentDataBuilder setAvoidability(int avoidability) {
      this.avoidability = avoidability;
      return getThis();
   }

   public EquipmentDataBuilder setHands(int hands) {
      this.hands = hands;
      return getThis();
   }

   public EquipmentDataBuilder setSpeed(int speed) {
      this.speed = speed;
      return getThis();
   }

   public EquipmentDataBuilder setJump(int jump) {
      this.jump = jump;
      return getThis();
   }

   public EquipmentDataBuilder setHp(int hp) {
      this.hp = hp;
      return getThis();
   }

   public EquipmentDataBuilder setMp(int mp) {
      this.mp = mp;
      return getThis();
   }

   public EquipmentDataBuilder setSlots(int slots) {
      this.slots = slots;
      return getThis();
   }
}
