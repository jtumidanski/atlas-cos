package com.atlas.cos.attribute;

import rest.AttributeResult;

public record InventoryAttributes(String type, Integer capacity) implements AttributeResult {
}
