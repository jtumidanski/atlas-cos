package com.atlas.cos.attribute;

import rest.AttributeResult;

public record EquipmentAttributes(Integer equipmentId, Short slot) implements AttributeResult {
}
