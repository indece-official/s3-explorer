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

all: collect_licenses build_linux build_windows

collect_licenses:
	rm -rf /tmp/s3-explorer-build-licenses
	cd ../backend/ && go-licenses save ./... --save_path /tmp/s3-explorer-build-licenses/ || true
	rm -rf /tmp/s3-explorer-build-licenses/github.com/indece-official
	OUTPUTFILE="$$PWD/assets/bin/licenses.zip" && \
	rm -f $$OUTPUTFILE && \
	cd /tmp/s3-explorer-build-licenses && \
	zip -r $$OUTPUTFILE ./*

build_windows:
	docker run --rm -i -v "$$PWD:/work" -e "BUILD_VERSION=${BUILD_VERSION}" amake/innosetup ./setup.iss

build_linux:
	tar -zcvf ./dist/s3-explorer_v${BUILD_VERSION}_linux64.tar.gz \
	-C ./assets/bin/ LICENSE.txt licenses.zip s3-explorer

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
	-out ./dist/s3-explorer-setup_v${BUILD_VERSION}_win64.exe