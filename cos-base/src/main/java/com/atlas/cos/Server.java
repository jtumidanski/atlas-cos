package com.atlas.cos;

import com.atlas.cos.event.consumer.ChangeMapCommandConsumer;
import com.atlas.cos.event.consumer.CharacterExperienceConsumer;
import com.atlas.cos.event.consumer.CharacterLevelConsumer;
import com.atlas.cos.event.consumer.CharacterMovementConsumer;
import com.atlas.cos.event.consumer.CharacterStatusConsumer;
import com.atlas.cos.event.consumer.KillMonsterConsumer;
import com.atlas.cos.processor.BlockedNameProcessor;
import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;
import database.PersistenceManager;

import java.net.URI;
import java.util.Arrays;
import java.util.List;

public class Server {
   public static void main(String[] args) {
      PersistenceManager.construct("atlas-cos");

      SimpleEventConsumerFactory.create(new ChangeMapCommandConsumer());
      SimpleEventConsumerFactory.create(new CharacterStatusConsumer());
      SimpleEventConsumerFactory.create(new CharacterMovementConsumer());
      SimpleEventConsumerFactory.create(new KillMonsterConsumer());
      SimpleEventConsumerFactory.create(new CharacterExperienceConsumer());
      SimpleEventConsumerFactory.create(new CharacterLevelConsumer());

      List<String> blockedNameList = Arrays.asList("admin", "owner", "moderator", "intern", "donor", "administrator", "FREDRICK",
            "help", "helper", "alert", "notice", "maplestory", "fuck", "wizet", "fucking", "negro", "fuk", "fuc", "penis", "pussy",
            "asshole", "gay", "nigger", "homo", "suck", "cum", "shit", "shitty", "condom", "security", "official", "rape", "nigga",
            "sex", "tit", "boner", "orgy", "clit", "asshole", "fatass", "bitch", "support", "gamemaster", "cock", "gaay", "gm",
            "operate", "master", "sysop", "party", "GameMaster", "community", "message", "event", "test", "meso", "Scania", "yata",
            "AsiaSoft", "henesys");
      BlockedNameProcessor.bulkAddBlockedNames(blockedNameList);

      URI uri = UriBuilder.host(RestService.CHARACTER).uri();
      RestServerFactory.create(uri, "com.atlas.cos.rest");
   }
}
