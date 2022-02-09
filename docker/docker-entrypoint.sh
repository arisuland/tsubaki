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

set -o errexit
set -o nounset
set -o pipefail

. /app/arisu/tsubaki/lib/liblog.sh

if ! [[ "${TSUBAKI_ENABLE_WELCOME_PROMPT:-yes}" =~ ^(no|false)$ ]]; then
    info ""
    info "   Welcome to the ${BOLD}Tsubaki${RESET} Docker container image!"
    info "   Core infrastructure + heart and soul of Arisu (https://arisu.land)"
    info ""
    info "   * Subscribe to the project for more updates: https://github.com/arisuland/tsubaki"
    info "   * Any issues or bugs occur while using? Report it here: https://github.com/arisuland/tsubaki/issues"
    info ""
fi

exec "$@"
