package com.atlas.cos.event.consumer;

import com.atlas.kafka.consumer.SimpleEventHandler;

public abstract class AbstractEventConsumer<T> implements SimpleEventHandler<T> {
   @Override
   public String getConsumerId() {
      return "Character Service";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }
}
