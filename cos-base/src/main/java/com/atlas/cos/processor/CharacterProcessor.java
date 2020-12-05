package com.atlas.cos.processor;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.event.producer.CharacterCreatedProducer;
import com.atlas.cos.event.producer.MapChangedProducer;
import com.atlas.cos.model.CharacterData;
import database.Connection;

import java.util.Optional;

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
      Connection.instance().with(entityManager -> CharacterAdministrator.updateSpawnPoint(entityManager, characterId, newSpawnPoint));
   }
}
