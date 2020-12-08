package com.atlas.cos.entity;

import java.io.Serializable;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.EnumType;
import javax.persistence.Enumerated;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Table;

@Entity
@Table(name = "savedlocations")
public class SavedLocation implements Serializable {
   private static final long serialVersionUID = 1L;

   @Id
   @GeneratedValue(strategy=GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer characterId;

   @Enumerated(EnumType.STRING)
   @Column(nullable = false)
   private LocationType locationType;

   @Column(nullable = false)
   private Integer map;

   @Column(nullable = false)
   private Integer portal;

   public SavedLocation() {
   }

   public Integer getId() {
      return id;
   }

   public void setId(Integer id) {
      this.id = id;
   }

   public Integer getCharacterId() {
      return characterId;
   }

   public void setCharacterId(Integer characterId) {
      this.characterId = characterId;
   }

   public LocationType getLocationType() {
      return locationType;
   }

   public void setLocationType(LocationType locationType) {
      this.locationType = locationType;
   }

   public Integer getMap() {
      return map;
   }

   public void setMap(Integer map) {
      this.map = map;
   }

   public Integer getPortal() {
      return portal;
   }

   public void setPortal(Integer portal) {
      this.portal = portal;
   }
}
