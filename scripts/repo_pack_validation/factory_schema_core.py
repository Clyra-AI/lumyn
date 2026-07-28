from __future__ import annotations

import re
from datetime import datetime
import json
from pathlib import Path
from typing import Any, Callable, Optional

FailFn = Callable[[str], None]
LoadJsonFn = Callable[[Path], Any]

SUPPORTED_SCHEMA_KEYWORDS = {
    "$comment",
    "$defs",
    "$id",
    "$ref",
    "$schema",
    "additionalProperties",
    "allOf",
    "anyOf",
    "const",
    "contains",
    "default",
    "dependentRequired",
    "description",
    "else",
    "enum",
    "examples",
    "format",
    "if",
    "items",
    "maxContains",
    "maxItems",
    "maxLength",
    "maxProperties",
    "maximum",
    "minContains",
    "minItems",
    "minLength",
    "minProperties",
    "minimum",
    "not",
    "oneOf",
    "pattern",
    "properties",
    "propertyNames",
    "required",
    "then",
    "title",
    "type",
    "uniqueItems",
    "x-factory-cross-field-constraints",
}

RFC3339_DATE_TIME_RE = re.compile(
    r"^(?P<date>\d{4}-\d{2}-\d{2})[Tt]"
    r"(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})"
    r"(?P<fraction>\.\d+)?(?P<offset>[Zz]|[+-]\d{2}:\d{2})$"
)


def validate_supported_schema_keywords(schema: dict[str, Any], path: str, *, fail: FailFn) -> None:
    unsupported = sorted(set(schema) - SUPPORTED_SCHEMA_KEYWORDS)
    if unsupported:
        fail(f"{path}: unsupported schema keywords {unsupported!r}")

    for container_key in ("properties", "$defs"):
        container = schema.get(container_key, {})
        if container is None:
            continue
        if not isinstance(container, dict):
            fail(f"{path}.{container_key}: expected object")
        for name, subschema in container.items():
            if not isinstance(subschema, dict):
                fail(f"{path}.{container_key}.{name}: expected schema object")
            validate_supported_schema_keywords(subschema, f"{path}.{container_key}.{name}", fail=fail)

    for list_key in ("allOf", "anyOf", "oneOf"):
        values = schema.get(list_key, [])
        if values is None:
            continue
        if not isinstance(values, list):
            fail(f"{path}.{list_key}: expected array")
        for index, subschema in enumerate(values):
            if not isinstance(subschema, dict):
                fail(f"{path}.{list_key}[{index}]: expected schema object")
            validate_supported_schema_keywords(subschema, f"{path}.{list_key}[{index}]", fail=fail)

    for schema_key in ("additionalProperties", "contains", "else", "if", "items", "not", "propertyNames", "then"):
        subschema = schema.get(schema_key)
        if isinstance(subschema, dict):
            validate_supported_schema_keywords(subschema, f"{path}.{schema_key}", fail=fail)


def validate_type(expected: str, value: Any, path: str, *, fail: FailFn) -> None:
    if expected == "object" and not isinstance(value, dict):
        fail(f"{path}: expected object")
    if expected == "array" and not isinstance(value, list):
        fail(f"{path}: expected array")
    if expected == "string" and not isinstance(value, str):
        fail(f"{path}: expected string")
    if expected == "integer" and not (isinstance(value, int) and not isinstance(value, bool)):
        fail(f"{path}: expected integer")
    if expected == "number" and not (
        (isinstance(value, int) or isinstance(value, float)) and not isinstance(value, bool)
    ):
        fail(f"{path}: expected number")
    if expected == "boolean" and not isinstance(value, bool):
        fail(f"{path}: expected boolean")


def validate_format(expected: str, value: str, path: str, *, fail: FailFn) -> None:
    if expected != "date-time":
        fail(f"{path}: unsupported schema format {expected!r}")
    match = RFC3339_DATE_TIME_RE.fullmatch(value)
    if match is None:
        fail(f"{path}: expected RFC3339 date-time")
    second = match.group("second")
    if int(second) > 60:
        fail(f"{path}: expected RFC3339 date-time")
    normalized_second = "59" if second == "60" else second
    offset = match.group("offset")
    normalized_offset = "+00:00" if offset.lower() == "z" else offset
    candidate = (
        f"{match.group('date')}T{match.group('hour')}:{match.group('minute')}:"
        f"{normalized_second}{match.group('fraction') or ''}{normalized_offset}"
    )
    try:
        parsed = datetime.fromisoformat(candidate)
    except ValueError:
        fail(f"{path}: expected RFC3339 date-time")
    if parsed.tzinfo is None:
        fail(f"{path}: date-time requires timezone")


def resolve_json_pointer(
    root: dict[str, Any],
    pointer: str,
    ref: str,
    path: str,
    *,
    fail: FailFn,
) -> dict[str, Any]:
    if pointer in {"", "#"}:
        return root
    if not pointer.startswith("#/"):
        fail(f"{path}: unsupported schema $ref {ref!r}")
    target: Any = root
    for raw_token in pointer[2:].split("/"):
        token = raw_token.replace("~1", "/").replace("~0", "~")
        if not isinstance(target, dict) or token not in target:
            fail(f"{path}: unresolved schema $ref {ref!r}")
        target = target[token]
    if not isinstance(target, dict):
        fail(f"{path}: schema $ref {ref!r} does not resolve to an object")
    return target


def resolve_schema_ref(
    ref: str,
    root_schema: dict[str, Any],
    path: str,
    *,
    root: Path,
    fail: FailFn,
    load_json: LoadJsonFn,
) -> tuple[dict[str, Any], dict[str, Any]]:
    if ref.startswith("#"):
        return resolve_json_pointer(root_schema, ref, ref, path, fail=fail), root_schema

    schema_ref, _, pointer = ref.partition("#")
    if not schema_ref:
        fail(f"{path}: unsupported schema $ref {ref!r}")
    if re.match(r"^[A-Za-z][A-Za-z0-9+.-]*:", schema_ref):
        fail(f"{path}: external schema $ref {ref!r} is not supported")
    # Lumyn vendors the two pinned Factory schemas in one isolated directory.
    # Resolution remains relative and fail-closed, matching Factory semantics
    # without making a sibling Factory checkout a CI prerequisite.
    schema_path = (root / schema_ref).resolve()
    artifact_schema_dir = root.resolve()
    if artifact_schema_dir not in schema_path.parents or not schema_path.is_file():
        fail(f"{path}: schema $ref {ref!r} does not resolve to an artifact schema")
    target_root = load_json(schema_path)
    target_pointer = f"#{pointer}" if pointer else "#"
    return resolve_json_pointer(target_root, target_pointer, ref, path, fail=fail), target_root


def validate_schema(
    schema: dict[str, Any],
    value: Any,
    path: str = "$",
    root_schema: Optional[dict[str, Any]] = None,
    *,
    root: Path,
    fail: FailFn,
    load_json: LoadJsonFn,
    validation_error_type: type[Exception],
) -> None:
    if root_schema is None:
        root_schema = schema
        validate_supported_schema_keywords(schema, path, fail=fail)

    ref = schema.get("$ref")
    if ref is not None:
        if not isinstance(ref, str):
            fail(f"{path}: schema $ref must be a string")
        resolved_schema, resolved_root = resolve_schema_ref(
            ref,
            root_schema,
            path,
            root=root,
            fail=fail,
            load_json=load_json,
        )
        validate_schema(
            resolved_schema,
            value,
            path,
            resolved_root,
            root=root,
            fail=fail,
            load_json=load_json,
            validation_error_type=validation_error_type,
        )
        schema = {key: item for key, item in schema.items() if key != "$ref"}
        if not schema:
            return

    expected_type = schema.get("type")
    if expected_type:
        validate_type(expected_type, value, path, fail=fail)

    if "enum" in schema and value not in schema["enum"]:
        fail(f"{path}: value {value!r} not in enum {schema['enum']!r}")

    if "const" in schema and value != schema["const"]:
        fail(f"{path}: value {value!r} does not match const {schema['const']!r}")

    if "pattern" in schema and isinstance(value, str):
        pattern = schema["pattern"]
        if isinstance(pattern, str) and re.search(pattern, value) is None:
            fail(f"{path}: string does not match pattern {pattern!r}")

    if "minLength" in schema and isinstance(value, str):
        min_length = schema["minLength"]
        if isinstance(min_length, int) and len(value) < min_length:
            fail(f"{path}: string shorter than minLength {min_length!r}")

    if "maxLength" in schema and isinstance(value, str):
        max_length = schema["maxLength"]
        if isinstance(max_length, int) and len(value) > max_length:
            fail(f"{path}: string longer than maxLength {max_length!r}")

    if "format" in schema and isinstance(value, str):
        expected_format = schema["format"]
        if not isinstance(expected_format, str):
            fail(f"{path}: format must be a string")
        validate_format(expected_format, value, path, fail=fail)

    if "minimum" in schema and isinstance(value, (int, float)) and not isinstance(value, bool):
        minimum = schema["minimum"]
        if isinstance(minimum, (int, float)) and value < minimum:
            fail(f"{path}: value {value!r} below minimum {minimum!r}")

    if "maximum" in schema and isinstance(value, (int, float)) and not isinstance(value, bool):
        maximum = schema["maximum"]
        if isinstance(maximum, (int, float)) and value > maximum:
            fail(f"{path}: value {value!r} above maximum {maximum!r}")

    for subschema in schema.get("allOf", []):
        if not isinstance(subschema, dict):
            fail(f"{path}: allOf entries must be schemas")
        validate_schema(
            subschema,
            value,
            path,
            root_schema,
            root=root,
            fail=fail,
            load_json=load_json,
            validation_error_type=validation_error_type,
        )

    any_of = schema.get("anyOf")
    if any_of is not None:
        if not isinstance(any_of, list) or not any_of:
            fail(f"{path}: anyOf must be a non-empty list")
        any_errors = []
        for subschema in any_of:
            if not isinstance(subschema, dict):
                fail(f"{path}: anyOf entries must be schemas")
            try:
                validate_schema(
                    subschema,
                    value,
                    path,
                    root_schema,
                    root=root,
                    fail=fail,
                    load_json=load_json,
                    validation_error_type=validation_error_type,
                )
                break
            except validation_error_type as exc:
                any_errors.append(str(exc))
        else:
            fail(f"{path}: value did not match anyOf schemas: {any_errors!r}")

    one_of = schema.get("oneOf")
    if one_of is not None:
        if not isinstance(one_of, list) or not one_of:
            fail(f"{path}: oneOf must be a non-empty list")
        matched = sum(
            1
            for subschema in one_of
            if isinstance(subschema, dict)
            and schema_matches(
                subschema,
                value,
                path,
                root_schema,
                root=root,
                fail=fail,
                load_json=load_json,
                validation_error_type=validation_error_type,
            )
        )
        if matched != 1:
            fail(f"{path}: value must match exactly one oneOf schema, matched {matched}")

    not_schema = schema.get("not")
    if not_schema is not None:
        if not isinstance(not_schema, dict):
            fail(f"{path}: not must be a schema")
        if schema_matches(
            not_schema,
            value,
            path,
            root_schema,
            root=root,
            fail=fail,
            load_json=load_json,
            validation_error_type=validation_error_type,
        ):
            fail(f"{path}: value matched forbidden schema")

    if_schema = schema.get("if")
    if if_schema is not None:
        if not isinstance(if_schema, dict):
            fail(f"{path}: if must be a schema")
        matched = schema_matches(
            if_schema,
            value,
            path,
            root_schema,
            root=root,
            fail=fail,
            load_json=load_json,
            validation_error_type=validation_error_type,
        )
        branch = schema.get("then") if matched else schema.get("else")
        if branch is not None:
            if not isinstance(branch, dict):
                fail(f"{path}: conditional branch must be a schema")
            validate_schema(
                branch,
                value,
                path,
                root_schema,
                root=root,
                fail=fail,
                load_json=load_json,
                validation_error_type=validation_error_type,
            )

    objectish = expected_type == "object" or (
        isinstance(value, dict)
        and any(key in schema for key in ["required", "dependentRequired", "properties", "additionalProperties"])
    )
    if objectish:
        min_properties = schema.get("minProperties")
        if isinstance(min_properties, int) and len(value) < min_properties:
            fail(f"{path}: expected at least {min_properties} properties")
        max_properties = schema.get("maxProperties")
        if isinstance(max_properties, int) and len(value) > max_properties:
            fail(f"{path}: expected at most {max_properties} properties")
        required = schema.get("required", [])
        for key in required:
            if key not in value:
                fail(f"{path}: missing required key {key!r}")
        for key, dependent_keys in schema.get("dependentRequired", {}).items():
            if key not in value:
                continue
            if not isinstance(dependent_keys, list):
                fail(f"{path}: dependentRequired for {key!r} must be an array")
            for dependent_key in dependent_keys:
                if dependent_key not in value:
                    fail(f"{path}: key {key!r} requires dependent key {dependent_key!r}")
        for key, subschema in schema.get("properties", {}).items():
            if key in value:
                validate_schema(
                    subschema,
                    value[key],
                    f"{path}.{key}",
                    root_schema,
                    root=root,
                    fail=fail,
                    load_json=load_json,
                    validation_error_type=validation_error_type,
                )
        allowed = set(schema.get("properties", {}).keys())
        extras = set(value.keys()) - allowed
        additional_properties = schema.get("additionalProperties")
        if additional_properties is False:
            if extras:
                fail(f"{path}: unexpected keys {sorted(extras)!r}")
        elif isinstance(additional_properties, dict):
            for key in extras:
                validate_schema(
                    additional_properties,
                    value[key],
                    f"{path}.{key}",
                    root_schema,
                    root=root,
                    fail=fail,
                    load_json=load_json,
                    validation_error_type=validation_error_type,
                )
        property_names = schema.get("propertyNames")
        if isinstance(property_names, dict):
            for key in value:
                validate_schema(
                    property_names,
                    key,
                    f"{path}.<property-name>",
                    root_schema,
                    root=root,
                    fail=fail,
                    load_json=load_json,
                    validation_error_type=validation_error_type,
                )

    arrayish = expected_type == "array" or (
        isinstance(value, list) and any(key in schema for key in ["minItems", "maxItems", "items", "contains"])
    )
    if arrayish:
        min_items = schema.get("minItems")
        if isinstance(min_items, int) and len(value) < min_items:
            fail(f"{path}: expected at least {min_items} items")
        max_items = schema.get("maxItems")
        if isinstance(max_items, int) and len(value) > max_items:
            fail(f"{path}: expected at most {max_items} items")
        contains_schema = schema.get("contains")
        if contains_schema is not None:
            if not isinstance(contains_schema, dict):
                fail(f"{path}: contains must be a schema")
            matched_count = sum(
                1
                for item in value
                if schema_matches(
                    contains_schema,
                    item,
                    path,
                    root_schema,
                    root=root,
                    fail=fail,
                    load_json=load_json,
                    validation_error_type=validation_error_type,
                )
            )
            min_contains = schema.get("minContains", 1)
            max_contains = schema.get("maxContains")
            if not isinstance(min_contains, int) or min_contains < 0:
                fail(f"{path}: minContains must be a non-negative integer")
            if matched_count < min_contains:
                fail(f"{path}: expected at least {min_contains} items matching contains schema")
            if isinstance(max_contains, int) and matched_count > max_contains:
                fail(f"{path}: expected at most {max_contains} items matching contains schema")
        if schema.get("uniqueItems") is True:
            canonical = [json.dumps(item, sort_keys=True, separators=(",", ":")) for item in value]
            if len(canonical) != len(set(canonical)):
                fail(f"{path}: expected unique array items")
        item_schema = schema.get("items")
        if item_schema is None:
            return
        for index, item in enumerate(value):
            validate_schema(
                item_schema,
                item,
                f"{path}[{index}]",
                root_schema,
                root=root,
                fail=fail,
                load_json=load_json,
                validation_error_type=validation_error_type,
            )


def schema_matches(
    schema: dict[str, Any],
    value: Any,
    path: str = "$",
    root_schema: Optional[dict[str, Any]] = None,
    *,
    root: Path,
    fail: FailFn,
    load_json: LoadJsonFn,
    validation_error_type: type[Exception],
) -> bool:
    try:
        validate_schema(
            schema,
            value,
            path,
            root_schema,
            root=root,
            fail=fail,
            load_json=load_json,
            validation_error_type=validation_error_type,
        )
        return True
    except validation_error_type:
        return False
