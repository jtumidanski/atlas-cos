package com.atlas.cos.database.provider;

import java.util.List;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;

import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

import accessor.AbstractQueryExecutor;

public class BlockedNameProvider extends AbstractQueryExecutor {
   private static final Object lock = new Object();

   private static volatile BlockedNameProvider instance;

   public static BlockedNameProvider getInstance() {
      BlockedNameProvider result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new BlockedNameProvider();
               instance = result;
            }
         }
      }
      return result;
   }

   public List<BlockedNameData> getNames(EntityManager entityManager) {
      TypedQuery<BlockedName> query = entityManager.createQuery("SELECT b FROM BlockedName b", BlockedName.class);
      return getResultList(query, blockedName -> new BlockedNameData(blockedName.getName()));
   }
}
