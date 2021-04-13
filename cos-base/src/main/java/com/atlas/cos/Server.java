package com.atlas.cos;

import java.net.URI;

import com.atlas.cos.constant.RestConstants;
import com.atlas.cos.event.consumer.AdjustHealthConsumer;
import com.atlas.cos.event.consumer.AdjustManaConsumer;
import com.atlas.cos.event.consumer.AdjustMesoConsumer;
import com.atlas.cos.event.consumer.AssignApConsumer;
import com.atlas.cos.event.consumer.AssignSpConsumer;
import com.atlas.cos.event.consumer.ChangeMapCommandConsumer;
import com.atlas.cos.event.consumer.CharacterExperienceConsumer;
import com.atlas.cos.event.consumer.CharacterLevelConsumer;
import com.atlas.cos.event.consumer.CharacterMovementConsumer;
import com.atlas.cos.event.consumer.CharacterStatusConsumer;
import com.atlas.cos.event.consumer.DropReservationEventConsumer;
import com.atlas.cos.event.consumer.GainMesoConsumer;
import com.atlas.cos.event.consumer.KillMonsterConsumer;
import com.atlas.kafka.consumer.SimpleEventConsumerBuilder;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.UriBuilder;

import database.PersistenceManager;

public class Server {
   public static void main(String[] args) {
      PersistenceManager.construct("atlas-cos");

      SimpleEventConsumerBuilder.builder()
            .addConsumer(new ChangeMapCommandConsumer())
            .addConsumer(new CharacterStatusConsumer())
            .addConsumer(new CharacterMovementConsumer())
            .addConsumer(new KillMonsterConsumer())
            .addConsumer(new CharacterExperienceConsumer())
            .addConsumer(new CharacterLevelConsumer())
            .addConsumer(new AssignApConsumer())
            .addConsumer(new AssignSpConsumer())
            .addConsumer(new DropReservationEventConsumer())
            .addConsumer(new GainMesoConsumer())
            .addConsumer(new AdjustMesoConsumer())
            .addConsumer(new AdjustHealthConsumer())
            .addConsumer(new AdjustManaConsumer())
            .initialize();

      URI uri = UriBuilder.host(RestConstants.SERVICE).uri();
      RestServerFactory.create(uri, "com.atlas.cos.rest");
   }
}
