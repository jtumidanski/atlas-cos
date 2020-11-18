package com.atlas.cos.attribute;

import rest.AttributeResult;

public record JobAttributes(String name, Integer createIndex) implements AttributeResult {
}
