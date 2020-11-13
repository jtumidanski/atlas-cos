package com.atlas.cos.entity;

import java.io.Serializable;
import java.sql.Timestamp;
import java.util.Calendar;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Index;
import javax.persistence.Table;

@Entity
@Table(name = "characters", indexes = {
      @Index(name = "accountId", columnList = "accountId"),
      @Index(name = "id_accountId_world", columnList = "id,accountId,world"),
      @Index(name = "id_accountId_name", columnList = "id,accountId,name")
})
public class Character implements Serializable {
   private static final long serialVersionUID = 1L;

   @Id
   @GeneratedValue(strategy = GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer accountId;

   @Column(nullable = false)
   private Integer world;

   @Column(nullable = false, length = 13)
   private String name;

   @Column(nullable = false)
   private Integer level = 1;

   @Column(nullable = false)
   private Integer exp = 0;

   @Column(nullable = false)
   private Integer gachaponExp = 0;

   @Column(nullable = false)
   private Integer strength = 12;

   @Column(nullable = false)
   private Integer dexterity = 5;

   @Column(nullable = false)
   private Integer luck = 4;

   @Column(nullable = false)
   private Integer intelligence = 4;

   @Column(nullable = false)
   private Integer hp = 50;

   @Column(nullable = false)
   private Integer mp = 5;

   @Column(nullable = false)
   private Integer maxHp = 50;

   @Column(nullable = false)
   private Integer maxMp = 5;

   @Column(nullable = false)
   private Integer meso = 0;

   @Column(nullable = false)
   private Integer hpMpUsed = 0;

   @Column(nullable = false)
   private Integer job = 0;

   @Column(nullable = false)
   private Integer skinColor = 0;

   @Column(nullable = false)
   private Byte gender = 0;

   @Column(nullable = false)
   private Integer fame = 0;

   @Column(nullable = false)
   private Integer hair = 0;

   @Column(nullable = false)
   private Integer face = 0;

   @Column(nullable = false)
   private Integer ap = 0;

   @Column(nullable = false)
   private String sp = "0,0,0,0,0,0,0,0,0,0";

   @Column(nullable = false)
   private Integer map = 0;

   @Column(nullable = false)
   private Integer spawnPoint = 0;

   @Column(nullable = false)
   private Integer gm = 0;

   @Column(nullable = false)
   private Timestamp createDate;

   @Column(nullable = false)
   private Timestamp lastLogin;

   public Character() {
      createDate = new Timestamp(Calendar.getInstance().getTimeInMillis());
      lastLogin = new Timestamp(0);
   }

   public Integer getId() {
      return id;
   }

   public Integer getAccountId() {
      return accountId;
   }

   public void setAccountId(Integer accountId) {
      this.accountId = accountId;
   }

   public Integer getWorld() {
      return world;
   }

   public void setWorld(Integer world) {
      this.world = world;
   }

   public String getName() {
      return name;
   }

   public void setName(String name) {
      this.name = name;
   }

   public Integer getLevel() {
      return level;
   }

   public void setLevel(Integer level) {
      this.level = level;
   }

   public Integer getExp() {
      return exp;
   }

   public void setExp(Integer exp) {
      this.exp = exp;
   }

   public Integer getGachaponExp() {
      return gachaponExp;
   }

   public void setGachaponExp(Integer gachaponExp) {
      this.gachaponExp = gachaponExp;
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

   public Integer getLuck() {
      return luck;
   }

   public void setLuck(Integer luck) {
      this.luck = luck;
   }

   public Integer getIntelligence() {
      return intelligence;
   }

   public void setIntelligence(Integer intelligence) {
      this.intelligence = intelligence;
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

   public Integer getMaxHp() {
      return maxHp;
   }

   public void setMaxHp(Integer maxHp) {
      this.maxHp = maxHp;
   }

   public Integer getMaxMp() {
      return maxMp;
   }

   public void setMaxMp(Integer maxMp) {
      this.maxMp = maxMp;
   }

   public Integer getMeso() {
      return meso;
   }

   public void setMeso(Integer meso) {
      this.meso = meso;
   }

   public Integer getHpMpUsed() {
      return hpMpUsed;
   }

   public void setHpMpUsed(Integer hpMpUsed) {
      this.hpMpUsed = hpMpUsed;
   }

   public Integer getJob() {
      return job;
   }

   public void setJob(Integer job) {
      this.job = job;
   }

   public Integer getSkinColor() {
      return skinColor;
   }

   public void setSkinColor(Integer skinColor) {
      this.skinColor = skinColor;
   }

   public Byte getGender() {
      return gender;
   }

   public void setGender(Byte gender) {
      this.gender = gender;
   }

   public Integer getFame() {
      return fame;
   }

   public void setFame(Integer fame) {
      this.fame = fame;
   }

   public Integer getHair() {
      return hair;
   }

   public void setHair(Integer hair) {
      this.hair = hair;
   }

   public Integer getFace() {
      return face;
   }

   public void setFace(Integer face) {
      this.face = face;
   }

   public Integer getAp() {
      return ap;
   }

   public void setAp(Integer ap) {
      this.ap = ap;
   }

   public String getSp() {
      return sp;
   }

   public void setSp(String sp) {
      this.sp = sp;
   }

   public Integer getMap() {
      return map;
   }

   public void setMap(Integer map) {
      this.map = map;
   }

   public Integer getSpawnPoint() {
      return spawnPoint;
   }

   public void setSpawnPoint(Integer spawnPoint) {
      this.spawnPoint = spawnPoint;
   }

   public Integer getGm() {
      return gm;
   }

   public void setGm(Integer gm) {
      this.gm = gm;
   }

   public Timestamp getCreateDate() {
      return createDate;
   }

   public void setCreateDate(Timestamp createDate) {
      this.createDate = createDate;
   }

   public Timestamp getLastLogin() {
      return lastLogin;
   }

   public void setLastLogin(Timestamp lastLogin) {
      this.lastLogin = lastLogin;
   }
}
