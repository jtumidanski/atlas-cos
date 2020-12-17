package com.atlas.cos.event;

public record InventoryModification(int mode, int inventoryType, int quantity, int position, int oldPosition) {
   public InventoryModification(int mode, int inventoryType, int quantity, int position) {
      this(mode, inventoryType, quantity, position, 0);
   }
}
