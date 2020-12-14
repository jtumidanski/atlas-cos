package com.atlas.cos.event.consumer;

import com.atlas.cos.processor.MonsterProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterKilledEvent;

public class KillMonsterConsumer implements SimpleEventHandler<MonsterKilledEvent> {
   @Override
   public void handle(Long key, MonsterKilledEvent event) {
      MonsterProcessor.getMonster(event.monsterId())
            .ifPresent(monster -> MonsterProcessor.distributeExperience(event.worldId(), event.channelId(), event.mapId(), monster,
                  event.damageEntries()));
      //         QuestProcessor.getInstance().raiseQuestMobCount(attacker, id());
   }

   @Override
   public Class<MonsterKilledEvent> getEventClass() {
      return MonsterKilledEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Character Service";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_MONSTER_KILLED_EVENT);
   }
}
