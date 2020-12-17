package com.atlas.cos.event;

public record InventoryModification(int mode, int inventoryType, int quantity, int oldPosition) {
   public InventoryModification(int mode, int inventoryType, int quantity) {
      this(mode, inventoryType, quantity, 0);
   }
}
