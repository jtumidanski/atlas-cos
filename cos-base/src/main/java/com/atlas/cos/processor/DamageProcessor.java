package com.atlas.cos.processor;

import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.Job;
import com.atlas.cos.model.WeaponType;

public final class DamageProcessor {
   private DamageProcessor() {
   }

   public static int getMaxBaseDamage(int characterId) {
      return CharacterProcessor.getById(characterId)
            .map(DamageProcessor::getMaxBaseDamage)
            .orElse(1);
   }

   protected static int getMaxBaseDamage(CharacterData character) {
      int weaponAttack = CharacterProcessor.getWeaponAttack(character);

      return ItemProcessor.getEquipedItemForCharacterBySlot(character.id(), (short) -11)
            .map(weapon -> getMaxBaseDamage(character, weaponAttack, WeaponType.getWeaponType(weapon.itemId())))
            .orElse(getMaxBaseDamageNoWeapon(character));
   }

   protected static int getMaxBaseDamage(CharacterData character, int weaponAttack, WeaponType weaponType) {
      int mainStat;
      int secondaryStat;
      WeaponType workingWeaponType = weaponType;

      Job job = Job.getById(character.jobId()).orElseThrow();
      if (job.isA(Job.THIEF) && workingWeaponType.equals(WeaponType.DAGGER_OTHER)) {
         workingWeaponType = WeaponType.DAGGER_THIEVES;
      }

      if (workingWeaponType.equals(WeaponType.BOW) || workingWeaponType.equals(WeaponType.CROSSBOW) || workingWeaponType
            .equals(WeaponType.GUN)) {
         mainStat = CharacterProcessor.getDexterity(character);
         secondaryStat = CharacterProcessor.getStrength(character);
      } else if (workingWeaponType.equals(WeaponType.CLAW) || workingWeaponType.equals(WeaponType.DAGGER_THIEVES)) {
         mainStat = CharacterProcessor.getLuck(character);
         secondaryStat = CharacterProcessor.getDexterity(character) + CharacterProcessor.getStrength(character);
      } else {
         mainStat = CharacterProcessor.getStrength(character);
         secondaryStat = CharacterProcessor.getDexterity(character);
      }
      return (int) Math.ceil(((weaponType.getDamageMultiplier() * mainStat + secondaryStat) / 100.0) * weaponAttack);
   }

   protected static Integer getMaxBaseDamageNoWeapon(CharacterData character) {
      Job job = Job.getById(character.jobId()).orElseThrow();
      if (job.isA(Job.PIRATE) || job.isA(Job.THUNDER_BREAKER_1)) {
         double weaponMultiplier = 3;
         if (job.getId() % 100 != 0) {
            weaponMultiplier = 4.2;
         }
         int attack = (int) Math.min(Math.floor((2 * character.level() + 31) / 3.0), 31);
         int strength = CharacterProcessor.getStrength(character);
         int dexterity = CharacterProcessor.getDexterity(character);
         return (int) Math.ceil((strength * weaponMultiplier + dexterity) * attack / 100.0);
      }
      return 1;
   }
}
