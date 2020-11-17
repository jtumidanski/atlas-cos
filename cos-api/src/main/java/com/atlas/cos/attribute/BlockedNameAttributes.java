package com.atlas.cos.attribute;

import rest.AttributeResult;

public class BlockedNameAttributes implements AttributeResult {
   private String name;

   public String getName() {
      return name;
   }

   public void setName(String name) {
      this.name = name;
   }
}
