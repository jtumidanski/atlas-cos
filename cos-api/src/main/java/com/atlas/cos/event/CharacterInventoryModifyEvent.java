package com.atlas.cos.event;

import java.util.List;

public record CharacterInventoryModifyEvent(int characterId, boolean updateTick, List<InventoryModification> modifications) {
}
