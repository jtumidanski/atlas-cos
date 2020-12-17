package com.atlas.cos.model;

import java.util.Arrays;
import java.util.Optional;

public enum InventoryType {
   EQUIP(1),
   USE(2),
   SETUP(3),
   ETC(4),
   CASH(5);

   final byte type;

   InventoryType(int type) {
      this.type = (byte) type;
   }

   public static Optional<InventoryType> getByType(byte type) {
      return Arrays.stream(InventoryType.values())
            .filter(inventoryType -> inventoryType.getType() == type)
            .findFirst();
   }

   public static Optional<InventoryType> getByName(String name) {
      return switch (name.toUpperCase()) {
         case "EQUIP" -> Optional.of(EQUIP);
         case "USE" -> Optional.of(USE);
         case "SETUP" -> Optional.of(SETUP);
         case "ETC" -> Optional.of(ETC);
         case "CASH" -> Optional.of(CASH);
         default -> Optional.empty();
      };
   }

   public byte getType() {
      return type;
   }
}
