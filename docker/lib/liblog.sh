#!/bin/bash

# ☔ Arisu: Translation made with simplicity, yet robust.
# Copyright (C) 2020-2022 Noelware
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

BLUE='\033[38;2;81;81;140m'
GREEN='\033[38;2;165;204;165m'
PINK='\033[38;2;241;204;209m'
RESET='\033[0m'
BOLD='\033[1m'
UNDERLINE='\033[4m'

info() {
  timestamp=$(date +"%D ~ %r")
  printf "%b\\n" "${GREEN}${BOLD}info${RESET}  | ${PINK}${BOLD}${timestamp}${RESET} ~ $1"
}

debug() {
  local debug="${TSUBAKI_DEBUG:-false}"
  shopt -s nocasematch
  timestamp=$(date +"%D ~%r")

  if ! [[ "$debug" = "1" || "$debug" =~ ^(no|false)$ ]]; then
    printf "%b\\n" "${BLUE}${BOLD}debug${RESET} | ${PINK}${BOLD}${timestamp}${RESET} ~ $1"
  fi
}