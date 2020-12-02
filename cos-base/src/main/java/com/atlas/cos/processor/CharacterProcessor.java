package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.model.CharacterData;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

import database.Connection;

public class CharacterProcessor {
   private static final Object lock = new Object();

   private static volatile CharacterProcessor instance;

   private final Producer<Long, CharacterCreatedEvent> producer;

   public static CharacterProcessor getInstance() {
      CharacterProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   private CharacterProcessor() {
      producer = KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS"));
   }

   public Optional<CharacterData> getByName(String name) {
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForName(entityManager, name))
            .stream().findFirst();
   }

   public Optional<CharacterData> getById(int id) {
      return Connection.instance()
            .element(entityManager -> CharacterProvider.getById(entityManager, id));
   }

   protected Optional<CharacterData> create(CharacterBuilder builder) {
      CharacterData original = builder.build();

      Optional<CharacterData> result = Connection.instance().element(entityManager ->
            CharacterAdministrator.create(entityManager,
                  original.accountId(), original.worldId(), original.name(), original.level(),
                  original.strength(), original.dexterity(), original.luck(), original.intelligence(),
                  original.maxHp(), original.maxMp(), original.jobId(), original.gender(), original.hair(),
                  original.face(), original.mapId())
      );

      result.ifPresent(this::notifyCharacterCreated);
      return result;
   }

   protected void notifyCharacterCreated(CharacterData data) {
      String topic = System.getenv(EventConstants.TOPIC_CHARACTER_CREATED_EVENT);
      long key = data.id();
      producer.send(new ProducerRecord<>(topic, key, new CharacterCreatedEvent(data.worldId(), data.id(), data.name())));
   }

   public Optional<CharacterData> createBeginner(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 10000);
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public Optional<CharacterData> createNoblesse(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 130030000);
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public Optional<CharacterData> createLegend(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 914000000);
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public Optional<CharacterData> createBeginner(CharacterBuilder builder) {
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      builder.setMapId(10000);
      return create(builder);
   }

   public Optional<CharacterData> createNoblesse(CharacterBuilder builder) {
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      builder.setMapId(130030000);
      return create(builder);
   }

   public Optional<CharacterData> createLegend(CharacterBuilder builder) {
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      builder.setMapId(914000000);
      return create(builder);
   }

   public void updateMap(int worldId, int channelId, int characterId, int mapId, int portalId) {
      Connection.instance().with(entityManager -> CharacterAdministrator.updateMap(entityManager, characterId, mapId));
      MapChangedProcessor.getInstance().notifyChange(worldId, channelId, characterId, mapId, portalId);
   }
}
