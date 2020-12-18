package com.atlas.cos.attribute;

import rest.AttributeResult;

public record ItemAttributes(Integer itemId, Integer quantity, Short slot) implements AttributeResult {
}
