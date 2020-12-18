package com.atlas.cos.model;

import java.util.List;

public record Inventory(byte id, String type, int capacity, List<InventoryItem> items) {
}
