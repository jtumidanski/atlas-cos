package com.atlas.cos.entity;

import java.io.Serializable;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.NamedQueries;
import javax.persistence.NamedQuery;

@Entity
@NamedQueries({
      @NamedQuery(name = Equipment.GET_BY_ID, query = "SELECT e FROM Equipment e WHERE e.id = :id"),
      @NamedQuery(name = Equipment.GET_FOR_CHARACTER, query = "SELECT e FROM Equipment e WHERE e.characterId = :characterId"),
      @NamedQuery(name = Equipment.GET_FOR_EQUIPMENT, query = "SELECT e FROM Equipment e WHERE e.equipmentId = :equipmentId"),
      @NamedQuery(name = Equipment.GET_FOR_CHARACTER_BY_SLOT, query = "SELECT e FROM Equipment e WHERE e.characterId = "
            + ":characterId AND e.slot = :slot")
})
public class Equipment implements Serializable {
   private static final long serialVersionUID = 1L;

   public static final String GET_FOR_CHARACTER = "Equipment.GET_FOR_CHARACTER";

   public static final String GET_FOR_EQUIPMENT = "Equipment.GET_FOR_EQUIPMENT";

   public static final String GET_BY_ID = "Equipment.GET_BY_ID";

   public static final String GET_FOR_CHARACTER_BY_SLOT = "Equipment.GET_FOR_CHARACTER_BY_SLOT";

   public static final String CHARACTER_ID = "characterId";

   public static final String EQUIPMENT_ID = "equipmentId";

   public static final String ID = "id";

   public static final String SLOT = "slot";

   @Id
   @GeneratedValue(strategy = GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer characterId;

   @Column(nullable = false)
   private Integer equipmentId;

   @Column(nullable = false)
   private Short slot;

   public Integer getId() {
      return id;
   }

   public Integer getCharacterId() {
      return characterId;
   }

   public void setCharacterId(Integer characterId) {
      this.characterId = characterId;
   }

   public Integer getEquipmentId() {
      return equipmentId;
   }

   public void setEquipmentId(Integer equipmentId) {
      this.equipmentId = equipmentId;
   }

   public Short getSlot() {
      return slot;
   }

   public void setSlot(Short slot) {
      this.slot = slot;
   }
}
