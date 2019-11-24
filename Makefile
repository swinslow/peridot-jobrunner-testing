# SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

test: FORCE
	docker-compose up --abort-on-container-exit

clean:
	docker-compose down

build:
	docker-compose build

FORCE:
