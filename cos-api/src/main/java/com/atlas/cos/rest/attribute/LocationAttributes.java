package v;

import rest.AttributeResult;

public record LocationAttributes(LocationType type, Integer mapId, Integer portalId) implements AttributeResult {
}
