package com.atlas.cos.model;

public enum WeaponType {
   NOT_A_WEAPON(0),
   GENERAL1H_SWING(4.4),
   GENERAL1H_STAB(3.2),
   GENERAL2H_SWING(4.8),
   GENERAL2H_STAB(3.4),
   BOW(3.4),
   CLAW(3.6),
   CROSSBOW(3.6),
   DAGGER_THIEVES(3.6),
   DAGGER_OTHER(4),
   GUN(3.6),
   KNUCKLE(4.8),
   POLE_ARM_SWING(5.0),
   POLE_ARM_STAB(3.0),
   SPEAR_STAB(5.0),
   SPEAR_SWING(3.0),
   STAFF(3.6),
   SWORD1H(4.0),
   SWORD2H(4.6),
   WAND(3.6);

   private final double damageMultiplier;

   WeaponType(double damageMultiplier) {
      this.damageMultiplier = damageMultiplier;
   }

   public double getDamageMultiplier() {
      return damageMultiplier;
   }

   public static WeaponType getWeaponType(int itemId) {
      int cat = (itemId / 10000) % 100;
      WeaponType[] type = {WeaponType.SWORD1H, WeaponType.GENERAL1H_SWING, WeaponType.GENERAL1H_SWING,
            WeaponType.DAGGER_OTHER, WeaponType.NOT_A_WEAPON, WeaponType.NOT_A_WEAPON, WeaponType.NOT_A_WEAPON,
            WeaponType.WAND, WeaponType.STAFF, WeaponType.NOT_A_WEAPON, WeaponType.SWORD2H,
            WeaponType.GENERAL2H_SWING, WeaponType.GENERAL2H_SWING, WeaponType.SPEAR_STAB,
            WeaponType.POLE_ARM_SWING, WeaponType.BOW, WeaponType.CROSSBOW, WeaponType.CLAW,
            WeaponType.KNUCKLE, WeaponType.GUN};
      if (cat < 30 || cat > 49) {
         return WeaponType.NOT_A_WEAPON;
      }
      return type[cat - 30];
   }
}
