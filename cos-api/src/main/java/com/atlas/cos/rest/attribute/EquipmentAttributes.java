package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record EquipmentAttributes(Integer equipmentId, Short slot) implements AttributeResult {
}
