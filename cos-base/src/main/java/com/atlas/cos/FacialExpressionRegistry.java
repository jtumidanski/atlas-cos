package com.atlas.cos;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

public class FacialExpressionRegistry {
   private static final Object lock = new Object();

   private static volatile FacialExpressionRegistry instance;

   private final Map<Integer, Integer> locks;

   private final Map<Integer, ScheduledFuture<?>> expressionTasks;

   public static FacialExpressionRegistry getInstance() {
      FacialExpressionRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new FacialExpressionRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private FacialExpressionRegistry() {
      locks = new ConcurrentHashMap<>();
      expressionTasks = new HashMap<>();
   }

   protected Object getLock(final Integer characterId) {
      locks.putIfAbsent(characterId, characterId);
      return locks.get(characterId);
   }

   public void registerChange(int characterId, Runnable returnToNormal) {
      synchronized (getLock(characterId)) {
         if (expressionTasks.containsKey(characterId)) {
            ScheduledFuture<?> future = expressionTasks.get(characterId);
            future.cancel(true);
            expressionTasks.remove(characterId);
         }

         Runnable wrappedRunnable = () -> {
            returnToNormal.run();
            synchronized (getLock(characterId)) {
               expressionTasks.remove(characterId);
            }
         };

         ScheduledExecutorService executor = Executors.newSingleThreadScheduledExecutor();
         ScheduledFuture<?> future = executor.schedule(wrappedRunnable, 5, TimeUnit.SECONDS);
         expressionTasks.put(characterId, future);
      }
   }

   public void cancelChange(int characterId) {
      if (expressionTasks.containsKey(characterId)) {
         synchronized (getLock(characterId)) {
            if (expressionTasks.containsKey(characterId)) {
               ScheduledFuture<?> future = expressionTasks.get(characterId);
               future.cancel(true);
               expressionTasks.remove(characterId);
            }
         }
      }
   }
}
