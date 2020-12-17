package com.atlas.cos.event;

public record InventoryModification(int mode, int itemId, int inventoryType, int quantity, int position, int oldPosition) {
   public InventoryModification(int mode, int itemId, int inventoryType, int quantity, int position) {
      this(mode, itemId, inventoryType, quantity, position, 0);
   }
}
