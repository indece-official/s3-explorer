# S3 Explorer
# Copyright (C) 2020  indece UG (haftungsbeschränkt)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License or any
# later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

all: build

build:
	docker run --rm -i -v "$$PWD:/work" -e "BUILD_VERSION=${BUILD_VERSION}" amake/innosetup ./setup.iss

sign:
	# Sign installer
	test -n "${CODESIGN_CERT}"
	@test -n "${CODESIGN_PASSWORD}"
	@osslsigncode sign \
	-pkcs12 "${CODESIGN_CERT}" \
	-pass "${CODESIGN_PASSWORD}" \
	-n "Installer for S3 Explorer" \
	-i https://www.indece.com/ \
	-in ./dist/s3-explorer-setup_win64-unsigned.exe \
	-out ./dist/s3-explorer-setup_win64.exe