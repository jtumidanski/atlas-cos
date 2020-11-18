package com.atlas.cos.database.provider;

import java.util.List;
import javax.persistence.EntityManager;

import com.app.database.provider.NamedQueryClient;
import com.atlas.cos.database.transformer.BlockedNameDataTransformer;
import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

public class BlockedNameProvider {
   private BlockedNameProvider() {
   }

   public static List<BlockedNameData> getNames(EntityManager entityManager) {
      return new NamedQueryClient<>(entityManager, BlockedName.GET_ALL, BlockedName.class)
            .list(new BlockedNameDataTransformer());
   }
}
