package com.atlas.cos.database.administrator;

import javax.persistence.EntityManager;

import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

import accessor.AbstractQueryExecutor;

public class BlockedNameAdministrator extends AbstractQueryExecutor {
   private static final Object lock = new Object();

   private static volatile BlockedNameAdministrator instance;

   public static BlockedNameAdministrator getInstance() {
      BlockedNameAdministrator result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new BlockedNameAdministrator();
               instance = result;
            }
         }
      }
      return result;
   }

   public BlockedNameData createBlockedName(EntityManager entityManager, String name) {
      BlockedName blockedName = new BlockedName();
      blockedName.setName(name);
      return insertAndReturn(entityManager, blockedName, result -> new BlockedNameData(result.getName()));
   }
}
