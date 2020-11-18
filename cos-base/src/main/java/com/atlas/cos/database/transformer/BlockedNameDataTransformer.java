package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

import transformer.SqlTransformer;

public class BlockedNameDataTransformer implements SqlTransformer<BlockedNameData, BlockedName> {
   @Override
   public BlockedNameData transform(BlockedName blockedName) {
      return new BlockedNameData(blockedName.getName());
   }
}
