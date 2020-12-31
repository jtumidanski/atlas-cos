package com.atlas.cos.processor;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.ConfigurationRegistry;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.builder.StatisticChangeSummaryBuilder;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.database.provider.SkillProvider;
import com.atlas.cos.event.StatUpdateType;
import com.atlas.cos.event.producer.CharacterCreatedProducer;
import com.atlas.cos.event.producer.CharacterLevelEventProducer;
import com.atlas.cos.event.producer.CharacterStatUpdateProducer;
import com.atlas.cos.event.producer.MapChangedProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.EquipmentStatistics;
import com.atlas.cos.model.HpMpSummary;
import com.atlas.cos.model.Job;
import com.atlas.cos.model.SkillData;
import com.atlas.cos.model.StatisticChangeSummary;
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
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 130030000);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 914000000);
      return create(builder);
   }

   public static Optional<CharacterData> createBeginner(CharacterBuilder builder) {
      builder.setJobId(Job.BEGINNER.getId());
      builder.setMapId(10000);
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterBuilder builder) {
      builder.setJobId(Job.NOBLESSE.getId());
      builder.setMapId(130030000);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterBuilder builder) {
      builder.setJobId(Job.LEGEND.getId());
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

      weaponAttack += EquipmentProcessor.getEquipmentForCharacter(character.id())
            .map(EquipmentData::equipmentId)
            .map(EquipmentStatisticsProcessor::getEquipData)
            .flatMap(Optional::stream)
            .mapToInt(EquipmentStatistics::weaponAttack)
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

      strength += EquipmentProcessor.getEquipmentForCharacter(character.id())
            .map(EquipmentData::equipmentId)
            .map(EquipmentStatisticsProcessor::getEquipData)
            .flatMap(Optional::stream)
            .mapToInt(EquipmentStatistics::strength)
            .sum();

      return strength;
   }

   public static int getDexterity(CharacterData character) {
      int dexterity = character.dexterity();

      //TODO
      // apply Maple Warrior

      dexterity += EquipmentProcessor.getEquipmentForCharacter(character.id())
            .map(EquipmentData::equipmentId)
            .map(EquipmentStatisticsProcessor::getEquipData)
            .flatMap(Optional::stream)
            .mapToInt(EquipmentStatistics::dexterity)
            .sum();

      return dexterity;
   }

   public static int getLuck(CharacterData character) {
      int luck = character.luck();

      //TODO
      // apply Maple Warrior

      luck += EquipmentProcessor.getEquipmentForCharacter(character.id())
            .map(EquipmentData::equipmentId)
            .map(EquipmentStatisticsProcessor::getEquipData)
            .flatMap(Optional::stream)
            .mapToInt(EquipmentStatistics::luck)
            .sum();

      return luck;
   }

   public static int getIntelligence(CharacterData character) {
      int intelligence = character.intelligence();

      //TODO
      // apply Maple Warrior

      intelligence += EquipmentProcessor.getEquipmentForCharacter(character.id())
            .map(EquipmentData::equipmentId)
            .map(EquipmentStatisticsProcessor::getEquipData)
            .flatMap(Optional::stream)
            .mapToInt(EquipmentStatistics::intelligence)
            .sum();

      return intelligence;
   }

   public static boolean inMap(int characterId, int mapId) {
      return getById(characterId)
            .map(CharacterData::mapId)
            .filter(id -> id == mapId)
            .isPresent();
   }

   public static void increaseExperience(int characterId, int amount) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.increaseExperience(entityManager, characterId, amount));
      CharacterStatUpdateProducer
            .statsUpdated(characterId, Collections.singleton(StatUpdateType.EXPERIENCE));
   }

   public static void increaseLevel(int characterId) {
      CharacterData runningCharacter = Connection.instance()
            .element(entityManager -> CharacterProvider.getById(entityManager, characterId)).orElseThrow();

      boolean isBeginner = runningCharacter.isBeginnerJob();

      StatisticChangeSummaryBuilder builder = new StatisticChangeSummaryBuilder();
      if (ConfigurationRegistry.getInstance().getConfiguration().useAutoAssignStartersAp && isBeginner
            && runningCharacter.level() < 11) {

         if (runningCharacter.level() < 6) {
            builder.setStrength(5);
         } else {
            builder.setStrength(4);
            builder.setStrength(1);
         }
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
         builder.setAp(remainingAp);
      }

      HpMpSummary hpMpSummary = levelUpHealthAndManaPoints(runningCharacter, isBeginner);
      builder.setHp(hpMpSummary.hp())
            .setMp(hpMpSummary.mp())
            .setMaxHp(hpMpSummary.hp())
            .setMaxMp(hpMpSummary.mp())
            .setLevel(1);

      StatisticChangeSummary summary = builder.build();
      processStatisticChangeSummary(characterId, summary);

      //levelUpGainSp();

      CharacterStatUpdateProducer.statsUpdated(characterId, Arrays.asList(
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

   protected static void processStatisticChangeSummary(int characterId, StatisticChangeSummary summary) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.update(entityManager, characterId, summary.hp(),
                  summary.maxHp(), summary.mp(), summary.maxMp(), summary.strength(), summary.dexterity(),
                  summary.intelligence(), summary.luck(), summary.ap(), summary.level()));
   }

   public static void assignStrDexIntLuk(CharacterData character, int deltaStr, int deltaDex,
                                         int deltaInt, int deltaLuk) {
      int apUsed = apAssigned(deltaStr) + apAssigned(deltaDex) + apAssigned(deltaInt) + apAssigned(deltaLuk);
      if (apUsed > character.ap()) {
         return;
      }

      int newStr = character.strength() + deltaStr;
      int newDex = character.dexterity() + deltaDex;
      int newInt = character.intelligence() + deltaInt;
      int newLuk = character.luck() + deltaLuk;
      if (outOfRange(newStr, deltaStr)) {
         return;
      }

      if (outOfRange(newDex, deltaDex)) {
         return;
      }

      if (outOfRange(newInt, deltaInt)) {
         return;
      }

      if (outOfRange(newLuk, deltaLuk)) {
         return;
      }

      StatisticChangeSummary summary = new StatisticChangeSummaryBuilder()
            .setStrength(deltaStr)
            .setDexterity(deltaDex)
            .setIntelligence(deltaInt)
            .setLuck(deltaLuk)
            .setAp(-apUsed)
            .build();
      processStatisticChangeSummary(character.id(), summary);

      List<StatUpdateType> statUpdateTypes = new ArrayList<>();
      if (deltaStr != 0) {
         statUpdateTypes.add(StatUpdateType.STRENGTH);
      }
      if (deltaDex != 0) {
         statUpdateTypes.add(StatUpdateType.DEXTERITY);
      }
      if (deltaInt != 0) {
         statUpdateTypes.add(StatUpdateType.INTELLIGENCE);
      }
      if (deltaLuk != 0) {
         statUpdateTypes.add(StatUpdateType.LUCK);
      }
      if (apUsed != 0) {
         statUpdateTypes.add(StatUpdateType.AVAILABLE_AP);
      }

      CharacterStatUpdateProducer.statsUpdated(character.id(), statUpdateTypes);
   }

   protected static int calcHpChange(CharacterData character, boolean usedAPReset) {
      Job job = Job.getById(character.jobId()).orElseThrow();
      int MaxHP = 0;

      if (job.isA(Job.WARRIOR) || job.isA(Job.DAWN_WARRIOR_1)) {
         //         if (!usedAPReset) {
         //            int skillId = job.isA(MapleJob.DAWN_WARRIOR_1) ? DawnWarrior.MAX_HP_INCREASE : Warrior.IMPROVED_MAX_HP;
         //            MaxHP += SkillFactory.applyIfHasSkill(player, skillId, (skill, skillLevel) -> skill.getEffect(skillLevel).getY(), 0);
         //         }

         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 20;
            } else {
               MaxHP += Randomizer.rand(18, 22);
            }
         } else {
            MaxHP += 20;
         }
      } else if (job.isA(Job.ARAN1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 20;
            } else {
               MaxHP += Randomizer.rand(26, 30);
            }
         } else {
            MaxHP += 28;
         }
      } else if (job.isA(Job.MAGICIAN) || job.isA(Job.BLAZE_WIZARD_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 6;
            } else {
               MaxHP += Randomizer.rand(5, 9);
            }
         } else {
            MaxHP += 6;
         }
      } else if (job.isA(Job.THIEF) || job.isA(Job.NIGHT_WALKER_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 16;
            } else {
               MaxHP += Randomizer.rand(14, 18);
            }
         } else {
            MaxHP += 16;
         }
      } else if (job.isA(Job.BOWMAN) || job.isA(Job.WIND_ARCHER_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 16;
            } else {
               MaxHP += Randomizer.rand(14, 18);
            }
         } else {
            MaxHP += 16;
         }
      } else if (job.isA(Job.PIRATE) || job.isA(Job.THUNDER_BREAKER_1)) {
         //         if (!usedAPReset) {
         //            MaxHP += SkillFactory.applyIfHasSkill(player, job.isA(MapleJob.PIRATE) ? Brawler.IMPROVE_MAX_HP : ThunderBreaker.IMPROVE_MAX_HP, (skill, skillLevel) -> skill.getEffect(skillLevel).getY(), 0);
         //         }

         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (usedAPReset) {
               MaxHP += 18;
            } else {
               MaxHP += Randomizer.rand(16, 20);
            }
         } else {
            MaxHP += 18;
         }
      } else if (usedAPReset) {
         MaxHP += 8;
      } else {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            MaxHP += Randomizer.rand(8, 12);
         } else {
            MaxHP += 10;
         }
      }

      return MaxHP;
   }

   protected static int calcMpChange(CharacterData character, boolean usedAPReset) {
      Job job = Job.getById(character.jobId()).orElseThrow();
      int maxMp = 0;

      if (job.isA(Job.WARRIOR) || job.isA(Job.DAWN_WARRIOR_1) || job.isA(Job.ARAN1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(2, 4) + (character.intelligence() / 10));
            } else {
               maxMp += 2;
            }
         } else {
            maxMp += 3;
         }
      } else if (job.isA(Job.MAGICIAN) || job.isA(Job.BLAZE_WIZARD_1)) {
         //         if (!usedAPReset) {
         //            int skillId = job.isA(MapleJob.BLAZE_WIZARD_1) ? BlazeWizard.INCREASING_MAX_MP : Magician.IMPROVED_MAX_MP_INCREASE;
         //            MaxMP += SkillFactory.applyIfHasSkill(player, skillId, (skill, skillLevel) -> skill.getEffect(skillLevel).getY(), 0);
         //         }

         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(12, 16) + (character.intelligence() / 20));
            } else {
               maxMp += 18;
            }
         } else {
            maxMp += 18;
         }
      } else if (job.isA(Job.BOWMAN) || job.isA(Job.WIND_ARCHER_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(6, 8) + (character.intelligence() / 10));
            } else {
               maxMp += 10;
            }
         } else {
            maxMp += 10;
         }
      } else if (job.isA(Job.THIEF) || job.isA(Job.NIGHT_WALKER_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(6, 8) + (character.intelligence() / 10));
            } else {
               maxMp += 10;
            }
         } else {
            maxMp += 10;
         }
      } else if (job.isA(Job.PIRATE) || job.isA(Job.THUNDER_BREAKER_1)) {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(7, 9) + (character.intelligence() / 10));
            } else {
               maxMp += 14;
            }
         } else {
            maxMp += 14;
         }
      } else {
         if (ConfigurationRegistry.getInstance().getConfiguration().useRandomizeHpMpGain) {
            if (!usedAPReset) {
               maxMp += (Randomizer.rand(4, 6) + (character.intelligence() / 10));
            } else {
               maxMp += 6;
            }
         } else {
            maxMp += 6;
         }
      }

      return maxMp;
   }

   protected static boolean outOfRange(int newAp, int deltaAp) {
      return newAp < 4 && deltaAp != Short.MIN_VALUE || newAp > ConfigurationRegistry.getInstance().getConfiguration().maxAp;
   }

   protected static int apAssigned(int x) {
      return x != Short.MIN_VALUE ? x : 0;
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

   protected static HpMpSummary levelUpHealthAndManaPoints(CharacterData character, boolean isBeginner) {
      Job job = Job.getById(character.jobId()).orElseThrow();

      //      int improvingMaxHPSkillId = -1;
      //      int improvingMaxMPSkillId = -1;

      HpMpSummary summary = new HpMpSummary(0, 0);
      if (isBeginner) {
         summary = summary.increaseHp(Randomizer.rand(12, 16));
         summary = summary.increaseMp(Randomizer.rand(10, 12));
      } else if (job.isA(Job.WARRIOR) || job.isA(Job.DAWN_WARRIOR_1)) {
         //         improvingMaxHPSkillId = isCygnus() ? DawnWarrior.MAX_HP_INCREASE : Warrior.IMPROVED_MAX_HP;
         //         if (job.isA(MapleJob.CRUSADER)) {
         //            improvingMaxMPSkillId = 1210000;
         //         } else if (job.isA(MapleJob.DAWN_WARRIOR_2)) {
         //            improvingMaxMPSkillId = 11110000;
         //         }
         summary = summary.increaseHp(Randomizer.rand(24, 28));
         summary = summary.increaseMp(Randomizer.rand(4, 6));
      } else if (job.isA(Job.MAGICIAN) || job.isA(Job.BLAZE_WIZARD_1)) {
         //         improvingMaxMPSkillId = isCygnus() ? BlazeWizard.INCREASING_MAX_MP : Magician.IMPROVED_MAX_MP_INCREASE;
         summary = summary.increaseHp(Randomizer.rand(10, 14));
         summary = summary.increaseMp(Randomizer.rand(22, 24));
      } else if (job.isA(Job.BOWMAN) || job.isA(Job.THIEF) || (job.getId() > 1299 && job.getId() < 1500)) {
         summary = summary.increaseHp(Randomizer.rand(20, 24));
         summary = summary.increaseMp(Randomizer.rand(14, 16));
      } else if (job.isA(Job.GM)) {
         summary = summary.increaseHp(30000);
         summary = summary.increaseMp(30000);
      } else if (job.isA(Job.PIRATE) || job.isA(Job.THUNDER_BREAKER_1)) {
         //         improvingMaxHPSkillId = isCygnus() ? ThunderBreaker.IMPROVE_MAX_HP : Brawler.IMPROVE_MAX_HP;
         summary = summary.increaseHp(Randomizer.rand(22, 28));
         summary = summary.increaseMp(Randomizer.rand(18, 23));
      } else if (job.isA(Job.ARAN1)) {
         summary = summary.increaseHp(Randomizer.rand(44, 48));
         int aids = Randomizer.rand(4, 8);
         summary = summary.increaseMp((int) (aids + Math.floor(aids * 0.1)));
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
            summary = summary.increaseMp(getIntelligence(character) / 20);
         } else {
            summary = summary.increaseMp(getIntelligence(character) / 10);
         }
      }

      return summary;
   }

   public static void setExperience(int characterId, int experience) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.setExperience(entityManager, characterId, experience));
      CharacterStatUpdateProducer
            .statsUpdated(characterId, Collections.singleton(StatUpdateType.EXPERIENCE));
   }

   public static void gainExperience(int characterId, int gain) {
      getById(characterId).ifPresent(character -> gainExperience(characterId,
            character.level(), character.maxClassLevel(), character.experience(), gain));
   }

   protected static void gainExperience(int characterId, int level,
                                        int maxLevel, int experience, int gain) {
      if (level < maxLevel) {
         int toNextLevel = ExpTable.getExpNeededForLevel(level) - experience;
         if (toNextLevel <= gain) {
            setExperience(characterId, 0);
            CharacterLevelEventProducer.gainLevel(characterId);
            gainExperience(characterId, level + 1, maxLevel, 0, gain - toNextLevel);
         } else {
            increaseExperience(characterId, gain);
         }
      } else {
         setExperience(characterId, 0);
      }
   }

   public static void assignHpMp(CharacterData character, int deltaHp, int deltaMp) {
      HpMpSummary hpMpSummary = new HpMpSummary(0, 0);
      if (deltaHp > 0) {
         hpMpSummary = hpMpSummary.increaseHp(calcHpChange(character, false));
      }
      if (deltaMp > 0) {
         hpMpSummary = hpMpSummary.increaseMp(calcMpChange(character, false));
      }

      StatisticChangeSummary summary = new StatisticChangeSummaryBuilder()
            .setMaxHp(hpMpSummary.hp())
            .setMaxMp(hpMpSummary.mp())
            .build();
      processStatisticChangeSummary(character.id(), summary);

      List<StatUpdateType> statUpdateTypes = new ArrayList<>();
      if (hpMpSummary.hp() != 0) {
         statUpdateTypes.add(StatUpdateType.MAX_HP);
      }
      if (hpMpSummary.mp() != 0) {
         statUpdateTypes.add(StatUpdateType.MAX_MP);
      }

      CharacterStatUpdateProducer.statsUpdated(character.id(), statUpdateTypes);
   }

   public static void gainMeso(int characterId, int meso) {
      Connection.instance().with(entityManager -> CharacterAdministrator.increaseMeso(entityManager, characterId, meso));
      CharacterStatUpdateProducer.statsUpdated(characterId, Collections.singleton(StatUpdateType.MESO));
   }

   public static List<SkillData> getSkills(int characterId) {
      return Connection.instance().list(entityManager -> SkillProvider.getSkills(entityManager, characterId));
   }
}
