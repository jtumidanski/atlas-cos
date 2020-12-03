package com.atlas.cos;

import java.net.URI;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Executors;

import com.atlas.cos.command.ChangeMapCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.consumer.ChangeMapCommandConsumer;
import com.atlas.cos.event.consumer.CharacterLoggedInConsumer;
import com.atlas.cos.processor.BlockedNameProcessor;
import com.atlas.csrv.event.CharacterLoggedInEvent;
import com.atlas.kafka.consumer.ConsumerBuilder;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import database.PersistenceManager;

public class Server {
   public static void main(String[] args) {
      PersistenceManager.construct("atlas-cos");

      List<String> blockedNameList = Arrays.asList("admin", "owner", "moderator", "intern", "donor", "administrator", "FREDRICK",
            "help", "helper", "alert", "notice", "maplestory", "fuck", "wizet", "fucking", "negro", "fuk", "fuc", "penis", "pussy",
            "asshole", "gay", "nigger", "homo", "suck", "cum", "shit", "shitty", "condom", "security", "official", "rape", "nigga",
            "sex", "tit", "boner", "orgy", "clit", "asshole", "fatass", "bitch", "support", "gamemaster", "cock", "gaay", "gm",
            "operate", "master", "sysop", "party", "GameMaster", "community", "message", "event", "test", "meso", "Scania", "yata",
            "AsiaSoft", "henesys");
      BlockedNameProcessor.getInstance().bulkAddBlockedNames(blockedNameList);

      URI uri = UriBuilder.host(RestService.CHARACTER).uri();
      RestServerFactory.create(uri, "com.atlas.cos.rest");

      Executors.newSingleThreadExecutor().execute(
            new ConsumerBuilder<>("Character Service", ChangeMapCommand.class)
                  .setBootstrapServers(System.getenv("BOOTSTRAP_SERVERS"))
                  .setTopic(System.getenv(EventConstants.TOPIC_CHANGE_MAP_COMMAND))
                  .setHandler(new ChangeMapCommandConsumer())
                  .build()
      );

      Executors.newSingleThreadExecutor().execute(
            new ConsumerBuilder<>("Character Service", CharacterLoggedInEvent.class)
                  .setBootstrapServers(System.getenv("BOOTSTRAP_SERVERS"))
                  .setTopic(System.getenv(com.atlas.csrv.constant.EventConstants.TOPIC_CHARACTER_LOGIN))
                  .setHandler(new CharacterLoggedInConsumer())
                  .build()
      );
   }
}
