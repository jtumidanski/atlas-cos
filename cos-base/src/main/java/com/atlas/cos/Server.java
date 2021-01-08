package com.atlas.cos;

import java.net.URI;
import java.util.Arrays;
import java.util.List;

import com.atlas.cos.constant.RestConstants;
import com.atlas.cos.event.consumer.AdjustHealthConsumer;
import com.atlas.cos.event.consumer.AdjustManaConsumer;
import com.atlas.cos.event.consumer.AdjustMesoConsumer;
import com.atlas.cos.event.consumer.AssignApConsumer;
import com.atlas.cos.event.consumer.AssignSpConsumer;
import com.atlas.cos.event.consumer.ChangeMapCommandConsumer;
import com.atlas.cos.event.consumer.CharacterExperienceConsumer;
import com.atlas.cos.event.consumer.CharacterExpressionCommandConsumer;
import com.atlas.cos.event.consumer.CharacterLevelConsumer;
import com.atlas.cos.event.consumer.CharacterMovementConsumer;
import com.atlas.cos.event.consumer.CharacterStatusConsumer;
import com.atlas.cos.event.consumer.DropReservationEventConsumer;
import com.atlas.cos.event.consumer.GainMesoConsumer;
import com.atlas.cos.event.consumer.KillMonsterConsumer;
import com.atlas.cos.processor.BlockedNameProcessor;
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
            .addConsumer(new CharacterExpressionCommandConsumer())
            .addConsumer(new AdjustMesoConsumer())
            .addConsumer(new AdjustHealthConsumer())
            .addConsumer(new AdjustManaConsumer())
            .initialize();

      List<String> blockedNameList = Arrays.asList("admin", "owner", "moderator", "intern", "donor", "administrator", "FREDRICK",
            "help", "helper", "alert", "notice", "maplestory", "fuck", "wizet", "fucking", "negro", "fuk", "fuc", "penis", "pussy",
            "asshole", "gay", "nigger", "homo", "suck", "cum", "shit", "shitty", "condom", "security", "official", "rape", "nigga",
            "sex", "tit", "boner", "orgy", "clit", "asshole", "fatass", "bitch", "support", "gamemaster", "cock", "gaay", "gm",
            "operate", "master", "sysop", "party", "GameMaster", "community", "message", "event", "test", "meso", "Scania", "yata",
            "AsiaSoft", "henesys");
      BlockedNameProcessor.bulkAddBlockedNames(blockedNameList);

      URI uri = UriBuilder.host(RestConstants.SERVICE).uri();
      RestServerFactory.create(uri, "com.atlas.cos.rest");
   }
}
