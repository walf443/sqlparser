
test: mysql_test

mysql_test:
	cd mysql && make test

.PHONY: test mysql_test
