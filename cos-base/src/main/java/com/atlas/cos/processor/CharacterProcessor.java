package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.Collections;
import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.ConfigurationRegistry;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.event.StatUpdateType;
import com.atlas.cos.event.producer.CharacterCreatedProducer;
import com.atlas.cos.event.producer.CharacterLevelEventProducer;
import com.atlas.cos.event.producer.CharacterStatUpdateProducer;
import com.atlas.cos.event.producer.MapChangedProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.Job;
import com.atlas.cos.util.ExpTable;
import com.atlas.cos.util.Randomizer;

import database.Connection;

public final class CharacterProcessor {
   private CharacterProcessor() {
   }

   public static Optional<CharacterData> getByName(String name) {
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForName(entityManager, name))
            .stream().findFirst();
   }

   public static Optional<CharacterData> getById(int id) {
      return Connection.instance()
            .element(entityManager -> CharacterProvider.getById(entityManager, id));
   }

   protected static Optional<CharacterData> create(CharacterBuilder builder) {
      CharacterData original = builder.build();

      Optional<CharacterData> result = Connection.instance().element(entityManager ->
            CharacterAdministrator.create(entityManager,
                  original.accountId(), original.worldId(), original.name(), original.level(),
                  original.strength(), original.dexterity(), original.luck(), original.intelligence(),
                  original.maxHp(), original.maxMp(), original.jobId(), original.gender(), original.hair(),
                  original.face(), original.mapId())
      );

      result.ifPresent(CharacterCreatedProducer::notifyCharacterCreated);
      return result;
   }

   public static Optional<CharacterData> createBeginner(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 10000);
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 130030000);
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 914000000);
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createBeginner(CharacterBuilder builder) {
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      builder.setMapId(10000);
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterBuilder builder) {
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      builder.setMapId(130030000);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterBuilder builder) {
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      builder.setMapId(914000000);
      return create(builder);
   }

   public static void updateMap(int worldId, int channelId, int characterId, int mapId, int portalId) {
      Connection.instance().with(entityManager -> CharacterAdministrator.updateMap(entityManager, characterId, mapId));

      PortalProcessor.getMapPortalById(mapId, portalId)
            .ifPresent(portal -> CharacterTemporalRegistry.getInstance()
                  .updatePosition(characterId, portal.x(), portal.y()));

      MapChangedProducer.notifyChange(worldId, channelId, characterId, mapId, portalId);
   }

   public static void updateSpawnPoint(int characterId, int newSpawnPoint) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.updateSpawnPoint(entityManager, characterId, newSpawnPoint));
   }

   public static int getWeaponAttack(CharacterData character) {
      int weaponAttack = 0;

      weaponAttack += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::weaponAttack)
            .sum();

      //TODO
      // apply Aran Combo
      // apply ThunderBreaker Marauder energy charge
      // apply Marksman Boost or Bowmaster Expert
      // apply weapon attack buffs
      // apply blessing
      // apply active projectile

      return weaponAttack;
   }

   public static int getStrength(CharacterData character) {
      int strength = character.strength();

      //TODO
      // apply Maple Warrior

      strength += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::strength)
            .sum();

      return strength;
   }

   public static int getDexterity(CharacterData character) {
      int dexterity = character.dexterity();

      //TODO
      // apply Maple Warrior

      dexterity += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::dexterity)
            .sum();

      return dexterity;
   }

   public static int getLuck(CharacterData character) {
      int luck = character.luck();

      //TODO
      // apply Maple Warrior

      luck += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::luck)
            .sum();

      return luck;
   }

   public static int getIntelligence(CharacterData character) {
      int intelligence = character.intelligence();

      //TODO
      // apply Maple Warrior

      intelligence += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::intelligence)
            .sum();

      return intelligence;
   }

   public static boolean inMap(int characterId, int mapId) {
      return getById(characterId)
            .map(CharacterData::mapId)
            .filter(id -> id == mapId)
            .isPresent();
   }

   public static void increaseExperience(int worldId, int channelId, int mapId, int characterId, int amount) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.increaseExperience(entityManager, characterId, amount));
      CharacterStatUpdateProducer
            .statsUpdated(worldId, channelId, mapId, characterId, Collections.singleton(StatUpdateType.EXPERIENCE));
   }

   public static void increaseLevel(int worldId, int channelId, int mapId, int characterId) {
      Connection.instance().with(entityManager ->
            QueryAdministratorUtil.inTransaction(entityManager, transaction -> {
               CharacterData runningCharacter = CharacterProvider.getById(transaction, characterId).orElseThrow();

               boolean isBeginner = runningCharacter.isBeginnerJob();
               if (ConfigurationRegistry.getInstance().getConfiguration().useAutoAssignStartersAp && isBeginner
                     && runningCharacter.level() < 11) {
                  gainAp(transaction, runningCharacter, 5, true);

                  int str = 0, dex = 0;
                  if (runningCharacter.level() < 6) {
                     str += 5;
                  } else {
                     str += 4;
                     dex += 1;
                  }

                  assignStrDexIntLuk(transaction, runningCharacter, str, dex, 0, 0);
               } else {
                  int remainingAp = 5;

                  if (runningCharacter.isCygnus()) {
                     if (runningCharacter.level() > 10) {
                        if (runningCharacter.level() <= 17) {
                           remainingAp += 2;
                        } else if (runningCharacter.level() < 77) {
                           remainingAp++;
                        }
                     }
                  }

                  gainAp(transaction, runningCharacter, remainingAp, true);
               }

               levelUpHealthAndManaPoints(transaction, runningCharacter, isBeginner);

               CharacterAdministrator.increaseLevel(transaction, characterId);

               //levelUpGainSp();
            }));
      //
      //      levelUpHealthAndManaPoints(isBeginner);

      //      levelUpGainSp();
      //
      //      effLock.lock();
      //      statWriteLock.lock();
      //      try {
      //         recalculateLocalStats();
      //         changeHpMp(localMaxHp, localMaxMp, true);
      //
      //         List<Pair<MapleStat, Integer>> statIncreases = new ArrayList<>(10);
      //         statIncreases.add(new Pair<>(MapleStat.AVAILABLE_AP, remainingAp));
      //         statIncreases.add(new Pair<>(MapleStat.AVAILABLE_SP, remainingSp[GameConstants.getSkillBook(job.getId())]));
      //         statIncreases.add(new Pair<>(MapleStat.HP, hp));
      //         statIncreases.add(new Pair<>(MapleStat.MP, mp));
      //         statIncreases.add(new Pair<>(MapleStat.EXP, exp.get()));
      //         statIncreases.add(new Pair<>(MapleStat.LEVEL, level));
      //         statIncreases.add(new Pair<>(MapleStat.MAX_HP, clientMaxHp));
      //         statIncreases.add(new Pair<>(MapleStat.MAX_MP, clientMaxMp));
      //         statIncreases.add(new Pair<>(MapleStat.STR, str));
      //         statIncreases.add(new Pair<>(MapleStat.DEX, dex));
      //
      //         PacketCreator.announce(client, new UpdatePlayerStats(statIncreases, true, this));
      //      } finally {
      //         statWriteLock.unlock();
      //         effLock.unlock();
      //      }

      CharacterStatUpdateProducer.statsUpdated(worldId, channelId, mapId, characterId, Arrays.asList(
            StatUpdateType.EXPERIENCE,
            StatUpdateType.LEVEL,
            StatUpdateType.AVAILABLE_AP,
            StatUpdateType.HP,
            StatUpdateType.MP,
            StatUpdateType.MAX_HP,
            StatUpdateType.MAX_MP,
            StatUpdateType.STRENGTH,
            StatUpdateType.DEXTERITY,
            StatUpdateType.LUCK,
            StatUpdateType.INTELLIGENCE
      ));
   }

   //   protected static void levelUpGainSp(Job job) {
   //      if (job.getJobBranch() == 0) {
   //         return;
   //      }
   //
   //      int spGain = 3;
   //      if (ConfigurationRegistry.getInstance().getConfiguration().useEnforceJobSpRange && !job.hasSpTable()) {
   //         spGain = getSpGain(spGain, job);
   //      }
   //
   //      if (spGain > 0) {
   //         gainSp(spGain, GameConstants.getSkillBook(job.getId()), true);
   //      }
   //   }

   //   protected static int getSpGain(int spGain, Job job) {
   //      int curSp = getUsedSp(job) + getJobRemainingSp(job);
   //      return getSpGain(spGain, curSp, job);
   //   }
   //
   //   protected static int getUsedSp(Job job) {
   //      int jobId = job.getId();
   //      int spUsed = 0;
   //
   ////      for (Map.Entry<Skill, SkillEntry> s : this.getSkills().entrySet()) {
   ////         Skill skill = s.getKey();
   ////         if (GameConstants.isInJobTree(skill.getId(), jobId) && !skill.isBeginnerSkill()) {
   ////            spUsed += s.getValue().skillLevel();
   ////         }
   ////      }
   //
   //      return spUsed;
   //   }
   //
   //   protected static int getJobRemainingSp(Job job) {
   //      int skillBook = GameConstants.getSkillBook(job.getId());
   //
   //      int ret = 0;
   //      for (int i = 0; i <= skillBook; i++) {
   //         ret += this.getRemainingSp(i);
   //      }
   //
   //      return ret;
   //   }

   //   protected static int getSpGain(int spGain, int curSp, Job job) {
   //      int maxSp = getJobMaxSp(job);
   //      spGain = Math.min(spGain, maxSp - curSp);
   //      return spGain;
   //   }

   //   protected static int getJobMaxSp(Job job) {
   //      int jobRange = job.getJobUpgradeLevelRange();
   //      return getJobLevelSp(jobRange, job);
   //   }

   //   protected static int getJobLevelSp(int level, Job job) {
   //      if (Job.getJobStyleInternal(job.getId(), (byte) 0x40) == Job.MAGICIAN) {
   //         level += 2;  // starts earlier, level 8
   //      }
   //
   //      return 3 * level + job.getChangeJobSpUpgrade();
   //   }

   protected static void gainAp(EntityManager entityManager, CharacterData character, int deltaAp, boolean silent) {
      changeRemainingAp(entityManager, character, Math.max(0, character.ap() + deltaAp), silent);
   }

   protected static boolean assignStrDexIntLuk(EntityManager entityManager, CharacterData character, int deltaStr, int deltaDex,
                                               int deltaInt, int deltaLuk) {
      int apUsed = apAssigned(deltaStr) + apAssigned(deltaDex) + apAssigned(deltaInt) + apAssigned(deltaLuk);
      if (apUsed > character.ap()) {
         return false;
      }

      int newStr = character.strength() + deltaStr;
      int newDex = character.dexterity() + deltaDex;
      int newInt = character.intelligence() + deltaInt;
      int newLuk = character.luck() + deltaLuk;

      if (outOfRange(newStr, deltaStr)) {
         return false;
      }

      if (outOfRange(newDex, deltaDex)) {
         return false;
      }

      if (outOfRange(newInt, deltaInt)) {
         return false;
      }

      if (outOfRange(newLuk, deltaLuk)) {
         return false;
      }

      int newAp = character.ap() - apUsed;
      updateStrDexIntLuk(entityManager, character.id(), newStr, newDex, newInt, newLuk, newAp);
      return true;
   }

   protected static void levelUpHealthAndManaPoints(EntityManager entityManager, CharacterData character, boolean isBeginner) {
      Job job = Job.getById(character.jobId()).orElseThrow();

      //      int improvingMaxHPSkillId = -1;
      //      int improvingMaxMPSkillId = -1;

      int addHp = 0, addMp = 0;
      if (isBeginner) {
         addHp += Randomizer.rand(12, 16);
         addMp += Randomizer.rand(10, 12);
      } else if (job.isA(Job.WARRIOR) || job.isA(Job.DAWN_WARRIOR_1)) {
         //         improvingMaxHPSkillId = isCygnus() ? DawnWarrior.MAX_HP_INCREASE : Warrior.IMPROVED_MAX_HP;
         //         if (job.isA(MapleJob.CRUSADER)) {
         //            improvingMaxMPSkillId = 1210000;
         //         } else if (job.isA(MapleJob.DAWN_WARRIOR_2)) {
         //            improvingMaxMPSkillId = 11110000;
         //         }
         addHp += Randomizer.rand(24, 28);
         addMp += Randomizer.rand(4, 6);
      } else if (job.isA(Job.MAGICIAN) || job.isA(Job.BLAZE_WIZARD_1)) {
         //         improvingMaxMPSkillId = isCygnus() ? BlazeWizard.INCREASING_MAX_MP : Magician.IMPROVED_MAX_MP_INCREASE;
         addHp += Randomizer.rand(10, 14);
         addMp += Randomizer.rand(22, 24);
      } else if (job.isA(Job.BOWMAN) || job.isA(Job.THIEF) || (job.getId() > 1299 && job.getId() < 1500)) {
         addHp += Randomizer.rand(20, 24);
         addMp += Randomizer.rand(14, 16);
      } else if (job.isA(Job.GM)) {
         addHp += 30000;
         addMp += 30000;
      } else if (job.isA(Job.PIRATE) || job.isA(Job.THUNDER_BREAKER_1)) {
         //         improvingMaxHPSkillId = isCygnus() ? ThunderBreaker.IMPROVE_MAX_HP : Brawler.IMPROVE_MAX_HP;
         addHp += Randomizer.rand(22, 28);
         addMp += Randomizer.rand(18, 23);
      } else if (job.isA(Job.ARAN1)) {
         addHp += Randomizer.rand(44, 48);
         int aids = Randomizer.rand(4, 8);
         addMp += aids + Math.floor(aids * 0.1);
      }

      //      Optional<Skill> improvingMaxHP = SkillFactory.getSkill(improvingMaxHPSkillId);
      //      if (improvingMaxHP.isPresent()) {
      //         int improvingMaxHPLevel = getSkillLevel(improvingMaxHP.get());
      //         if (improvingMaxHPLevel > 0 && (job.isA(MapleJob.WARRIOR) || job.isA(MapleJob.PIRATE) || job.isA(MapleJob.DAWN_WARRIOR_1)
      //               || job.isA(MapleJob.THUNDER_BREAKER_1))) {
      //            addHp += improvingMaxHP.get().getEffect(improvingMaxHPLevel).getX();
      //         }
      //      }

      //      Optional<Skill> improvingMaxMP = SkillFactory.getSkill(improvingMaxMPSkillId);
      //      if (improvingMaxMP.isPresent()) {
      //         int improvingMaxMPLevel = getSkillLevel(improvingMaxMP.get());
      //         if (improvingMaxMPLevel > 0 && (job.isA(MapleJob.MAGICIAN) || job.isA(MapleJob.CRUSADER) || job
      //               .isA(MapleJob.BLAZE_WIZARD_1))) {
      //            addMp += improvingMaxMP.get().getEffect(improvingMaxMPLevel).getX();
      //         }
      //      }

      if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
         if (job.getJobStyle(character.strength(), character.dexterity()) == Job.MAGICIAN) {
            addMp += getIntelligence(character) / 20;
         } else {
            addMp += getIntelligence(character) / 10;
         }
      }

      addMaxMPMaxHP(entityManager, character, addHp, addMp);
   }

   protected static void addMaxMPMaxHP(EntityManager entityManager, CharacterData character, int hpDelta, int mpDelta) {
      changeHpMpPool(entityManager, character.id(), Short.MIN_VALUE, Short.MIN_VALUE, character.maxHp() + hpDelta,
            character.maxMp() + mpDelta);
   }

   protected static void changeHpMpPool(EntityManager entityManager, int characterId, int hp, int mp, int maxHp, int maxMp) {
      long hpMpPool = calcStatPoolLong(hp, mp, maxHp, maxMp);
      changeStatPool(entityManager, characterId, hpMpPool, null, null, -1);
   }

   protected static void updateStrDexIntLuk(EntityManager entityManager, int characterId, int str, int dex, int int_, int luk,
                                            int remainingAp) {
      changeStrDexIntLuk(entityManager, characterId, str, dex, int_, luk, remainingAp);
   }

   protected static int apAssigned(int x) {
      return x != Short.MIN_VALUE ? x : 0;
   }

   protected static boolean outOfRange(int newAp, int deltaAp) {
      return newAp < 4 && deltaAp != Short.MIN_VALUE || newAp > ConfigurationRegistry.getInstance().getConfiguration().maxAp;
   }

   protected static void changeRemainingAp(EntityManager entityManager, CharacterData character, int x, boolean silent) {
      changeStrDexIntLuk(entityManager, character.id(), character.strength(), character.dexterity(), character.intelligence(),
            character.luck(), x
      );
   }

   protected static void changeStrDexIntLuk(EntityManager entityManager, int characterId, int str, int dex, int int_, int luk,
                                            int remainingAp) {
      long strDexIntLuk = calcStatPoolLong(str, dex, int_, luk);
      changeStatPool(entityManager, characterId, null, strDexIntLuk, null, remainingAp);
   }

   protected static void changeStatPool(EntityManager entityManager, int characterId, Long hpMpPool, Long strDexIntLuk, Long newSp
         , int newAp) {
      if (hpMpPool != null) {
         short newHp = (short) (hpMpPool >> 48);
         short newMp = (short) (hpMpPool >> 32);
         short newMaxHp = (short) (hpMpPool >> 16);
         short newMaxMp = hpMpPool.shortValue();

         if (newMaxHp != Short.MIN_VALUE) {
            if (newMaxHp < 50) {
               newMaxHp = 50;
            }
            int finalNewMaxHp = newMaxHp;
            CharacterAdministrator.setMaxHp(entityManager, characterId, finalNewMaxHp);
         }

         if (newHp != Short.MIN_VALUE) {
            CharacterAdministrator.setHp(entityManager, characterId, newHp);
         }

         if (newMaxMp != Short.MIN_VALUE) {
            if (newMaxMp < 5) {
               newMaxMp = 5;
            }
            int finalNewMaxMp = newMaxMp;
            CharacterAdministrator.setMaxMp(entityManager, characterId, finalNewMaxMp);
         }

         if (newMp != Short.MIN_VALUE) {
            CharacterAdministrator.setMp(entityManager, characterId, newMp);
         }
      }

      if (strDexIntLuk != null) {
         short newStr = (short) (strDexIntLuk >> 48);
         short newDex = (short) (strDexIntLuk >> 32);
         short newInt = (short) (strDexIntLuk >> 16);
         short newLuk = strDexIntLuk.shortValue();

         if (newStr >= 4) {
            CharacterAdministrator.setStrength(entityManager, characterId, newStr);
         }

         if (newDex >= 4) {
            CharacterAdministrator.setDexterity(entityManager, characterId, newDex);
         }

         if (newInt >= 4) {
            CharacterAdministrator.setIntelligence(entityManager, characterId, newInt);
         }

         if (newLuk >= 4) {
            CharacterAdministrator.setLuck(entityManager, characterId, newLuk);
         }

         if (newAp >= 0) {
            CharacterAdministrator.setAp(entityManager, characterId, newAp);
         }
      }

      if (newSp != null) {
         short sp = (short) (newSp >> 16);
         short skillBook = newSp.shortValue();

         setRemainingSp(sp, skillBook);
      }
   }

   protected static long calcStatPoolLong(int v1, int v2, int v3, int v4) {
      long ret = 0;
      ret |= calcStatPoolNode(v1, 48);
      ret |= calcStatPoolNode(v2, 32);
      ret |= calcStatPoolNode(v3, 16);
      ret |= calcStatPoolNode(v4, 0);
      return ret;
   }

   protected static void setRemainingSp(int remainingSp, int skillBook) {
      //this.remainingSp[skillBook] = remainingSp;
   }

   protected static long calcStatPoolNode(long v, int displacement) {
      if (v > Short.MAX_VALUE) {
         v = Short.MAX_VALUE;
      } else if (v < Short.MIN_VALUE) {
         v = Short.MIN_VALUE;
      }
      return ((v & 0x0FFFF) << displacement);
   }

   public static void setExperience(int worldId, int channelId, int mapId, int characterId, int experience) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.setExperience(entityManager, characterId, experience));
      CharacterStatUpdateProducer
            .statsUpdated(worldId, channelId, mapId, characterId, Collections.singleton(StatUpdateType.EXPERIENCE));
   }

   public static void gainExperience(int worldId, int channelId, int mapId, int characterId, int gain) {
      getById(characterId).ifPresent(character -> gainExperience(worldId, channelId, mapId, characterId,
            character.level(), character.maxClassLevel(), character.experience(), gain));
   }

   protected static void gainExperience(int worldId, int channelId, int mapId, int characterId, int level,
                                        int maxLevel, int experience, int gain) {
      if (level < maxLevel) {
         int toNextLevel = ExpTable.getExpNeededForLevel(level) - experience;
         if (toNextLevel <= gain) {
            setExperience(worldId, channelId, mapId, characterId, 0);
            CharacterLevelEventProducer.gainLevel(worldId, channelId, mapId, characterId);
            gainExperience(worldId, channelId, mapId, characterId, level + 1, maxLevel, 0, gain - toNextLevel);
         } else {
            increaseExperience(worldId, channelId, mapId, characterId, gain);
         }
      } else {
         setExperience(worldId, channelId, mapId, characterId, 0);
      }
   }
}
