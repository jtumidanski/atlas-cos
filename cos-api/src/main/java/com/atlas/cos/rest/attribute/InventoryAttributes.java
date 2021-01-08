package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record InventoryAttributes(String type, Integer capacity) implements AttributeResult {
}
