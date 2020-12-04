package com.atlas.cos.processor;

import java.awt.*;
import java.util.Optional;

import com.atlas.cos.model.Portal;

public class MapProcessor {
   private static final Object lock = new Object();

   private static volatile MapProcessor instance;

   public static MapProcessor getInstance() {
      MapProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new MapProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public Optional<Integer> findClosestSpawnPoint(int mapId, int x, int y) {
      Point from = new Point(x, y);
      return PortalProcessor.getInstance()
            .getMapPortals(mapId).stream()
            .filter(this::isSpawnPoint)
            .min((o1, o2) -> compareDistanceFromPoint(from, o1, o2))
            .map(Portal::id);
   }

   protected boolean isSpawnPoint(Portal portal) {
      return (portal.type() == 0 || portal.type() == 1) && portal.targetMap() == 999999999;
   }

   protected int compareDistanceFromPoint(Point from, Portal o1, Portal o2) {
      double o1Distance = new Point(o1.x(), o1.y()).distanceSq(from);
      double o2Distance = new Point(o2.x(), o2.y()).distanceSq(from);
      return Double.compare(o1Distance, o2Distance);
   }
}
