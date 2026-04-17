#!/bin/sh
set -eu

TEMPLATE_FILE="/app/etc/chat.yaml.template"
OUTPUT_FILE="/app/etc/chat.yaml"

cp "$TEMPLATE_FILE" "$OUTPUT_FILE"

replace_or_fail() {
  key="$1"
  value="$2"
  if [ -z "$value" ]; then
    echo "Missing required environment variable for $key" >&2
    exit 1
  fi
  escaped_value=$(printf '%s' "$value" | sed 's/[\/&]/\\&/g')
  sed -i "s|\${$key}|$escaped_value|g" "$OUTPUT_FILE"
}

replace_allow_empty() {
  key="$1"
  value="$2"
  escaped_value=$(printf '%s' "$value" | sed 's/[\/&]/\\&/g')
  sed -i "s|\${$key}|$escaped_value|g" "$OUTPUT_FILE"
}

replace_or_fail OPENAI_API_KEY "${OPENAI_API_KEY:-}"
replace_or_fail OPENAI_BASE_URL "${OPENAI_BASE_URL:-}"
replace_or_fail OPENAI_MODEL "${OPENAI_MODEL:-}"
replace_or_fail VECTOR_DB_HOST "${VECTOR_DB_HOST:-}"
replace_or_fail VECTOR_DB_PORT "${VECTOR_DB_PORT:-}"
replace_or_fail VECTOR_DB_NAME "${VECTOR_DB_NAME:-}"
replace_or_fail VECTOR_DB_USER "${VECTOR_DB_USER:-}"
replace_or_fail VECTOR_DB_PASSWORD "${VECTOR_DB_PASSWORD:-}"
replace_or_fail VECTOR_DB_TABLE "${VECTOR_DB_TABLE:-}"
replace_or_fail VECTOR_DB_MAX_CONN "${VECTOR_DB_MAX_CONN:-}"
replace_or_fail VECTOR_DB_EMBEDDING_MODEL "${VECTOR_DB_EMBEDDING_MODEL:-}"
replace_or_fail VECTOR_DB_KNOWLEDGE_MAX_CHUNK_SIZE "${VECTOR_DB_KNOWLEDGE_MAX_CHUNK_SIZE:-}"
replace_or_fail VECTOR_DB_KNOWLEDGE_TOP_K "${VECTOR_DB_KNOWLEDGE_TOP_K:-}"
replace_or_fail VECTOR_DB_KNOWLEDGE_MAX_CONTEXT_LENGTH "${VECTOR_DB_KNOWLEDGE_MAX_CONTEXT_LENGTH:-}"
replace_or_fail MCP_ENDPOINT "${MCP_ENDPOINT:-}"
replace_or_fail REDIS_HOST "${REDIS_HOST:-}"
replace_or_fail REDIS_PORT "${REDIS_PORT:-}"
replace_allow_empty REDIS_PASSWORD "${REDIS_PASSWORD-}"
replace_or_fail REDIS_DB "${REDIS_DB:-}"

replace_or_fail OPENAI_MAX_TOKENS "${OPENAI_MAX_TOKENS:-2048}"
replace_or_fail OPENAI_TEMPERATURE "${OPENAI_TEMPERATURE:-0.7}"

exec /app/api -f /app/etc/chat.yaml
