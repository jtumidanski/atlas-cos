package com.atlas.cos.model;

public record Drop(int id, int itemId, int quantity, int meso, long dropTime, int dropType, int ownerId, boolean playerDrop) {

   public final boolean isFFADrop() {
      return dropType == 2 || dropType == 3 || hasExpiredOwnershipTime();
   }

   public final boolean hasExpiredOwnershipTime() {
      return System.currentTimeMillis() - dropTime >= 15 * 1000;
   }
}
