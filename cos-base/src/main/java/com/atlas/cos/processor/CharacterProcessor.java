package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.Collections;
import java.util.Optional;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.ConfigurationRegistry;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.builder.StatisticChangeSummaryBuilder;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.event.StatUpdateType;
import com.atlas.cos.event.producer.CharacterCreatedProducer;
import com.atlas.cos.event.producer.CharacterLevelEventProducer;
import com.atlas.cos.event.producer.CharacterStatUpdateProducer;
import com.atlas.cos.event.producer.MapChangedProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.HpMpSummary;
import com.atlas.cos.model.Job;
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

      Connection.instance()
            .with(entityManager -> CharacterAdministrator.update(entityManager, characterId, summary.hp(),
                  summary.maxHp(), summary.mp(), summary.maxMp(), summary.strength(), summary.dexterity(),
                  summary.intelligence(), summary.luck(), summary.ap(), summary.level()));

      //levelUpGainSp();
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
