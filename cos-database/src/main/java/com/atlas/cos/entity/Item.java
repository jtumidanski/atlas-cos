package com.atlas.cos.entity;

import java.io.Serializable;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class Item implements Serializable {
   private static final long serialVersionUID = 1L;

   @Id
   @GeneratedValue(strategy = GenerationType.IDENTITY)
   private Integer id;

   @Column(nullable = false)
   private Integer characterId;

   @Column(nullable = false)
   private Byte inventoryType;

   @Column(nullable = false)
   private Integer itemId;

   @Column(nullable = false)
   private Integer quantity;

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

   public Integer getItemId() {
      return itemId;
   }

   public void setItemId(Integer itemId) {
      this.itemId = itemId;
   }

   public Byte getInventoryType() {
      return inventoryType;
   }

   public void setInventoryType(Byte inventoryType) {
      this.inventoryType = inventoryType;
   }

   public Integer getQuantity() {
      return quantity;
   }

   public void setQuantity(Integer quantity) {
      this.quantity = quantity;
   }

   public Short getSlot() {
      return slot;
   }

   public void setSlot(Short slot) {
      this.slot = slot;
   }
}
