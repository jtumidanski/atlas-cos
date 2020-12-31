package com.atlas.cos.entity;

import java.io.Serializable;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Index;
import javax.persistence.Table;

@Entity
@Table(indexes = {
      @Index(name = "skillPair", columnList = "skillId,characterId", unique = true)
})
public class Skill implements Serializable {
   private static final long serialVersionUID = 1L;

   @Id
   @GeneratedValue(strategy=GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer skillId;

   @Column(nullable = false)
   private Integer characterId;

   @Column(nullable = false)
   private Integer skillLevel;

   @Column(nullable = false)
   private Integer masterLevel;

   @Column(nullable = false)
   private Long expiration = -1L;

   public Skill() {
   }

   public Integer getId() {
      return id;
   }

   public void setId(Integer id) {
      this.id = id;
   }

   public Integer getSkillId() {
      return skillId;
   }

   public void setSkillId(Integer skillId) {
      this.skillId = skillId;
   }

   public Integer getCharacterId() {
      return characterId;
   }

   public void setCharacterId(Integer characterId) {
      this.characterId = characterId;
   }

   public Integer getSkillLevel() {
      return skillLevel;
   }

   public void setSkillLevel(Integer skillLevel) {
      this.skillLevel = skillLevel;
   }

   public Integer getMasterLevel() {
      return masterLevel;
   }

   public void setMasterLevel(Integer masterLevel) {
      this.masterLevel = masterLevel;
   }

   public Long getExpiration() {
      return expiration;
   }

   public void setExpiration(Long expiration) {
      this.expiration = expiration;
   }
}
