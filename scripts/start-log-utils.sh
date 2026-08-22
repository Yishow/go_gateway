#!/usr/bin/env bash

# dashboard 輪詢的 datalink 實體清單（單一來源，兩處比對共用；與 start-log-utils.ps1 對齊）
DATALINK_DASHBOARD_ENTITIES='devices|polling-groups|points|mappings|tags'

get_log_noise_category() {
  local line="$1"
  if [[ "$line" =~ ^watching|^building\.\.\.|^\!exclude|^[[:space:]]*/|^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
    printf "air-watcher"
    return
  fi
  if [[ "$line" == *"/api/v1/datalink/modbus-share/status"* ]]; then
    printf "status-polling"
    return
  fi
  if [[ "$line" =~ /api/v1/datalink/(${DATALINK_DASHBOARD_ENTITIES}) ]]; then
    printf "dashboard-refresh"
    return
  fi
  if [[ "$line" == *"資料庫路徑"* || "$line" == *"Executing SQLite migration"* || "$line" == *"ConnectionManager 已初始化"* || "$line" == *"已註冊的協議"* ]]; then
    printf "startup-detail"
    return
  fi
  printf ""
}

runtime_log_line() {
  local line="$1"
  [[ -z "$line" ]] && return

  if [[ "$line" =~ \[([^\]]+)\][[:space:]]+[^[:space:]]+[[:space:]]+(GET|POST|PUT|DELETE|PATCH)[[:space:]]+([^[:space:]]+)[[:space:]]+([0-9]{3})[[:space:]]+(.+) ]]; then
    local ts="${BASH_REMATCH[1]}"
    local method="${BASH_REMATCH[2]}"
    local path="${BASH_REMATCH[3]}"
    local status="${BASH_REMATCH[4]}"
    local latency="${BASH_REMATCH[5]}"
    local time_part
    if [[ "$ts" =~ ([0-9]{2}:[0-9]{2}:[0-9]{2})$ ]]; then
      time_part="${BASH_REMATCH[1]}"
    else
      time_part="--:--:--"
    fi

    if [[ "$method" == "GET" && "$status" == "200" && "$path" =~ ^/api/v1/datalink/(${DATALINK_DASHBOARD_ENTITIES}|modbus-share/status)$ ]]; then
      if [[ "$VERBOSE" != true ]]; then
        LOG_NOISE_COUNTERS["dashboard-refresh"]=$(( ${LOG_NOISE_COUNTERS["dashboard-refresh"]:-0} + 1 ))
        return
      fi
    fi

    local formatted
    formatted=$(printf "[HTTP] %s  %-6s %s %3s %9s" "$time_part" "$method" "$path" "$status" "$latency")
    if (( status >= 500 )); then
      printf "\033[31m%s\033[0m\n" "$formatted"
    elif (( status >= 400 )); then
      printf "\033[33m%s\033[0m\n" "$formatted"
    else
      printf "\033[90m%s\033[0m\n" "$formatted"
    fi
    return
  fi

  if [[ "$line" =~ ^[0-9]{4}/[0-9]{2}/[0-9]{2}[[:space:]]+[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+[^:]+:[0-9]+:[[:space:]]+(.+)$ ]]; then
    local msg="${BASH_REMATCH[1]}"
    local category
    category="$(get_log_noise_category "$line")"
    if [[ -n "$category" && "$VERBOSE" != true ]]; then
      LOG_NOISE_COUNTERS["$category"]=$(( ${LOG_NOISE_COUNTERS["$category"]:-0} + 1 ))
      return
    fi
    if [[ "$msg" == *"127.0.0.1:5020"* ]]; then
      printf "\033[32m[BOOT] Local Modbus share started on 127.0.0.1:5020\033[0m\n"
    elif [[ "$msg" == *"測試工具伺服器啟動於 "* ]]; then
      local startup_url="${msg##*啟動於 }"
      printf "\033[32m[BOOT] Server started at %s\033[0m\n" "$startup_url"
    elif [[ "$msg" == *"資料庫路徑"* ]]; then
      printf "\033[90m[BOOT] Database initialized\033[0m\n"
    elif [[ "$VERBOSE" == true ]]; then
      printf "\033[90m[BOOT] %s\033[0m\n" "$msg"
    fi
    return
  fi

  local category
  category="$(get_log_noise_category "$line")"
  if [[ -n "$category" && "$VERBOSE" != true ]]; then
    LOG_NOISE_COUNTERS["$category"]=$(( ${LOG_NOISE_COUNTERS["$category"]:-0} + 1 ))
    return
  fi

  if [[ "$line" =~ ERROR|Error|panic|FATAL ]]; then
    printf "\033[31m%s\033[0m\n" "$line"
  elif [[ "$line" =~ WARN|Warning ]]; then
    printf "\033[33m%s\033[0m\n" "$line"
  elif [[ "$line" =~ 啟動於|本機\ Modbus\ 分享服務已啟動 ]]; then
    printf "\033[32m%s\033[0m\n" "$line"
  else
    printf "\033[90m%s\033[0m\n" "$line"
  fi
}

show_log_noise_summary() {
  [[ ${#LOG_NOISE_COUNTERS[@]} -eq 0 ]] && return
  [[ "$VERBOSE" == true ]] && return
  printf "\n\033[36m============================================\033[0m\n"
  printf "\033[36m   已隱藏雜訊日誌\033[0m\n"
  printf "\033[36m============================================\033[0m\n"
  for key in "${!LOG_NOISE_COUNTERS[@]}"; do
    printf "\033[33m- %s: %s 行\033[0m\n" "$key" "${LOG_NOISE_COUNTERS[$key]}"
  done
  info "可加上 --verbose 顯示全部原始日誌。"
}
