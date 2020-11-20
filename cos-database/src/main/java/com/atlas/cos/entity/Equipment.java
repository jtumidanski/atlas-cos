package com.atlas.cos.entity;

import java.io.Serializable;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.NamedQueries;
import javax.persistence.NamedQuery;
import javax.persistence.Table;

@Entity
@NamedQueries({
      @NamedQuery(name = Equipment.GET_BY_ID, query = "SELECT e FROM Equipment e WHERE e.id = :id"),
      @NamedQuery(name = Equipment.GET_FOR_CHARACTER, query = "SELECT e FROM Equipment e WHERE e.characterId = :characterId"),
      @NamedQuery(name = Equipment.GET_FOR_CHARACTER_BY_SLOT, query = "SELECT e FROM Equipment e WHERE e.characterId = "
            + ":characterId AND e.slot = :slot")
})
@Table(name = "equipment")
public class Equipment implements Serializable {
   private static final long serialVersionUID = 1L;

   public static final String GET_FOR_CHARACTER = "Equipment.GET_FOR_CHARACTER";

   public static final String GET_BY_ID = "Equipment.GET_BY_ID";

   public static final String GET_FOR_CHARACTER_BY_SLOT = "Equipment.GET_FOR_CHARACTER_BY_SLOT";

   public static final String CHARACTER_ID = "characterId";

   public static final String ID = "id";

   public static final String SLOT = "slot";

   @Id
   @GeneratedValue(strategy = GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer characterId;

   @Column(nullable = false)
   private Short slot;

   @Column(nullable = false)
   private Integer itemId;

   @Column(nullable = false)
   private Integer strength;

   @Column(nullable = false)
   private Integer dexterity;

   @Column(nullable = false)
   private Integer intelligence;

   @Column(nullable = false)
   private Integer luck;

   @Column(nullable = false)
   private Integer hp;

   @Column(nullable = false)
   private Integer mp;

   @Column(nullable = false)
   private Integer weaponAttack;

   @Column(nullable = false)
   private Integer magicAttack;

   @Column(nullable = false)
   private Integer weaponDefense;

   @Column(nullable = false)
   private Integer magicDefense;

   @Column(nullable = false)
   private Integer accuracy;

   @Column(nullable = false)
   private Integer avoidability;

   @Column(nullable = false)
   private Integer hands;

   @Column(nullable = false)
   private Integer speed;

   @Column(nullable = false)
   private Integer jump;

   @Column(nullable = false)
   private Integer slots;

   public Integer getId() {
      return id;
   }

   public Integer getCharacterId() {
      return characterId;
   }

   public void setCharacterId(Integer characterId) {
      this.characterId = characterId;
   }

   public Short getSlot() {
      return slot;
   }

   public void setSlot(Short slot) {
      this.slot = slot;
   }

   public Integer getItemId() {
      return itemId;
   }

   public void setItemId(Integer itemId) {
      this.itemId = itemId;
   }

   public Integer getStrength() {
      return strength;
   }

   public void setStrength(Integer strength) {
      this.strength = strength;
   }

   public Integer getDexterity() {
      return dexterity;
   }

   public void setDexterity(Integer dexterity) {
      this.dexterity = dexterity;
   }

   public Integer getIntelligence() {
      return intelligence;
   }

   public void setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
   }

   public Integer getLuck() {
      return luck;
   }

   public void setLuck(Integer luck) {
      this.luck = luck;
   }

   public Integer getHp() {
      return hp;
   }

   public void setHp(Integer hp) {
      this.hp = hp;
   }

   public Integer getMp() {
      return mp;
   }

   public void setMp(Integer mp) {
      this.mp = mp;
   }

   public Integer getWeaponAttack() {
      return weaponAttack;
   }

   public void setWeaponAttack(Integer weaponAttack) {
      this.weaponAttack = weaponAttack;
   }

   public Integer getMagicAttack() {
      return magicAttack;
   }

   public void setMagicAttack(Integer magicAttack) {
      this.magicAttack = magicAttack;
   }

   public Integer getWeaponDefense() {
      return weaponDefense;
   }

   public void setWeaponDefense(Integer weaponDefense) {
      this.weaponDefense = weaponDefense;
   }

   public Integer getMagicDefense() {
      return magicDefense;
   }

   public void setMagicDefense(Integer magicDefense) {
      this.magicDefense = magicDefense;
   }

   public Integer getAccuracy() {
      return accuracy;
   }

   public void setAccuracy(Integer accuracy) {
      this.accuracy = accuracy;
   }

   public Integer getAvoidability() {
      return avoidability;
   }

   public void setAvoidability(Integer avoidability) {
      this.avoidability = avoidability;
   }

   public Integer getHands() {
      return hands;
   }

   public void setHands(Integer hands) {
      this.hands = hands;
   }

   public Integer getSpeed() {
      return speed;
   }

   public void setSpeed(Integer speed) {
      this.speed = speed;
   }

   public Integer getJump() {
      return jump;
   }

   public void setJump(Integer jump) {
      this.jump = jump;
   }

   public Integer getSlots() {
      return slots;
   }

   public void setSlots(Integer slots) {
      this.slots = slots;
   }
}
