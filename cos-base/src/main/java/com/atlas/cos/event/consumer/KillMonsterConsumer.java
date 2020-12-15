package com.atlas.cos.event.consumer;

import com.atlas.cos.processor.MonsterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterKilledEvent;

public class KillMonsterConsumer extends AbstractEventConsumer<MonsterKilledEvent> {
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
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_MONSTER_KILLED_EVENT);
   }
}
